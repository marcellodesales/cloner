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
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v5"
	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/plumbing/transport"
	sshgit "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
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
* Executes an auth with the given private key file
  https://stackoverflow.com/questions/44269142/golang-ssh-getting-must-specify-hoskeycallback-error-despite-setting-it-to-n/63308243#63308243
  https://skarlso.github.io/2019/02/17/go-ssh-with-host-key-verification/
  https://github.com/src-d/go-git/issues/637#issuecomment-543015125
  Needs to execute the agent https://github.com/src-d/go-git/issues/550#issuecomment-323075887
*/
func sshAuth(privateKeyFilePath string) (transport.AuthMethod, error) {
	// Verification on the private key https://github.com/src-d/go-git/issues/637#issuecomment-404851019, not complete though
	isPrivateKey := func(privateKeyContent string) bool {
		if len(privateKeyContent) > 1000 && strings.HasPrefix(privateKeyContent, "-----") {
			return true
		}
		return false
	}

	privateKey, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	if isPrivateKey(string(privateKey)) {
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, err
		}

		// Let's assume that if the file exists, there were connections to the host before
		knownHostsFilePath := fmt.Sprintf("%s/.ssh/known_hosts", os.Getenv("HOME"))
		if util.FileExists(knownHostsFilePath) {
			// This whole piece fails because of the missing .ssh/known_hosts file.
			// TODO Need a way to generate the known_hosts when the file does not exist
			// https://skarlso.github.io/2019/02/17/go-ssh-with-host-key-verification/
			knownHostsCallback, err := kh.New(fmt.Sprintf("%s/.ssh/known_hosts", os.Getenv("HOME")))
			if err != nil {
				return nil, fmt.Errorf("could not create hostkeycallback function: %v", err)
			}

			// Verify the key with the content of the known_hosts file
			return &sshgit.PublicKeys{
				User:   "git",
				Signer: signer,
				HostKeyCallbackHelper: sshgit.HostKeyCallbackHelper{
					HostKeyCallback: knownHostsCallback,
				},
			}, nil
		}

		// In case the known hosts is not available, just don't verify
		// This is specific on containers, CI, tests, https://github.com/src-d/go-git/issues/637#issuecomment-404851019
		return &sshgit.PublicKeys{
			User:   "git",
			Signer: signer,
			HostKeyCallbackHelper: sshgit.HostKeyCallbackHelper{
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			},
		}, nil
	}
	return nil, nil
}

/**
 * Clone the git repo to the clone location using go-git based on the request
 */
func (service GitServiceType) GoCloneRepo(gitRepoCloneRequest *CloneGitRepoRequest, config *config.Configuration) error {
	var cloneOptions *git.CloneOptions
	var err error

	// https://git-scm.com/book/en/v2/Appendix-B%3A-Embedding-Git-in-your-Applications-go-git
	log.Debugf("Attempting to clone repo '%s' => '%s'", gitRepoCloneRequest.Url, gitRepoCloneRequest.CloneLocation)

	if gitRepoCloneRequest.PrivateKeyFile != "" {
		log.Debugf("Authenticating using the key ")
		auth, _ := sshAuth(gitRepoCloneRequest.PrivateKeyFile)

		// https://github.com/go-git/go-git/issues/169
		cloneOptions = &git.CloneOptions{
			URL:      gitRepoCloneRequest.Url,
			Auth:     auth,
			Progress: os.Stdout,
		}

	} else {
		cloneOptions = &git.CloneOptions{
			URL:      gitRepoCloneRequest.Url,
			Progress: os.Stdout,
		}
	}

	_, err = git.PlainClone(gitRepoCloneRequest.CloneLocation, false, cloneOptions)
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
