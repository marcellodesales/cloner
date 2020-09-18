package config

import (
	"fmt"
	"strings"
)

var VersionBuildTime string
var VersionBuildGitSha string
var VersionBuildGoModule string
var VersionBuildNumber string

func ShowVersionDetails() {
	// Print the default docker user when in debug mode
	params := []interface{}{VersionBuildGoModule, VersionBuildGitSha, VersionBuildNumber, strings.ReplaceAll(VersionBuildTime, "_", " ")}
	fmt.Println()
	fmt.Printf("cloner %s@%s version %s built on %s", params...)
	fmt.Println()
}
