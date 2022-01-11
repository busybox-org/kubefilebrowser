package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(mkdirCmd)
}

// mkdirCmd represents the mkdir command
var mkdirCmd = &cobra.Command{
	Use:     "mkdir",
	Aliases: []string{"mk"},
	Short:   "Create the DIRECTORY(ies), if they do not already exist.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll(args[0], os.ModeDir)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
	},
}
