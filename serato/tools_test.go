package serato

import "testing"

type TestTable struct {
	path      string
	seratoDir string
	result    string
}

func TestGetFilePath(t *testing.T) {
	defaultHomeDirGetter = stubHomeDirGetter{"/Users/test"}

	for _, table := range getTableForFilePathTest() {
		s, e := GetFilePath(table.path, table.seratoDir)
		if e != nil {
			t.Errorf("Error: %s", e.Error())
		} else if s != table.result {
			t.Errorf("File path for: %s with serato dir: %s. Got %s, expected %s", table.path, table.seratoDir, s, table.result)
		}
	}
}

type getDarwinVolumeTestCase struct {
	path   string
	result string
}

func TestGetDarwinVolume(t *testing.T) {
	for _, test := range []getDarwinVolumeTestCase{
		{"/Volumes/TestVolume1", "/Volumes/TestVolume1"},
		{"/Volumes/TestVolume1/", "/Volumes/TestVolume1"},
		{"/Volumes/Test-Volume1", "/Volumes/Test-Volume1"},
		{"/Volumes/Test-Volume1/", "/Volumes/Test-Volume1"},
		{"/Volumes/TestVolume1/example.mp3", "/Volumes/TestVolume1"},
		{"/Users/test/Desktop/example.mp3", "/"},
		{"SomeRelativePath/example.mp3", ""},
	} {
		result := GetDarwinVolume(test.path)
		if result != test.result {
			t.Errorf("expected '%s', got '%s'", test.result, result)
		}
	}
}

type removeVolumeTestCase struct {
	path   string
	result string
}

func TestRemoveVolumeFromPath(t *testing.T) {
	for _, test := range removeVolumeTestCaseExpect {
		result, _ := RemoveVolumeFromPath(test.path)
		if result != test.result {
			t.Errorf("expected '%s', got '%s'", test.result, result)
		}
	}
}

type stubHomeDirGetter struct {
	homeDir string
}

func (dirGetter stubHomeDirGetter) getHomeDir() string {
	return dirGetter.homeDir
}

type getSeratoDirTestCase struct {
	path   string
	result string
}

func TestGetSeratoDir(t *testing.T) {
	defaultHomeDirGetter = stubHomeDirGetter{"/Users/test/Music"}

	for _, test := range []getSeratoDirTestCase{
		{"/Volumes/TestVolume1/Music/example.mp3", "/Volumes/TestVolume1/_Serato_"},
		{"/Users/test/Desktop/example.mp3", "/Users/test/Music/_Serato_"},
		{"/", "/Users/test/Music/_Serato_"},
		{"/SomeDir/example.mp3", "/Users/test/Music/_Serato_"},
		{"SomeRelativePath/example.mp3", ""},
	} {
		result, _ := GetSeratoDir(&Config{MusicPath: test.path})
		if result != test.result {
			t.Errorf("expected '%s', got '%s'", test.result, result)
		}
	}
}
