package bootstrap

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type loggerOption struct {
	LogPath         string
	ApplicationName string
	Debug           bool
}

func NewLoggerOption(conf *viper.Viper) loggerOption {
	return loggerOption{
		LogPath:         conf.GetString("log.log_path"),
		ApplicationName: conf.GetString("log.application_name"),
		Debug:           conf.GetBool("log.debug"),
	}
}

func NewLogger(option loggerOption) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   option.LogPath,
		MaxSize:    128,
		MaxAge:     7,
		MaxBackups: 30,
		Compress:   false,
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "file",
		StacktraceKey: "stacktrace",

		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	if option.Debug {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(writes...), atomicLevel)

	caller := zap.AddCaller()

	callerSkip := zap.AddCallerSkip(1)

	development := zap.Development()

	field := zap.Fields(zap.String("ApplicationName", option.ApplicationName))

	zapLogger := zap.New(core, caller, callerSkip, development, field)

	zapLogger.Info("log init success")

	return zapLogger
}
