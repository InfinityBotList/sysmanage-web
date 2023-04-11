package main

import (
	"encoding/json"
	"html/template"
	"io"
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

		go func() {
			defer markLogDone(reqId)

			addToLog(reqId, "Waiting for other builds to finish...", true)

			inDeploy.Lock()
			defer inDeploy.Unlock()

			addToLog(reqId, "Starting build process to convert service templates to systemd services...", true)

			servicesToEnable := []string{}
			servicesToDisable := []string{}
			for _, folderDef := range config.ServiceDefinitions {
				// First load in the _meta.yaml file from the folder
				open, err := os.Open(folderDef + "/_meta.yaml")

				if err != nil {
					addToLog(reqId, "ERROR: Failed to open _meta.yaml file in "+folderDef, true)
					return
				}

				var meta types.MetaYAML

				err = yaml.NewDecoder(open).Decode(&meta)

				if err != nil {
					addToLog(reqId, "ERROR: Failed to decode _meta.yaml file in "+folderDef+": "+err.Error(), true)
					return
				}

				// Validate meta
				err = v.Struct(meta)

				if err != nil {
					addToLog(reqId, "ERROR: Failed to validate _meta.yaml file in "+folderDef+": "+err.Error(), true)
					return
				}

				var targetTemplate = template.Must(template.New("target").Parse(targetTemplate))

				for _, target := range meta.Targets {
					outFile := config.ServiceOut + "/" + target.Name + ".target"

					// Create file
					out, err := os.Create(outFile)

					if err != nil {
						addToLog(reqId, "ERROR: Failed to create target file "+outFile+": "+err.Error(), true)
						return
					}

					defer out.Close()

					err = targetTemplate.Execute(out, target)

					if err != nil {
						addToLog(reqId, "ERROR: Failed to execute target template "+outFile+": "+err.Error(), true)
						return
					}

					addToLog(reqId, "Created target file "+outFile, true)
				}

				// Next load every service definition in the folder
				fsd, err := os.ReadDir(folderDef)

				if err != nil {
					addToLog(reqId, "ERROR: Failed to read service definition "+err.Error(), true)
					return
				}

				var serviceTemplate = template.Must(template.New("service").Parse(serviceTemplate))

				for _, file := range fsd {
					if file.Name() == "_meta.yaml" {
						addToLog(reqId, "Skipping _meta.yaml as already parsed", true)
						continue // Skip _meta.yaml
					}

					if strings.HasSuffix(file.Name(), ".service") {
						// Copy to service out
						addToLog(reqId, "Copying "+file.Name()+" to service out", true)

						// Open source
						src, err := os.Open(folderDef + "/" + file.Name())

						if err != nil {
							addToLog(reqId, "ERROR: Failed to open service definition "+file.Name()+": "+err.Error(), true)
							return
						}

						defer src.Close()

						// Open destination
						dst, err := os.Create(config.ServiceOut + "/" + file.Name())

						if err != nil {
							addToLog(reqId, "ERROR: Failed to create service definition "+file.Name()+": "+err.Error(), true)
							return
						}

						defer dst.Close()

						// Copy
						_, err = io.Copy(dst, src)

						if err != nil {
							addToLog(reqId, "ERROR: Failed to copy service definition "+file.Name()+": "+err.Error(), true)
							return
						}

						// Enable service in systemd
						servicesToEnable = append(servicesToEnable, file.Name())
						continue
					}

					open, err := os.Open(folderDef + "/" + file.Name())

					if err != nil {
						addToLog(reqId, "ERROR: Failed to open service definition "+file.Name()+": "+err.Error(), true)
						return
					}

					defer open.Close()

					var service types.TemplateYaml

					err = yaml.NewDecoder(open).Decode(&service)

					if err != nil {
						addToLog(reqId, "ERROR: Failed to decode service definition "+file.Name()+": "+err.Error(), true)
						return
					}

					// Validate service
					err = v.Struct(service)

					if err != nil {
						addToLog(reqId, "ERROR: Failed to validate service definition "+file.Name()+": "+err.Error(), true)
						return
					}

					if strings.Contains(service.Target, ".") {
						addToLog(reqId, "ERROR: Target name cannot contain a period, not adding service...", true)
						continue
					}

					if strings.Contains(service.After, ".") {
						addToLog(reqId, "ERROR: After name cannot contain a period, not adding service...", true)
						continue
					}

					targetNames := []string{}

					for _, target := range meta.Targets {
						targetNames = append(targetNames, target.Name)
					}

					if !slices.Contains(targetNames, service.Target) {
						addToLog(reqId, "ERROR: Target "+service.Target+" does not exist, not adding service...", true)
						continue
					}

					outFile := config.ServiceOut + "/" + strings.TrimSuffix(file.Name(), ".yaml") + ".service"

					// Create file
					out, err := os.Create(outFile)

					if err != nil {
						addToLog(reqId, "ERROR: Failed to create service file "+outFile+": "+err.Error(), true)
						return
					}

					defer out.Close()

					err = serviceTemplate.Execute(out, service)

					if err != nil {
						addToLog(reqId, "ERROR: Failed to execute service template "+outFile+": "+err.Error(), true)
						return
					}

					addToLog(reqId, "Created service file "+outFile, true)

					// Enable service in systemd
					if service.Broken {
						servicesToDisable = append(servicesToDisable, strings.TrimSuffix(file.Name(), ".yaml")+".service")
					} else {
						servicesToEnable = append(servicesToEnable, strings.TrimSuffix(file.Name(), ".yaml")+".service")
					}
				}
			}

			addToLog(reqId, "Finished building services.", true)

			// Now we need to reload systemd
			addToLog(reqId, "Reloading systemd...", true)

			cmd := exec.Command("systemctl", "daemon-reload")
			cmd.Stdout = autoLogger{id: reqId}
			cmd.Stderr = autoLogger{id: reqId, Error: true}

			err := cmd.Run()

			if err != nil {
				addToLog(reqId, "ERROR: Failed to reload systemd: "+err.Error(), true)
				return
			}

			addToLog(reqId, "Finished reloading systemd.", true)

			// Now we need to enable the services
			addToLog(reqId, "Enabling services...", true)

			// Prepend "enable" to the list of
			servicesToEnable = append([]string{"enable"}, servicesToEnable...)

			cmd = exec.Command("systemctl", servicesToEnable...)
			cmd.Stdout = autoLogger{id: reqId}
			cmd.Stderr = autoLogger{id: reqId, Error: true}

			err = cmd.Run()

			if err != nil {
				addToLog(reqId, "ERROR: Failed to enable services "+strings.Join(servicesToEnable, ","), true)
				return
			}

			addToLog(reqId, "Finished enabling services.", true)

			// Now we need to disable the services
			addToLog(reqId, "Disabling broken services...", true)

			// Prepend "disable" to the list of
			servicesToDisable = append([]string{"disable"}, servicesToDisable...)

			cmd = exec.Command("systemctl", servicesToDisable...)
			cmd.Stdout = autoLogger{id: reqId}
			cmd.Stderr = autoLogger{id: reqId, Error: true}

			err = cmd.Run()

			if err != nil {
				addToLog(reqId, "ERROR: Failed to disable services "+strings.Join(servicesToDisable, ","), true)
				return
			}

			addToLog(reqId, "Finished disabling broken services.", true)

			addToLog(reqId, "Finished building services.", true)
		}()

		w.Write([]byte(reqId))
	})
}
