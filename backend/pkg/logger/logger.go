package logger

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Option struct {
	LogPath         string
	ApplicationName string
	Debug           bool
}

func NewOption(conf *viper.Viper) Option {
	return Option{
		LogPath:         conf.GetString("log.log_path"),
		ApplicationName: conf.GetString("log.application_name"),
		Debug:           conf.GetBool("log.debug"),
	}
}

func NewZapLogger(Option Option) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   Option.LogPath,
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
	if Option.Debug {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(writes...), atomicLevel)

	caller := zap.AddCaller()

	callerSkip := zap.AddCallerSkip(1)

	development := zap.Development()

	field := zap.Fields(zap.String("ApplicationName", Option.ApplicationName))

	zapLogger := zap.New(core, caller, callerSkip, development, field)

	zapLogger.Info("log init success")

	return zapLogger
}
