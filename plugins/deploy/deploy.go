package deploy

import (
	"context"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/infinitybotlist/sysmanage-web/core/logger"
)

const scriptTmpl = `
#!/bin/bash

{{range $val := .}}
echo "> {{$val}}"
{{$val}}
{{end}}
`

var templ = template.Must(template.New("script").Parse(scriptTmpl))

func InitDeploy(logId string, d *DeployMeta) {
	if d.Src == nil {
		logger.LogMap.Add(logId, "FATAL: Deploy does not have an associated source setup.", true)
		return
	}

	if d.Broken {
		logger.LogMap.Add(logId, "FATAL: Deploy is marked as broken.", true)
		return
	}

	defer logger.LogMap.MarkDone(logId)

	logger.LogMap.Add(logId, "Started deploy on: "+time.Now().Format(time.RFC3339), true)
	logger.LogMap.Add(logId, "Deploy Source:"+d.Src.String(), true)
	logger.LogMap.Add(logId, "Waiting for builds to finish...", true)

	maxConcurrency++

	waits := 0
	for len(builds) > maxConcurrency {
		waits++

		if waits%5 == 0 {
			// Print log entries
			b := []string{}

			for k := range builds {
				b = append(b, k+": "+builds[k].String())
			}

			logger.LogMap.Add(logId, strings.Join(b, "\n"), true)
		}

		time.Sleep(5 * time.Second)
	}

	breakpoint.Lock()
	builds[logId] = &DeployStatus{
		Source:    d.Src,
		CreatedAt: time.Now(),
	}
	breakpoint.Unlock()

	buildDir := "/tmp/deploys/" + logId + "/output"

	err := os.MkdirAll(buildDir, 0755)

	if err != nil {
		logger.LogMap.Add(logId, "FATAL: could not create build folder ["+buildDir+"]: "+err.Error(), true)
		return
	}

	defer os.RemoveAll(buildDir)

	logger.LogMap.Add(logId, "Output path: "+buildDir, true)

	srcFn, ok := DeploySources[d.Src.Type]

	if !ok {
		logger.LogMap.Add(logId, "FATAL: Unknown deploy source type: "+d.Src.Type, true)
		return
	}

	err = srcFn(logId, buildDir, d)

	if err != nil {
		logger.LogMap.Add(logId, "Error loading source "+d.Src.Type+": "+err.Error(), true)
		return
	}

	// Create script
	f, err := os.Create(buildDir + "/builder")

	if err != nil {
		logger.LogMap.Add(logId, "Error creating script: "+err.Error(), true)
		return
	}

	defer f.Close()

	// Write script
	err = templ.Execute(f, d.Commands)

	if err != nil {
		logger.LogMap.Add(logId, "Error writing script: "+err.Error(), true)
		return
	}

	// Run script using bash as a seperate contained process
	// This is to prevent the script from doing anything malicious
	// to the system

	// Create a new bash process
	ctx := context.Background()
	if d.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(d.Timeout)*time.Second)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, "bash", buildDir+"/builder")
	cmd.Dir = buildDir
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	for k, v := range d.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Stdout = logger.AutoLogger{ID: logId}
	cmd.Stderr = logger.AutoLogger{ID: logId, Error: true}

	err = cmd.Run()

	if err != nil {
		logger.LogMap.Add(logId, "Error running command: "+err.Error(), true)
		return
	}

	// Copy any potential config files to deploy folder
	logger.LogMap.Add(logId, "Copying config files to new deploy", true)
	for _, file := range d.ConfigFiles {
		logger.LogMap.Add(logId, "=> "+file, true)

		f, err := os.Open(d.OutputPath + "/" + file)
		if err != nil {
			logger.LogMap.Add(logId, "WARNING: Could not open config file "+file, true)
			continue
		}
		defer f.Close()

		newF, err := os.Create(buildDir + "/" + file)
		if err != nil {
			logger.LogMap.Add(logId, "Error creating config file: "+err.Error(), true)
			return
		}
		defer newF.Close()

		_, err = newF.ReadFrom(f)

		if err != nil {
			logger.LogMap.Add(logId, "Error copying config file: "+err.Error(), true)
			return
		}
	}

	breakpoint.Lock()
	defer breakpoint.Unlock()

	// Ensure output path exists first before continuing
	err = os.MkdirAll(d.OutputPath, 0755)

	if err != nil {
		logger.LogMap.Add(logId, "Error validating service folder: "+err.Error(), true)
		return
	}

	// Rename service folder
	err = os.Rename(d.OutputPath, d.OutputPath+"-old")

	if err != nil {
		logger.LogMap.Add(logId, "Error renaming service folder: "+err.Error(), true)
		return
	}

	// Move deploy folder to service folder
	err = os.Rename(buildDir, d.OutputPath)

	if err != nil {
		logger.LogMap.Add(logId, "Error moving build directory: "+err.Error(), true)

		// Move old service folder back
		err = os.Rename(d.OutputPath+"-old", d.OutputPath)

		if err != nil {
			logger.LogMap.Add(logId, "Error moving old service folder back: "+err.Error(), true)
		}

		return
	}

	// Remove old service folder
	err = os.RemoveAll(d.OutputPath + "-old")

	if err != nil {
		logger.LogMap.Add(logId, "Error removing old service folder: "+err.Error(), true)
		return
	}

	logger.LogMap.Add(logId, "Deploy finished on: "+time.Now().Format(time.RFC3339), true)
}
