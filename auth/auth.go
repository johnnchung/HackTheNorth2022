package auth

import (
	"net/http"
	"os"

	"github.com/johnnchung/HackTheNorth2022/helpers"
)

func ReqAPIKey(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	reqToken := r.Header.Get("Authorization")
	apiToken := os.Getenv("API_KEY")

	if reqToken != apiToken {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Authentication Failed! Need API key")
		return
	}

	next(w, r)
}
