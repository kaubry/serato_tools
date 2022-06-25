package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"path/filepath"
)

var normalizeFilenameCommand = &cobra.Command{
	Use:   "normalize-filename",
	Short: "Convert the files name to composing form",
	Long:  "Convert all files name with decomposing form to composing form (2 unicode to 1 unicode code point for special characters)  ",
	Run:   normalizeFilename,
}

func init() {
	normalizeFilenameCommand.Flags().StringVarP(&directory, "dir", "d", "", "Directory (recursive) to be normalized (required)")
	normalizeFilenameCommand.MarkFlagRequired("dir")
	normalizeFilenameCommand.Flags().BoolVar(&dryRun, "dryrun", false, "Run dry, it doesn't modify the files. (Default: false)")
	rootCmd.AddCommand(normalizeFilenameCommand)
}

func normalizeFilename(cmd *cobra.Command, args []string) {
	filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() {
			fileName := d.Name()

			normalized := norm.NFC.Bytes([]byte(fileName))

			if string(normalized) != fileName {

				fmt.Printf("%s needs to be normalized", path)
				fmt.Printf("\n")

				fmt.Printf("normalized string: ")
				fmt.Printf("%s", normalized)
				fmt.Printf("\n")

				fmt.Printf("quoted normalized string: ")
				fmt.Printf("%+q", normalized)
				fmt.Printf("\n")

				fmt.Printf("\n")

				oldPath := fmt.Sprintf("%s", path)
				currentDirectory := filepath.Dir(path)

				newPath := filepath.Join(currentDirectory, string(normalized))

				fmt.Printf("Old Path: %s: ", oldPath)
				fmt.Printf("\n")
				fmt.Printf("New Path: %s: ", newPath)

				fmt.Printf("\n")

				if !dryRun {
					e := os.Rename(oldPath, newPath)
					if e != nil {
						log.Fatal(e)
					}
				}
			}
		}
		return nil
	})
}
