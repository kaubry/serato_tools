package serato

func getTableForFilePathTest() []TestTable {
	tables := []TestTable{
		{"\\GDrive\\DJ\\2012\\Electro\\Drum n Bass\\test.mp3", "E:\\_Serato_", "E:\\GDrive\\DJ\\2012\\Electro\\Drum n Bass\\test.mp3"},
	}
	return tables
}

var removeVolumeTestCaseExpect = []removeVolumeTestCase{
	{"C:\\example.mp3", "example.mp3"},
	{"C:\\Music\\example.mp3", "Music\\example.mp3"},
	{"C:\\Users\\test\\Desktop\\example.mp3", "Users\\test\\Desktop\\example.mp3"},
	{"SomeRelativePath\\example.mp3", ""},
}
