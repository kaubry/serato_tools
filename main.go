package main

import "github.com/kaubry/serato_tools/cmd"

func main() {
	cmd.Execute()

	//f, err := os.Open("./database V2")
	//check(err)
	//d := serato.NewDatabase(f)
	//logger.Logger.Debug(d.String())
	//for _, t := range d.Dmfs {
	//	logger.Logger.Debug(t.String())
	//}
	//
	//files.WriteToFile("./database V2_1", d.GetBytes())
}

//func readCrate(path string) {
//	f, err := os.Open(path)
//	check(err)
//	NewCrate(f)
//track := crate.TrackList()[0]

//fmt.Printf("%s\n", crate)
//for _, c := range crate.columns {
//	fmt.Printf("%s\n", c)
//}
//fmt.Printf("---- Tracks ----\n")
//for _, t := range crate.tracks {
//	fmt.Printf("%s\n", t)
//
//}
//fmt.Printf("%s", crate.tracks[0])

//t := NewTrack("Users/kevin/Downloads/Haila - De Donde Vengo.mp3")
//fmt.Printf("%v", t.Equals(crate.tracks[0]))
//display := crate.brev
//fmt.Printf("%d bytes: %s\n", len(display), string(display))
//}

//func createCrateWithTracks(path string, columns []ColumnName, files []string) {
//	crate
//}

//func getCrates(map[string][]string) map[string][]Crate {
// return nil
//}
