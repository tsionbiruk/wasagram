package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) putBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_name := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	target_id, err := rt.db.GetUserIdFromUserName(target_name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if user_id == target_id {
		http.Error(w, "Cannot ban yourself!", http.StatusBadRequest)
		return
	}

	err = rt.db.UserBan(user_id, target_id)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_name := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	target_id, err := rt.db.GetUserIdFromUserName(target_name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if user_id == target_id {
		http.Error(w, "Cannot unban yourself!", http.StatusBadRequest)
		return
	}

	err = rt.db.UserUnban(user_id, target_id)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	banned, err := rt.db.UserGetBanned(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get banned users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(banned)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal banned users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(jsonstr))
}
