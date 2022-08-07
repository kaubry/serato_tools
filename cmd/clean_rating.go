package cmd

import (
	"fmt"
	"github.com/bogem/id3v2"
	fileTools "github.com/kaubry/serato_tools/files"
	"github.com/spf13/cobra"
	"gopkg.in/fatih/set.v0"
	"log"
)

var cleanRatingCommand = &cobra.Command{
	Use:   "clean-rating",
	Short: "Clean your ratings",
	Long:  "Remove the BOM characters from the rating tag",
	Run:   cleanRatings,
}

func init() {
	cleanRatingCommand.Flags().StringVarP(&musicDir, "dir", "d", "", "Root directory for your music (required)")
	cleanRatingCommand.MarkFlagRequired("dir")

	cleanRatingCommand.Flags().BoolVar(&dryRun, "dryrun", false, "Run dry, it doesn't modify the files. (Default: false)")

	rootCmd.AddCommand(cleanRatingCommand)
}

func cleanRatings(cmd *cobra.Command, args []string) {
	supportedExtension := set.New(set.ThreadSafe)
	supportedExtension.Add(".mp3")
	folders := fileTools.ListFiles(musicDir, supportedExtension)
	counter := 0
	for _, files := range folders {
		for _, file := range files {
			tag, err := id3v2.Open(file, id3v2.Options{Parse: true})
			if err != nil {
				log.Printf("Error while opening mp3 file: ", err, file)
			}
			if tag != nil {
				defer tag.Close()
				tf := tag.GetTextFrame(tag.CommonID("Composer"))
				composerBytes := []byte(tf.Text)
				if startWithBom(composerBytes) {
					//cleanedComposer := composerBytes[3:]
					fmt.Println(file)
					//tag.AddTextFrame(tag.CommonID("Composer"), tf.Encoding, string(cleanedComposer))
					//if err = tag.Save(); err != nil {
					//	log.Fatal("Error while saving a tag: ", err)
					//} else {
					//	counter++
					//}
				}
			}
		}
	}
	fmt.Printf("Cleaned %d files", counter)
}

func startWithBom(text []byte) bool {
	return  len(text) > 3 && text[0] == 239 && text[1] == 187 && text[2] == 191
}
