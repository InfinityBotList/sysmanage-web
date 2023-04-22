package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/crypto"
)

func loadNginxApi(r *chi.Mux) {
	r.Post("/api/nginx/buildNginx", func(w http.ResponseWriter, r *http.Request) {
		reqId := crypto.RandString(64)

		go buildNginx(reqId)

		w.Write([]byte(reqId))

	})
}
