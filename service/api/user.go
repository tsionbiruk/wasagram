package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) GetProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_username := ps.ByName("target_user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	profile, err := rt.db.UserProfile(target_username)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve profile information: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal profile data: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonstr)
}

func (rt *_router) GetStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	target_username := ps.ByName("target_user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	photos, stream, err := rt.db.GetStream(target_username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve stream: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"stream": stream,
		"photos": photos,
	}

	// Marshal the combined response into JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal combined response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	_, _ = w.Write(jsonResponse)

}

func (rt *_router) Rename(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var newname string

	if newname == "" {
		http.Error(w, "Invalid username!", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &newname)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal username from POST request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = rt.db.UpdateUserName(username, newname)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to rename user, can not rename other users: %s", err.Error()), http.StatusInternalServerError)
		return

	}

	// Create the success response
	response := map[string]string{
		"message": "Username updated successfully",
	}

	// Encode the response to JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

}

func (rt *_router) GetUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	users, err := rt.db.GetAllUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonstr, err := json.Marshal(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal users array: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonstr)
}

func (rt *_router) DoLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var username string
	err = json.Unmarshal(body, &username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal username from POST request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if username == "" {
		http.Error(w, "Invalid username!", http.StatusBadRequest)
		return
	}

	var token int
	token, err = rt.db.Creatuser_Getuserfromtoken(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to either retrieve or create user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	//MARSHAL THETOKEN NOT THE USERNAME

	jsonstr, err := json.Marshal(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal token %s: %s", token, err.Error()), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(jsonstr))
}
