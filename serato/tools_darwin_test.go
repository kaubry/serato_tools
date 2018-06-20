package serato

func getTableForFilePathTest() []TestTable {
	tables := []TestTable{
		{"/GDrive/DJ/Hip-Hop/Latin/Test.mp3", "/Volumes/128Go SD/GDrive/_Serato_", "/Volumes/128Go SD/GDrive/DJ/Hip-Hop/Latin/Test.mp3"},
		{"/Users/test/Desktop/Test.mp3", "/Users/kevin/Music/_Serato_", "/Users/test/Desktop/Test.mp3"},
	}
	return tables
}
