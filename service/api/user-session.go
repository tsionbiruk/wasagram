package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var username string
	err = json.Unmarshal(body, &username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal username from POST request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if username == "" {
		http.Error(w, "Invalid username!", http.StatusBadRequest)
		return
	}

	err = rt.db.CreateIfNoUser(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to either retrieve or create user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, "Failed to log-in", http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(user_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal UserId %d: %s", user_id, err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(jsonstr))
	if err != nil {

		return
	}
}

func (rt *_router) Authorize(w http.ResponseWriter, r *http.Request, username string) bool {
	tokenstr, ok := r.Header["Authorization"]
	if !ok {
		http.Error(w, "Action unauthorised - authorization not found in header!", http.StatusUnauthorized)
		return false
	}

	token, err := strconv.ParseInt(tokenstr[0], 10, 64)
	if err != nil {
		http.Error(w, "Action unauthorised - corrupted user token!", http.StatusUnauthorized)
		return false
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, "Action unauthorised - wrong user token!", http.StatusUnauthorized)
		return false
	}

	if user_id != token {
		http.Error(w, "Action unauthorised - wrong user token!", http.StatusUnauthorized)
		return false
	}
	return true
}
