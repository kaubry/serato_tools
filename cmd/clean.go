package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"github.com/watershine/serato_tools/logger"
	"path/filepath"
	"go.uber.org/zap"
	"github.com/watershine/serato_tools/serato"
	"os"
	"github.com/watershine/serato_tools/files"
)

var seratoDir string
var database, dryRun bool

var cleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "Clean your crates",
	Long:  "Remove the song in the crates that are missing from your computer",
	Run:   cleanCrates,
}

func init() {
	cleanCommand.Flags().StringVarP(&seratoDir, "dir", "d", "", "Serato directory to be cleaned (required)")
	cleanCommand.MarkFlagRequired("dir")

	cleanCommand.Flags().BoolVar(&database, "database", true, "Clean the database file. (Default: true)")
	cleanCommand.Flags().BoolVar(&dryRun, "dryrun", false, "Run dry, it doesn't modify the files. (Default: false)")

	rootCmd.AddCommand(cleanCommand)
}

func cleanCrates(cmd *cobra.Command, args []string) {
	crates := getCrates(filepath.Join(seratoDir, "Subcrates"))
	for _, c := range crates {
		f, _ := os.Open(c)
		crate := serato.NewCrate(f)
		logger.Logger.Info("Reading crate", zap.String("crate", c))
		before := crate.NumberOfTracks()
		cleanCrate(crate)
		if before != crate.NumberOfTracks() {
			logger.Logger.Info("Updating crate", zap.String("crate", c))
			files.WriteToFile(c, crate.GetCrateBytes())
		} else {
			logger.Logger.Info("Nothing to clean", zap.String("crate", c))
		}
	}
	if database {
		cleanDatabase()
	}
}

func cleanCrate(crate *serato.Crate) {
	tracks := crate.TrackList()
	for _, t := range tracks {
		filePath, _ := serato.GetFilePath(string(os.PathSeparator)+t, seratoDir)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			logger.Logger.Info("Removing track from crate", zap.String("track", filePath))
			if !dryRun {
				crate.RemoveTrack(t)
			}
		}
	}
}

func cleanDatabase() {
	f, err := os.Open(seratoDir + string(os.PathSeparator) + "database V2")
	if err != nil {
		logger.Logger.Error(err.Error())
	} else {
		db := serato.NewDatabase(f)
		logger.Logger.Info("Cleaning Serato Database")
		dmfPaths := db.GetMusicFiles()
		before := len(dmfPaths)
		for _, p := range dmfPaths {
			filePath, _ := serato.GetFilePath(string(os.PathSeparator)+p, seratoDir)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				logger.Logger.Info("Removing music file from database", zap.String("music file", filePath))
				if !dryRun {
					db.RemoveMusicFile(p)
				}
			}
		}
		if before != len(db.GetMusicFiles()) {
			logger.Logger.Info("Updating Database")
			files.WriteToFile(f.Name(), db.GetBytes())
		} else {
			logger.Logger.Info("Nothing to clean")
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
