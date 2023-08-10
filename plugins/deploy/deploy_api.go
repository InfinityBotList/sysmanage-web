package deploy

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/crypto"
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

		var flag bool
		for _, webh := range cfg.Webhooks {
			if webh.Type != "api" {
				continue
			}

			if webh.Token == token {
				flag = true
				break
			}
		}

		if !flag {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid token"))
			return
		}

		reqId := crypto.RandString(64)

		go initDeploy(reqId, cfg)

		w.Write([]byte(reqId))
	})
}
