package persist

import (
	"errors"
	"fmt"
	"strings"
	"sysmanage-web/core/logger"
	"sysmanage-web/core/plugins"
	"sysmanage-web/core/state"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type persistOutput struct {
	buf *strings.Builder
}

func (po persistOutput) Write(p []byte) (n int, err error) {
	po.buf.Write(p)
	return len(p), nil
}

func PersistToGit(logId string) error {
	if !plugins.Enabled("persist") {
		return errors.New("not persisting, persist plugin is disabled")
	}

	// Open current directory as git repo
	if logId != "" {
		logger.LogMap.Add(logId, "Persisting changes to git", true)
	}

	if state.Config.GithubPat == "" {
		if logId != "" {
			logger.LogMap.Add(logId, "WARNING: Github PAT not set. Git operations are disabled", true)
		}
		return nil
	}

	repo, err := git.PlainOpen(".")

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get the working directory for the repository
	w, err := repo.Worktree()

	if err != nil {
		fmt.Println(err)
		return err
	}

	if logId != "" {
		logger.LogMap.Add(logId, "Loaded working directory", true)
	}

	// First try pulling
	err = w.Pull(&git.PullOptions{
		Auth: &githttp.BasicAuth{
			Username: state.Config.GithubPat,
			Password: state.Config.GithubPat,
		},
	})

	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			if logId != "" {
				logger.LogMap.Add(logId, "PASS: No changes to pull", true)
			}
		} else {
			fmt.Println(err)
			return errors.New("FATAL: " + err.Error())
		}
	} else {
		if logId != "" {
			logger.LogMap.Add(logId, "Pulled changes", true)
		}
	}

	if status, err := w.Status(); err == nil {
		if status.IsClean() {
			if logId != "" {
				logger.LogMap.Add(logId, "No changes to persist", true)
			}

			return nil
		}
	} else {
		logger.LogMap.Add(logId, "FATAL: Error getting git status - "+err.Error(), true)
	}

	// Add all changes to the staging area
	_, err = w.Add(".")

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Commit the changes
	_, err = w.Commit("ci(update): persist changes to git", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name: "sysmanage-web[auto]",
			When: time.Now(),
		},
	})

	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			if logId != "" {
				logger.LogMap.Add(logId, "No changes to commit", true)
			}

			return nil
		} else if errors.Is(err, git.ErrEmptyCommit) {
			if logId != "" {
				logger.LogMap.Add(logId, "No changes to commit [doing so would create a empty commit]", true)
			}

			return nil
		}

		fmt.Println(err)
		return err
	}

	if logId != "" {
		logger.LogMap.Add(logId, "Committed changes", true)
	}

	outp := persistOutput{buf: &strings.Builder{}}

	err = repo.Push(&git.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: state.Config.GithubPat,
			Password: state.Config.GithubPat,
		},
		Progress: outp,
	})

	// Force push if we get an error
	if err != nil {
		err = repo.Push(&git.PushOptions{
			Force: true,
			Auth: &githttp.BasicAuth{
				Username: state.Config.GithubPat,
				Password: state.Config.GithubPat,
			},
			Progress: outp,
		})

		if err != nil {
			fmt.Println(err)

			if logId != "" {
				logger.LogMap.Add(logId, outp.buf.String()+": "+err.Error(), true)
			}

			return err
		}

		if logId != "" {
			logger.LogMap.Add(logId, "Pushed changes (force-push): "+outp.buf.String(), true)
		}
	} else {
		if logId != "" {
			logger.LogMap.Add(logId, "Pushed (no force-push):"+outp.buf.String(), true)
		}
	}

	return nil
}