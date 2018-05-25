package cmd

import (
	"github.com/spf13/cobra"
	"github.com/watershine/serato_crates/files"
	"github.com/watershine/serato_crates/serato"
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
	syncCommand.Flags().StringVarP(&volume, "volume", "v", "", "Volume where your music directory is (required)")
	syncCommand.MarkFlagRequired("volume")

	syncCommand.Flags().StringVarP(&musicDir, "dir", "d", "", "Root directory for your music (required)")
	syncCommand.MarkFlagRequired("dir")

	syncCommand.Flags().StringVarP(&rootCrate, "crate", "c", "", "Parent crate name (optional).")

	rootCmd.AddCommand(syncCommand)
}

func sync(cmd *cobra.Command, args []string) {
	//@Todo fix issue with volume and music path
	//path := filepath.Join(volume, musicDir)
	f := files.ListFiles(musicDir)
	config := &serato.Config{
		VolumePath: volume,
		MusicPath: musicDir,
		RootCrate: rootCrate,
	}
	serato.CreateCrates(f, config)
}
