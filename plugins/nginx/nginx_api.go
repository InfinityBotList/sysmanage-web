package nginx

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/plugins/persist"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/crypto"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

func loadNginxApi(r chi.Router) {
	r.Post("/buildNginx", func(w http.ResponseWriter, r *http.Request) {
		reqId := crypto.RandString(64)

		go buildNginx(reqId)

		w.Write([]byte(reqId))
	})

	r.Post("/updateDnsRecordCf", func(w http.ResponseWriter, r *http.Request) {
		reqId := crypto.RandString(64)

		go updateDnsRecordCf(reqId)

		w.Write([]byte(reqId))
	})

	r.Post("/getDomainList", func(w http.ResponseWriter, r *http.Request) {
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

	r.Post("/publishCerts", func(w http.ResponseWriter, r *http.Request) {
		var req NginxAPIPublishCert

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Validate request
		err = state.Validator.Struct(req)

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

		// Remove any http/https prefix
		if strings.Contains(req.Domain, "/") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Domain cannot contain http/https prefix or slashes"))
			return
		}

		domainSplit := strings.Split(req.Domain, ".")

		if len(domainSplit) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Domain must contain a dot and must not be a subdomain"))
			return
		}

		if cf != nil {
			// Ensure domain is on cloudflare
			if _, ok := zoneMap[req.Domain]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Domain must be on Cloudflare"))
				return
			}
		}

		certFile := meta.NginxCertPath + "/cert-" + req.Domain + ".pem"
		keyFile := meta.NginxCertPath + "/key-" + req.Domain + ".pem"

		// Check that the cert and key files do not already exists
		if r.URL.Query().Get("force") != "true" {
			exists := true

			for _, f := range []string{certFile, keyFile} {
				if _, err := os.Stat(f); errors.Is(err, fs.ErrNotExist) {
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

	r.Post("/getCertList", func(w http.ResponseWriter, r *http.Request) {
		// Load meta
		meta, err := loadNginxMeta()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		fsd, err := os.ReadDir(meta.NginxCertPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var certList []string

		for _, f := range fsd {
			if strings.Contains(f.Name(), "cert-") {
				certList = append(certList, f.Name())
			}
		}

		bytes, err := json.Marshal(certList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(bytes)
	})

	r.Post("/addDomain", func(w http.ResponseWriter, r *http.Request) {
		domainName := r.URL.Query().Get("domain")

		if domainName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Domain name cannot be empty"))
			return
		}

		// Load meta
		meta, err := loadNginxMeta()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Check that cert and key exists
		certFile := meta.NginxCertPath + "/cert-" + domainName + ".pem"
		keyFile := meta.NginxCertPath + "/key-" + domainName + ".pem"

		_, err = tls.LoadX509KeyPair(certFile, keyFile)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Check that the domain does not already exists
		domList, err := getNginxDomainList()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		for _, d := range domList {
			if d.Domain == domainName {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Domain already exists"))
				return
			}
		}

		// Add domain
		f, err := os.Create(nginxDefinitions + "/" + domainName + ".yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		defer f.Close()

		_, err = f.WriteString("servers:")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		go persist.PersistToGit("")

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/updateDomain", func(w http.ResponseWriter, r *http.Request) {
		var req NginxServerManage

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if len(req.Server.Servers) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No servers found?"))
		}

		err = state.Validator.Struct(req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		getSub := []string{} // Used to check for duplicate subdomains
		for _, srv := range req.Server.Servers {
			if strings.Contains(srv.ID, " ") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Server ID cannot contain spaces"))
				return
			}

			for i := range srv.Names {
				if srv.Names[i] == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Emoty subdomains are not allowed"))
					return
				}

				srv.Names[i] = strings.Replace(srv.Names[i], "."+req.Domain, "", 1)

				if strings.Contains(srv.Names[i], " ") {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Subdomains should not spaces"))
					return
				}

				if slices.Contains(getSub, srv.Names[i]) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("All subdomains must be unique"))
					return
				}

				getSub = append(getSub, srv.Names[i])
			}

			if len(srv.Locations) > 0 {
				gotRoot := false
				gotPaths := []string{}

				for _, loc := range srv.Locations {
					if loc.Path == "/" {
						gotRoot = true
					}

					if strings.Contains(loc.Proxy, ";") || strings.Contains(loc.Proxy, " ") {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte("Proxy cannot contain spaces or semicolons"))
						return
					}

					gotPaths = append(gotPaths, loc.Path)
				}

				// TODO: Use gotPaths (???) or remove it
				fmt.Println(gotPaths)

				if !gotRoot {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Atleast one location named '/' must be specified"))
					return
				}
			}
		}

		// Check that the domain exists
		_, err = os.Stat(nginxDefinitions + "/" + req.Domain + ".yaml")

		if errors.Is(err, fs.ErrNotExist) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Domain does not exist"))
			return
		}

		// Update domain
		f, err := os.Create(nginxDefinitions + "/" + req.Domain + ".yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = yaml.NewEncoder(f).Encode(req.Server)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = f.Close()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		go persist.PersistToGit("")

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/deleteDomain", func(w http.ResponseWriter, r *http.Request) {
		domainName := r.URL.Query().Get("domain")

		if domainName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Domain must be specified"))
		}

		// create task id
		reqId := crypto.RandString(64)

		go deleteDomain(reqId, domainName)

		w.Write([]byte(reqId))
	})
}
