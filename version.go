package kubefilebrowser

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

var (
	GoVersion string
	GitUrl    string
	GitBranch string
	GitCommit string
	BuildTime string
	title     = figure.NewFigure("KubeFileBrowser", "doom", true).String()
)

func VersionIfo() string {
	return fmt.Sprintf("GoVersion: %s\nGitUrl: %s\nGitBranch: %s\nGitCommit: %s\nBuildTime: %s",
		GoVersion, GitUrl, GitBranch, GitCommit, BuildTime)
}

func PrintHeadInfo() {
	fmt.Println(title)
	fmt.Println(VersionIfo())
}
