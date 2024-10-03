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

func (rt *_router) PostComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	photo_id_str := ps.ByName("photoid")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	var text string
	err = json.Unmarshal(body, &text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal comment from POST request body: %v", err.Error()), http.StatusBadRequest)
		return
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	err = rt.db.Comment(username, PhotoId, text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post comment: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	// Retrieve comments
	comments, err := rt.db.Getcomment(PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve comments: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Marshal comments to JSON
	responseData, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal comments to JSON: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseData); err != nil {
		// Handle the error
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) DeleteComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	photo_id_str := ps.ByName("photoid")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	comment_id_str := ps.ByName("commentid")
	CommentId, err := strconv.ParseInt(comment_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse comment ID: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	var og_commenter []string
	og_commenter, err = rt.db.GetAuthorcommenter(PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post comment: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	found := false // Flag to track if username is found

	for _, v := range og_commenter {
		if v == username {
			err = rt.db.Uncomment(CommentId)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to unlike the photo, can't unlike what you didn't like: %v", err.Error()), http.StatusInternalServerError)
				return
			}
			found = true // Mark that username was found
			break        // Exit the loop since the user has been found
		}
	}

	// Execute the else-like condition if the username wasn't found
	if !found {
		http.Error(w, "You cannot unlike a comment that you didn't make", http.StatusBadRequest)
	}

	// Retrieve comments
	comments, err := rt.db.Getcomment(PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve comments: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Marshal comments to JSON
	responseData, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal comments to JSON: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseData); err != nil {
		// Handle the error
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
