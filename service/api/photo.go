package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tsionbiruk/wasagram/service/api/reqcontext"
)

func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user")

	if token := rt.Authorize(w, r, username); !token {
		return
	}

	photo, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read POST request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if http.DetectContentType(photo) != "image/png" {
		http.Error(w, "Invalid photo type! Only .png files accepted.", http.StatusInternalServerError)
		return
	}

	err = rt.db.PhotoInsert(username, photo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	err = rt.db.PhotoDelete(username, photo_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "image/png")

	photo_id_str := ps.ByName("photoid")
	photo_id, err := strconv.ParseInt(photo_id_str, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	photo, err := rt.db.PhotoGet(photo_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve photo: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(photo)
}
