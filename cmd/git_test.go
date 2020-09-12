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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/**
 * Testing Golang CLI
 * https://stackoverflow.com/questions/59709345/how-to-implement-unit-tests-for-cli-commands-in-go/59714127#59714127
 */
func TestGitCloneWithWrongURLs(t *testing.T) {

	// https://github.com/smartystreets/goconvey/wiki#your-first-goconvey-test
	// Only pass t into top-level Convey calls
	Convey("Given an empty URL to clone", t, func() {
		url := ""
		forceClone := false
		exitCode, errors := executeGitClone(url, forceClone)

		Convey("The exit code should be 1", func() {
			So(exitCode, ShouldEqual, 1)
		})
		Convey("With 1 error message", func() {
			So(len(errors), ShouldEqual, 1)
			So(errors[0].Error(), ShouldEqual, "git URL invalid: you must provide the repo URL")
		})
	})

	// https://github.com/smartystreets/goconvey/wiki#your-first-goconvey-test
	// Only pass t into top-level Convey calls
	Convey("Given an invalid git URL to clone", t, func() {
		url := "abc-not-url"
		forceClone := false
		exitCode, errors := executeGitClone(url, forceClone)

		Convey("The exit code should be 1", func() {
			So(exitCode, ShouldEqual, 1)
		})
		Convey("With 1 error message", func() {
			So(len(errors), ShouldEqual, 1)
			So(errors[0].Error(), ShouldContainSubstring, "git URL invalid:")
			So(errors[0].Error(), ShouldContainSubstring, url)
		})
	})

}
