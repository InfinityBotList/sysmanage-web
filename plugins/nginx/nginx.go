package nginx

import (
	"crypto/tls"
	"errors"
	"html/template"
	"net"
	"os"
	"os/exec"
	"strings"
	"sysmanage-web/core/logger"
	"sysmanage-web/core/state"

	"gopkg.in/yaml.v3"
)

func loadNginxMeta() (NginxMeta, error) {
	open, err := os.Open(nginxDefinitions + "/_meta.yaml")

	if err != nil {
		return NginxMeta{}, errors.New("Failed to open _meta.yaml file in " + nginxDefinitions + ": " + err.Error())
	}

	var meta NginxMeta

	err = yaml.NewDecoder(open).Decode(&meta)

	if err != nil {
		return NginxMeta{}, errors.New("Failed to decode _meta.yaml file in " + nginxDefinitions + ": " + err.Error())
	}

	// Validate meta
	err = state.Validator.Struct(meta)

	if err != nil {
		return NginxMeta{}, errors.New("Failed to validate _meta.yaml file in " + nginxDefinitions + ": " + err.Error())
	}

	return meta, nil
}

func getNginxDomainList() ([]NginxServerManage, error) {
	// Get all files in path
	fsd, err := os.ReadDir(nginxDefinitions)

	if err != nil {
		return nil, errors.New("Failed to read nginx definition " + err.Error())
	}

	servers := make([]NginxServerManage, 0)

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
		f, err := os.Open(nginxDefinitions + "/" + file.Name())

		if err != nil {
			return nil, errors.New("Failed to read nginx definition " + err.Error() + file.Name())
		}

		// Read file into NginxServer
		var server NginxYaml

		err = yaml.NewDecoder(f).Decode(&server)

		if err != nil {
			return nil, errors.New("Failed to decode nginx definition " + err.Error() + file.Name())
		}

		if len(server.Servers) == 0 {
			server.Servers = []NginxServer{}
		}

		servers = append(servers, NginxServerManage{
			Domain: strings.ReplaceAll(strings.TrimSuffix(file.Name(), ".yaml"), "-", "."),
			Server: server,
		})
	}

	return servers, nil
}

