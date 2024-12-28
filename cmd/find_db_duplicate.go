package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kaubry/serato_tools/logger"
	"github.com/kaubry/serato_tools/serato"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var fdbCommand = &cobra.Command{
	Use:   "fdb",
	Short: "Find duplicate entries in Serato database file",
	Run:   findDatabaseDuplicate,
}

func init() {
	fdbCommand.Flags().StringVarP(&musicDir, "dir", "d", "", "Root directory for your music")
	rootCmd.AddCommand(fdbCommand)
}

func findDatabaseDuplicate(cmd *cobra.Command, args []string) {
	config := &serato.Config{
		MusicPath: musicDir,
	}

	// Get the database file path
	seratoDir, err := serato.GetSeratoDir(config)
	if err != nil {
		logger.Logger.Error("Error getting Serato directory", zap.Error(err))
		return
	}

	// Open and read database
	dbFile, err := os.Open(seratoDir + string(os.PathSeparator) + "database V2")
	if err != nil {
		logger.Logger.Error("Error opening database file", zap.Error(err))
		return
	}
	defer dbFile.Close()

	db, err := serato.NewDatabase(dbFile)
	if err != nil {
		logger.Logger.Error("Error reading database", zap.Error(err))
		return
	}

	// Create a map to track paths and their occurrences
	pathCount := make(map[string][]int) // map[filepath][]index

	// Find duplicates
	for i, dmf := range db.Dmfs {
		path, err := dmf.GetFilePath()
		if strings.Contains(path, "El Secreto") {
			fmt.Printf("Path: %s\n", path)
			fmt.Printf("Byte array: %v\n", dmf.GetFilePathByteArray())
			fmt.Printf("Added date: %v\n", dmf.GetFilePathAddedDate())

		}
		if err != nil {
			logger.Logger.Error("Error getting file path", zap.Error(err))
			continue
		}
		pathCount[path] = append(pathCount[path], i)
	}

	// Print duplicates
	duplicateCount := 0
	for path, indices := range pathCount {
		if len(indices) > 1 {
			duplicateCount++
			fmt.Printf("\n----- Found duplicate for path: %s ------\n", path)
			fmt.Printf("Appears %d times at indices: %v\n", len(indices), indices)
		}
	}

	fmt.Printf("\nTotal number of duplicate entries found: %d\n", duplicateCount)
	fmt.Printf("Total number of entries in database: %d\n", len(db.Dmfs))
}
