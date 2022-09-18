package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/johnnchung/HackTheNorth2022/auth"
	"github.com/johnnchung/HackTheNorth2022/helpers"
	models "github.com/johnnchung/HackTheNorth2022/models"
	"github.com/rs/cors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/urfave/negroni"
)

var (
	YTDL_TIME_FORMAT = "15:04:05.00"
)

// ********************
// Setup stuff
// ********************

type Repo struct {
	muxClient *mux.Router
	db        *models.Db
}

type textReq struct {
	Text string `json:"text"`
}

// ********************
// Setup stuff
// ********************
func (r *Repo) HandlerInit() error {
	r.muxClient = mux.NewRouter()

	// connect DB
	r.db = models.NewDbInstance()
	if err := r.db.DbInit(); err != nil {
		return err
	}

	// add routes
	apiRoutes := mux.NewRouter()
	apiRoutes.HandleFunc("/api/internal/{term}", r.addDataFromCategoryHandler)

	// Core endpoint
	r.muxClient.HandleFunc("/api/v1/process", r.getVideoFromText)

	// add middlewares
	authMiddleware := negroni.HandlerFunc(auth.ReqAPIKey)

	// combine
	r.muxClient.PathPrefix("/api/internal").Handler(negroni.New(
		authMiddleware,
		negroni.Wrap(apiRoutes),
	))

	return nil
}

func (r *Repo) Run() error {
	n := negroni.Classic() // Includes some default middlewares

	n.UseHandler(r.muxClient)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},           // All origins
		AllowedMethods: []string{"GET", "POST"}, // Allowing only get, just an example
	})

	if err := http.ListenAndServe(":8080", c.Handler(n)); err != nil {
		return err
	}

	return nil
}

// ********************
// Handler stuff
// ********************

func (r *Repo) addDataFromCategoryHandler(w http.ResponseWriter, req *http.Request) {
	// gets a word
	vars := mux.Vars(req)

	// query youtube API for word
	videoIds, err := r.db.GetSearchFromYoutube(vars["term"])
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadGateway, err.Error())
		return
	}

	// send them to AssemblyAI; get a list of transcribed word with timestamps
	for _, id := range videoIds {
		// get words
		ytRsp, err := r.db.GetTranscribedYtVid(id)
		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// // add words as rows to Cockroach
		for _, word := range ytRsp.Words {
			word.Text = helpers.CleanPunctuation([]byte(word.Text))
			if err := r.db.InsertWords(word, ytRsp.URL); err != nil {
				helpers.RespondWithError(w, http.StatusBadGateway, err.Error())
				return
			}
		}

	}

	helpers.RespondWithJSON(w, http.StatusAccepted, nil)
	return

}

func (r *Repo) getVideoFromText(w http.ResponseWriter, req *http.Request) {

	// get text
	var reqBody textReq
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "")
		return
	}

	var ytLink []models.DbResp

	for _, word := range strings.Fields(reqBody.Text) {
		// get rid of punctuation
		cleanWord := helpers.CleanPunctuation([]byte(word))
		// lookup in DB => get results
		links, err := r.db.GetLinkFromWord(cleanWord)
		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// pick one randomly
		linksLen := len(links)
		choice := rand.Intn(linksLen)
		fmt.Println(choice, links[choice], links)

		// append to ytLink
		ytLink = append(ytLink, links[choice])

	}

	var videoList []string
	fmt.Println("Got all ytLinks ready: ", ytLink)

	for i, link := range ytLink {

		fmt.Println("Began processing link: ", link)

		// convert milliseconds to time Format
		var t time.Time
		start_time := t.Add(time.Duration(link.Start) * time.Millisecond).Format(YTDL_TIME_FORMAT)
		t = time.Time{}
		end_time := t.Add(time.Duration(link.End-link.Start) * time.Millisecond).Format(YTDL_TIME_FORMAT)

		// download video
		outFileName := fmt.Sprintf("temp/%d.mp4", i)
		cmd := exec.Command("yt-dlp", "-S", "res,ext:mp4:m4a", "--recode", "mp4", "-o", outFileName, link.Url)
		fmt.Println(cmd)

		if err := cmd.Run(); err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("Editing file....")
		// cut file``
		cutFileName := fmt.Sprintf("temp/cut_%d.mp4", i)
		if err := ffmpeg.Input(outFileName, ffmpeg.KwArgs{"ss": start_time}).
			Output(cutFileName, ffmpeg.KwArgs{"t": end_time}).
			OverWriteOutput().Run(); err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// append to list
		videoList = append(videoList, cutFileName)
	}

	fmt.Println("Stitching all videos...")

	f, err := os.Create("list.txt")
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer f.Close()
	for _, val := range videoList {
		_, err := f.WriteString(fmt.Sprintf("file '%s'\n", val))
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "list.txt", "-c", "copy", "final.mp4")
	if err := cmd.Run(); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// push it to cloudshiney
	url, err := r.db.UploadToCloudinary("final.mp4")
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	os.RemoveAll("temp")
	os.Remove("final.mp4")
	os.Remove("list.txt")

	// return link
	helpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"video_url": url})

}
