package cmd

import (
	"fmt"
	"os"

	"github.com/kaubry/serato_tools/logger"
	"github.com/kaubry/serato_tools/serato"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	databaseFileName = "database V2"
)

var fdbCommand = &cobra.Command{
	Use:   "fdb",
	Short: "Find duplicate entries in Serato database file",
	Run:   findDatabaseDuplicate,
}

//var musicDir string

func init() {
	fdbCommand.Flags().StringVarP(&musicDir, "dir", "d", "", "Root directory for your music")
	rootCmd.AddCommand(fdbCommand)
}

// handleErr logs the error and returns whether to continue execution
func handleErr(err error, message string) bool {
	if err != nil {
		logger.Logger.Error(message, zap.Error(err))
		return false // Stop execution
	}
	return true // Continue execution
}

func findDatabaseDuplicate(cmd *cobra.Command, args []string) {
	config := &serato.Config{
		MusicPath: musicDir,
	}

	// Get the Serato directory
	seratoDir, err := serato.GetSeratoDir(config)
	if !handleErr(err, "Error getting Serato directory") {
		return
	}

	// Construct the database file path
	dbPath := seratoDir + string(os.PathSeparator) + databaseFileName

	// Open the database file
	dbFile, err := os.Open(dbPath)
	if !handleErr(err, fmt.Sprintf("Error opening database file: %s", dbPath)) {
		return
	}
	defer dbFile.Close()

	// Read the database
	db, err := serato.NewDatabase(dbFile)
	if !handleErr(err, "Error reading database") {
		return
	}

	// Create a map to track paths and their occurrences
	pathCount := make(map[string][]int) // map[filepath][]index

	// Find duplicates
	for i, dmf := range db.Dmfs {
		path, err := dmf.GetFilePath()
		if !handleErr(err, "Error getting file path") {
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
