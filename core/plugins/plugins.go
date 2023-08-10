// Common functions for plugins
package plugins

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/infinitybotlist/sysmanage-web/core/state"

	"golang.org/x/exp/slices"
)

// Update this list when adding new plugins
var officialPlugins = []string{
	"acl",
	"actions",
	"authdp",
	"deploy",
	"frontend",
	"logger",
	"nginx",
	"persist",
	"systemd",
}

// Returns the list of official plugins
func GetOfficialPluginList() []string {
	return officialPlugins
}

type OpaqueConfig struct {
	inner map[string]any
}

func (i OpaqueConfig) GetString(key string) (string, error) {
	v, ok := i.inner[key]

	if !ok {
		return "", errors.New("key not found: " + key)
	}

	switch v := v.(type) {
	case string:
		return v, nil
	case nil:
		return "", nil
	}

	return "", errors.New("key not a string: " + key)
}

func (i OpaqueConfig) GetBool(key string) (bool, error) {
	v, ok := i.inner[key]

	if !ok {
		return false, errors.New("key not found: " + key)
	}

	switch v := v.(type) {
	case bool:
		return v, nil
	case nil:
		return false, nil
	}

	return false, errors.New("key not a boolean: " + key)
}

func (i OpaqueConfig) GetInt(key string) (int, error) {
	v, ok := i.inner[key]

	if !ok {
		return 0, errors.New("key not found: " + key)
	}

	switch v := v.(type) {
	case uint8, uint16, uint32, uint64, uint, int8, int16, int32, int64, int:
		return int(v.(int)), nil
	case nil:
		return 0, nil
	}

	return 0, errors.New("key not a integer: " + key)
}

func (i OpaqueConfig) GetStringArray(key string) ([]string, error) {
	v, ok := i.inner[key]

	if !ok {
		return nil, errors.New("key not found: " + key)
	}

	switch v := v.(type) {
	case []string:
		return v, nil
	case []int:
		// Convert to string array
		val := make([]string, len(v))

		for i, v := range v {
			val[i] = strconv.Itoa(v)
		}

		return val, nil
	case []any:
		// Convert to string array
		val := make([]string, len(v))

		for i, v := range v {
			val[i] = fmt.Sprint(v)
		}

		return val, nil
	case nil:
		return []string{}, nil
	}

	return nil, errors.New("key not a string array: " + key + "type: " + fmt.Sprintf("%s", v))
}

func Enabled(plugin string) bool {
	return slices.Contains(state.LoadedPlugins, plugin)
}

func GetConfig(plugin string) (*OpaqueConfig, error) {
	cfg, ok := state.Config.Plugins[plugin]

	if !ok {
		return nil, errors.New("plugin not enabled")
	}

	return &OpaqueConfig{inner: cfg}, nil
}
