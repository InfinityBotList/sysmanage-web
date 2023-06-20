package systemd

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/crypto"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/infinitybotlist/sysmanage-web/core/logger"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/plugins/persist"
)

func loadServiceApi(r chi.Router) {
	// Returns the list of services
	r.Post("/getServiceList", func(w http.ResponseWriter, r *http.Request) {
		serviceList, err := getServiceList(true)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get serviceList."))
			return
		}

		// JSON encode defines
		jsonStr, err := json.Marshal(serviceList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/getDefinitionFolders", func(w http.ResponseWriter, r *http.Request) {
		jsonStr, err := json.Marshal(serviceDefinitions)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/getMeta", func(w http.ResponseWriter, r *http.Request) {
		// Open _meta.yaml
		f, err := os.Open(serviceDefinitions + "/_meta.yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read service definition." + err.Error()))
			return
		}

		// Read file into TemplateYaml
		var metaYaml MetaYAML

		err = yaml.NewDecoder(f).Decode(&metaYaml)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read service definition" + err.Error()))
			return
		}

		// JSON encode defines
		jsonStr, err := json.Marshal(metaYaml)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service definitions."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	r.Post("/createService", func(w http.ResponseWriter, r *http.Request) {
		var createService CreateTemplate

		err := yaml.NewDecoder(r.Body).Decode(&createService)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to decode create service request."))
			return
		}

		// validate createService
		err = state.Validator.Struct(createService)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Open _meta.yaml
		f, err := os.Open(serviceDefinitions + "/_meta.yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read service definition." + err.Error()))
			return
		}

		// Read file into TemplateYaml
		var metaYaml MetaYAML

		err = yaml.NewDecoder(f).Decode(&metaYaml)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read service definition" + err.Error()))
			return
		}

		// ensure service target is in _meta.yaml
		flag := false
		for _, target := range metaYaml.Targets {
			if target.Name == createService.Service.Target {
				flag = true
			}
		}

		if !flag {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service target is not in _meta.yaml."))
			return
		}

		if strings.Contains(createService.Name, " ") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service name cannot contain spaces."))
			return
		}

		if strings.Contains(createService.Name, "/") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service name cannot contain slashes."))
			return
		}

		if strings.Contains(createService.Name, ".") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service name cannot contain dots."))
			return
		}

		// Check if service already exists
		if r.URL.Query().Get("update") != "true" {
			if _, err := os.Stat(serviceDefinitions + "/" + createService.Name + ".yaml"); errors.Is(err, os.ErrExist) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Service already exists:" + err.Error()))
				return
			} else if !errors.Is(err, os.ErrNotExist) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to check if service exists:" + err.Error()))
				return
			}
		} else {
			// Open file and copy git integration into it
			f, err := os.Open(serviceDefinitions + "/" + createService.Name + ".yaml")

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to read service definition." + err.Error()))
				return
			}

			// Read file into TemplateYaml
			var serviceYaml TemplateYaml

			err = yaml.NewDecoder(f).Decode(&serviceYaml)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to read service definition" + err.Error()))
				return
			}

			createService.Service.Git = serviceYaml.Git
		}

		// Create file
		f, err = os.Create(serviceDefinitions + "/" + createService.Name + ".yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to create service file."))
			return
		}

		// Create service
		err = yaml.NewEncoder(f).Encode(createService.Service)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode service."))
			return
		}

		go persist.PersistToGit("")

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/deleteService", func(w http.ResponseWriter, r *http.Request) {
		var deleteService DeleteTemplate

		err := yaml.NewDecoder(r.Body).Decode(&deleteService)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to decode delete service request."))
			return
		}

		// validate deleteService
		err = state.Validator.Struct(deleteService)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if strings.Contains(deleteService.Name, " ") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service name cannot contain spaces."))
			return
		}

		if strings.Contains(deleteService.Name, "/") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service name cannot contain slashes."))
			return
		}

		if strings.Contains(deleteService.Name, ".") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service name cannot contain dots."))
			return
		}

		logId := crypto.RandString(32)

		go func() {
			defer logger.LogMap.MarkDone(logId)

			// delete yaml file, ignore if it doesn't exist
			logger.LogMap.Add(logId, "Deleting service file...", true)

			err = os.Remove(serviceDefinitions + "/" + deleteService.Name + ".yaml")

			if err != nil {
				logger.LogMap.Add(logId, "Failed to delete service file: "+err.Error(), true)
			} else {
				logger.LogMap.Add(logId, "Deleted service file successfully.", true)
			}

			// disable service, ignore if it doesn't exist
			err = exec.Command("systemctl", "disable", deleteService.Name).Run()

			if err != nil {
				logger.LogMap.Add(logId, "Failed to disable service: "+err.Error(), true)
			} else {
				logger.LogMap.Add(logId, "Disabled service successfully.", true)
			}

			// stop service, ignore if it doesn't exist
			err = exec.Command("systemctl", "stop", deleteService.Name).Run()

			if err != nil {
				logger.LogMap.Add(logId, "Failed to stop service: "+err.Error(), true)
			} else {
				logger.LogMap.Add(logId, "Stopped service successfully.", true)
			}

			// delete service file, ignore if it doesn't exist
			err = os.Remove("/etc/systemd/system/" + deleteService.Name + ".service")

			if err != nil {
				logger.LogMap.Add(logId, "Failed to delete service file: "+err.Error(), true)
			} else {
				logger.LogMap.Add(logId, "Deleted service file successfully.", true)
			}

			// reload systemd
			err = exec.Command("systemctl", "daemon-reload").Run()

			if err != nil {
				logger.LogMap.Add(logId, "Failed to reload systemd: "+err.Error(), true)
			} else {
				logger.LogMap.Add(logId, "Reloaded systemd successfully.", true)
			}

			err := persist.PersistToGit(logId)

			if err != nil {
				logger.LogMap.Add(logId, "Failed to persist to git: "+err.Error(), true)
			}
		}()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(logId))
	})

	r.Post("/systemctl", func(w http.ResponseWriter, r *http.Request) {
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

	// Simple goroutine to clean up open entries
	const maxOpenTime = time.Minute * 5 // for now, 5 minutes

	r.Post("/getServiceLogs", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("id")

		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing id parameter."))
			return
		}

		logId := crypto.RandString(64)

		cmd := exec.Command("journalctl", "-u", name, "-n", "50", "-f")
		cmd.Stdout = logger.AutoLogger{ID: logId}
		cmd.Stderr = logger.AutoLogger{ID: logId}
		cmd.Stdin = nil

		logger.LogMap.Add(logId, "goro:start", true)

		go func() {
			logger.LogMap.Add(logId, "Starting logger for "+name, true)

			err := cmd.Run()

			if err != nil {
				logger.LogMap.Add(logId, "Failed to get logs: "+err.Error(), true)
			}

			logger.LogMap.Add(logId, "Logger died:", true)

			logger.LogMap.MarkDone(logId)
		}()

		go func() {
			time.Sleep(maxOpenTime)

			cmd.Process.Kill()

			logger.LogMap.Add(logId, "Max open time reached, closing log.", true)

			logger.LogMap.MarkDone(logId)
			logger.LogMap.Set(logId, logger.LogEntry{LastUpdate: time.Now()})
		}()

		w.Write([]byte(logId))
	})

	r.Post("/restartServer", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			cmd := exec.Command("reboot")
			_ = cmd.Run()
		}()

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/buildServices", func(w http.ResponseWriter, r *http.Request) {
		reqId := crypto.RandString(64)

		go buildServices(reqId)

		w.Write([]byte(reqId))
	})

	r.Post("/serviceMod", func(w http.ResponseWriter, r *http.Request) {
		act := r.URL.Query().Get("act")

		switch act {
		case "killall":
			// List all services and kill them
			services := []string{"stop"}
			fsd, err := os.ReadDir(serviceDefinitions)

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

				if slices.Contains(srvModBypass, strings.TrimSuffix(sname, ".service")) {
					continue
				}

				services = append(services, sname)
			}

			cmd := exec.Command("systemctl", services...)
			err = cmd.Run()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to kill services."))
				return
			}
		case "startall":
			// List all services and start them
			services := []string{"start"}
			fsd, err := os.ReadDir(serviceDefinitions)

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
					var service TemplateYaml

					f, err := os.Open(serviceDefinitions + "/" + file.Name())

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

				if slices.Contains(srvModBypass, strings.TrimSuffix(sname, ".service")) {
					continue
				}

				services = append(services, sname)
			}

			cmd := exec.Command("systemctl", services...)

			err = cmd.Run()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to start services."))
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/initDeploy", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("id")

		services, err := getServiceList(false)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get service list."))
			return
		}

		var gotService *ServiceManage
		for _, service := range services {
			if service.ID == name {
				gotService = &service
				break
			}
		}

		if gotService == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Service not found."))
			return
		}

		logId := crypto.RandString(64)

		go initDeploy(logId, *gotService)
		go persist.PersistToGit("")

		w.Write([]byte(logId))
	})

	r.Post("/createGit", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		services, err := getServiceList(false)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get service list."))
			return
		}

		// Check if service exists
		var gotService *ServiceManage

		for _, service := range services {
			if service.ID == id {
				gotService = &service
				break
			}
		}

		if gotService == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Service not found."))
			return
		}

		var git *Git

		err = yaml.NewDecoder(r.Body).Decode(&git)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to decode git data."))
			return
		}

		err = state.Validator.Struct(git)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid git data."))
			return
		}

		if !strings.HasPrefix(git.Repo, "https://") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid git URL."))
			return
		}

		if !strings.HasPrefix(git.Ref, "refs/heads/") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Currently, only refs starting with refs/heads/ are supported."))
			return
		}

		if len(git.BuildCommands) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No build commands specified."))
			return
		}

		gotService.Service.Git = git

		// Save service
		f, err := os.Create(serviceDefinitions + "/" + gotService.ID + ".yaml-1")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/" + gotService.ID + ".yaml-1")

			return
		}

		err = yaml.NewEncoder(f).Encode(gotService.Service)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/" + gotService.ID + ".yaml-1")

			return
		}

		err = f.Close()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/" + gotService.ID + ".yaml-1")

			return
		}

		err = os.Rename(serviceDefinitions+"/"+gotService.ID+".yaml-1", serviceDefinitions+"/"+gotService.ID+".yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/" + gotService.ID + ".yaml-1")

			return
		}

		go persist.PersistToGit("")
	})

	r.Post("/updateMeta", func(w http.ResponseWriter, r *http.Request) {
		var target MetaTarget

		err := json.NewDecoder(r.Body).Decode(&target)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to decode target data."))
			return
		}

		err = state.Validator.Struct(target)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Open _meta.yaml
		f, err := os.Open(serviceDefinitions + "/_meta.yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read service definition." + err.Error()))
			return
		}

		// Read file into TemplateYaml
		var metaYaml MetaYAML

		err = yaml.NewDecoder(f).Decode(&metaYaml)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to read service definition" + err.Error()))
			return
		}

		f.Close()

		switch r.URL.Query().Get("action") {
		case "create":
			flag := false

			for _, t := range metaYaml.Targets {
				if t.Name == target.Name {
					flag = true
					break
				}
			}

			if flag {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Target with this name already exists."))
				return
			}

			metaYaml.Targets = append(metaYaml.Targets, target)
		case "update":
			name := r.URL.Query().Get("name")

			if name == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("No target name specified."))
				return
			}

			flag := false

			for i, t := range metaYaml.Targets {
				if t.Name == name {
					metaYaml.Targets[i] = target
					flag = true
					break
				}
			}

			if !flag {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Target with this name does not exist."))
				return
			}

		case "delete":
			flag := false

			// Ensure no services are using this target
			serviceList, err := getServiceList(false)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to read service list:" + err.Error()))
			}

			for _, s := range serviceList {
				if s.Service.Target == target.Name {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Target is in use by service " + s.ID + "."))
					return
				}
			}

			// Then delete the target
			for i, t := range metaYaml.Targets {
				if t.Name == target.Name {
					metaYaml.Targets = append(metaYaml.Targets[:i], metaYaml.Targets[i+1:]...)
					flag = true
					break
				}
			}

			if !flag {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Target with this name does not exist."))
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid action"))
			return
		}

		// Save service
		f, err = os.Create(serviceDefinitions + "/_meta.yaml-1")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/_meta.yaml-1")

			return
		}

		err = yaml.NewEncoder(f).Encode(metaYaml)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/_meta.yaml-1")

			return
		}

		err = f.Close()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/_meta.yaml-1")

			return
		}

		err = os.Rename(serviceDefinitions+"/_meta.yaml-1", serviceDefinitions+"/_meta.yaml")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save service."))

			// Close file and delete it
			f.Close()
			os.Remove(serviceDefinitions + "/_meta.yaml-1")

			return
		}

		go persist.PersistToGit("")

		w.WriteHeader(http.StatusOK)
	})
}
