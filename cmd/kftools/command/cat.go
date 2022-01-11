package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	rootCmd.AddCommand(catCmd)
}

// catCmd represents the cat command
var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "Concatenate FILE(s) to standard output.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
	},
}
