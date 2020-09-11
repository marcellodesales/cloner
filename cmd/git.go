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
func GitCloneCmd(cmd *cobra.Command, args []string) {
	log.Debugf("config.git.cloneBaseDir=%s", config.INSTANCE.Git.CloneBaseDir)
	repo, _ := cmd.Flags().GetString("repo")
	forceClone, _ := cmd.Flags().GetBool("force")

	exitCode, errors := executeGitClone(repo, forceClone)

	// Show any errors if any
	if len(errors) > 0 {
		for _, error := range errors {
			log.Error(error)
		}
	}

	// Exit as part of the CLI contract
	os.Exit(exitCode)
}

func executeGitClone(repo string, forceClone bool) (int, []error) {
	// Make a clone repo request based on the input from the CLI
	cloneRepoRequest := git.GitService.Init(repo, forceClone)

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

	var verbose = false
	// https://github.com/spf13/cobra/issues/818#issuecomment-489021216
	initCmd.Flags().BoolVarP(&verbose, "force", "f", false, "Forces cloning by deleting existing dir")
}
