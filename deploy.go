package main

import (
	"os"
	"os/exec"
	"strings"
	"sysmanage-web/types"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func initDeploy(logId string, srv types.ServiceManage) {
	if srv.Service.Git == nil {
		logMap.Add(logId, "FATAL: Service does not have git setup.", true)
		return
	}

	if srv.Service.Broken {
		logMap.Add(logId, "FATAL: Service is marked as broken.", true)
		return
	}

	defer logMap.Persist(logId) // Persist log on exit
	defer logMap.MarkDone(logId)

	logMap.Add(logId, "Started deploy on: "+time.Now().Format(time.RFC3339), true)

	logMap.Add(logId, "Service Repo:"+srv.Service.Git.Repo, true)

	logMap.Add(logId, "Waiting for other builds to finish...", true)

	lsOp.Lock()
	defer lsOp.Unlock()

	deployFolder := "deploys/" + logId

	err := os.MkdirAll(deployFolder, 0755)

	if err != nil {
		logMap.Add(logId, "Error creating deploy folder: "+err.Error(), true)
		return
	}

	logMap.Add(logId, "Cloning repo to "+deployFolder, true)

	_, err = git.PlainClone(deployFolder, false, &git.CloneOptions{
		URL: srv.Service.Git.Repo,
		Auth: &githttp.BasicAuth{
			Username: config.GithubPat,
			Password: config.GithubPat,
		},
		Progress:      autoLogger{id: logId},
		ReferenceName: plumbing.ReferenceName(srv.Service.Git.Ref),
	})

	if err != nil {
		logMap.Add(logId, "Error cloning repo: "+err.Error(), true)
		return
	}

	for _, command := range srv.Service.Git.BuildCommands {
		// Split cmd into args
		args := strings.Split(command, " ")

		// Run the command
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = "deploys/" + logId
		cmd.Env = os.Environ()

		for k, v := range srv.Service.Git.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}

		cmd.Stdout = autoLogger{id: logId}
		cmd.Stderr = autoLogger{id: logId, Error: true}

		err = cmd.Run()

		if err != nil {
			logMap.Add(logId, "Error running command: "+err.Error(), true)
			return
		}
	}

	// Remove deploy.Git.Path and move deploys/deployID to deploy.Git.Path
	err = os.Rename(srv.Service.Directory, srv.Service.Directory+"-old")

	if err != nil {
		logMap.Add(logId, "Error moving old directory: "+err.Error(), true)
		return
	}

	err = os.Rename(deployFolder, srv.Service.Directory)

	if err != nil {
		logMap.Add(logId, "Error moving new directory: "+err.Error(), true)

		// Move old directory back
		os.RemoveAll(srv.Service.Directory)
		os.Rename(srv.Service.Directory+"-old", srv.Service.Directory)

		return
	}

	// Remove old directory
	os.RemoveAll(srv.Service.Directory + "-old")

	logMap.Add(logId, "Deploy finished on: "+time.Now().Format(time.RFC3339), true)

	// Run systemctl restart deploy.Git.Service
	cmd := exec.Command("systemctl", "restart", srv.ID)
	cmd.Env = os.Environ()
	cmd.Stdout = autoLogger{id: logId}
	cmd.Stderr = autoLogger{id: logId, Error: true}

	err = cmd.Run()

	if err != nil {
		logMap.Add(logId, "Error restarting service: "+err.Error(), true)
	}

	logMap.Add(logId, "Service restarted on: "+time.Now().Format(time.RFC3339), false)
}

/*
func loadDeployRoutes(r *chi.Mux) {
	r.Post("/__external__/github", func(w http.ResponseWriter, r *http.Request) {})
}
*/
