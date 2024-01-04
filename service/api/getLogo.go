package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getLogo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// this function receives a user id and returns the logo

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if errTok != nil || token == 0 {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var id = ps.ByName("id")
	if id == "" {
		// ID not found
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var i, e = strconv.Atoi(id)

	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Searchs for the user to get its logo
	var p Photo
	var u User
	u.ID = i

	dbphoto, dbuser, err := rt.db.GetLogo(p.ToDatabase(), u.ToDatabase())

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// Here we can re-use `photo` as FromDatabase is overwriting every variable in the structure.
	p.FromDatabase(dbphoto)
	u.FromDatabase(dbuser)

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	err2 := json.NewEncoder(w).Encode(p)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
