package cmd

import (
	"github.com/spf13/cobra"
	"github.com/watershine/serato_tools/files"
	"github.com/watershine/serato_tools/serato"
)

var volume string
var musicDir string
var rootCrate string

var syncCommand = &cobra.Command{
	Use: "sync",
	Short: "Sync a folder (and all his subfolder with Serato)",
	Run: sync,
}

func init() {

	syncCommand.Flags().StringVarP(&musicDir, "dir", "d", "", "Root directory for your music (required)")
	syncCommand.MarkFlagRequired("dir")

	syncCommand.Flags().StringVarP(&rootCrate, "crate", "c", "", "Parent crate name")
	syncCommand.MarkFlagRequired("crate")

	rootCmd.AddCommand(syncCommand)
}

func sync(cmd *cobra.Command, args []string) {
	f := files.ListFiles(musicDir, serato.GetSupportedExtension())
	config := &serato.Config{
		MusicPath: musicDir,
		RootCrate: rootCrate,
	}
	serato.CreateCrates(f, config)
}
