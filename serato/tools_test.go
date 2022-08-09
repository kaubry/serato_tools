package serato

import "testing"

type TestTable struct {
	path      string
	seratoDir string
	result    string
}

type pathTestCase struct {
	path   string
	result string
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

func TestGetVolume(t *testing.T) {
	for _, test := range getVolumeExpect {
		result := GetVolume(test.path)
		if result != test.result {
			t.Errorf("expected '%s', got '%s'", test.result, result)
		}
	}
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

func TestGetSeratoDir(t *testing.T) {
	setHomeDir()

	for _, test := range seratoDirExpect {
		result, _ := GetSeratoDir(&Config{MusicPath: test.path})
		if result != test.result {
			t.Errorf("expected '%s', got '%s'", test.result, result)
		}
	}
}
