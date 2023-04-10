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
)

func getServiceStatus(id string) string {
	cmd := exec.Command("systemctl", "check", id)
	out, _ := cmd.CombinedOutput()

	return strings.ReplaceAll(string(out), "\n", "")
}

func loadApi(r *chi.Mux) {
	// Returns the list of services
	r.Post("/api/getServiceList", func(w http.ResponseWriter, r *http.Request) {
		defines := []ServiceManage{}

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
				var service TemplateYaml

				err = yaml.NewDecoder(f).Decode(&service)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Failed to read service definition" + err.Error() + file.Name()))
					return
				}

				// Service name is the name without .yaml
				sname := strings.TrimSuffix(file.Name(), ".yaml")

				defines = append(defines, ServiceManage{
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

	r.Post("/api/buildServices", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("consoleOf") != "" {
			// Fetch from redis
			console := rdb.Get(ctx, logPrefix+r.URL.Query().Get("consoleOf")).Val()

			isDone := rdb.Get(ctx, logPrefix+r.URL.Query().Get("consoleOf")+markerDoneSuffix).Val()

			if isDone == "1" {
				w.Header().Set("X-Is-Done", "1")
			}

			w.Write([]byte(console))
			return
		}

		reqId := crypto.RandString(64)

		go func() {
			addToLog(reqId, "Waiting for other builds to finish...", true)

			inDeploy.Lock()
			defer inDeploy.Unlock()

			addToLog(reqId, "Starting build process to convert service templates to systemd services...", true)

			cmd := exec.Command("make", "systemd")
			cmd.Dir = config.InfraFolder
			cmd.Env = os.Environ()
			cmd.Stdout = autoLogger{id: reqId}
			cmd.Stderr = autoLogger{id: reqId, Error: true}

			err := cmd.Run()

			if err != nil {
				addToLog(reqId, "Failed to build services: "+err.Error(), true)
				return
			}

			addToLog(reqId, "Finished building services.", false)
			markLogDone(reqId)
		}()

		w.Write([]byte(reqId))
	})
}
