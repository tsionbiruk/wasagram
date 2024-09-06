package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) followuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) string {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("username")
	target_username := ps.ByName("target_username")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	if token := rt.Authorize(w, r, userClaims.Issuer); !token {
		return ""
	}

	_, err = rt.db.FollowUser(username, target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to follow target: %s", err.Error()), http.StatusInternalServerError)
		return "can not follow yourself"
	}
	return "user followed succesfully"
}

func (rt *_router) unfolloweruser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) string {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("username")
	target_username := ps.ByName("target_username")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return ""
	}

	_, err = rt.db.UnFollowUser(username, target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unfollow user target: %s", err.Error()), http.StatusInternalServerError)
		return "can not unfollow yourself!"
	}
	return "user unfollowed!"
}

func (rt *_router) Ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("username")
	target_username := ps.ByName("target_username")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return
	}

	_, err = rt.db.BanUsers(username, target_username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) unban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) string {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("username")
	target_username := ps.ByName("target_username")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return ""
	}

	_, err = rt.db.UnBanUser(username, target_username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unban target: %s", err.Error()), http.StatusInternalServerError)
		return "you can not unban yourself"
	}

	return "user banned succesfully!"
}
