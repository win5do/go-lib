package logx

import (
	"log"
	"sync/atomic"
	"unsafe"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate sh -c "go run ./generator >zap_sugar_generated.go"

var (
	_globalL = NewLogger(zapcore.InfoLevel)
	_globalS = _globalL.Sugar()
)

const skip = 1

func SetLogger(l *zap.Logger) {
	lPtr := (*unsafe.Pointer)(unsafe.Pointer(&_globalL))
	atomic.StorePointer(lPtr, unsafe.Pointer(l))
	sPtr := (*unsafe.Pointer)(unsafe.Pointer(&_globalS))
	atomic.StorePointer(sPtr, unsafe.Pointer(l.Sugar()))
}

func GetLogger() *zap.Logger {
	l := _globalL
	return l.WithOptions(zap.AddCallerSkip(-skip)) // unwrap
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
