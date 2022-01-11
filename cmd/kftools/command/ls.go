package command

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type File struct {
	Name    string      `json:"Name"`
	Path    string      `json:"Path"`
	Size    int64       `json:"Size"`
	Mode    string      `json:"Mode"`
	ModTime time.Time   `json:"ModTime"`
	IsDir   bool        `json:"IsDir"`
	Sys     interface{} `json:"SysInfo"`
}

var denyFileOrList = []string{
	"/kftools", "/kftools.exe",
	"/tools/kftools", "/tools/kftools.exe",
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List information about the FILEs (the current directory by default).",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) == 0 {
			path = "."
		} else {
			path = args[0]
		}
		dir, err := ioutil.ReadDir(path)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			os.Exit(253)
		}
		path = strings.Replace(path, "\\", "/", -1)
		path = strings.TrimRight(path, "/")
		var files []File
		for _, d := range dir {
			p := fmt.Sprintf("%s/%s", path, d.Name())
			if inSlice(p, denyFileOrList) {
				continue
			}
			isDir := d.IsDir()
			size := d.Size()
			// check link is file or dir
			if d.Mode()&os.ModeSymlink != 0 {
				s, err := os.Stat(p)
				if err == nil {
					isDir = s.IsDir()
					size = s.Size()
				}
			}
			f := File{
				Name:    d.Name(),
				Path:    p,
				Size:    size,
				Mode:    d.Mode().String(),
				ModTime: d.ModTime(),
				IsDir:   isDir,
				Sys:     d.Sys(),
			}
			files = append(files, f)
		}
		s, err := json.Marshal(files)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			os.Exit(253)
		}
		fmt.Println(string(s))
	},
}

func inSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
