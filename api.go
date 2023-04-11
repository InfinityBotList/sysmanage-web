package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/crypto"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"sysmanage-web/types"
)

func getServiceStatus(id string) string {
	cmd := exec.Command("systemctl", "check", id)
	out, _ := cmd.CombinedOutput()

	return strings.ReplaceAll(string(out), "\n", "")
}

func loadApi(r *chi.Mux) {
	// Returns the list of services
	r.Post("/api/getServiceList", func(w http.ResponseWriter, r *http.Request) {
		defines := []types.ServiceManage{}

		for _, path := range config.ServiceDefinitions {
			// Get all files in path
			fsd, err := os.ReadDir(path)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to read service definition " + err.Error()))
				return
			}

			for _, file := range fsd {
				if file.Name() == "_meta.yaml" {
					continue // Skip _meta.yaml
				}

				if file.IsDir() {
					continue // Skip directories
				}

				if !strings.HasSuffix(file.Name(), ".yaml") {
					continue // Skip non-yaml files
				}

				// Read file into TemplateYaml
				f, err := os.Open(path + "/" + file.Name())

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Failed to read service definition." + err.Error() + file.Name()))
					return
				}

				// Read file into TemplateYaml
				var service types.TemplateYaml

				err = yaml.NewDecoder(f).Decode(&service)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Failed to read service definition" + err.Error() + file.Name()))
					return
				}

				// Service name is the name without .yaml
				sname := strings.TrimSuffix(file.Name(), ".yaml")

				defines = append(defines, types.ServiceManage{
					DefinitionFolder: path,
					Service:          service,
					ID:               sname,
					Status:           getServiceStatus(sname),
				})
			}
		}

		// JSON encode defines
		jsonStr, err := json.Marshal(defines)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/api/systemctl", func(w http.ResponseWriter, r *http.Request) {
		tgt := r.URL.Query().Get("tgt")
		act := r.URL.Query().Get("act")

		if !slices.Contains([]string{"start", "stop", "restart", "list-dependencies"}, act) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid action."))
			return
		}

		if tgt == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing tgt parameter."))
			return
		}

		cmd := exec.Command("systemctl", act, tgt)
		out, _ := cmd.CombinedOutput()

		w.Write(out)
	})

	r.Post("/api/getLogEntry", func(w http.ResponseWriter, r *http.Request) {
		// Fetch from redis
		console := rdb.Get(ctx, logPrefix+r.URL.Query().Get("id")).Val()

		isDone := rdb.Get(ctx, logPrefix+r.URL.Query().Get("id")+markerDoneSuffix).Val()

		if isDone == "1" {
			w.Header().Set("X-Is-Done", "1")
		}

		w.Write([]byte(console))
	})

	r.Post("/api/buildServices", func(w http.ResponseWriter, r *http.Request) {
		reqId := crypto.RandString(64)

		go buildServices(reqId)

		w.Write([]byte(reqId))
	})
}
