package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var imageId int

	reqToken := r.Header.Get("Authorization")

	token, errTok := strconv.Atoi(reqToken)

	if errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// This function receives an image ID and returns the image.

	i := ps.ByName("id")

	if i == "" {
		// Empty ID
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageId, e2 := strconv.Atoi(i)

	if e2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// Searchs for the user to get its logo

	var p Photo

	p.ID = imageId
	dbphoto, err3 := rt.db.GetImage(p.ToDatabase())

	if err3 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Here we can re-use `photo` as FromDatabase is overwriting every variable in the structure.
	p.FromDatabase(dbphoto)

	if token != p.UserId {
		// Unauthorized request
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	err4 := json.NewEncoder(w).Encode(p)
	if err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
