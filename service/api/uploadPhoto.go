package api

import (
	"bytes"
	"encoding/json"
	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Takes the userId and the path of the photo, and uploads it (updates the stream of photos)

	intId, err := checkId(ps)

	if err != nil {
		// error on database
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// path to the image
	buf := new(bytes.Buffer)
	n, err2 := buf.ReadFrom(r.Body)
	if err2 != nil || n == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newStr := buf.String()

	var p Photo

	p.UserId = intId
	p.Path = newStr

	var u User
	u.ID = intId

	if u.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// update info from database
	dbphoto, dbuser, err3 := rt.db.UploadPhoto(p.ToDatabase(), u.ToDatabase())

	if err3 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err3).Error("can't upload the photo")
		w.WriteHeader(http.StatusBadRequest) // 500
		return
	}
	// Here we can re-use `user` as FromDatabase is overwriting every variable in the structure.
	p.FromDatabase(dbphoto)
	u.FromDatabase(dbuser)

	// Send the output to the user.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err4 := json.NewEncoder(w).Encode(p)
	if err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
