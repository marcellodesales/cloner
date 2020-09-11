/*
Copyright © 2019 Marcello de Sales <marcello.desales@gmail.com>

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
package cmd

import (
	"os"
	"testing"

	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	// Sets up Viper with the given config file name
	util.SetupViperHomeConfig(cfgFile, ".cloner", "yaml", "CLONER")

	// Setups up the configuration parsed from Viper config (maps -> structs)
	config.Setup()

	if err := util.SetUpLogs(os.Stdout, logrus.InfoLevel.String()); err != nil {
		os.Exit(1)
	}
}
