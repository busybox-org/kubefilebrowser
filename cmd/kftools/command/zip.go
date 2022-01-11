package command

import (
	"archive/zip"
	"bufio"
	"github.com/spf13/cobra"
	"github.com/xmapst/kubefilebrowser/utils"
	"github.com/xmapst/kubefilebrowser/utils/ratelimit"
	"github.com/xmapst/kubefilebrowser/utils/symwalk"
	"io"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(zipCmd)
}

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "The default action is to add or replace zipfile entries from list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zw := zip.NewWriter(ratelimit.Writer(os.Stdout, ratelimit.New(2000*1024))) // 限制输出 2000KB/s
		defer func() {
			_ = zw.Close()
		}()
		// 监听标准输入, 获取是否退出
		go closeEvent()
		for _, p := range args {
			if !utils.FileOrPathExist(p) {
				continue
			}
			_ = makeZip(p, zw)
		}
	},
}

func makeZip(inFilepath string, zw *zip.Writer) error {
	return symwalk.Walk(inFilepath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		// 目录拉平
		//relPath := strings.TrimPrefix(filePath, filepath.Dir(inFilepath))
		var zwPath = utils.ToLinuxPath(filePath)
		zipFile, err := zw.Create(strings.TrimPrefix(zwPath, "/"))
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer func() {
			_ = fsFile.Close()
		}()
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
}

func closeEvent() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			continue
		}
		switch string(line) {
		case "close":
			os.Exit(0)
		}
	}
}
