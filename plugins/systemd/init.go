package systemd

import (
	"errors"
	"fmt"
	"os"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/plugins/frontend"
	"github.com/infinitybotlist/sysmanage-web/types"
)

var (
	targetTemplate                string
	serviceTemplate               string
	serviceDefinitions            string
	serviceOut                    string
	ignoreSuffixForCopy           bool
	ignoreSuffixForGetServiceList bool
	trimSuffixForManualUnits      bool
	srvModBypass                  []string
)

const ID = "systemd"

func InitPlugin(c *types.PluginConfig) error {
	// Register links
	frontend.AddLink(c, frontend.Link{
		Title:       "Systemd Service Management",
		Description: "Manage systemd services on the system.",
		LinkText:    "View Service List",
		Href:        "@root",
	})

	frontend.AddLink(c, frontend.Link{
		Title:       "Systemd Metadata Editor",
		Description: "Edit the metadata of systemd targets on the system.",
		LinkText:    "Edit Metadata",
		Href:        "@root/meta",
	})

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

	ignoreSuffixForCopy, err = cfgData.GetBool("ignore_suffix_for_copy")

	if err != nil {
		fmt.Println("Ignoring config option ignore_suffix_for_copy: " + err.Error())
	}

	ignoreSuffixForGetServiceList, err = cfgData.GetBool("ignore_suffix_for_get_service_list")

	if err != nil {
		fmt.Println("Ignoring config option ignore_suffix_for_get_service_list: " + err.Error())
	}

	trimSuffixForManualUnits, err = cfgData.GetBool("trim_suffix_for_manual_units")

	if err != nil {
		fmt.Println("Ignoring config option trim_suffix_for_manual_units: " + err.Error())
	}

	loadServiceApi(c.Mux)

	return nil
}
