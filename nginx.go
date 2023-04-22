package main

import (
	"crypto/tls"
	"html/template"
	"os"
	"strings"
	"sysmanage-web/types"

	"gopkg.in/yaml.v3"
)

func buildNginx(reqId string) {
	defer logMap.MarkDone(reqId)

	logMap.Add(reqId, "Waiting for other builds to finish...", true)

	lsOp.Lock()
	defer lsOp.Unlock()

	logMap.Add(reqId, "Starting build process to convert nginx templates to nginx config files...", true)

	// First load in the _meta.yaml file from the folder
	open, err := os.Open(config.NginxDefinitions + "/_meta.yaml")

	if err != nil {
		logMap.Add(reqId, "ERROR: Failed to open _meta.yaml file in "+config.NginxDefinitions, true)
		return
	}

	var meta types.NginxMeta

	err = yaml.NewDecoder(open).Decode(&meta)

	if err != nil {
		logMap.Add(reqId, "ERROR: Failed to decode _meta.yaml file in "+config.NginxDefinitions+": "+err.Error(), true)
		return
	}

	// Validate meta
	err = v.Struct(meta)

	if err != nil {
		logMap.Add(reqId, "ERROR: Failed to validate _meta.yaml file in "+config.NginxDefinitions+": "+err.Error(), true)
		return
	}

	// Next load every nginx definition in the folder
	fsd, err := os.ReadDir(config.NginxDefinitions)

	if err != nil {
		logMap.Add(reqId, "ERROR: Failed to read nginx definition "+err.Error(), true)
		return
	}

	var nginxTemplate = template.Must(
		template.New("nginx").Funcs(template.FuncMap{
			"ConcatNames": func(domain string, s []string) string {
				var parsedSlice []string

				for _, v := range s {
					parsedSlice = append(parsedSlice, v+"."+domain)
				}

				return strings.Join(parsedSlice, ", ")
			},
		}).Parse(nginxTemplate),
	)

	// Create a temp folder for nginx to use
	ngxDir, err := os.MkdirTemp("", ".ngx-temp")

	if err != nil {
		logMap.Add(reqId, "ERROR: Failed to create temp folder for nginx: "+err.Error(), true)
		return
	}

	defer os.RemoveAll(ngxDir)

	logMap.Add(reqId, "Created temp folder for nginx: "+ngxDir, true)

	for _, file := range fsd {
		if file.Name() == "_meta.yaml" {
			logMap.Add(reqId, "Skipping _meta.yaml as already parsed", true)
			continue // Skip _meta.yaml
		}

		if !strings.HasSuffix(file.Name(), ".yaml") {
			logMap.Add(reqId, "Skipping "+file.Name()+" as not a yaml file", true)
			continue // Skip non-yaml files
		}

		// Create certfile and keyfile from file.Name
		certFile := meta.NginxCertPath + "/cert-" + strings.TrimSuffix(file.Name(), ".yaml") + ".pem"
		keyFile := meta.NginxCertPath + "/key-" + strings.TrimSuffix(file.Name(), ".yaml") + ".pem"

		// Ensure certfile and keyfile exist
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			logMap.Add(reqId, "SANITY FAILED: Failed to find required certfile "+certFile, true)
			return
		}

		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			logMap.Add(reqId, "SANITY FAILED: Failed to find required keyfile "+keyFile, true)
			return
		}

		// Try parsing certfile and keyfile
		_, err = tls.LoadX509KeyPair(certFile, keyFile)

		if err != nil {
			logMap.Add(reqId, "SANITY FAILED: Failed to parse certfile "+certFile+" and keyfile "+keyFile+": "+err.Error(), true)
			return
		}

		open, err := os.Open(config.NginxDefinitions + "/" + file.Name())

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to open nginx definition "+file.Name()+": "+err.Error(), true)
			return
		}

		defer open.Close()

		var nginxCfg types.NginxYaml

		err = yaml.NewDecoder(open).Decode(&nginxCfg)

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to decode nginx definition "+file.Name()+": "+err.Error(), true)
			return
		}

		// Validate nginx definition
		err = v.Struct(nginxCfg)

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to validate nginx definition "+file.Name()+": "+err.Error(), true)
			return
		}

		// Create file
		outFile := strings.TrimSuffix(file.Name(), ".yaml") + ".conf"
		out, err := os.Create(ngxDir + "/" + outFile)

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to create config file "+outFile+": "+err.Error(), true)
			return
		}

		defer out.Close()

		err = nginxTemplate.Execute(out, types.NginxTemplate{
			Servers:  nginxCfg.Servers,
			Meta:     meta,
			Domain:   strings.ReplaceAll(strings.TrimSuffix(file.Name(), ".yaml"), "-", "."),
			CertFile: certFile,
			KeyFile:  keyFile,
		})

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to execute nginx template "+outFile+": "+err.Error(), true)
			return
		}

		logMap.Add(reqId, "Created nginx file "+outFile, true)

		// DEBUG: Move to .debug folder
		err = os.Rename(ngxDir+"/"+outFile, ".debug/"+outFile)

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to move nginx file to debug folder "+outFile+": "+err.Error(), true)
			return
		}
	}
}
