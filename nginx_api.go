package main

import (
	"encoding/json"
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

	r.Get("/api/nginx/getDomainList", func(w http.ResponseWriter, r *http.Request) {
		domList, err := getNginxDomainList()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		bytes, err := json.Marshal(domList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(bytes)
	})
}
