package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger
var WarnColor = "\033[1;31m%s\033[0m\n"

func InitCfg() {

	/*
			viper.SetConfigName("config")
			viper.SetConfigType("toml")
			viper.AddConfigPath(".")
			err := viper.ReadInConfig()

		LogInit(viper.GetString("log-level"), viper.GetString("log-name")) //"log-path"
		if err != nil {
			//Logger.Sugar().Error(err.Error())
			Logger.Errorf(err.Error())
		}
	*/
	//0504---Logger.Infof("bsc balance start")

	//Logger.Sugar().Info("bsc balance start")
	//Logger.Sugar().Info("bsc balance start")

}

// func LogInit() {
func LogInit(LogLevel string, logPath string) {

	// writeSyncer := getLogWriter()
	writeSyncer := getLogWriter(logPath)

	encoder := getEncoder()
	colorEncoder := getColorEncoder()
	//level := config.Conf.Service.LogLevel
	level := LogLevel

	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	StdoutCore := zapcore.NewCore(colorEncoder, zapcore.Lock(os.Stdout), l)
	teeCore := zapcore.NewTee(
		core,
		StdoutCore,
	)
	Logger = zap.New(teeCore, zap.AddCaller()).Sugar()

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getColorEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(LogPath string) zapcore.WriteSyncer {
	// logPath := os.Getenv("LogPath")
	// if logPath == "" {
	// 	fmt.Printf(WarnColor, "logpath is empty,use log path in ./update.log")
	// 	logPath = "./update.log"
	// }
	//20230407sgj
	//logPath := config.Conf.Service.LogPath
	logPath := LogPath
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     24,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}
func Info(args ...interface{}) {
	Logger.Info(args...)
}
func Error(args ...interface{}) {
	Logger.Error(args...)
}
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}
