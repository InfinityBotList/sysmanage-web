// Common functions for plugins
package plugins

import (
	"errors"
	"fmt"
	"strconv"
	"sysmanage-web/core/state"

	"golang.org/x/exp/slices"
)

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
			val[i] = fmt.Sprintf("%s", v)
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
