package systemd

import (
	"errors"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/infinitybotlist/sysmanage-web/core/logger"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/plugins/persist"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

func GetServiceStatus(ids []string) []string {
	ids = append([]string{"check"}, ids...)
	cmd := exec.Command("systemctl", ids...)
	out, _ := cmd.CombinedOutput()

	return strings.Split(string(out), "\n")
}

func GetServiceList(getStatus bool) ([]ServiceManage, error) {
	// Get all files in path
	fsd, err := os.ReadDir(serviceDefinitions)

	if err != nil {
		return nil, errors.New("Failed to read service definition " + err.Error())
	}

	services := make([]ServiceManage, 0)
	ids := make([]string, 0)

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
		f, err := os.Open(serviceDefinitions + "/" + file.Name())

		if err != nil {
			return nil, errors.New("Failed to read service definition " + err.Error() + file.Name())
		}

		// Read file into TemplateYaml
		var service TemplateYaml

		err = yaml.NewDecoder(f).Decode(&service)

		if err != nil {
			return nil, errors.New("Failed to read service definition " + err.Error() + file.Name())
		}

		// Service name is the name without .yaml
		sname := strings.TrimSuffix(file.Name(), ".yaml")

		services = append(services, ServiceManage{
			Service: service,
			ID:      sname,
		})

		ids = append(ids, sname)
	}

	// Get status of services
	if getStatus {
		statuses := GetServiceStatus(ids)

		for i := range services {
			services[i].Status = statuses[i]
		}
	}

	return services, nil
}

