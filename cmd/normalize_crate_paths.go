package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaubry/serato_tools/logger"
	"github.com/kaubry/serato_tools/serato"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/text/unicode/norm"
)

var normalizeCratePathsCommand = &cobra.Command{
	Use:   "normalize-crate-paths",
	Short: "Normalize file paths in crate files",
	Long:  "Iterate through crate files and normalize all file paths within them",
	Run:   normalizeCratePaths,
}

var crateDirectory string

func init() {
	normalizeCratePathsCommand.Flags().StringVarP(&crateDirectory, "dir", "d", "", "Directory containing crate files")
	normalizeCratePathsCommand.MarkFlagRequired("dir")
	normalizeCratePathsCommand.Flags().BoolVar(&dryRun, "dryrun", false, "Run dry, it doesn't modify the files. (Default: false)")
	rootCmd.AddCommand(normalizeCratePathsCommand)
}

func normalizeCratePaths(cmd *cobra.Command, args []string) {
	totalCrates := 0
	cratesWithNonNormalizedPaths := 0
	totalTracksToNormalize := 0

	// Walk through the directory
	err := filepath.WalkDir(crateDirectory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if file is a crate file
		if filepath.Ext(path) != ".crate" {
			return nil
		}

		totalCrates++
		logger.Logger.Info("Processing crate file", zap.String("path", path))

		// Open and read the crate file
		f, err := os.Open(path)
		if err != nil {
			logger.Logger.Error("Failed to open crate file", zap.Error(err))
			return nil
		}
		defer f.Close()

		crate, err := serato.NewCrate(f)
		if err != nil {
			logger.Logger.Error("Failed to read crate file", zap.Error(err))
			return nil
		}

		// Get all tracks
		tracks := crate.TrackList()
		modified := false
		crateNonNormalizedCount := 0

		logger.Logger.Debug("Crate details", 
			zap.String("crate", path),
			zap.Int("track_count", len(tracks)))

		// Process each track
		for _, track := range tracks {
			normalizedPath := norm.NFC.String(track)
			
			// Debug: print hex representation of both strings
			originalBytes := []byte(track)
			normalizedBytes := []byte(normalizedPath)
			
			if normalizedPath != track {
				logger.Logger.Debug("Track comparison",
					zap.String("original", track),
					zap.String("normalized", normalizedPath),
					zap.String("original_hex", fmt.Sprintf("%x", originalBytes)),
					zap.String("normalized_hex", fmt.Sprintf("%x", normalizedBytes)))

				logger.Logger.Info("Normalizing track path",
					zap.String("file", path),
					zap.String("from", track),
					zap.String("to", normalizedPath))

				crateNonNormalizedCount++

				if !dryRun {
					// Remove old track and add normalized one
					crate.RemoveTrack(track)
					crate.AddTrack(normalizedPath)
					modified = true
				}
			}
		}

		if crateNonNormalizedCount > 0 {
			cratesWithNonNormalizedPaths++
			totalTracksToNormalize += crateNonNormalizedCount
		}

		// If changes were made and not in dry run mode, save the crate
		if modified && !dryRun {
			logger.Logger.Info("Saving modified crate", 
				zap.String("path", path),
				zap.Int("changes", crateNonNormalizedCount))

			// Create a backup of the original file
			backupPath := path + ".bak"
			err = os.Rename(path, backupPath)
			if err != nil {
				logger.Logger.Error("Failed to create backup", zap.Error(err))
				return nil
			}

			// Write the modified crate
			newFile, err := os.Create(path)
			if err != nil {
				logger.Logger.Error("Failed to create new crate file", zap.Error(err))
				return nil
			}
			defer newFile.Close()

			_, err = newFile.Write(crate.GetCrateBytes())
			if err != nil {
				logger.Logger.Error("Failed to write new crate file", zap.Error(err))
				// Try to restore from backup
				os.Rename(backupPath, path)
				return nil
			}

			// Remove backup if everything succeeded
			os.Remove(backupPath)
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error("Error walking directory", zap.Error(err))
	}

	// Log summary
	logger.Logger.Info("Normalization scan complete",
		zap.Int("total_crates", totalCrates),
		zap.Int("crates_with_non_normalized_paths", cratesWithNonNormalizedPaths),
		zap.Int("total_tracks_to_normalize", totalTracksToNormalize))
}
