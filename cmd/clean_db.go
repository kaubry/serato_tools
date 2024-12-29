package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/kaubry/serato_tools/logger"
	"github.com/kaubry/serato_tools/serato"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/text/unicode/norm"
)

var cleanDBCommand = &cobra.Command{
	Use:   "clean-db",
	Short: "Clean duplicate entries in Serato database file",
	Long:  "Remove duplicate entries from Serato database file, keeping the oldest version of each track based on import date (tadd)",
	Run:   cleanDuplicates,
}

var db *serato.Database

func init() {
	cleanDBCommand.Flags().StringVarP(&musicDir, "dir", "d", "", "Root directory for your music (required)")
	cleanDBCommand.MarkFlagRequired("dir")
	cleanDBCommand.Flags().BoolVar(&dryRun, "dryrun", false, "Run dry, don't modify the database (Default: false)")
	rootCmd.AddCommand(cleanDBCommand)
}

type duplicateEntry struct {
	path      string
	timestamp []byte // tadd field
	index     int    // Index into db.Dmfs
	toRemove  bool   // Flag to mark entries for removal
	normalizedPath string
}

// Helper function to get timestamp as time.Time
func (de *duplicateEntry) getAddedTime() time.Time {
	if de == nil || db == nil || de.index < 0 || de.index >= len(db.Dmfs) {
		return time.Time{} // Return zero time for invalid entries
	}
	t, err := db.Dmfs[de.index].GetFilePathAddedDateTime()
	if err != nil || t.IsZero() {
		return time.Time{} // Return zero time for entries with missing or invalid timestamps
	}
	return t
}

