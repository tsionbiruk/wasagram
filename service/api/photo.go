package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) []byte {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil
	}

	if !rt.Authorize(w, r, userClaims.Username) {
		return nil
	}

	// Parse the form to retrieve the file and text
	err = r.ParseMultipartForm(10 << 20) // Limit the size to 10 MB
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %s", err.Error()), http.StatusInternalServerError)
		return nil
	}

	// Retrieve the photo file
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get photo: %s", err.Error()), http.StatusBadRequest)
		return nil
	}
	defer file.Close()

	// Retrieve the text
	caption := r.FormValue("text")

	// Read the photo data
	photo, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read photo: %s", err.Error()), http.StatusInternalServerError)
		return nil
	}

	// Validate the photo content type
	contentType := http.DetectContentType(photo)
	if contentType != "image/png" {
		http.Error(w, "Invalid photo type! Only .png files accepted.", http.StatusBadRequest)
		return nil
	}

	// Insert the photo and text into the database
	err = rt.db.UploadPhoto(username, caption, photo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post photo: %s", err.Error()), http.StatusInternalServerError)
		return nil
	}

	return photo

}

func (rt *_router) Photolike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) ([]string, error) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("username")
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

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return nil, nil
	}

	err = rt.db.Photolike(username, PhotoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to like the photo: %s", err.Error()), http.StatusInternalServerError)
		return nil, nil
	}
	return rt.db.Getlikes(PhotoId)
}

func (rt *_router) Photounlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) ([]string, error) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("username")
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

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return nil, nil
	}

	if username == userClaims.Username {
		err = rt.db.Photounlike(username, PhotoId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to unlike the photo, cant unlike what you didn't like: %s", err.Error()), http.StatusInternalServerError)
			return nil, nil
		}
	}

	return rt.db.Getlikes(PhotoId)
}

func (rt *_router) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) (string, error) {
	w.Header().Set("Content-Type", "application/json")
	//username of the persons post you want to delete
	username := ps.ByName("username")
	photo_id_str := ps.ByName("PhotoId")

	userClaims, err := rt.getUserInfoFromRequest(r)
	if err != nil {
		// Handle error (e.g., invalid or missing token)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return "", nil
	}

	if token := rt.Authorize(w, r, userClaims.Username); !token {
		return "", nil
	}

	PhotoId, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return "", nil
	}

	if username == userClaims.Username {
		_, err = rt.db.DeletePost(PhotoId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve photo: %s", err.Error()), http.StatusInternalServerError)
			return "", nil
		}
	}
	return "Only owner of the photo can delete the post", nil

}
