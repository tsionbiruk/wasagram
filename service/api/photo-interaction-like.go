package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) putLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	photo_id_str := ps.ByName("photoid")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	photo_id, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to like the photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = rt.db.PhotoLike(user_id, photo_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to like the photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	photo_id_str := ps.ByName("photoid")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	photo_id, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unlike the photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = rt.db.PhotoUnlike(user_id, photo_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unlike the photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
