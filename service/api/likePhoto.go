package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Takes the photo Id and updates its like in the photos table
	// user id
	i := ps.ByName("id")

	if i == "" {
		// Empty ID
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intId, err := strconv.Atoi(i)
	if err != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// photo id
	photoId := ps.ByName("photoId")

	if photoId == "" {
		// Empty Photo ID
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intPhoto, err := strconv.Atoi(photoId)
	if err != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var p Photo
	var u User

	p.ID = intPhoto
	p.UserId = intId // only for sending it to the database function
	u.ID = intId

	if u.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// update info from database
	dbphoto, dbuser, err2 := rt.db.LikePhoto(p.ToDatabase(), u.ToDatabase())
	if err2 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err2).Error("can't upload the like to the photo")
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// Here we can re-use `photo ` as FromDatabase is overwriting every variable in the structure.
	p.FromDatabase(dbphoto)
	u.FromDatabase(dbuser)

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")

	if p.ID == 0 { // user not found
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// Send the output to the user.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err3 := json.NewEncoder(w).Encode(p)
	if err3 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
