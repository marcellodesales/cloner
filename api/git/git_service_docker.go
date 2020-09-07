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
	"errors"
	"fmt"
	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	log "github.com/sirupsen/logrus"
	"os"
)

func (service GitServiceType) VerifyCloneDir(gitRepoClone *GitRepoClone, forceClone bool, config *config.Configuration) error {
	// The location is provided by the api
	gitRepoClone.CloneLocation = service.GetRepoLocalPath(gitRepoClone, config)
	log.Debugf("Verifying if the clone path '%s' exists or exists and is empty", gitRepoClone.CloneLocation)

	if util.DirExists(gitRepoClone.CloneLocation) {
		if forceClone {
			util.DeleteDir(gitRepoClone.CloneLocation)
			return nil
		}

		// Verify if the user repo is not empty
		dirIsEmpty, _ := util.IsDirEmpty(gitRepoClone.CloneLocation)
		if !dirIsEmpty {
			return errors.New(fmt.Sprintf("clone location '%s' exists and it's not empty", gitRepoClone.CloneLocation))
		}
	}
	return nil
}

/**
 * Regular git
 */
func (service GitServiceType) DockerGitClone(gitRepoClone *GitRepoClone, config *config.Configuration) (string, error) {
	// The location is provided by the api
	var cloneLocation = service.GetOrgLocalPath(gitRepoClone, config)

	userHomeDir, _ := os.UserHomeDir()

	var dockerCommandExecutor *util.DockerExecutor

	// Only for branch or tag, full paths
	//var selectedRevision string
	//if gitRepoClone.Branch != "" {
	//	selectedRevision = gitRepoClone.Branch
	//
	//} else {
	//	selectedRevision = gitRepoClone.Tag
	//}

	// Set the depth if it is set
	depthSetting := ""
	if gitRepoClone.Depth > 0 {
		depthSetting = fmt.Sprintf("--depth=%d ", gitRepoClone.Depth)
	}

	dockerImage := config.Git.DockerImage
	dockerCommandArgs := []string{userHomeDir, cloneLocation, dockerImage, depthSetting, gitRepoClone.Url}
	// The docker command to be passed
	//dockerCommand := "docker run --rm -v %s/.ssh:/root/.ssh -v %s:/git %s clone %s %s --branch %s"
	dockerCommand := "docker run --rm -v %s/.ssh:/root/.ssh -v %s:/git %s clone %s%s"

	// Create a new Executor
	dockerCommandExecutor = util.NewDockerExecutor(dockerCommand, dockerCommandArgs)

	// Execute the command
	stdout, err := dockerCommandExecutor.Execute()
	// For now, just exit from errors
	if err != nil {
		return stdout, err
	}
	return stdout, nil
}

/**
 * Show the tree of files generated
 */
func (service GitServiceType) DockerFilesTree(gitRepoClone *GitRepoClone) (string, error) {
	workingDir := service.GetRepoCloneDir(gitRepoClone)
	dockerCommandArgs := []string{workingDir, workingDir, workingDir}

	// The docker command to be passed
	dockerCommand := "docker container run --rm -v %s:%s iankoulski/tree %s"
	// Create a new Executor
	dockerCommandExecutor := util.NewDockerExecutor(dockerCommand, dockerCommandArgs)

	// Execute the command
	stdout, err := dockerCommandExecutor.Execute()
	// For now, just exit from errors
	if err != nil {
		return stdout, err
	}

	return stdout, nil
}