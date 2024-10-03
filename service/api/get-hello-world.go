package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	if _, err := w.Write([]byte("Hello World!")); err != nil {
		// Handle the error
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
