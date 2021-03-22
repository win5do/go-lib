// Code generated by log-gen. DO NOT EDIT.
package logx

func Desugar() *Logger {
	return _globalS.Desugar()
}

func Named(name string) *SugaredLogger {
	return _globalS.Named(name)
}

func With(args ...interface{}) *SugaredLogger {
	return _globalS.With(args...)
}

func Debug(args ...interface{}) {
	_globalS.Debug(args...)
}

func Info(args ...interface{}) {
	_globalS.Info(args...)
}

func Warn(args ...interface{}) {
	_globalS.Warn(args...)
}

func Error(args ...interface{}) {
	_globalS.Error(args...)
}

func DPanic(args ...interface{}) {
	_globalS.DPanic(args...)
}

func Panic(args ...interface{}) {
	_globalS.Panic(args...)
}

func Fatal(args ...interface{}) {
	_globalS.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	_globalS.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	_globalS.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	_globalS.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	_globalS.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	_globalS.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	_globalS.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	_globalS.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	_globalS.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	_globalS.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	_globalS.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	_globalS.Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	_globalS.DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	_globalS.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	_globalS.Fatalw(msg, keysAndValues...)
}

func Sync() error {
	return _globalS.Sync()
}
