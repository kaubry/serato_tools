package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"github.com/watershine/serato_crates/logger"
	"path/filepath"
	"go.uber.org/zap"
	"github.com/watershine/serato_crates/serato"
	"os"
	"github.com/watershine/serato_crates/files"
)

var seratoDir string

var cleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "Clean your crates",
	Long:  "Remove the song in the crates that are missing from your computer",
	Run:   cleanCrates,
}

func init() {
	cleanCommand.Flags().StringVarP(&seratoDir, "dir", "d", "", "Serato directory to be cleaned (required)")
	cleanCommand.MarkFlagRequired("dir")

	rootCmd.AddCommand(cleanCommand)
}

func cleanCrates(cmd *cobra.Command, args []string) {
	crates := getCrates(filepath.Join(seratoDir, "Subcrates"))
	for _, c := range crates {
		f, _ := os.Open(c)
		crate := serato.NewCrate(f)
		logger.Logger.Debug("Cleaning crate", zap.String("crate", c))
		before := crate.NumberOfTracks()
		cleanCrate(crate)
		if before != crate.NumberOfTracks() {
			logger.Logger.Debug("Updating crate", zap.String("crate", c))
			files.WriteToFile(c, crate.GetCrateBytes())
		}

	}

}

func cleanCrate(crate *serato.Crate) {
	tracks := crate.TrackList()
	for _, t := range tracks {
		//@TODO check on Windows
		if _, err := os.Stat(string(os.PathSeparator)+t); os.IsNotExist(err) {
			logger.Logger.Info("Removing track from crate", zap.String("track", t))
			crate.RemoveTrack(t)
		}
	}
}

func getCrates(dir string) []string {
	var crates []string
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Logger.Error(err.Error())
	} else {
		for _, f := range fileInfos {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".crate" {
				crates = append(crates, filepath.Join(dir, f.Name()))
			}
		}
	}
	return crates
}


