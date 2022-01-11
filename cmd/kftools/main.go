package main

import (
	"fmt"
	"github.com/xmapst/kubefilebrowser/cmd/kftools/command"
)

var (
	BuildAt string
	GitHash string
)

func main() {
	command.Execute(fmt.Sprintf("Hash: %s\nBuildDate: %s", GitHash, BuildAt))
}
