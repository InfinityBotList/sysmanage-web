package logger

import "github.com/infinitybotlist/sysmanage-web/types"

const ID = "logger"

func InitPlugin(c *types.PluginConfig) error {
	loadLoggerApi(c.Mux)
	return nil
}
