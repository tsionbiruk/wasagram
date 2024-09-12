package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if !rt.Authorize(w, r, username) {
		return
	}

	// Parse the form to retrieve the file and text
	err := r.ParseMultipartForm(10 << 20) // Limit the size to 10 MB
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Retrieve the photo file
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get photo: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Retrieve the text
	caption := r.FormValue("text")

	// Read the photo data
	photo, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Validate the photo content type
	contentType := http.DetectContentType(photo)
	if contentType != "image/png" {
		http.Error(w, "Invalid photo type! Only .png files accepted.", http.StatusBadRequest)
		return
	}

	// Insert the photo and text into the database
	err = rt.db.UploadPhoto(username, caption, photo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Encode the photo to base64
	encodedPhoto := base64.StdEncoding.EncodeToString(photo)

	// Create the response map
	response := map[string]string{
		"photo":   encodedPhoto,
		"caption": caption,
	}

	// Marshal the response to JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response to JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

}

func (rt *_router) Photolike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")
	//target_username := ps.ByName("target_user")
	photo_id_str := ps.ByName("photoid")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = rt.db.Photolike(username, PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to like the photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	likes, err := rt.db.Getlikes(PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve likes: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Marshal the likes to JSON
	responseData, err := json.Marshal(likes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal likes to JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (rt *_router) Photounlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	photo_id_str := ps.ByName("photoid")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var og_liker string
	og_liker, err = rt.db.GetAuthorliker(PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post comment: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if username == og_liker {
		err = rt.db.Photounlike(username, PhotoId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to unlike the photo, cant unlike what you didn't like: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}

	likes, err := rt.db.Getlikes(PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve likes: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Marshal the likes to JSON
	responseData, err := json.Marshal(likes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal likes to JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (rt *_router) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	//username of the persons post you want to delete
	username := ps.ByName("user")
	photo_id_str := ps.ByName("photoid")
	fmt.Print(photo_id_str)

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var og_poster string
	og_poster, err = rt.db.GetAuthorId(PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if username == og_poster {
		_, err = rt.db.DeletePost(PhotoId)
		if err != nil {
			http.Error(w, fmt.Sprintf("only posts author can delete the post: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
	// Success message
	response := map[string]string{
		"message": "Photo deleted successfully",
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response to JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

}
