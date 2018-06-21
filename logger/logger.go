package logger

import (
	"go.uber.org/zap"
	"encoding/json"
)

var Logger *zap.Logger

func init() {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "console",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "capitalColor"
	  }
	}`)
	var cfg zap.Config


	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
	defer logger.Sync()
}

