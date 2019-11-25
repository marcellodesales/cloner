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
	"github.com/marcellodesales/cloner/api/git"
	"github.com/marcellodesales/cloner/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "git",
	Short: "Clones a given git repo",
	Long: `Clones a given github repo using the handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			log.Errorf("You must provide the github repo clone URL")
			os.Exit(1)
		}
		log.Debugf("Setting up for %s", repo)

		gitRepo, err := git.GitService.ParseRepoString(repo)
		if err != nil {
			log.Errorf("Can't parse the Git Repo: %v", err)
			os.Exit(2)
		}

		orgDir := git.GitService.GetRepoUserDir(gitRepo)
		log.Infof("Cloning the provided repo at %s", orgDir)

		err = git.GitService.MakeRepoUserDir(gitRepo)
		if err != nil {
			log.Errorf("Can't create the base clone repo '%s': %v", gitRepo.Type.GetUserDir(), err)
			os.Exit(3)
		}

		cloneStdout, err := git.GitService.DockerGitClone(gitRepo)
		if err != nil {
			log.Errorf("Can't clone the repo at '%s': %v", gitRepo.Type.GetRepoDir(), cloneStdout)
		} else {
			log.Info("Finished cloning...")
		}

		// Show the files cloned
		stdout, err := git.GitService.DockerFilesTree(gitRepo)
		if err != nil {
			log.Errorf("Can't show the repo tree '%s': %v", gitRepo.Type.GetRepoDir(), err)
			os.Exit(4)
		}

		if !util.IsLogInDebug() {
			log.Infof("\n%s", stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringP("repo", "r", "", "The repo to clone")
}
