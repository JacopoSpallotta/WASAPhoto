package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	reqToken := r.Header.Get("Authorization")
	token, errTok := strconv.Atoi(reqToken)
	if errTok != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusBadRequest)
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

	// followedUser
	idFollowed := ps.ByName("id2")

	if idFollowed == "" {
		// Empty Followed ID
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intFollowed, err2 := strconv.Atoi(idFollowed)
	if err2 != nil {
		// id was not properly cast
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// User 1
	var u1 User
	u1.ID = intId

	// User 2
	var u2 User
	u2.ID = intFollowed

	if u1.ID != token {
		// Error: the authorization header is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if u2.ID == u1.ID {
		// You can't ban yourself
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// update info from database
	dbuser1, err3 := rt.db.BanUser(u1.ToDatabase(), u2.ToDatabase())
	if err3 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err3).Error("can't update the banned list")
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	u1.FromDatabase(dbuser1)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err4 := json.NewEncoder(w).Encode(u1)
	if err4 != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return

	}

	defer r.Body.Close()

}
