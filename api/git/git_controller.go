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
package git

import (
	"fmt"

	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	log "github.com/sirupsen/logrus"
)

/**
 * Clones the repo implementation, returning the exit code and any error message
 */
func CloneGitRepo(cloneRequest *CloneGitRepoRequest, config *config.Configuration) (int, []error) {
	log.Debugf("config.git.cloneBaseDir=%s", config.Git.CloneBaseDir)
	err := GitService.ParseRepoString(cloneRequest)
	if err != nil {
		return 1, []error{fmt.Errorf("git URL invalid: %v", err)}
	}

	if cloneRequest.Force {
		log.Info("Forcing clone...")
	}

	deletedExistingDir, err := GitService.VerifyCloneDir(cloneRequest, config)
	if deletedExistingDir {
		log.Infof("Deleted dir '%s'", cloneRequest.CloneLocation)
	}
	if err != nil {
		return 2, []error{fmt.Errorf("can't clone repo: %v", err),
			fmt.Errorf("you can specify --force or -f to delete the existing dir and clone again. " +
				"Make sure there are no panding changes")}
	}

	err = GitService.MakeCloneDir(cloneRequest, config)
	if err != nil {
		return 3, []error{fmt.Errorf("can't create the base clone repo '%s': %v", cloneRequest.Type.GetUserDir(), err)}
	}

	log.Infof("Cloning into '%s'", cloneRequest.CloneLocation)
	err = GitService.GoCloneRepo(cloneRequest, config)
	if err != nil {
		return 4, []error{fmt.Errorf("can't clone the repo at '%s': %v", cloneRequest.Type.GetRepoDir(), err)}

	} else {
		log.Info("Done...")
	}

	// Show the files cloned
	filesListTree, err := GitService.GoPrintTree(cloneRequest)
	if err != nil {
		return 5, []error{fmt.Errorf("can't show the cloned repo tree '%s': %v", cloneRequest.Type.GetRepoDir(), err)}
	}

	if util.IsLogInDebug() {
		log.Infof("\n%s", filesListTree)
	}

	return 0, nil
}
