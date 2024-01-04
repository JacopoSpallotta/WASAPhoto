package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Takes the userId and the comment, and uploads it (updates the comments table)

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

	// Comment
	buf := new(bytes.Buffer)
	n, err2 := buf.ReadFrom(r.Body)
	if err2 != nil || n == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	comment := buf.String()

	// create a Photo Struct
	var c Comment

	c.UserId = intId
	c.Content = comment
	c.PhotoId = intPhoto

	if c.Content == "" || !c.RightComment() { // empty comment or not valid (len not in range (3,144) )
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	var p Photo
	p.ID = intPhoto

	var u User
	u.ID = intId

	if u.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// update info from database
	dbcomment, dbphoto, dbuser, err3 := rt.db.CommentPhoto(c.ToDatabase(), p.ToDatabase(), u.ToDatabase())
	if err3 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't upload the comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Here we can re-use `comment` as FromDatabase is overwriting every variable in the structure.
	c.FromDatabase(dbcomment)
	p.FromDatabase(dbphoto)
	u.FromDatabase(dbuser)

	if c.ID == 0 { // user not found
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// Send the output to the user.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err4 := json.NewEncoder(w).Encode(c)
	if err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
