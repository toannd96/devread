package log

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:   "message",
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		EncodeLevel:  CustomLevelEncoder,         //Format cách hiển thị level log
		EncodeTime:   SyslogTimeEncoder,          //Format hiển thị thời điểm log
		EncodeCaller: zapcore.ShortCallerEncoder, //Format dòng code bắt đầu log
	})
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func logErrorWriter() (zapcore.WriteSyncer, error) {
	writerError, err := os.OpenFile("./log_files/"+os.Getenv("APP_NAME")+"_error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&lumberjack.Logger{
			Filename: writerError.Name(),
			MaxSize:  10, // megabytes
			MaxAge:   10, // days
		}),
		zapcore.AddSync(os.Stdout)), nil
}

func logInfoWriter() (zapcore.WriteSyncer, error) {
	writerInfo, err := os.OpenFile("./log_files/"+os.Getenv("APP_NAME")+"_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&lumberjack.Logger{
			Filename: writerInfo.Name(),
			MaxSize:  10,
			MaxAge:   10,
		}),
		zapcore.AddSync(os.Stdout)), nil
}

// Write log to file by level log and console
func WriteLog() *zap.Logger {
	highWriteSyncer, _ := logErrorWriter()
	lowWriteSyncer, _ := logInfoWriter()
	encoder := getEncoder()

	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	lowCore := zapcore.NewCore(encoder, lowWriteSyncer, lowPriority)
	highCore := zapcore.NewCore(encoder, highWriteSyncer, highPriority)

	logger := zap.New(zapcore.NewTee(lowCore, highCore), zap.AddCaller())
	return logger
}
