package serato

import (
	"github.com/kaubry/serato_tools/logger"
	"os"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	f, _ := os.Open("../database V2")
	db, _ := NewDatabase(f)
	logger.Logger.Info(db.String())
}
