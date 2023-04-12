package main

import (
	"errors"
	"os"
	"strings"
	"sync"
	"sysmanage-web/types"
	"time"

	"gopkg.in/yaml.v3"
)

const expiryTime = 10 * time.Minute

var cache ICache
var inOp sync.Mutex

type ICache struct {
	services    []types.ServiceManage
	ids         []string
	valid       bool
	lastUpdated time.Time
}

// Similar to ICache but exported and not opaque
type Cache struct {
	Services []types.ServiceManage
	Ids      []string
	Valid    bool
}

func (c *ICache) GetCache() (*Cache, error) {
	if time.Since(c.lastUpdated) > expiryTime || !c.valid {
		err := c.Load()

		if err != nil {
			return nil, err
		}
	}

	return &Cache{
		Services: c.services,
		Ids:      c.ids,
		Valid:    c.valid,
	}, nil
}

func (c *ICache) Load() error {
	inOp.Lock()
	defer inOp.Unlock()

	c.services = []types.ServiceManage{}
	c.ids = []string{}

	// Get all files in path
	fsd, err := os.ReadDir(config.ServiceDefinitions)

	if err != nil {
		return errors.New("Failed to read service definition " + err.Error())
	}

	for _, file := range fsd {
		if file.Name() == "_meta.yaml" {
			continue // Skip _meta.yaml
		}

		if file.IsDir() {
			continue // Skip directories
		}

		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue // Skip non-yaml files
		}

		// Read file into TemplateYaml
		f, err := os.Open(config.ServiceDefinitions + "/" + file.Name())

		if err != nil {
			return errors.New("Failed to read service definition " + err.Error() + file.Name())
		}

		// Read file into TemplateYaml
		var service types.TemplateYaml

		err = yaml.NewDecoder(f).Decode(&service)

		if err != nil {
			return errors.New("Failed to read service definition " + err.Error() + file.Name())
		}

		// Service name is the name without .yaml
		sname := strings.TrimSuffix(file.Name(), ".yaml")

		c.ids = append(c.ids, sname)

		c.services = append(c.services, types.ServiceManage{
			Service: service,
			ID:      sname,
		})
	}

	// Get status of services
	statuses := getServiceStatus(c.ids)

	for i := range c.services {
		c.services[i].Status = statuses[i]
	}

	c.valid = true
	c.lastUpdated = time.Now()

	return nil
}

func (c *ICache) Invalidate() {
	c.valid = false
}