func buildNginx(reqId string) {
	defer logger.LogMap.MarkDone(reqId)

	logger.LogMap.Add(reqId, "Waiting for other builds to finish...", true)

	state.LsOp.Lock()
	defer state.LsOp.Unlock()

	logger.LogMap.Add(reqId, "Starting build process to convert nginx templates to nginx config files...", true)

	// First load in the _meta.yaml file from the folder
	open, err := os.Open(nginxDefinitions + "/_meta.yaml")

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to open _meta.yaml file in "+nginxDefinitions, true)
		return
	}

	var meta NginxMeta

	err = yaml.NewDecoder(open).Decode(&meta)

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to decode _meta.yaml file in "+nginxDefinitions+": "+err.Error(), true)
		return
	}

	// Validate meta
	err = state.Validator.Struct(meta)

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to validate _meta.yaml file in "+nginxDefinitions+": "+err.Error(), true)
		return
	}

	// Next load every nginx definition in the folder
	fsd, err := os.ReadDir(nginxDefinitions)

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to read nginx definition "+err.Error(), true)
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
			logger.LogMap.Add(reqId, "Skipping _meta.yaml as already parsed", true)
			continue // Skip _meta.yaml
		}

		if !strings.HasSuffix(file.Name(), ".yaml") {
			logger.LogMap.Add(reqId, "Skipping "+file.Name()+" as not a yaml file", true)
			continue // Skip non-yaml files
		}

		// Create certfile and keyfile from file.Name
		certFile := meta.NginxCertPath + "/cert-" + strings.TrimSuffix(file.Name(), ".yaml") + ".pem"
		keyFile := meta.NginxCertPath + "/key-" + strings.TrimSuffix(file.Name(), ".yaml") + ".pem"

		// Ensure certfile and keyfile exist
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			logger.LogMap.Add(reqId, "SANITY FAILED: Failed to find required certfile "+certFile, true)
			return
		}

		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			logger.LogMap.Add(reqId, "SANITY FAILED: Failed to find required keyfile "+keyFile, true)
			return
		}

		// Try parsing certfile and keyfile
		_, err = tls.LoadX509KeyPair(certFile, keyFile)

		if err != nil {
			logger.LogMap.Add(reqId, "SANITY FAILED: Failed to parse certfile "+certFile+" and keyfile "+keyFile+": "+err.Error(), true)
			return
		}

		open, err := os.Open(nginxDefinitions + "/" + file.Name())

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to open nginx definition "+file.Name()+": "+err.Error(), true)
			return
		}

		defer open.Close()

		var nginxCfg NginxYaml

		err = yaml.NewDecoder(open).Decode(&nginxCfg)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to decode nginx definition "+file.Name()+": "+err.Error(), true)
			return
		}

		if len(nginxCfg.Servers) == 0 {
			logger.LogMap.Add(reqId, "ERROR: Failed to find servers in nginx definition "+file.Name()+", skipping...", true)
			continue
		}

		// Validate nginx definition
		err = state.Validator.Struct(nginxCfg)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to validate nginx definition "+file.Name()+": "+err.Error(), true)
			return
		}

		// Create file
		outFile := strings.TrimSuffix(file.Name(), ".yaml") + ".conf"
		out, err := os.Create("/etc/nginx/conf.d/" + outFile)

		if err != nil {
			logger.LogMap.Add(reqId, "ERROR: Failed to create config file "+outFile+": "+err.Error(), true)
			return
		}

		defer out.Close()

		err = nginxTemplate.Execute(out, NginxTemplate{
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
			logger.LogMap.Add(reqId, "ERROR: Failed to execute nginx template "+outFile+": "+err.Error(), true)
			return
		}

		logger.LogMap.Add(reqId, "Created nginx file /etc/nginx/conf.d/"+outFile, true)
	}

	// Run nginx -t to validate config
	cmd := exec.Command("nginx", "-t")

	cmd.Stdout = logger.AutoLogger{ID: reqId}
	cmd.Stderr = logger.AutoLogger{ID: reqId, Error: true}

	err = cmd.Run()

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to validate nginx config: "+err.Error(), true)
		return
	}

	// Restart nginx
	cmd = exec.Command("systemctl", "restart", "nginx")

	cmd.Stdout = logger.AutoLogger{ID: reqId}
	cmd.Stderr = logger.AutoLogger{ID: reqId, Error: true}

	err = cmd.Run()

	if err != nil {
		logger.LogMap.Add(reqId, "ERROR: Failed to restart nginx: "+err.Error(), true)
		return
	}

	logger.LogMap.Add(reqId, "Restarted nginx", true)
}

// Get preferred outbound ip of this machine
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func updateDnsRecordCf(reqId string) {
	if cf == nil {
		logger.LogMap.Add(reqId, "Not updating DNS, CF is disabled!", true)
		return
	}

	logger.LogMap.Add(reqId, "Updating DNS record", true)

	// Get current IP
	ip, err := getOutboundIP()

	if err != nil {
		logger.LogMap.Add(reqId, "Failed to get IP address:"+err.Error(), true)
	}

	logger.LogMap.Add(reqId, "Current IP address is "+ip.String(), true)

	// Get servers to update
	srv, err := getNginxDomainList()

	if err != nil {
		logger.LogMap.Add(reqId, "Failed to get nginx domain list:"+err.Error(), true)
	}

	for _, s := range srv {
		for _, serverName := range s.Server.Servers {
			for _, name := range serverName.Names {
				domExpanded := name + "." + s.Domain
				logger.LogMap.Add(reqId, "Updating IP of "+domExpanded+" to "+ip.String(), true)

				// Find any existing records

			}
		}
	}
}
