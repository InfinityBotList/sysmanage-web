package systemd

import (
	"errors"
	"os"
	"sysmanage-web/core/plugins"
	"sysmanage-web/types"
)

var (
	targetTemplate     string
	serviceTemplate    string
	serviceDefinitions string
	serviceOut         string
	srvModBypass       []string
)

func InitPlugin(c *types.PluginConfig) error {
	// Open data/servicegen/target.tmpl
	bytes, err := os.ReadFile("data/servicegen/target.tmpl")

	if err != nil {
		return errors.New("Failed to read target.tmpl: " + err.Error())
	}

	targetTemplate = string(bytes)

	// Open data/servicegen/service.tmpl
	bytes, err = os.ReadFile("data/servicegen/service.tmpl")

	if err != nil {
		return errors.New("Failed to read service.tmpl: " + err.Error())
	}

	serviceTemplate = string(bytes)

	cfgData, err := plugins.GetConfig(c.Name)

	if err != nil {
		return errors.New("Failed to get nginx config: " + err.Error())
	}

	serviceDefinitions, err = cfgData.GetString("service_definitions")

	if err != nil {
		return err
	}

	serviceOut, err = cfgData.GetString("service_out")

	if err != nil {
		return err
	}

	srvModBypass, err = cfgData.GetStringArray("srv_mod_bypass")

	if err != nil {
		return err
	}

	loadServiceApi(c.Mux)

	return nil
}
