package api

import (
	"encoding/json"
	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// comment id
	commentId := ps.ByName("commentId")

	if commentId == "" {
		// Empty Photo ID
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intComment, err2 := strconv.Atoi(commentId)
	if err2 != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create a Photo Struct
	var c Comment
	c.ID = intComment // default
	c.UserId = intId

	// We pass this photo because we will have to update its comments in the database
	var p Photo

	var u User
	u.ID = intId

	if u.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// update info from database
	dbcomment, dbphoto, dbuser, err3 := rt.db.UncommentPhoto(c.ToDatabase(), p.ToDatabase(), u.ToDatabase())
	if err3 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err3).Error("can't delete the comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Here we can re-use `comment` as FromDatabase is overwriting every variable in the structure.
	c.FromDatabase(dbcomment)
	p.FromDatabase(dbphoto)
	u.FromDatabase(dbuser)

	// Send the output to the user.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err4 := json.NewEncoder(w).Encode(p) // Return the photo
	if err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

}
