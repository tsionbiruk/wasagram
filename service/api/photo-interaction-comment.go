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

func (rt *_router) postComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	photo_id_str := ps.ByName("photoid")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var comment_text string
	err = json.Unmarshal(body, &comment_text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal comment from POST request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	photo_id, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = rt.db.PhotoComment(username, photo_id, comment_text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post comment: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	comment_id_str := ps.ByName("commentid")
	comment_id, err := strconv.ParseInt(comment_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse comment ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = rt.db.PhotoUncomment(username, photo_id, comment_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete comment: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
