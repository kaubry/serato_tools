package cmd

import (
	"github.com/kaubry/serato_tools/serato"
	"github.com/spf13/cobra"
	"log"
)

var sortCratesCommand = &cobra.Command{
	Use:   "sort-crates",
	Short: "Sort crates alphabetically",
	Long:  "Sort the creates under a specific path alphabetically",
	Run:   sort,
}

func init() {

	sortCratesCommand.Flags().StringVar(&prefFilePath, "file", "", "neworder.pref filepath (required)")
	sortCratesCommand.MarkFlagRequired("file")

	sortCratesCommand.Flags().StringVar(&path, "path", "", "Root path under which to sort")
	sortCratesCommand.MarkFlagRequired("path")

	rootCmd.AddCommand(sortCratesCommand)
}

func sort(cmd *cobra.Command, args []string) {
	log.Printf("File to be sorted: %s\n", prefFilePath)
	log.Printf("Path to be sorted: %s\n", path)

	err := serato.SortPrefixLineInFile(prefFilePath, path)
	if err != nil {
		log.Fatal(err)
	}
}
