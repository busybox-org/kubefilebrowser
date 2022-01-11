package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "del", "delete"},
	Short:   "Remove (unlink) the FILE(s).",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, p := range args {
			err := os.RemoveAll(p)
			if err != nil {
				_, _ = fmt.Fprint(os.Stderr, err)
			}
		}
	},
}
