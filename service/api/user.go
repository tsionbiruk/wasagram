package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	stream, err := rt.db.UserStream(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve stream: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(stream)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal stream data: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonstr)
}
func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	profile, err := rt.db.UserProfile(username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve profile information: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal profile data: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonstr)
}

func (rt *_router) postUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var newname string

	// if newname == "" {
	// http.Error(w, "Invalid username!", http.StatusBadRequest)
	// return
	// }

	err = json.Unmarshal(body, &newname)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal username from POST request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = rt.db.UserRename(username, newname)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to rename user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	users, err := rt.db.GetAllUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal users array: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonstr)
}
