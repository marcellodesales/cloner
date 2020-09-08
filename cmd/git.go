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
	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "git",
	Short: "Clones a given git repo",
	Long: `Clones a given git repo URL`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("config.git.cloneBaseDir=%s", config.INSTANCE.Git.CloneBaseDir)
		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			log.Errorf("You must provide the repo URL")
			os.Exit(1)
		}
		log.Debugf("Setting up for %s", repo)

		gitRepo, err := git.GitService.ParseRepoString(repo)
		if err != nil {
			log.Errorf("git URL invalid: %v", err)
			os.Exit(2)
		}

		forceClone, _ := cmd.Flags().GetBool("force")
		if forceClone {
			log.Info("Forcing clone...")
		}

		deletedExistingDir, err := git.GitService.VerifyCloneDir(gitRepo, forceClone, config.INSTANCE)
		if deletedExistingDir {
			log.Infof("Deleted dir '%s'", gitRepo.CloneLocation)
		}
		if err != nil {
			log.Errorf("Can't clone repo: %v", err)
			log.Errorf("You can specify --force or -f to delete the existing dir and clone again. " +
				"Make sure there are no panding changes!")
			os.Exit(3)
		}

		err = git.GitService.MakeCloneDir(gitRepo, config.INSTANCE)
		if err != nil {
			log.Errorf("Can't create the base clone repo '%s': %v", gitRepo.Type.GetUserDir(), err)
			os.Exit(4)
		}

		log.Infof("Cloning into '%s'", gitRepo.CloneLocation)
		err = git.GitService.GoCloneRepo(gitRepo, config.INSTANCE)
		if err != nil {
			log.Errorf("Can't clone the repo at '%s': %v", gitRepo.Type.GetRepoDir(), err)
			os.Exit(5)

		} else {
			log.Info("Done...")
		}

		// Show the files cloned
		stdout, err := git.GitService.DockerFilesTree(gitRepo, config.INSTANCE)
		if err != nil {
			log.Errorf("Can't show the cloned repo tree '%s': %v", gitRepo.Type.GetRepoDir(), err)
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
	initCmd.Flags().StringP("repo", "r", "", "The repo URL to clone")

	var verbose = false
	// https://github.com/spf13/cobra/issues/818#issuecomment-489021216
	initCmd.Flags().BoolVarP(&verbose, "force", "f", false, "Forces cloning by deleting existing dir")
}
