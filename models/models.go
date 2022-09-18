package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	_ "github.com/lib/pq"
)

var (
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
)

// *********************
//	Structs
// *********************

type Db struct {
	db        *sql.DB
	ytClient  *youtube.Service
	cldClient *cloudinary.Cloudinary
}

// AssemblyAPI->Youtube
type YtResp struct {
	Words  []YtWord `json:"words"`
	Status string   `json:"status"`
	URL    string   `json:"url,omit_empty"`
}

type YtWord struct {
	End   int64  `json:"end"`
	Start int64  `json:"start"`
	Text  string `json:"text"`
}

type DbResp struct {
	End   int64  `json:"end"`
	Start int64  `json:"start"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

// *********************
// Functions
// *********************

func NewDbInstance() *Db {
	return &Db{}
}

func (d *Db) DbInit() error {
	// setup DB
	connStr := os.Getenv("CONNECTION_STRING")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	d.db = db

	// setup YT
	ytDevKey := os.Getenv("YOUTUBE_API_KEY")
	client := &http.Client{
		Transport: &transport.APIKey{Key: ytDevKey},
	}

	ytService, err := youtube.New(client)
	if err != nil {
		return fmt.Errorf("Error creating new YouTube client: %v", err)
	}
	d.ytClient = ytService

	// setup cloudinary
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLD_NAME"), os.Getenv("CLD_KEY"), os.Getenv("CLD_SECRET"))
	d.cldClient = cld

	return nil
}

// **********************
// CockroachDB Endpoints
// **********************

func (d *Db) InsertWords(audioWord YtWord, url string) error {
	if _, err := d.db.Exec(
		`INSERT INTO ytdata (start_sec, end_sec, word, url) VALUES ($1, $2, $3, $4)`,
		audioWord.Start,
		audioWord.End,
		audioWord.Text,
		url,
	); err != nil {
		return err
	}

	return nil
}

func (d *Db) GetWords(word string) ([]*DbResp, error) {
	var result []*DbResp

	rows, err := d.db.Query(
		`SELECT text FROM ytdata WHERE text=$1;`,
		word,
	)
	if err != nil {
		return nil, err
	}

	tmpRow := &DbResp{}
	for rows.Next() {
		err := rows.Scan(tmpRow)
		if err != nil {
			return nil, err
		}

		result = append(result, tmpRow)
	}
	return result, nil
}

// word -> list of links

// func wordToLinks(word string) []string {

// 	return
// }

// url -> download videos

//

// **********************
//	API Handlers
// **********************

func (d *Db) GetTranscribedYtVid(ytID string) (*YtResp, error) {
	// 1.convert to audio file and download
	ytUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", ytID)
	cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", ytUrl)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 1.1 Find mp3
	filename := findMP3(".")

	// 2.upload to AssemblyAI
	resp, err := d.cldClient.Upload.Upload(context.Background(), filename, uploader.UploadParams{})
	if err != nil {
		return nil, err
	}

	// delete file
	os.Remove(filename)

	// setup url and body
	values := map[string]string{"audio_url": fmt.Sprintf("%s", resp.URL)}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	// Call requeust
	client := http.Client{}
	req, _ := http.NewRequest("POST", os.Getenv("TRANSCRIPT_URL"), bytes.NewBuffer(jsonData))
	req.Header.Set("authorization", os.Getenv("ASSEMBLY_API_KEY"))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	// get ytWords
	ytRsp, err := d.getAudioFileFromID(fmt.Sprintf("%s", result["id"]))
	if err != nil {
		return nil, err
	}

	ytRsp.URL = ytUrl

	return ytRsp, nil
}

func (d *Db) getAudioFileFromID(audioFileID string) (*YtResp, error) {
	var POLLING_URL = os.Getenv("TRANSCRIPT_URL") + "/" + audioFileID

	var res *http.Response
	var result YtResp
	var err error

	client := http.Client{}
	req, _ := http.NewRequest("GET", POLLING_URL, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", os.Getenv("ASSEMBLY_API_KEY"))

	// poll
	for {
		fmt.Println("Polling for Completed Status")
		res, err = client.Do(req)
		if err != nil {
			return nil, err
		}

		// decode body
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return nil, err
		}

		if result.Status == "completed" {
			break
		}

		time.Sleep(5)
	}
	fmt.Println("Done polling for Completed Status")

	return &result, nil
}

func (d *Db) GetSearchFromYoutube(word string) ([]string, error) {
	// create query
	query := word

	// get list
	call := d.ytClient.Search.List([]string{"id,snippet"}).
		Q(query).
		MaxResults(*maxResults)

	// execute req
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	// extract data
	respArr := []string{}
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			respArr = append(respArr, item.Id.VideoId)
		default:
			continue
		}
	}

	return respArr, nil
}

// **************************
// Helpers
// **************************

func findMP3(root string) string {
	var file string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if filepath.Ext(s) == ".mp3" {
			file = s
		}
		return nil
	})
	return file
}
