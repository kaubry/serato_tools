package serato

import (
	"testing"
)

type MockHomeDirGetter struct {
	homeDir string
}

func (m MockHomeDirGetter) GetHomeDir() string {
	return m.homeDir
}

func TestGetSeratoDir(t *testing.T) {
	folderTool := FolderTool{}
	for _, test := range seratoDirExpect {
		result, _ := folderTool.GetSeratoDir(&Config{MusicPath: test.path}, MockHomeDirGetter{homeDir: "/Users/test/Music"})
		if result != test.result {
			t.Errorf("expected '%s', got '%s'", test.result, result)
		}
	}
}
