package api

import (
	"bytes"
	"encoding/json"
	"github.com/jorgeferrerhn/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// función que recibe username. Busca en la base de datos, si no existe lo crea.
	var user User

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(r.Body)
	if err != nil || n == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Name = buf.String()

	if !user.IsValid() {
		// Here we validated the user structure content (correct name), and we discovered that the user data is not valid
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbuser, err2 := rt.db.DoLogin(user.ToDatabase())

	if err2 != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err2).Error("can't create the user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Here we can re-use `user` as FromDatabase is overwriting every variable in the structure.
	user.FromDatabase(dbuser)

	// Send the output to the user.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err3 := json.NewEncoder(w).Encode(user)
	if err3 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	defer r.Body.Close()
}
