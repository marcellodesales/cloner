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
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"

	log "github.com/sirupsen/logrus"
)

type GitServiceType struct{}

var GitService GitServiceType

// The regex to use to parse the repo https://play.golang.org/p/iBeytKqRDBX
var githubRepoUrlRe = regexp.MustCompile(`(?m)^(?P<protocol>https|git)(:\/\/|@)(?P<host>[^\/:]+)[\/:](?P<user>[^\/:]+)\/(?P<repo>.+).git$`)

// Init from the CLI inputs
func (service GitServiceType) Init(repoUrl, privateKeyPath string, forceClone bool) (*CloneGitRepoRequest, error) {
	if repoUrl == "" {
		return nil, fmt.Errorf("you must provide the repo URL")
	}

	if !strings.HasSuffix(repoUrl, ".git") {
		repoUrl += ".git"
	}

	// Make a clone repo request based on the input from the CLI
	if strings.HasPrefix(repoUrl, "git@") && privateKeyPath == "" {
		privateKeyPath = fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
		log.Debugf("Assuming default private key at '%s': repo URL is prefixed by 'git@'", privateKeyPath)
	}

	// for git@host:org/repo, SSH keys are used as auth form.
	if privateKeyPath != "" && !util.FileExists(privateKeyPath) {
		return nil, fmt.Errorf("private key '%s' does NOT exist. Make sure to provide for ssh-related clones (git@) repos", privateKeyPath)
	}

	return &CloneGitRepoRequest{
		Url:            repoUrl,
		Force:          forceClone,
		PrivateKeyFile: privateKeyPath,
	}, nil
}

/**
 * @return a new instance of the GitRepoType
 */
func (service GitServiceType) ParseRepoString(gitRepoClone *CloneGitRepoRequest) error {
	// Parse the regex
	gitRepoValues, err := util.RegexProcessString(githubRepoUrlRe, gitRepoClone.Url)
	if err != nil {
		return err
	}

	gitRepoClone.Type = &GitRepoType{
		Protocol: gitRepoValues["protocol"],
		Host:     gitRepoValues["host"],
		User:     gitRepoValues["user"],
		Repo:     gitRepoValues["repo"], // may contain "/" when in repo is in gitlab folders format
	}

	return nil
}

/**
 * Get the org
 */
func (service GitServiceType) GetRepoLocalPath(gitRepoClone *CloneGitRepoRequest, config *config.Configuration) string {
	return path.Join(config.Git.CloneBaseDir, gitRepoClone.Type.GetRepoDir())
}

/**
 * Get the repo clone dir expected to be created
 */
func (service GitServiceType) GetOrgLocalPath(gitRepoClone *CloneGitRepoRequest, config *config.Configuration) string {
	return path.Join(config.Git.CloneBaseDir, gitRepoClone.Type.GetUserDir())
}

/**
 * Make sure the clone dir does not exist. If it does, force clone must be provided to not fail.
 */
func (service GitServiceType) VerifyCloneDir(gitRepoClone *CloneGitRepoRequest, config *config.Configuration) (bool, error) {
	// The location is provided by the api
	gitRepoClone.CloneLocation = service.GetRepoLocalPath(gitRepoClone, config)
	//log.Debugf("Verifying if the clone path '%s' exists or exists and is empty", gitRepoClone.CloneLocation)

	if util.DirExists(gitRepoClone.CloneLocation) {
		if gitRepoClone.Force {
			err := util.DeleteDir(gitRepoClone.CloneLocation)
			if err != nil {
				return false, fmt.Errorf("can't force deletion of dir '%s': %v", gitRepoClone.CloneLocation, err)
			}
			return true, nil
		}

		// Verify if the user repo is not empty
		dirIsEmpty, _ := util.IsDirEmpty(gitRepoClone.CloneLocation)
		if !dirIsEmpty {
			return false, fmt.Errorf("clone location '%s' exists and it's not empty", gitRepoClone.CloneLocation)
		}
	}
	return false, nil
}

/**
 * Makes the base clone dir is the org and parts of the repo dir since the clone service creates the repo dir
 */
func (service GitServiceType) MakeCloneDir(gitRepoClone *CloneGitRepoRequest, config *config.Configuration) error {
	baseCloneDir := service.GetOrgLocalPath(gitRepoClone, config)

	// For gitlab and other hosts with more than a single folder
	if strings.Contains(gitRepoClone.Type.Repo, "/") {
		split := strings.Split(gitRepoClone.Type.Repo, "/")
		for i := 0; i <= len(split)-2; i++ {
			baseCloneDir += "/" + split[i]
		}
	}

	err := os.MkdirAll(gitRepoClone.CloneLocation, 0755)
	if err != nil {
		return err
	}

	gitRepoClone.CloneLocation = service.GetRepoLocalPath(gitRepoClone, config)
	return nil
}

/**
 * Clone the git repo to the clone location using go-git
 */
func (service GitServiceType) GoCloneRepo(gitRepoClone *CloneGitRepoRequest, config *config.Configuration) error {
	var publicKeys *ssh.PublicKeys
	var err error
	if strings.HasPrefix(gitRepoClone.Url, "git@") {
		// default private key is
		privateKeyFile := fmt.Sprintf("%s/.ssh/gmail_rsa", os.Getenv("HOME"))
		publicKeys, err = ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
		if err != nil {
			return err
		}
		log.Debugf("Using public keys %s", publicKeys.String())
	}

	// https://git-scm.com/book/en/v2/Appendix-B%3A-Embedding-Git-in-your-Applications-go-git
	log.Debugf("Attempting to clone repo '%s' => '%s'", gitRepoClone.Url, gitRepoClone.CloneLocation)
	_, err = git.PlainClone(gitRepoClone.CloneLocation, false, &git.CloneOptions{
		URL:      gitRepoClone.Url,
		Auth:     publicKeys,
		Progress: os.Stdout,
	})

	return err
}

/**
 * Print the list of the files in the dir just like the "tree" unix command
 */
func (service GitServiceType) GoPrintTree(gitRepoClone *CloneGitRepoRequest) (string, error) {
	// Exclude the .git dir from the list
	excludedList := []string{".git"}
	return util.GetDirTree(gitRepoClone.CloneLocation, excludedList)
}
