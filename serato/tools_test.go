package serato

import "testing"

type TestTable struct {
	path      string
	seratoDir string
	result    string
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


