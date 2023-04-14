package main

import (
	"errors"
	"fmt"
	"strings"
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

func persistToGit(logId string) error {
	// Open current directory as git repo

	if logId != "" {
		logMap.Add(logId, "Persisting changes to git", true)
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
		logMap.Add(logId, "Loaded working directory", true)
	}

	// Add all changes to the staging area
	_, err = w.Add(".")

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Commit the changes
	_, err = w.Commit("Persist changes to git", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name: "sysmanage-web[auto]",
			When: time.Now(),
		},
	})

	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			if logId != "" {
				logMap.Add(logId, "No changes to commit", true)
			}

			return nil
		} else if errors.Is(err, git.ErrEmptyCommit) {
			if logId != "" {
				logMap.Add(logId, "No changes to commit [doing so would create a empty commit]", true)
			}

			return nil
		}

		fmt.Println(err)
		return err
	}

	if logId != "" {
		logMap.Add(logId, "Committed changes", true)
	}

	outp := persistOutput{buf: &strings.Builder{}}

	err = repo.Push(&git.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: config.GithubPat,
			Password: config.GithubPat,
		},
		Progress: outp,
	})

	// Force push if we get an error
	if err != nil {
		err = repo.Push(&git.PushOptions{
			Force: true,
			Auth: &githttp.BasicAuth{
				Username: config.GithubPat,
				Password: config.GithubPat,
			},
			Progress: outp,
		})

		if err != nil {
			fmt.Println(err)

			if logId != "" {
				logMap.Add(logId, outp.buf.String()+": "+err.Error(), true)
			}

			return err
		}

		if logId != "" {
			logMap.Add(logId, "Pushed changes (force-push): "+outp.buf.String(), true)
		}
	} else {
		if logId != "" {
			logMap.Add(logId, "Pushed (no force-push):"+outp.buf.String(), true)
		}
	}

	return nil
}
