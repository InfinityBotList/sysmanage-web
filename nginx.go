package main

import (
	"crypto/tls"
	"errors"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"sysmanage-web/types"

	"gopkg.in/yaml.v3"
)

func loadNginxMeta() (types.NginxMeta, error) {
	open, err := os.Open(config.NginxDefinitions + "/_meta.yaml")

	if err != nil {
		return types.NginxMeta{}, errors.New("Failed to open _meta.yaml file in " + config.NginxDefinitions + ": " + err.Error())
	}

	var meta types.NginxMeta

	err = yaml.NewDecoder(open).Decode(&meta)

	if err != nil {
		return types.NginxMeta{}, errors.New("Failed to decode _meta.yaml file in " + config.NginxDefinitions + ": " + err.Error())
	}

	// Validate meta
	err = v.Struct(meta)

	if err != nil {
		return types.NginxMeta{}, errors.New("Failed to validate _meta.yaml file in " + config.NginxDefinitions + ": " + err.Error())
	}

	return meta, nil
}

func getNginxDomainList() ([]types.NginxServerManage, error) {
	// Get all files in path
	fsd, err := os.ReadDir(config.NginxDefinitions)

	if err != nil {
		return nil, errors.New("Failed to read nginx definition " + err.Error())
	}

	servers := make([]types.NginxServerManage, 0)

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

		// Read file into NginxServer
		f, err := os.Open(config.NginxDefinitions + "/" + file.Name())

		if err != nil {
			return nil, errors.New("Failed to read nginx definition " + err.Error() + file.Name())
		}

		// Read file into NginxServer
		var server types.NginxYaml

		err = yaml.NewDecoder(f).Decode(&server)

		if err != nil {
			return nil, errors.New("Failed to decode nginx definition " + err.Error() + file.Name())
		}

		if len(server.Servers) == 0 {
			server.Servers = []types.NginxServer{}
		}

		servers = append(servers, types.NginxServerManage{
			Domain: strings.ReplaceAll(strings.TrimSuffix(file.Name(), ".yaml"), "-", "."),
			Server: server,
		})
	}

	return servers, nil
}

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
					if v == "@root" {
						parsedSlice = append(parsedSlice, domain)
						continue
					}
					parsedSlice = append(parsedSlice, v+"."+domain)
				}

				return strings.Join(parsedSlice, " ")
			},
			"ParseOpts": func(opts []string) string {
				if len(opts) == 0 {
					return ""
				}

				var parsedSlice []string

				for _, v := range opts {
					if strings.HasSuffix(v, ";") {
						parsedSlice = append(parsedSlice, v)
					} else {
						parsedSlice = append(parsedSlice, v+";")
					}
				}

				return "\n\t\t" + strings.Join(parsedSlice, "\n\t\t")
			},
		}).Parse(nginxTemplate),
	)

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

		if len(nginxCfg.Servers) == 0 {
			logMap.Add(reqId, "ERROR: Failed to find servers in nginx definition "+file.Name()+", skipping...", true)
			continue
		}

		// Validate nginx definition
		err = v.Struct(nginxCfg)

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to validate nginx definition "+file.Name()+": "+err.Error(), true)
			return
		}

		// Create file
		outFile := strings.TrimSuffix(file.Name(), ".yaml") + ".conf"
		out, err := os.Create("/etc/nginx/conf.d/" + outFile)

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
			MetaCommon: func() string {
				return strings.Join(strings.Split(meta.Common, "\n"), "\n\t")
			}(),
		})

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to execute nginx template "+outFile+": "+err.Error(), true)
			return
		}

		logMap.Add(reqId, "Created nginx file /etc/nginx/conf.d/"+outFile, true)

		// Run nginx -t to validate config
		cmd := exec.Command("nginx", "-t")

		cmd.Stdout = autoLogger{id: reqId}
		cmd.Stderr = autoLogger{id: reqId, Error: true}

		err = cmd.Run()

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to validate nginx config: "+err.Error(), true)
			return
		}

		// Restart nginx
		cmd = exec.Command("systemctl", "restart", "nginx")

		cmd.Stdout = autoLogger{id: reqId}
		cmd.Stderr = autoLogger{id: reqId, Error: true}

		err = cmd.Run()

		if err != nil {
			logMap.Add(reqId, "ERROR: Failed to restart nginx: "+err.Error(), true)
			return
		}

		logMap.Add(reqId, "Restarted nginx", true)
	}
}
