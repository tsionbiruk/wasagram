package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) putFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_name := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to follow target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	target_id, err := rt.db.GetUserIdFromUserName(target_name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to follow target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if user_id == target_id {
		http.Error(w, "Cannot follow yourself!", http.StatusBadRequest)
		return
	}

	err = rt.db.UserFollow(user_id, target_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to follow target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_name := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unfollow target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	target_id, err := rt.db.GetUserIdFromUserName(target_name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unfollow target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if user_id == target_id {
		http.Error(w, "Cannot unfollow yourself!", http.StatusBadRequest)
		return
	}

	err = rt.db.UserUnfollow(user_id, target_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unfollow user target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	followed, err := rt.db.UserGetFollowed(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followed users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // Explicitly set status 200 OK

	err = json.NewEncoder(w).Encode(followed)
	if err != nil {
		http.Error(w, "Failed to encode followed users", http.StatusInternalServerError)
		return
	}
}
