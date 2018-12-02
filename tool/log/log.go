package log

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger //global logger

func InitLog() {
	logLevel := viper.GetString("Log.LogLevel")
	//logPath := viper.GetString("Log.LogPath")
	logName := viper.GetString("Log.LogName")
	logAge := viper.GetInt("Log.LogAge")
	logSize := viper.GetInt("Log.LogSize")

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logName,
		MaxSize:    logSize, // megabytes
		MaxBackups: 3,
		LocalTime:  true,
		MaxAge:     logAge, // days
	})
	zapLogLevel := zap.NewAtomicLevel()
	if err := zapLogLevel.UnmarshalText([]byte(strings.ToLower(logLevel))); err != nil {
		panic(fmt.Errorf("get config log level:%v config error: %v", logLevel, err))
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		w,
		zapLogLevel,
	)
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
	Logger.Info("logger init successful!")
}
