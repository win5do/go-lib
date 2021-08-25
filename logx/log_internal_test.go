package logx

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func setDiscardLogger() {
	SetLogger(
		zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			&zaptest.Discarder{},
			zap.DebugLevel,
		)),
	)
}

func BenchmarkLock(b *testing.B) {
	setDiscardLogger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Info("info")
		}
	})
}

func BenchmarkUnlock(b *testing.B) {
	setDiscardLogger()

	unlockLog := func(args ...interface{}) {
		_globalS.Info(args...)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			unlockLog("info")
		}
	})
}
