/*
Copyright © 2019 Marcello de Sales <marcello.gmail@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config

import (
	"errors"
	"github.com/marcellodesales/cloner/util"
	"os"
	"path"
	"strings"
)

/**
 * https://godoc.org/github.com/mitchellh/mapstructure
 */
type Git struct {
	// The name of the protos used by the generated Stubs
	DockerImage string `mapstructure:"dockerImage"`

	CloneBaseDir string `mapstructure:"cloneBaseDir"`

	// The list of plugins to initialize
	Properties map[string]string `mapstructure:"properties"`
}

/**
 * Loads the config from file
 */
func parseGitConfig() (*Git, error) {
	// Make the dependencies reference
	config := &Git{}

	// Read the properties from Viper
	if err := util.ReadConfig("git", config); err != nil {
		return nil, err
	}

	// Set default values
	setDefaultCliValues(config)

	return config, nil
}

func setDefaultCliValues(git *Git) {
	if git.DockerImage == "" {
		git.DockerImage = "alpine/git"
	}
	homeDir, _ := os.UserHomeDir()
	if git.CloneBaseDir == "" {
		git.CloneBaseDir = path.Join(homeDir, "cloner")

	} else if strings.Contains(git.CloneBaseDir, "~") {
		// mkdir in go works, but Since it will be used by docker, we need to expland it
		git.CloneBaseDir = strings.ReplaceAll(git.CloneBaseDir, "~", homeDir)
	}
}

/**
 * Validate the initialization
 */
func validateInitConfig(git *Git) error {
	if git.DockerImage == "" {
		git.DockerImage = "alpine/git"
		return nil
	}
	if git.CloneBaseDir == "" {
		return errors.New("you must provide the clone base dir")
	}
	return nil
}

/**
 * Whether the dependencies contain locked resources
 */
func (init *Git) getPropertyNames() []string {
	var keys []string
	for k := range init.Properties {
		keys = append(keys, k)
	}
	return keys
}
