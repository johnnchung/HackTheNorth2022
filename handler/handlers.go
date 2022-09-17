package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johnnchung/HackTheNorth2022/helpers"
	models "github.com/johnnchung/HackTheNorth2022/models"
	"github.com/urfave/negroni"
)

type Repo struct {
	muxClient *mux.Router
	db        *models.Db
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
	r.muxClient.HandleFunc("/{term}", r.addDataFromCategoryHandler)

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
			if err := r.db.InsertWords(word, ytRsp.URL); err != nil {
				helpers.RespondWithError(w, http.StatusBadGateway, err.Error())
				return
			}
		}

	}

	helpers.RespondWithJSON(w, http.StatusAccepted, nil)
	return

}
