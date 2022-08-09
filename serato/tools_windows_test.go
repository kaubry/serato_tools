package serato

func getTableForFilePathTest() []TestTable {
	tables := []TestTable{
		{"\\GDrive\\DJ\\2012\\Electro\\Drum n Bass\\test.mp3", "E:\\_Serato_", "E:\\GDrive\\DJ\\2012\\Electro\\Drum n Bass\\test.mp3"},
	}
	return tables
}

func setHomeDir() {
	defaultHomeDirGetter = stubHomeDirGetter{"C:\\Users\\TestUser\\Music"}
}

var removeVolumeTestCaseExpect = []pathTestCase{
	{"C:\\example.mp3", "example.mp3"},
	{"C:\\Music\\example.mp3", "Music\\example.mp3"},
	{"C:\\Users\\test\\Desktop\\example.mp3", "Users\\test\\Desktop\\example.mp3"},
	{"SomeRelativePath\\example.mp3", ""},
}

var getVolumeExpect = []pathTestCase{
	{"C:\\example.mp3", "C:"},
	{"C:\\Music\\example.mp3", "C:"},
	{"C:\\Users\\test\\Desktop\\example.mp3", "C:"},
	{"D:\\Users\\test\\Desktop\\example.mp3", "D:"},
	{"\\Users\\test\\Desktop\\example.mp3", ""},
	{"SomeRelativePath\\example.mp3", ""},
}

var seratoDirExpect = []pathTestCase{
	{"C:\\example.mp3", "C:\\Users\\TestUser\\Music\\_Serato_"},
	{"C:\\Music\\example.mp3", "C:\\Users\\TestUser\\Music\\_Serato_"},
	{"C:\\Users\\test\\Desktop\\example.mp3", "C:\\Users\\TestUser\\Music\\_Serato_"},
	{"D:\\Users\\test\\Desktop\\example.mp3", "D:\\_Serato_"},
	{"\\Users\\test\\Desktop\\example.mp3", ""},
}
