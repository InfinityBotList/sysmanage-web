package deploy

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/infinitybotlist/sysmanage-web/core/logger"
)

// A deploy source is a public API provided for use by other sysmanage plugins
//
// # It accepts a string action to clone the source to the folder desired.
//
// E.g: func(logId, buildDir string, d *DeployMeta) error
var DeploySources = map[string]func(logId, buildDir string, d *DeployMeta) error{
	"git": func(logId, buildDir string, d *DeployMeta) error {
		logger.LogMap.Add(logId, "Cloning "+d.Src.Url, true)
		_, err := git.PlainClone(buildDir, false, &git.CloneOptions{
			URL: d.Src.Url,
			Auth: &githttp.BasicAuth{
				Username: d.Src.Token,
				Password: d.Src.Token,
			},
			Progress:      logger.AutoLogger{ID: logId},
			ReferenceName: plumbing.ReferenceName(d.Src.Ref),
		})

		if err != nil {
			return err
		}

		return nil
	},
}
