package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(mvCmd)
}

// mvCmd represents the mv command
var mvCmd = &cobra.Command{
	Use:     "mv",
	Aliases: []string{"move", "rename"},
	Short:   "Rename SOURCE to DEST, or move SOURCE(s) to DIRECTORY.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Rename(args[0], args[1])
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			os.Exit(253)
		}
	},
}
