package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) Followuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var target_username string

	err = json.Unmarshal(body, &target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal username from POST request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	if username == target_username {
		response := map[string]string{
			"message": "Can not follow youself",
		}

		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the response
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(responseData); err != nil {
			// Handle the error
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	}

	_, err = rt.db.FollowUser(username, target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to follow target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Create the success response
	response := map[string]string{
		"message": "User followed succesfully",
	}

	// Encode the response to JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
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

func (rt *_router) Unfollowuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("user")
	target_username := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	_, err := rt.db.UnFollowUser(username, target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unfollow user target: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Create the success response
	response := map[string]string{
		"message": "User unfollowed succesfully",
	}

	// Encode the response to JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
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

func (rt *_router) Ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("user")
	target_username := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	_, err := rt.db.BanUsers(username, target_username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ban target: %s", err.Error()), http.StatusInternalServerError)
		return

	}

	if username == target_username {
		response := map[string]string{
			"message": "Cannot ban yourself",
		}

		// Encode the response to JSON
		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the response
		w.WriteHeader(http.StatusBadRequest) // 400 status code for invalid request
		if _, err := w.Write(responseData); err != nil {
			// Handle the error
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	}

	response := map[string]string{
		"message": "User banned succesfully",
	}

	// Encode the response to JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
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

func (rt *_router) Unban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_username := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	if username == target_username {
		response := map[string]string{
			"message": "Cannot unban yourself",
		}
		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the response
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(responseData); err != nil {
			// Handle the error
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	}

	_, err := rt.db.UnBanUser(username, target_username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unban target: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "User unbanned succesfully",
	}

	// Encode the response to JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
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

func (rt *_router) GetFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var target_username string

	err = json.Unmarshal(body, &target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal username from POST request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	followers, err := rt.db.Getfollowers(target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followers: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(followers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal followed users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write([]byte(jsonstr)); err != nil {
		// Handle the error
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) Getfollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_username := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	following, err := rt.db.Getfollowing(target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followers: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(following)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal followed users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write([]byte(jsonstr)); err != nil {
		// Handle the error
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) GetBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_username := ps.ByName("targetuser")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	banned, err := rt.db.Getbannedusers(target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get banned users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(banned)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal followed users: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write([]byte(jsonstr)); err != nil {
		// Handle the error
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
