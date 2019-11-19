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
	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	"os"
	"path"
	"regexp"
)

type GitServiceType struct {}

var GitService GitServiceType

// The regex to use to parse the repo https://play.golang.org/p/iBeytKqRDBX
var githubRepoUrlRe = regexp.MustCompile(`(?m)^(?P<protocol>https|git)(:\/\/|@)(?P<host>[^\/:]+)[\/:](?P<user>[^\/:]+)\/(?P<repo>.+).git$`)

/**
 * @return a new instance of the GitRepoType
 */
func (serviceType GitServiceType) ParseRepoString(repo string) (*GitRepoClone, error) {
	if repo == "" {
		return nil, errors.New("you must provide the repo URL")
	}

	// Parse the regex
	gitRepoValues, err := util.RegexProcessString(githubRepoUrlRe, repo)
	if err != nil {
		return nil, err
	}

	gitRepoInstance := GitRepoType{
		Protocol: gitRepoValues["protocol"],
		Host: gitRepoValues["host"],
		User: gitRepoValues["user"],
		Repo: gitRepoValues["repo"],
	}

	gitRepoClone := &GitRepoClone{
		Url: repo,
		Type: gitRepoInstance,
	}

	return gitRepoClone, nil
}

/**
 * Get the repo clone dir expected to be created
 */
func (serviceType GitServiceType) GetRepoCloneDir(gitRepoClone *GitRepoClone) string {
	return path.Join(config.INSTANCE.Git.CloneBaseDir, gitRepoClone.Type.GetRepoDir())
}

/**
 * Get the repo clone dir expected to be created
 */
func (serviceType GitServiceType) GetRepoUserDir(gitRepoClone *GitRepoClone) string {
	return path.Join(config.INSTANCE.Git.CloneBaseDir, gitRepoClone.Type.GetUserDir())
}

/**
 * Makes the base clone dir
 */
func (serviceType GitServiceType) MakeRepoUserDir(gitRepoClone *GitRepoClone) error {
	baseCloneDir := serviceType.GetRepoUserDir(gitRepoClone)
	err := os.MkdirAll(baseCloneDir, 0755)
	if err != nil {
		return err
	}
	return nil
}