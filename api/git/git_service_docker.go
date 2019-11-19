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
	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	"os"
)

/**
 * List the gitman dependencies by calling Gitman itself
 */
func (service GitServiceType) DockerGitClone(gitRepoClone *GitRepoClone) (string, error) {
	// The default Git Docker Image to use
	gitmanDockerImage := config.INSTANCE.Git.DockerImage

	workingDir := service.GetRepoUserDir(gitRepoClone)
	userHomeDir, _ := os.UserHomeDir()
	dockerCommandArgs := []string{userHomeDir, workingDir, gitmanDockerImage, gitRepoClone.Url}
	// The docker command to be passed
	dockerCommand := "docker run --rm -v %s/.ssh:/root/.ssh -v %s:/git %s clone %s"

	// Create a new Executor
	dockerCommandExecutor := util.NewDockerExecutor(dockerCommand, dockerCommandArgs)
	// Execute the command
	stdout, err := dockerCommandExecutor.Execute()
	// For now, just exit from errors
	if err != nil {
		return "", err
	}
	return stdout, nil
}