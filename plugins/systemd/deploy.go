package systemd

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/infinitybotlist/sysmanage-web/core/logger"
	"github.com/infinitybotlist/sysmanage-web/core/state"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func initDeploy(logId string, srv ServiceManage) {
	if srv.Service.Git == nil {
		logger.LogMap.Add(logId, "FATAL: Service does not have git setup.", true)
		return
	}

	if srv.Service.Broken {
		logger.LogMap.Add(logId, "FATAL: Service is marked as broken.", true)
		return
	}

	if state.Config.GithubPat == "" {
		logger.LogMap.Add(logId, "FATAL: Github PAT not set. Git operations are disabled", true)
		return
	}

	defer logger.LogMap.MarkDone(logId)

	logger.LogMap.Add(logId, "Started deploy on: "+time.Now().Format(time.RFC3339), true)

	logger.LogMap.Add(logId, "Service Repo:"+srv.Service.Git.Repo, true)

	logger.LogMap.Add(logId, "Waiting for other builds to finish...", true)

	state.LsOp.Lock()
	defer state.LsOp.Unlock()

	deployFolder := srv.Service.Directory

	// Check that the deploy folder exists
	deployViaClone := false

	if _, err := os.Stat(deployFolder); os.IsNotExist(err) {
		logger.LogMap.Add(logId, "Deploy folder does not exist: "+deployFolder+", cloning it", true)

		err := os.MkdirAll(deployFolder, 0755)

		if err != nil {
			logger.LogMap.Add(logId, "Error creating deploy folder: "+err.Error(), true)
			return
		}

		deployViaClone = true
	} else if !srv.Service.Git.AllowDirty {
		logger.LogMap.Add(logId, "Dirty builds not allowed, performing fresh clone", true)
		deployViaClone = true
	}

	if deployViaClone {
		deployFolder = "deploys/" + logId

		logger.LogMap.Add(logId, "Cloning "+srv.Service.Git.Repo, true)
		_, err := git.PlainClone(deployFolder, false, &git.CloneOptions{
			URL: srv.Service.Git.Repo,
			Auth: &githttp.BasicAuth{
				Username: state.Config.GithubPat,
				Password: state.Config.GithubPat,
			},
			Progress:      logger.AutoLogger{ID: logId},
			ReferenceName: plumbing.ReferenceName(srv.Service.Git.Ref),
		})

		if err != nil {
			logger.LogMap.Add(logId, "Error cloning repo: "+err.Error(), true)
			return
		}
	} else {
		logger.LogMap.Add(logId, "Pulling into "+deployFolder, true)

		// Pull repo
		repo, err := git.PlainOpen(deployFolder)

		if err != nil {
			logger.LogMap.Add(logId, "Error opening repo: "+err.Error(), true)
			return
		}

		w, err := repo.Worktree()

		if err != nil {
			logger.LogMap.Add(logId, "Error getting worktree: "+err.Error(), true)
			return
		}

		// If there are unstaged changes, add+commit+push them
		status, err := w.Status()

		if err != nil {
			logger.LogMap.Add(logId, "WARNING: Error getting status: "+err.Error(), true)
		} else if err == nil {
			if !status.IsClean() {
				logger.LogMap.Add(logId, "Unstaged changes detected, committing them", true)

				_, err = w.Add(".")

				if err != nil {
					logger.LogMap.Add(logId, "Error adding unstaged changes: "+err.Error(), true)
					return
				}

				_, err = w.Commit("ci(update): Auto commit from github.com/infinitybotlist/sysmanage-web", &git.CommitOptions{
					All:               true,
					AllowEmptyCommits: true,
					Author: &object.Signature{
						Name: "github.com/infinitybotlist/sysmanage-web[auto]",
						When: time.Now(),
					},
				})

				if err != nil {
					logger.LogMap.Add(logId, "Error committing unstaged changes: "+err.Error(), true)
					return
				}

				err = repo.Push(&git.PushOptions{
					Auth: &githttp.BasicAuth{
						Username: state.Config.GithubPat,
						Password: state.Config.GithubPat,
					},
					Progress: logger.AutoLogger{ID: logId},
				})

				// Force push if we get an error
				if err != nil {
					logger.LogMap.Add(logId, "Error pushing unstaged changes: "+err.Error(), true)
					logger.LogMap.Add(logId, "Force pushing unstaged changes", true)

					err = repo.Push(&git.PushOptions{
						Force: true,
						Auth: &githttp.BasicAuth{
							Username: state.Config.GithubPat,
							Password: state.Config.GithubPat,
						},
						Progress: logger.AutoLogger{ID: logId},
					})

					if err != nil {
						logger.LogMap.Add(logId, "Error force pushing unstaged changes: "+err.Error(), true)
						return
					}
				}
			}
		}

		err = w.Pull(&git.PullOptions{
			ReferenceName: plumbing.ReferenceName(srv.Service.Git.Ref),
			Auth: &githttp.BasicAuth{
				Username: state.Config.GithubPat,
				Password: state.Config.GithubPat,
			},
			Progress: logger.AutoLogger{ID: logId},
		})

		if err != nil && err != git.NoErrAlreadyUpToDate {
			logger.LogMap.Add(logId, "Error pulling repo: "+err.Error(), true)
			return
		}
	}

	logger.LogMap.Add(logId, "Deploy folder: "+deployFolder, true)

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
					logger.LogMap.Add(logId, "WARN: ", true)
				}
			} else {
				curDir = curDir + "/" + args[1]
			}

			logger.LogMap.Add(logId, "Changed directory to "+curDir, true)
			continue // Ignore rest of command
		}

		// Run the command
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = curDir
		cmd.Env = os.Environ()

		for k, v := range srv.Service.Git.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}

		cmd.Stdout = logger.AutoLogger{ID: logId}
		cmd.Stderr = logger.AutoLogger{ID: logId, Error: true}

		err := cmd.Run()

		if err != nil {
			logger.LogMap.Add(logId, "Error running command: "+err.Error(), true)
			return
		}
	}

	if deployViaClone {
		// Copy any potential config files to deploy folder
		logger.LogMap.Add(logId, "Copying config files to new deploy", true)
		for _, file := range srv.Service.Git.ConfigFiles {
			logger.LogMap.Add(logId, "Copying config file "+file, true)

			f, err := os.Open(srv.Service.Directory + "/" + file)

			if err != nil {
				logger.LogMap.Add(logId, "WARNING: Could not open config file "+file, true)
				continue
			}

			defer f.Close()

			newF, err := os.Create(deployFolder + "/" + file)

			if err != nil {
				logger.LogMap.Add(logId, "Error creating config file: "+err.Error(), true)
				return
			}

			defer newF.Close()

			_, err = io.Copy(newF, f)

			if err != nil {
				logger.LogMap.Add(logId, "Error copying config file: "+err.Error(), true)
				return
			}

			err = newF.Sync()

			if err != nil {
				logger.LogMap.Add(logId, "Error syncing config file: "+err.Error(), true)
				return
			}
		}

		// Rename service folder
		err := os.Rename(srv.Service.Directory, srv.Service.Directory+"-old")

		if err != nil {
			logger.LogMap.Add(logId, "Error renaming service folder: "+err.Error(), true)
			return
		}

		// Move deploy folder to service folder
		err = os.Rename(deployFolder, srv.Service.Directory)

		if err != nil {
			logger.LogMap.Add(logId, "Error moving deploy folder: "+err.Error(), true)

			// Move old service folder back
			err = os.Rename(srv.Service.Directory+"-old", srv.Service.Directory)

			if err != nil {
				logger.LogMap.Add(logId, "Error moving old service folder back: "+err.Error(), true)
			}

			return
		}

		// Remove old service folder
		err = os.RemoveAll(srv.Service.Directory + "-old")

		if err != nil {
			logger.LogMap.Add(logId, "Error removing old service folder: "+err.Error(), true)
			return
		}
	}

	logger.LogMap.Add(logId, "Deploy finished on: "+time.Now().Format(time.RFC3339), true)

	time.Sleep(5 * time.Second)

	// Run systemctl restart deploy.Git.Service
	cmd := exec.Command("systemctl", "restart", srv.ID)
	cmd.Env = os.Environ()
	cmd.Stdout = logger.AutoLogger{ID: logId}
	cmd.Stderr = logger.AutoLogger{ID: logId, Error: true}

	err := cmd.Run()

	if err != nil {
		logger.LogMap.Add(logId, "Error restarting service: "+err.Error(), true)
	}

	logger.LogMap.Add(logId, "Service restarted on: "+time.Now().Format(time.RFC3339), false)
}
