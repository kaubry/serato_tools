package serato

import (
	"testing"
)

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

func TestRemoveMusicPathFromPath(t *testing.T) {
	file := "/home/user/music/song.mp3"
	expected := "song.mp3"

	// Create a sample Config
	config := &Config{
		MusicPath: "/home/user/music/",
	}

	result := removeMusicPathFromPath(file, config)

	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

type MockSeratoDirGetter struct {
	Result string
	Err    error
}

func (m MockSeratoDirGetter) GetSeratoDir(c *Config, hdg HomeDirGetter) (string, error) {
	return m.Result, m.Err
}

func TestGetSubcrateFolder(t *testing.T) {
	config := Config{
		MusicPath: "/path/to/music",
		RootCrate: "Root",
	}

	mockSeratoDirGetter := MockSeratoDirGetter{
		Result: "/path/to/serato",
		Err:    nil,
	}

	subcrateFolder, err := GetSubcrateFolder(&config, mockSeratoDirGetter)
	if err != nil {
		t.Errorf("GetSubcrateFolder returned an error: %v", err)
	}

	expectedSubcrateFolder := "/path/to/serato/Subcrates"
	if subcrateFolder != expectedSubcrateFolder {
		t.Errorf("Subcrate folder does not match. Expected: %q, Got: %q", expectedSubcrateFolder, subcrateFolder)
	}
}
