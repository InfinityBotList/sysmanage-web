package frontend

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func loadFrontendApi(r chi.Router) {
	r.Post("/getRegisteredLinks", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(RegisteredLinks)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not marshal registered links"))
			return
		}

		w.Write(bytes)
	})
}
