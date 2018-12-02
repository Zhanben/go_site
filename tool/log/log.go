package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger //global logger

func InitLog() {
	//logLevel := viper.GetString("Log.LogLevel")
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

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		w,
		zap.InfoLevel,
	)
	Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger)
}
