package logx

import (
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate sh -c "go run ./generator >zap_sugar_generated.go"

var (
	_globalMu = sync.RWMutex{}
	_globalL  = NewLogger(zapcore.InfoLevel)
	_globalS  = _globalL.Sugar()
)

const skip = 1

func SetLogger(l *zap.Logger) {
	_globalMu.Lock()
	_globalL = l
	_globalS = l.Sugar()
	_globalMu.Unlock()
}

func GetLogger() *zap.Logger {
	_globalMu.RLock()
	l := _globalL
	_globalMu.RUnlock()
	return l.WithOptions(zap.AddCallerSkip(-skip)) // unwrap
}

func getSugarLogger() *zap.SugaredLogger {
	_globalMu.RLock()
	s := _globalS
	_globalMu.RUnlock()
	return s
}

func NewLogger(lv zapcore.Level) *zap.Logger {
	var zapConfig zap.Config

	switch lv {
	case zapcore.DebugLevel:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		zapConfig = zap.NewProductionConfig()
	}

	zapConfig.Level = zap.NewAtomicLevelAt(lv)
	logger, err := zapConfig.Build(
		// AddCallerSkip because we wrapped a layer.
		zap.AddCallerSkip(skip),
	)
	if err != nil {
		log.Panic(err)
	}

	return logger
}