func cleanDuplicates(cmd *cobra.Command, args []string) {
	logger.Logger.Info("Starting database duplicate cleaning process",
		zap.String("music_dir", musicDir),
		zap.Bool("dry_run", dryRun))

	// Open and read database
	dbPath := filepath.Join(musicDir, "database V2")
	logger.Logger.Debug("Opening database file", zap.String("path", dbPath))

	dbFile, err := os.Open(dbPath)
	if err != nil {
		logger.Logger.Error("Error opening database file", zap.Error(err))
		return
	}
	defer dbFile.Close()

	var dbErr error
	db, dbErr = serato.NewDatabase(dbFile)  // Store in global variable
	if dbErr != nil {
		logger.Logger.Error("Error reading database", zap.Error(dbErr))
		return
	}

	if db == nil {
		logger.Logger.Error("Database is nil after initialization")
		return
	}

	logger.Logger.Info("Successfully loaded database",
		zap.Int("total_entries", len(db.Dmfs)))

	// Initialize maps for duplicate detection
	pathMap := make(map[string][]*duplicateEntry)
	duplicateCount := 0
	totalDuplicates := 0
	totalMissingTimestamps := 0

	// Process each entry
	for i := range db.Dmfs {
		if i < 0 || i >= len(db.Dmfs) {
			logger.Logger.Error("Invalid index while processing entries",
				zap.Int("index", i),
				zap.Int("total_entries", len(db.Dmfs)))
			continue
		}

		path, err := db.Dmfs[i].GetFilePath()
		if err != nil {
			logger.Logger.Error("Error getting file path from entry", zap.Error(err))
			continue
		}

		// Normalize the path
		normalizedPath := string(norm.NFC.Bytes([]byte(path)))

		// Create duplicate entry
		entry := &duplicateEntry{
			path:      path,
			timestamp: db.Dmfs[i].GetFilePathAddedDate(), // Store tadd field for later comparison
			index:     i,  // Store index into db.Dmfs
			normalizedPath: normalizedPath,
		}

		// Validate index before adding to map
		if entry.index < 0 || entry.index >= len(db.Dmfs) {
			logger.Logger.Error("Invalid index in entry",
				zap.Int("index", entry.index),
				zap.Int("total_entries", len(db.Dmfs)))
			continue
		}

		// Add to path map
		pathMap[normalizedPath] = append(pathMap[normalizedPath], entry)
	}

	// Log duplicate summary
	logger.Logger.Info("Starting duplicate analysis")
	var totalToRemove int
	for normalizedPath, entries := range pathMap {
		if len(entries) > 1 {
			logger.Logger.Debug("Processing duplicates",
				zap.String("path", normalizedPath),
				zap.Int("total_entries", len(entries)))

			duplicateCount++
			totalDuplicates += len(entries) - 1 // subtract 1 as one entry is the original
			
			// Validate all entries before sorting
			validEntries := make([]*duplicateEntry, 0, len(entries))
			for _, e := range entries {
				if e == nil {
					logger.Logger.Error("Nil entry found before sorting",
						zap.String("path", normalizedPath))
					continue
				}
				if e.index < 0 || e.index >= len(db.Dmfs) {
					logger.Logger.Error("Invalid index in entry before sorting",
						zap.Int("index", e.index),
						zap.Int("total_entries", len(db.Dmfs)))
					continue
				}
				validEntries = append(validEntries, e)
			}

			// Only process if we have valid entries
			if len(validEntries) == 0 {
				logger.Logger.Error("No valid entries found",
					zap.String("path", normalizedPath))
				continue
			}

			// Sort valid entries by timestamp (oldest first)
			sort.Slice(validEntries, func(i, j int) bool {
				if i < 0 || i >= len(validEntries) || j < 0 || j >= len(validEntries) {
					logger.Logger.Error("Invalid index in sort function",
						zap.Int("i", i),
						zap.Int("j", j),
						zap.Int("entries_len", len(validEntries)))
					return false
				}
				
				timeI := validEntries[i].getAddedTime()
				timeJ := validEntries[j].getAddedTime()
				
				// If either time is zero (missing/invalid timestamp), put it at the end
				if timeI.IsZero() {
					return false
				}
				if timeJ.IsZero() {
					return true
				}
				return timeI.Before(timeJ)
			})

			// Process valid entries
			var oldestTime time.Time
			for i, entry := range validEntries {
				if entry.index < 0 || entry.index >= len(db.Dmfs) {
					logger.Logger.Error("Invalid index while processing entry",
						zap.Int("index", entry.index),
						zap.Int("total_entries", len(db.Dmfs)))
					continue
				}

				addedTime := entry.getAddedTime()
				// Validate time before using it
				hasValidTime := !addedTime.IsZero()
				timeStatus := "valid"
				if !hasValidTime {
					timeStatus = "missing or invalid"
				}

				if hasValidTime {
					if i == 0 {
						oldestTime = addedTime
					}
					// Mark all but the oldest entry (first one) for removal
					entry.toRemove = i > 0
					if entry.toRemove {
						totalToRemove++
						if !dryRun {
							// Create log fields slice with valid fields
							logFields := []zap.Field{
								zap.String("path", entry.normalizedPath),
								zap.Int("index", entry.index),
								zap.String("time_status", timeStatus),
							}
							if hasValidTime {
								logFields = append(logFields, zap.String("added_time", addedTime.String()))
							}

							logger.Logger.Debug("Marking entry for removal", logFields...)
							db.Dmfs[entry.index].RemoveFromDatabase()
							
							// Log removal status with time if available
							statusFields := []zap.Field{
								zap.String("path", entry.normalizedPath),
								zap.Int("index", entry.index),
								zap.Bool("removed_flag", db.Dmfs[entry.index].IsRemoved()),
								zap.String("time_status", timeStatus),
							}
							if hasValidTime {
								statusFields = append(statusFields, zap.String("added_time", addedTime.String()))
							}
							logger.Logger.Debug("Entry marked as removed", statusFields...)
						}
						if dryRun {
							// Show timestamp difference for entries being removed
							dryRunFields := []zap.Field{
								zap.String("path", entry.normalizedPath),
								zap.String("time_status", timeStatus),
							}
							if hasValidTime && !oldestTime.IsZero() {
								timeDiff := addedTime.Sub(oldestTime)
								dryRunFields = append(dryRunFields, 
									zap.String("added_time", addedTime.String()),
									zap.String("newer_by", timeDiff.String()))
							}
							logger.Logger.Info("Would remove duplicate", dryRunFields...)
						}
					} else {
						if dryRun {
							keepFields := []zap.Field{
								zap.String("path", entry.normalizedPath),
							}
							if hasValidTime {
								keepFields = append(keepFields, zap.String("added_time", addedTime.String()))
							}
							logger.Logger.Info("Would keep entry (oldest)", keepFields...)
						} else {
							// Normalize the path of the entry we're keeping to ensure NFC form
							db.Dmfs[entry.index].SetFilePath(entry.normalizedPath)
							logger.Logger.Debug("Normalized path of kept entry",
								zap.String("path", entry.normalizedPath))
						}
					}
				} else {
					// Keep entries with missing timestamps
					entry.toRemove = false
					totalMissingTimestamps++
					if dryRun {
						logger.Logger.Info("Would keep entry (missing timestamp)",
							zap.String("path", entry.normalizedPath))
					}
				}
			}

			logger.Logger.Info(fmt.Sprintf("Found %d duplicates for path:", len(entries)-1),
				zap.String("path", normalizedPath))
			
			for i, entry := range validEntries {
				var taddInfo string
				if entry.timestamp != nil {
					addedTime, err := db.Dmfs[entry.index].GetFilePathAddedDateTime()
					if err != nil {
						taddInfo = fmt.Sprintf("tadd: error parsing timestamp - %v", err)
					} else if !addedTime.IsZero() {
						taddInfo = fmt.Sprintf("added: %v", addedTime.Format(time.RFC3339))
					} else {
						taddInfo = "tadd: missing (will keep)"
					}
				} else {
					taddInfo = "tadd: missing (will keep)"
				}

				status := "will keep"
				if entry.toRemove {
					status = "will remove"
				}

				if i == 0 {
					logger.Logger.Info(fmt.Sprintf("  Entry %d: %s (oldest - %s)", i+1, taddInfo, status))
				} else {
					logger.Logger.Info(fmt.Sprintf("  Entry %d: %s (%s)", i+1, taddInfo, status))
				}
			}
		}
	}

	// Log summary statistics
	logger.Logger.Info("Duplicate cleanup summary",
		zap.Int("total_duplicates_found", totalDuplicates),
		zap.Int("entries_to_remove", totalToRemove),
		zap.Int("entries_with_missing_timestamps", totalMissingTimestamps))

	if dryRun {
		logger.Logger.Info("Dry run completed. No changes were made to the database.",
			zap.String("database_path", dbPath))
		return
	}

	// Create backup before making changes
	backupPath := dbPath + ".bak"
	logger.Logger.Info("Creating database backup", zap.String("backup_path", backupPath))
	
	dbBytes, err := os.ReadFile(dbPath)
	if err != nil {
		logger.Logger.Error("Failed to read database for backup", zap.Error(err))
		return
	}
	
	err = os.WriteFile(backupPath, dbBytes, 0644)
	if err != nil {
		logger.Logger.Error("Failed to write database backup", zap.Error(err))
		return
	}

	// Write changes to database
	logger.Logger.Info("Writing changes to database")
	modifiedDbBytes := db.GetBytes()
	err = os.WriteFile(dbPath, modifiedDbBytes, 0644)
	if err != nil {
		logger.Logger.Error("Failed to write changes to database", zap.Error(err))
		logger.Logger.Info("Original database preserved at", zap.String("backup_path", backupPath))
		return
	}

	// Verify the changes
	logger.Logger.Info("Verifying database changes")
	verifyDb, err := os.Open(dbPath)
	if err != nil {
		logger.Logger.Error("Failed to open database for verification", zap.Error(err))
		logger.Logger.Info("Original database preserved at", zap.String("backup_path", backupPath))
		return
	}
	defer verifyDb.Close()

	verifyDatabase, err := serato.NewDatabase(verifyDb)
	if err != nil {
		logger.Logger.Error("Failed to read modified database for verification", zap.Error(err))
		logger.Logger.Info("Original database preserved at", zap.String("backup_path", backupPath))
		return
	}

	// Verify entry count matches expected
	expectedEntries := len(db.Dmfs) - totalToRemove
	actualEntries := len(verifyDatabase.Dmfs)
	if actualEntries != expectedEntries {
		logger.Logger.Error("Database verification failed: entry count mismatch",
			zap.Int("expected_entries", expectedEntries),
			zap.Int("actual_entries", actualEntries))
		logger.Logger.Info("Original database preserved at", zap.String("backup_path", backupPath))
		return
	}

	logger.Logger.Info("Database cleanup completed successfully",
		zap.Int("entries_removed", totalToRemove),
		zap.Int("entries_remaining", actualEntries),
		zap.String("database_path", dbPath),
		zap.String("backup_path", backupPath))
}
