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
	"github.com/marcellodesales/cloner/config"
	"github.com/marcellodesales/cloner/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

//The verbose flag value
var v string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloner",
	Short: "Managing github clones with a single command",
	Long: `cloner knows how to clone version-control urls by simply making sure
it can place software in specific location designed.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Setup Viper Configuration file type
	rootCmd.PersistentFlags().StringVar(&cfgFile, "cloner", "cloner", "Config file (default is $HOME/.cloner.yaml)")

	// Callbacks to initilize, in order, for cobra-viper-anything else
	cobra.OnInitialize(
		func() {
			// Sets up Viper with the given config file name
			util.SetupViperHomeConfig(cfgFile, ".cloner", "yaml", "CLONER")
		},
		func() {
			// Setups up the configuration parsed from Viper config (maps -> structs)
			config.Setup()
		},
	)

	// Setup Logger https://le-gall.bzh/post/go/integrating-logrus-with-cobra/
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := util.SetUpLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}

	//Here is where we bind the verbose flag
	//Default value is the info level
	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
