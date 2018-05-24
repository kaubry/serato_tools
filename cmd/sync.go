package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	rootCmd.AddCommand(syncCommand)
}

var syncCommand = &cobra.Command{
	Use: "sync",
	Short: "Sync a folder (and all his subfolder with Serato)",
	Run: sync,
}

func sync(cmd *cobra.Command, args []string) {
		fmt.Println("hello there")
}
