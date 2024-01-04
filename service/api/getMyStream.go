package api

import (
	"encoding/json"
	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// This function receives a user id and returns the stream of photos of that user

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	intId, err := checkId(ps)
	if err != nil {
		// error on database
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Search for the user to get the stream of photos
	var user User
	user.ID = intId

	if user.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photos, err2 := rt.db.GetMyStream(user.ToDatabase())
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	err4 := json.NewEncoder(w).Encode(photos)
	if err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
