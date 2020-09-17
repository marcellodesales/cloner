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

	"github.com/marcellodesales/cloner/api/git"
	"github.com/marcellodesales/cloner/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "git",
	Short: "Clones a given git repo",
	Long:  `Clones a given git repo URL`,
	Run:   GitCloneCmd,
}

// Exposed command for testing
// https://stackoverflow.com/questions/59709345/how-to-implement-unit-tests-for-cli-commands-in-go/59714127#59714127
func GitCloneCmd(cmd *cobra.Command, args []string) {
	repo, _ := cmd.Flags().GetString("repo")
	forceClone, _ := cmd.Flags().GetBool("force")
	privateKey, _ := cmd.Flags().GetString("privateKey")

	exitCode, errors := executeGitClone(repo, privateKey, forceClone)

	// Show any errors if any
	if len(errors) > 0 {
		for _, error := range errors {
			log.Error(error)
		}
	}

	// Exit as part of the CLI contract
	os.Exit(exitCode)
}

func executeGitClone(repo, privateKeyPath string, forceClone bool) (int, []error) {
	cloneRepoRequest, err := git.GitService.Init(repo, privateKeyPath, forceClone)
	if err != nil {
		return -1, []error{fmt.Errorf("can't initialize due to parameters' values: %v", err)}
	}

	// Execute the implementation, getting the exit code and any error
	return git.CloneGitRepo(cloneRepoRequest, config.INSTANCE)
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringP("repo", "r", "", "The repo URL to clone")
	initCmd.Flags().StringP("privateKey", "k", "", "The private key associated to the public key to clone 'git@' repos")

	var verbose = false
	// https://github.com/spf13/cobra/issues/818#issuecomment-489021216
	initCmd.Flags().BoolVarP(&verbose, "force", "f", false, "Forces cloning by deleting existing dir")
}
