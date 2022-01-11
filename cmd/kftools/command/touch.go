package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
)

func init() {
	rootCmd.AddCommand(touchCmd)
}

// touchCmd represents the touch command
var touchCmd = &cobra.Command{
	Use:     "touch",
	Aliases: []string{"echo"},
	Short:   "Get content from standard input and write it to file.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll(path.Dir(args[0]), os.ModeDir)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
		f, err := os.Create(args[0])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
		_, err = io.Copy(f, os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
	},
}
