package api

import (
	"encoding/json"
	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
	// Searchs for the user to get the profile and returns the information
	var user User
	user.ID = intId

	if user.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbuser, err2 := rt.db.GetUserProfile(user.ToDatabase())

	if err2 != nil {

		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err2).Error("can't create the user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user.FromDatabase(dbuser)

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	err3 := json.NewEncoder(w).Encode(user)
	if err3 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

}
