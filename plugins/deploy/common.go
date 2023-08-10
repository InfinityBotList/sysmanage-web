package deploy

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func LoadConfig(name string) (*DeployMeta, error) {
	// Read file into *DeployMeta
	f, err := os.Open(deployConfigPath + "/" + name)

	if err != nil {
		return nil, errors.New("Failed to read deploy config " + err.Error() + f.Name())
	}

	// Read file into *DeployMeta
	var meta *DeployMeta

	err = yaml.NewDecoder(f).Decode(&meta)

	if err != nil {
		return nil, errors.New("Failed to read deploy config " + err.Error() + f.Name())
	}

	return meta, nil
}

func GetDeployList() ([]*DeployMeta, error) {
	// Get all files in path
	fsd, err := os.ReadDir(deployConfigPath)

	if err != nil {
		return nil, errors.New("Failed to load deployConfigPath " + err.Error())
	}

	dms := make([]*DeployMeta, 0)

	for _, file := range fsd {
		if file.IsDir() {
			continue // Skip directories
		}

		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue // Skip non-yaml files
		}

		// Read file into *DeployMeta
		f, err := os.Open(deployConfigPath + "/" + file.Name())

		if err != nil {
			return nil, errors.New("Failed to read deploy config " + err.Error() + file.Name())
		}

		// Read file into *DeployMeta
		var meta *DeployMeta

		err = yaml.NewDecoder(f).Decode(&meta)

		if err != nil {
			return nil, errors.New("Failed to read deploy config " + err.Error() + file.Name())
		}

		dms = append(dms, meta)
	}

	return dms, nil
}
