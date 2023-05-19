package actions

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func loadActionsApi(r *chi.Mux) {
	r.Post("/api/getActionList", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(Actions)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not marshal actions"))
			return
		}

		w.Write(bytes)
	})

	r.Post("/api/executeAction", func(w http.ResponseWriter, r *http.Request) {
		actionName := r.URL.Query().Get("actionName")

		if actionName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing actionName"))
			return
		}

		action, ok := Actions.Find(actionName)

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Action not found"))
			return
		}

		response, err := action.Handler(&ActionContext{Request: r})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		if response.TaskID != "" {
			w.Header().Add("X-Task-ID", response.TaskID)
		}

		w.WriteHeader(response.StatusCode)
		w.Write([]byte(response.Resp))
	})
}
