package logx_test

import (
	"testing"

	"go.uber.org/zap/zapcore"

	log "github.com/win5do/go-lib/logx"
)

func TestLog(t *testing.T) {
	log.SetLogger(log.NewLogger(zapcore.DebugLevel))
	log.Debug("debug")
	log.Info("info")
	l := log.GetLogger().Sugar()
	l.Debug("debug")
	l.Info("info")
}
