package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sysmanage-web/types"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	deployViaClone := false

	if _, err := os.Stat(deployFolder); os.IsNotExist(err) {
		logMap.Add(logId, "Deploy folder does not exist: "+deployFolder+", cloning it", true)
		deployViaClone = true
	} else if !srv.Service.Git.AllowDirty {
		logMap.Add(logId, "Dirty builds not allowed, performing fresh clone", true)
		deployViaClone = true
	}

	if deployViaClone {
		deployFolder = "deploys/" + logId

		logMap.Add(logId, "Cloning "+srv.Service.Git.Repo, true)
		_, err := git.PlainClone(deployFolder, false, &git.CloneOptions{
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
		logMap.Add(logId, "Pulling into "+deployFolder, true)

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

		// If there are unstaged changes, add+commit+push them
		if status, err := w.Status(); err == nil {
			if !status.IsClean() {
				logMap.Add(logId, "Unstaged changes detected, committing them", true)

				_, err = w.Add(".")

				if err != nil {
					logMap.Add(logId, "Error adding unstaged changes: "+err.Error(), true)
					return
				}

				_, err = w.Commit("Auto commit from sysmanage-web", &git.CommitOptions{
					All:               true,
					AllowEmptyCommits: true,
					Author: &object.Signature{
						Name: "sysmanage-web[auto]",
						When: time.Now(),
					},
				})

				if err != nil {
					logMap.Add(logId, "Error committing unstaged changes: "+err.Error(), true)
					return
				}

				err = repo.Push(&git.PushOptions{
					Auth: &githttp.BasicAuth{
						Username: config.GithubPat,
						Password: config.GithubPat,
					},
					Progress: autoLogger{id: logId},
				})

				// Force push if we get an error
				if err != nil {
					logMap.Add(logId, "Error pushing unstaged changes: "+err.Error(), true)
					logMap.Add(logId, "Force pushing unstaged changes", true)

					err = repo.Push(&git.PushOptions{
						Force: true,
						Auth: &githttp.BasicAuth{
							Username: config.GithubPat,
							Password: config.GithubPat,
						},
						Progress: autoLogger{id: logId},
					})

					if err != nil {
						logMap.Add(logId, "Error force pushing unstaged changes: "+err.Error(), true)
						return
					}
				}
			}
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

	logMap.Add(logId, "Deploy folder: "+deployFolder, true)

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

	if deployViaClone {
		// Copy any potential config files to deploy folder
		logMap.Add(logId, "Copying config files to new deploy", true)
		for _, file := range srv.Service.Git.ConfigFiles {
			logMap.Add(logId, "Copying config file "+file, true)

			f, err := os.Open(srv.Service.Directory + "/" + file)

			if err != nil {
				logMap.Add(logId, "WARNING: Could not open config file "+file, true)
				continue
			}

			defer f.Close()

			newF, err := os.Create(deployFolder + "/" + file)

			if err != nil {
				logMap.Add(logId, "Error creating config file: "+err.Error(), true)
				return
			}

			defer newF.Close()

			_, err = io.Copy(newF, f)

			if err != nil {
				logMap.Add(logId, "Error copying config file: "+err.Error(), true)
				return
			}

			err = newF.Sync()

			if err != nil {
				logMap.Add(logId, "Error syncing config file: "+err.Error(), true)
				return
			}
		}

		// Rename service folder
		err := os.Rename(srv.Service.Directory, srv.Service.Directory+"-old")

		if err != nil {
			logMap.Add(logId, "Error renaming service folder: "+err.Error(), true)
			return
		}

		// Move deploy folder to service folder
		err = os.Rename(deployFolder, srv.Service.Directory)

		if err != nil {
			logMap.Add(logId, "Error moving deploy folder: "+err.Error(), true)

			// Move old service folder back
			err = os.Rename(srv.Service.Directory+"-old", srv.Service.Directory)

			if err != nil {
				logMap.Add(logId, "Error moving old service folder back: "+err.Error(), true)
			}

			return
		}

		// Remove old service folder
		err = os.RemoveAll(srv.Service.Directory + "-old")

		if err != nil {
			logMap.Add(logId, "Error removing old service folder: "+err.Error(), true)
			return
		}
	}

	logMap.Add(logId, "Deploy finished on: "+time.Now().Format(time.RFC3339), true)

	err := logMap.Persist(logId)

	if err != nil {
		logMap.Add(logId, "Error persisting log: "+err.Error(), true)
	}

	time.Sleep(5 * time.Second)

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
