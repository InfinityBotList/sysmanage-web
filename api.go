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

func getServiceStatus(ids []string) []string {
	ids = append([]string{"check"}, ids...)
	cmd := exec.Command("systemctl", ids...)
	out, _ := cmd.CombinedOutput()

	return strings.Split(string(out), "\n")
}

func loadApi(r *chi.Mux) {
	// Returns the list of services
	r.Post("/api/getServiceList", func(w http.ResponseWriter, r *http.Request) {
		defines := []types.ServiceManage{}

		ids := []string{}
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

				ids = append(ids, sname)

				defines = append(defines, types.ServiceManage{
					DefinitionFolder: path,
					Service:          service,
					ID:               sname,
				})
			}
		}

		// Get status of services
		statuses := getServiceStatus(ids)

		for i := range defines {
			defines[i].Status = statuses[i]
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

	r.Post("/api/getDefinitionFolders", func(w http.ResponseWriter, r *http.Request) {
		jsonStr, err := json.Marshal(config.ServiceDefinitions)

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

	r.Post("/api/restartServer", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			cmd := exec.Command("reboot")
			_ = cmd.Run()
		}()

		w.WriteHeader(http.StatusNoContent)
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

	r.Post("/api/serviceMod", func(w http.ResponseWriter, r *http.Request) {
		act := r.URL.Query().Get("act")

		switch act {
		case "killall":
			// List all services and kill them
			services := []string{"stop"}
			for _, path := range config.ServiceDefinitions {
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

					var sname = strings.TrimSuffix(file.Name(), ".yaml") + ".service"
					if !strings.HasSuffix(file.Name(), ".yaml") {
						sname = file.Name()
					}

					if slices.Contains(config.SrvModBypass, strings.TrimSuffix(sname, ".service")) {
						continue
					}

					services = append(services, sname)
				}
			}

			cmd := exec.Command("systemctl", services...)
			err := cmd.Run()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to kill services."))
				return
			}
		case "startall":
			// List all services and start them
			services := []string{"start"}
			for _, path := range config.ServiceDefinitions {
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

					var sname = strings.TrimSuffix(file.Name(), ".yaml") + ".service"
					if !strings.HasSuffix(file.Name(), ".yaml") {
						sname = file.Name()
					} else {
						var service types.TemplateYaml

						f, err := os.Open(path + "/" + file.Name())

						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte("Failed to read service definition." + err.Error() + file.Name()))
							return
						}

						err = yaml.NewDecoder(f).Decode(&service)

						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte("Failed to read service definition" + err.Error() + file.Name()))
							return
						}

						if service.Broken {
							continue
						}
					}

					if slices.Contains(config.SrvModBypass, strings.TrimSuffix(sname, ".service")) {
						continue
					}

					services = append(services, sname)
				}
			}

			cmd := exec.Command("systemctl", services...)

			err := cmd.Run()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to start services."))
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
