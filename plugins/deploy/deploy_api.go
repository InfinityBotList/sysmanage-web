package deploy

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slices"
)

func loadDeployApi(r chi.Router) {
	r.Post("/getDeployMeta", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing id"))
			return
		}

		cfg, err := LoadConfig(id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to load config: " + err.Error()))
			return
		}

		// JSON encode defines
		jsonStr, err := json.Marshal(cfg)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/getDeployList", func(w http.ResponseWriter, r *http.Request) {
		cfg, err := GetDeployList()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to load config: " + err.Error()))
			return
		}

		// JSON encode defines
		jsonStr, err := json.Marshal(cfg)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/getDeploySourceTypes", func(w http.ResponseWriter, r *http.Request) {
		var srcs []string

		for k := range DeploySources {
			srcs = append(srcs, k)
		}

		// JSON encode defines
		jsonStr, err := json.Marshal(srcs)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/getDeployWebhookSourceTypes", func(w http.ResponseWriter, r *http.Request) {
		var srcs []string

		for k := range DeployWebhookSources {
			srcs = append(srcs, k)
		}

		// JSON encode defines
		jsonStr, err := json.Marshal(srcs)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/createDeploy", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing id"))
			return
		}

		token := r.URL.Query().Get("token")

		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing token"))
			return
		}

		cfg, err := LoadConfig(id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to load config: " + err.Error()))
			return
		}

		if len(cfg.AllowedIps) > 0 && !slices.Contains(cfg.AllowedIps, r.RemoteAddr) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("ip not allowed"))
			return
		}

		typ := r.URL.Query().Get("type")

		if typ == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing type"))
			return
		}

		wid := r.URL.Query().Get("wid")

		if wid == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing wid"))
			return
		}

		fn, ok := DeployWebhookSources[typ]

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid type"))
			return
		}

		logId, err := fn(cfg, wid, id, token)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(logId))
	})
}
