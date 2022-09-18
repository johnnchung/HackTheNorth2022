package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/johnnchung/HackTheNorth2022/auth"
	"github.com/johnnchung/HackTheNorth2022/helpers"
	models "github.com/johnnchung/HackTheNorth2022/models"
	"github.com/urfave/negroni"
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

	n.Run(":8080")

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

func (r *Repo) getLinksFromText(w http.ResponseWriter, req *http.Request) {

	// get text
	var reqBody textReq
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "")
		return
	}

	var ytLink []*models.DbResp

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

		// append to ytLink
		ytLink = append(ytLink, links[choice])

	}

	// for _, link := range ytLink {
	// 	// download video
	// 	// convert milliseconds to seconds
	// 	// cut video
	// 	// add to folder
	// }

	// combine video

	// push it to cloudshiney

	// return link

}
