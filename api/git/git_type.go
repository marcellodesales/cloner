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
	"path"
)

type GitRepoType struct {
	Protocol string
	Host     string
	User     string
	Repo     string
}

type CloneGitRepoRequest struct {
	Url  string
	Type *GitRepoType
	//Tag string
	//Branch string
	//Revision string
	Depth uint
	//SparsePaths []string
	CloneLocation  string
	Force          bool
	PrivateKeyFile string
}

func (gitRepo GitRepoType) GetRepoDir() string {
	return path.Join(gitRepo.Host, gitRepo.User, gitRepo.Repo)
}

func (gitRepo GitRepoType) GetUserDir() string {
	return path.Join(gitRepo.Host, gitRepo.User)
}
