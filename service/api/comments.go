package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) PostComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) ([]string, error) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	photo_id_str := ps.ByName("PhotoId")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, nil
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return nil, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return nil, nil
	}

	var text string
	err = json.Unmarshal(body, &text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal comment from POST request body: %s", err.Error()), http.StatusBadRequest)
		return nil, nil
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return nil, nil
	}

	err = rt.db.Comment(username, PhotoId, text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post comment: %s", err.Error()), http.StatusInternalServerError)
		return nil, err
	}
	return rt.db.Getcomment(PhotoId)
}

func (rt *_router) DeleteComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) ([]string, error) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	photo_id_str := ps.ByName("PhotoId")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, nil
	}

	if token := rt.Authorize(w, r, userClaims.Subject); !token {
		return nil, nil
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return nil, nil
	}
	comment_id_str := ps.ByName("CommentId")
	CommentId, err := strconv.ParseInt(comment_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse comment ID: %s", err.Error()), http.StatusInternalServerError)
		return nil, nil
	}

	if username == userClaims.Username {
		err = rt.db.Uncomment(CommentId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete comment: %s", err.Error()), http.StatusInternalServerError)
			return nil, nil
		}
	}

	return rt.db.Getcomment(PhotoId)
}
