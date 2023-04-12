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

	deployFolder := srv.Service.Directory

	// Check that the deploy folder exists
	if _, err := os.Stat(deployFolder); os.IsNotExist(err) {
		logMap.Add(logId, "Deploy folder does not exist: "+deployFolder+", cloning it", true)

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
	} else {
		logMap.Add(logId, "Pulling "+deployFolder, true)

		// Pull repo
		repo, err := git.PlainOpen(deployFolder)

		if err != nil {
			logMap.Add(logId, "Error opening repo: "+err.Error(), true)
			return
		}

		w, err := repo.Worktree()

		if err != nil {
			logMap.Add(logId, "Error getting worktree: "+err.Error(), true)
			return
		}

		err = w.Pull(&git.PullOptions{
			ReferenceName: plumbing.ReferenceName(srv.Service.Git.Ref),
			Auth: &githttp.BasicAuth{
				Username: config.GithubPat,
				Password: config.GithubPat,
			},
			Progress: autoLogger{id: logId},
		})

		if err != nil && err != git.NoErrAlreadyUpToDate {
			logMap.Add(logId, "Error pulling repo: "+err.Error(), true)
			return
		}
	}

	curDir := deployFolder

	for _, command := range srv.Service.Git.BuildCommands {
		// Split cmd into args
		args := strings.Split(command, " ")

		if args[0] == "cd" && len(args) > 1 {
			if args[1] == ".." {
				split := strings.Split(curDir, "/")

				if len(split) > 1 {
					curDir = strings.Join(split[:len(split)-1], "/")
				} else {
					logMap.Add(logId, "WARN: ", true)
				}
			} else {
				curDir = curDir + "/" + args[1]
			}

			logMap.Add(logId, "Changed directory to "+curDir, true)
			continue // Ignore rest of command
		}

		// Run the command
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = curDir
		cmd.Env = os.Environ()

		for k, v := range srv.Service.Git.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}

		cmd.Stdout = autoLogger{id: logId}
		cmd.Stderr = autoLogger{id: logId, Error: true}

		err := cmd.Run()

		if err != nil {
			logMap.Add(logId, "Error running command: "+err.Error(), true)
			return
		}
	}

	logMap.Add(logId, "Deploy finished on: "+time.Now().Format(time.RFC3339), true)

	// Run systemctl restart deploy.Git.Service
	cmd := exec.Command("systemctl", "restart", srv.ID)
	cmd.Env = os.Environ()
	cmd.Stdout = autoLogger{id: logId}
	cmd.Stderr = autoLogger{id: logId, Error: true}

	err := cmd.Run()

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