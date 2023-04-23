package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"sysmanage-web/types"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/crypto"
)

func loadNginxApi(r *chi.Mux) {
	r.Post("/api/nginx/buildNginx", func(w http.ResponseWriter, r *http.Request) {
		reqId := crypto.RandString(64)

		go buildNginx(reqId)

		w.Write([]byte(reqId))
	})

	r.Post("/api/nginx/getDomainList", func(w http.ResponseWriter, r *http.Request) {
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

	r.Post("/api/nginx/publishCerts", func(w http.ResponseWriter, r *http.Request) {
		var req types.NginxAPIPublishCert

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Validate request
		err = v.Struct(req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Load meta
		meta, err := loadNginxMeta()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Check cert and key
		_, err = tls.X509KeyPair([]byte(req.Cert), []byte(req.Key))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		domain := strings.ReplaceAll(req.Domain, ".", "-")

		// Remove any http/https prefix
		if strings.Contains(domain, "/") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Domain cannot contain http/https prefix or slashes"))
			return
		}

		// Ensure not subdomain
		domainSplit := strings.Split(domain, ".")

		if len(domainSplit) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Domain must contain a dot and must not be a subdomain"))
			return
		}

		certFile := meta.NginxCertPath + "/cert-" + domain + ".pem"
		keyFile := meta.NginxCertPath + "/key-" + domain + ".pem"

		// Check that the cert and key files do not already exists
		if r.URL.Query().Get("force") != "true" {
			exists := true

			for _, f := range []string{certFile, keyFile} {
				if _, err := os.Stat(f); err == nil || errors.Is(err, fs.ErrNotExist) {
					exists = false
					break
				}
			}

			if exists {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("ALREADY_EXISTS"))
				return
			}
		}

		// Write cert and key
		err = os.WriteFile(certFile, []byte(req.Cert), 0644)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = os.WriteFile(keyFile, []byte(req.Key), 0644)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
