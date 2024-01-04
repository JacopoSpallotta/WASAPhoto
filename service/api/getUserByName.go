package api

import (
	"encoding/json"
	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getUserByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if token == 0 || errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	name := ps.ByName("id") // Altough is a name, we need to cast it by id

	if name == "" {
		// Empty ID
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// Search for the user to get the profile and returns the information
	var user User
	user.Name = name

	dbuser, err2 := rt.db.GetUserByName(user.ToDatabase())

	if err2 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
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
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

}
