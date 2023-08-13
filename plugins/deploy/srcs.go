package deploy

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/infinitybotlist/eureka/crypto"
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

// Public API to allow plugins to define custom webhook sources
var DeployWebhookSources = map[string]func(cfg *DeployMeta, wid, id, token string) (logId string, err error){
	"api": func(cfg *DeployMeta, wid, id, token string) (logId string, err error) {
		var flag bool
		for _, webh := range cfg.Webhooks {
			if webh.Type != "api" {
				continue
			}

			if wid == webh.Id && webh.Token == token {
				flag = true
				break
			}
		}

		if !flag {
			return "", errors.New("invalid token")
		}

		reqId := crypto.RandString(64)

		go InitDeploy(reqId, cfg)

		return reqId, nil
	},
}
