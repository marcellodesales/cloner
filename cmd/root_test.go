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
	"fmt"
	"os"
	"testing"

	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

/**
 * Testing root and setup methods
 * https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go/34102842#34102842
 */
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

// the current CLI input to be tested
type TestInput struct {
	url        string
	forceClone bool
}

type CliTests struct {
	title                string
	input                TestInput
	exitCode             int
	exactErrorMessages   map[string]bool // https://yourbasic.org/golang/implement-set/
	errorMessageContains []string
}

/**
 * Execute the CLI tests using our strategy
 */
func ExecuteCliTestsStrategy(t *testing.T, cliTests []CliTests) {
	for _, tst := range cliTests {
		t.Run(tst.title, func(t *testing.T) {
			// https://github.com/smartystreets/goconvey/wiki#your-first-goconvey-test
			// Only pass t into top-level Convey calls
			Convey(tst.title, t, func() {

				// Actual clone command to be tested
				exitCode, errors := executeGitClone(tst.input.url, tst.input.forceClone)

				Convey(fmt.Sprintf("The exit code should be %d", tst.exitCode), func() {
					So(exitCode, ShouldEqual, tst.exitCode)
				})

				// Test exact error message cases
				if len(tst.exactErrorMessages) > 0 {
					Convey(fmt.Sprintf("With %d error message", len(tst.exactErrorMessages)), func() {
						So(len(errors), ShouldEqual, len(tst.exactErrorMessages))
						for _, errorMsg := range errors {
							So(tst.exactErrorMessages[errorMsg.Error()], ShouldBeTrue)
						}
					})
				}

				// Test whether the error message contains strings
				if len(tst.errorMessageContains) > 0 {
					Convey(fmt.Sprintf("With %d error message", len(tst.errorMessageContains)), func() {
						for _, errorMessageToken := range tst.errorMessageContains {
							So(errors[0].Error(), ShouldContainSubstring, errorMessageToken)
						}
					})
				}
			})
		})
	}
}
