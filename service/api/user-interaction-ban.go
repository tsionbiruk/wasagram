package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) putBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("<i am here")

	username := ps.ByName("user")
	target_name := ps.ByName("targetuser")

	log.Printf("Ban request received: %s banning %s", username, target_name)

	// Check Authorization
	if token := rt.Authorize(w, r, username); !token {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get user IDs
	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		log.Printf("Error fetching user ID for %s: %v", username, err)
		http.Error(w, `{"error": "Failed to retrieve user ID"}`, http.StatusInternalServerError)
		return
	}

	target_id, err := rt.db.GetUserIdFromUserName(target_name)
	if err != nil {
		log.Printf("Error fetching target ID for %s: %v", target_name, err)
		http.Error(w, `{"error": "Failed to retrieve target ID"}`, http.StatusInternalServerError)
		return
	}

	// Prevent banning yourself
	if user_id == target_id {
		http.Error(w, `{"error": "You cannot ban yourself"}`, http.StatusBadRequest)
		return
	}

	// Attempt to ban user
	err = rt.db.UserBan(user_id, target_id)
	if err != nil {
		log.Printf("Failed to ban user %d -> %d: %v", user_id, target_id, err)
		http.Error(w, `{"error": "Database error banning user"}`, http.StatusInternalServerError)
		return
	}

	// unfollow after banning
	err = rt.db.UserUnfollow(user_id, target_id)
	if err != nil {
		log.Printf("Failed to unfollow user %d -> %d: %v", user_id, target_id, err)
		http.Error(w, `{"error": "Database error unfollowing user afetr banning"}`, http.StatusInternalServerError)
		return
	}

	// remove banned user from followers list
	err = rt.db.UserUnfollow(target_id, user_id)
	if err != nil {
		log.Printf("Failed to remove user fromm follower list %d -> %d: %v", target_id, user_id, err)
		http.Error(w, `{"error": "Database error unfollowing user afetr banning"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("User %s successfully banned %s", username, target_name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User banned successfully"})
}

func (rt *_router) deleteBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("user")
	target_name := ps.ByName("targetuser")

	log.Printf("Unban request received: %s unbanning %s", username, target_name)

	// Check Authorization
	if token := rt.Authorize(w, r, username); !token {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get user IDs
	user_id, err := rt.db.GetUserIdFromUserName(username)
	if err != nil {
		log.Printf("Error fetching user ID for %s: %v", username, err)
		http.Error(w, `{"error": "Failed to retrieve user ID"}`, http.StatusInternalServerError)
		return
	}

	target_id, err := rt.db.GetUserIdFromUserName(target_name)
	if err != nil {
		log.Printf("Error fetching target ID for %s: %v", target_name, err)
		http.Error(w, `{"error": "Failed to retrieve target ID"}`, http.StatusInternalServerError)
		return
	}

	// Prevent unbanning yourself
	if user_id == target_id {
		http.Error(w, `{"error": "You cannot unban yourself"}`, http.StatusBadRequest)
		return
	}

	// Attempt to unban user
	err = rt.db.UserUnban(user_id, target_id)
	if err != nil {
		log.Printf("Failed to unban user %d -> %d: %v", user_id, target_id, err)
		http.Error(w, `{"error": "Database error unbanning user"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("User %s successfully unbanned %s", username, target_name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User unbanned successfully"})
}

func (rt *_router) getBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("user")

	log.Printf("Fetching banned users for %s", username)

	// Check Authorization
	if token := rt.Authorize(w, r, username); !token {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Fetch banned users
	banned, err := rt.db.UserGetBanned(username)
	if err != nil {
		log.Printf("Failed to get banned users for %s: %v", username, err)
		http.Error(w, `{"error": "Failed to fetch banned users"}`, http.StatusInternalServerError)
		return
	}

	// Convert to JSON
	response, err := json.Marshal(banned)
	if err != nil {
		log.Printf("Error marshalling banned users for %s: %v", username, err)
		http.Error(w, `{"error": "Failed to process banned users"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully retrieved banned users for %s", username)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
