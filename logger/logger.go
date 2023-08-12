package logger

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"runtime/debug"
)

var zapLog *zap.Logger

func init() {
	var err error
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	//enccoderConfig.StacktraceKey = "" // to hide stacktrace info
	config.EncoderConfig = enccoderConfig

	zapLog, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	grpc_zap.ReplaceGrpcLoggerV2(zapLog)
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}

// FatalfIf log to level error
func FatalfIf(cond bool, fmt string, args ...interface{}) {
	if !cond {
		return
	}
	debug.PrintStack()
	log.Fatalf(fmt, args...)
}

// FatalIfError if err is not nil, then log to level fatal and call os.Exit
func FatalIfError(err error) {
	FatalfIf(err != nil, "fatal error: %v", err)
}
