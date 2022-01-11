package kubefilebrowser

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

var (
	Version      string
	GoVersion    string
	GitUrl       string
	GitBranch    string
	GitCommit    string
	GitLatestTag string
	BuildTime    string
	title        = figure.NewFigure("KubeFileBrowser", "doom", true).String()
)

func VersionIfo() string {
	return fmt.Sprintf("Version: %s\nGoVersion: %s\nGitUrl: %s\nGitBranch: %s\nGitCommit: %s\nGitLatestTag: %s\nBuildTime: %s",
		Version, GoVersion, GitUrl, GitBranch, GitCommit, GitLatestTag, BuildTime)
}

func PrintHeadInfo() {
	fmt.Println(title)
	fmt.Println(VersionIfo())
}
