package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) Followuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) string {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return ""
	}

	_, err = rt.db.FollowUser(userClaims.Username, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to follow target: %s", err.Error()), http.StatusInternalServerError)
		return "can not follow yourself"
	}
	return "user followed succesfully"
}

func (rt *_router) Unfolloweruser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) string {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return ""
	}

	_, err = rt.db.UnFollowUser(userClaims.Username, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unfollow user target: %s", err.Error()), http.StatusInternalServerError)
		return "can not unfollow yourself!"
	}
	return "user unfollowed!"
}

func (rt *_router) Ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return
	}

	_, err = rt.db.BanUsers(userClaims.Username, username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) Unban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) string {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return ""
	}

	_, err = rt.db.UnBanUser(userClaims.Username, username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unban target: %s", err.Error()), http.StatusInternalServerError)
		return "you can not unban yourself"
	}

	return "user banned succesfully!"
}

func (rt *_router) GetFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return
	}

	followers, err := rt.db.Getfollowers(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followers: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(followers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal followed users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(jsonstr))
}

func (rt *_router) Getfollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return
	}

	following, err := rt.db.Getfollowing(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followers: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(following)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal followed users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(jsonstr))
}

func (rt *_router) GetBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return
	}

	banned, err := rt.db.Getbannedusers(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followers: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(banned)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal followed users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(jsonstr))
}