func BuildServices(reqId string) {
	defer logger.LogMap.MarkDone(reqId)

	logger.LogMap.Add(reqId, "Waiting for other builds to finish...", true)

	state.LsOp.Lock()
	defer state.LsOp.Unlock()

	logger.LogMap.Add(reqId, "Starting build process to convert service templates to systemd services...", true)

	servicesToEnable := []string{}
	servicesToDisable := []string{}
	// First load in the _meta.yaml file from the folder
	open, err := os.Open(serviceDefinitions + "/_meta.yaml")

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to open _meta.yaml file in "+serviceDefinitions, true)
		return
	}

	var meta MetaYAML

	err = yaml.NewDecoder(open).Decode(&meta)

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to decode _meta.yaml file in "+serviceDefinitions+": "+err.Error(), true)
		return
	}

	// Validate meta
	err = state.Validator.Struct(meta)

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to validate _meta.yaml file in "+serviceDefinitions+": "+err.Error(), true)
		return
	}

	var targetTemplate = template.Must(template.New("target").Parse(targetTemplate))

	for _, target := range meta.Targets {
		outFile := serviceOut + "/" + target.Name + ".target"

		// Create file
		out, err := os.Create(outFile)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to create target file "+outFile+": "+err.Error(), true)
			return
		}

		defer out.Close()

		err = targetTemplate.Execute(out, target)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to execute target template "+outFile+": "+err.Error(), true)
			return
		}

		logger.LogMap.Add(reqId, "Created target file "+outFile, true)
	}

	// Next load every service definition in the folder
	fsd, err := os.ReadDir(serviceDefinitions)

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to read service definition "+err.Error(), true)
		return
	}

	var serviceTemplate = template.Must(template.New("service").Parse(serviceTemplate))

	for _, file := range fsd {
		if file.Name() == "_meta.yaml" {
			logger.LogMap.Add(reqId, "Skipping _meta.yaml as already parsed", true)
			continue // Skip _meta.yaml
		}

		if strings.HasSuffix(file.Name(), ".service") {
			// Copy to service out
			logger.LogMap.Add(reqId, "Copying "+file.Name()+" to service out", true)

			// Open source
			src, err := os.Open(serviceDefinitions + "/" + file.Name())

			if err != nil {
				logger.LogMap.Add(reqId, "ERROR: Failed to open service definition "+file.Name()+": "+err.Error(), true)
				return
			}

			defer src.Close()

			// Open destination
			dst, err := os.Create(serviceOut + "/" + file.Name())

			if err != nil {
				logger.LogMap.Add(reqId, "ERROR: Failed to create service definition "+file.Name()+": "+err.Error(), true)
				return
			}

			defer dst.Close()

			// Copy
			_, err = io.Copy(dst, src)

			if err != nil {
				logger.LogMap.Add(reqId, "ERROR: Failed to copy service definition "+file.Name()+": "+err.Error(), true)
				return
			}

			// Enable service in systemd
			servicesToEnable = append(servicesToEnable, file.Name())
			continue
		}

		open, err := os.Open(serviceDefinitions + "/" + file.Name())

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to open service definition "+file.Name()+": "+err.Error(), true)
			return
		}

		defer open.Close()

		var service TemplateYaml

		err = yaml.NewDecoder(open).Decode(&service)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to decode service definition "+file.Name()+": "+err.Error(), true)
			return
		}

		// Validate service
		err = state.Validator.Struct(service)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to validate service definition "+file.Name()+": "+err.Error(), true)
			return
		}

		if strings.Contains(service.Target, ".") {
			logger.LogMap.Add(reqId, "ERROR: Target name cannot contain a period, not adding service...", true)
			continue
		}

		if strings.Contains(service.After, ".") {
			logger.LogMap.Add(reqId, "ERROR: After name cannot contain a period, not adding service...", true)
			continue
		}

		if service.User == "" {
			service.User = "root"
		}

		if service.Group == "" {
			service.Group = "root"
		}

		targetNames := []string{}

		for _, target := range meta.Targets {
			targetNames = append(targetNames, target.Name)
		}

		if !slices.Contains(targetNames, service.Target) {
			logger.LogMap.Add(reqId, "ERROR: Target "+service.Target+" does not exist, not adding service...", true)
			continue
		}

		outFile := serviceOut + "/" + strings.TrimSuffix(file.Name(), ".yaml") + ".service"

		// Create file
		out, err := os.Create(outFile)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to create service file "+outFile+": "+err.Error(), true)
			return
		}

		defer out.Close()

		err = serviceTemplate.Execute(out, service)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to execute service template "+outFile+": "+err.Error(), true)
			return
		}

		logger.LogMap.Add(reqId, "Created service file "+outFile, true)

		// Enable service in systemd
		if service.Broken {
			servicesToDisable = append(servicesToDisable, strings.TrimSuffix(file.Name(), ".yaml")+".service")
		} else {
			servicesToEnable = append(servicesToEnable, strings.TrimSuffix(file.Name(), ".yaml")+".service")
		}
	}

	logger.LogMap.Add(reqId, "Finished building services.", true)

	// Now we need to reload systemd
	logger.LogMap.Add(reqId, "Reloading systemd...", true)

	cmd := exec.Command("systemctl", "daemon-reload")
	cmd.Stdout = logger.AutoLogger{ID: reqId}
	cmd.Stderr = logger.AutoLogger{ID: reqId, Error: true}

	err = cmd.Run()

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to reload systemd: "+err.Error(), true)
		return
	}

	logger.LogMap.Add(reqId, "Finished reloading systemd.", true)

	// Now we need to enable the services
	logger.LogMap.Add(reqId, "Enabling services...: "+strings.Join(servicesToEnable, ","), true)

	// Prepend "enable" to the list of
	servicesToEnable = append([]string{"enable"}, servicesToEnable...)

	cmd = exec.Command("systemctl", servicesToEnable...)
	cmd.Stdout = logger.AutoLogger{ID: reqId}
	cmd.Stderr = logger.AutoLogger{ID: reqId, Error: true}

	err = cmd.Run()

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to enable services "+strings.Join(servicesToEnable, ","), true)
		return
	}

	logger.LogMap.Add(reqId, "Finished enabling services.", true)

	// Now we need to disable the services
	logger.LogMap.Add(reqId, "Disabling broken services...", true)

	// Prepend "disable" to the list of
	servicesToDisable = append([]string{"disable"}, servicesToDisable...)

	cmd = exec.Command("systemctl", servicesToDisable...)
	cmd.Stdout = logger.AutoLogger{ID: reqId}
	cmd.Stderr = logger.AutoLogger{ID: reqId, Error: true}

	err = cmd.Run()

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to disable services "+strings.Join(servicesToDisable, ","), true)
		return
	}

	logger.LogMap.Add(reqId, "Finished disabling broken services.", true)

	err = persist.PersistToGit(reqId)

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to persist to git: "+err.Error(), true)
		return
	}

	logger.LogMap.Add(reqId, "Finished building services.", true)
}
