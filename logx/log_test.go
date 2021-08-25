package logx_test

import (
	"testing"
	"time"

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

func TestDataRace(t *testing.T) {
	go func() {
		for {
			log.Info("info")
			time.Sleep(10 * time.Millisecond)
		}
	}()

	cancel := time.After(1 * time.Second)
	for {
		select {
		case <-cancel:
			return
		default:
			log.SetLogger(log.NewLogger(zapcore.DebugLevel))
			time.Sleep(30 * time.Millisecond)
		}
	}
}
