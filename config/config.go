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
	"fmt"
	"github.com/marcellodesales/cloner/util"
	log "github.com/sirupsen/logrus"
)

// Configuration represents the application configuration
type Configuration struct {
	FilePath string
	Git      Git `mapstructure:"git" yaml:"git"`
}

// Singleton instance of the configuration
var INSTANCE *Configuration

/**
 * Loads the config from file
 */
func initConfig() (*Configuration, error) {
	// Make the dependencies reference
	config := &Configuration{}

	//  ############# Initialize the init config
	git, err := parseGitConfig()
	if err != nil {
		return nil, fmt.Errorf("config: Can't parse init config: %g", err)
	}

	// Validate
	err = validateInitConfig(git)
	if err != nil {
		log.Error("Can't init because of config errors: %v", err)
	}

	config.Git = *git
	log.Debugf("Initializing config for git handler %v", config.Git)

	// The path for config
	config.FilePath = util.GetConfigFilePath()

	return config, nil
}

// Creates the singleton of the configs to be available to all modules
func Setup() {
	config, err := initConfig()
	if err != nil {
		log.Fatalf("Can't load config: %v", err)
	}
	INSTANCE = config
}