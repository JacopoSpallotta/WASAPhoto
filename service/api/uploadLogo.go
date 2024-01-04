package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadLogo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Takes the userId and the path of the photo, and uploads it (updates the stream of photos)

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

	intId2, err2 := strconv.Atoi(i)
	if err2 != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// path to the image
	buf := new(bytes.Buffer)
	n, err3 := buf.ReadFrom(r.Body)
	if err3 != nil || n == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newStr := buf.String()

	// create a Photo Struct
	var p Photo
	p.UserId = intId2
	p.Path = newStr

	// create a User Struct
	var u User
	u.ID = intId

	if u.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// update info from database
	dbphoto, dbuser, err4 := rt.db.UploadLogo(p.ToDatabase(), u.ToDatabase())

	if err4 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err4).Error("can't upload the logo")
		w.WriteHeader(http.StatusBadRequest) // 500
		return
	}
	// Here we can re-use `user` as FromDatabase is overwriting every variable in the structure.
	u.FromDatabase(dbuser)
	p.FromDatabase(dbphoto)

	// Send the output to the user.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err5 := json.NewEncoder(w).Encode(p)
	if err5 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
