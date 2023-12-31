package serato

import (
	"github.com/kaubry/serato_tools/files"
	"github.com/kaubry/serato_tools/logger"
	"os"
	"testing"
)

func TestNewCrate(t *testing.T) {
	f, _ := os.Open("../Bailamos.crate")
	test, _ := files.ReadBytesWithOffset(f, 0, 4)

	//crate, _ := NewCrate(f)
	logger.Logger.Info(string(test))
}
