package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

type customEncoder struct {
	zapcore.Encoder
}

func (c *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 格式化日志条目，包括时间、日志级别、调用者信息和消息
	formatted := entry.Time.Format("2006/01/02 - 15:04:05") + " | " +
		entry.Level.CapitalString() + " | " +
		entry.Caller.TrimmedPath() + " | " +
		entry.Message

	// 使用编码器对日志条目和字段进行编码
	buf, err := c.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		// 如果编码失败，则返回错误
		return nil, err
	}

	// 重置缓冲区
	buf.Reset()
	// 将格式化后的日志条目追加到缓冲区中
	buf.AppendString(formatted)
	return buf, nil
}

func newCustomEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	// 创建一个新的自定义编码器
	return &customEncoder{
		// 使用 zapcore.NewConsoleEncoder 函数创建一个控制台编码器，并传入配置
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}

func init() {
	// 创建一个生产环境的编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置时间格式为ISO8601
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 设置日志级别编码为大写形式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 设置调用者信息编码为简短形式
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 创建一个新的核心对象
	core := zapcore.NewCore(
		// 使用自定义编码器
		newCustomEncoder(encoderConfig),
		// 将输出锁定到标准输出
		zapcore.Lock(os.Stdout),
		// 设置日志级别为Debug
		zap.DebugLevel,
	)

	// 创建一个新的日志记录器，并添加调用者信息和跳过一层调用栈
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

// Info logs a message at InfoLevel. The message includes any fields passed.
func Info(args ...interface{}) {
	// 调用 logger 的 Info 方法，将 args 参数传递给该方法
	logger.Info(args...)
}

// Infof formats and logs a message at InfoLevel.
func Infof(template string, args ...interface{}) {
	// 调用logger的Infof方法，将模板和参数传递给它
	logger.Infof(template, args...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed.
func Error(args ...interface{}) {
	// 使用logger的Error方法记录错误日志
	logger.Error(args...)
}

// Errorf formats and logs a message at ErrorLevel.
func Errorf(template string, args ...interface{}) {
	// 调用 logger 的 Errorf 方法，传入模板字符串和参数列表
	logger.Errorf(template, args...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed.
func Debug(args ...interface{}) {
	// 调用 logger 的 Debug 方法，并传递 args 参数
	logger.Debug(args...)
}

// Debugf formats and logs a message at DebugLevel.
func Debugf(template string, args ...interface{}) {
	// 使用logger的Debugf方法输出调试信息
	// logger.Debugf(template, args...)
	logger.Debugf(template, args...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed.
func Warn(args ...interface{}) {
	// 使用logger记录警告信息
	logger.Warn(args...)
}

// Warnf formats and logs a message at WarnLevel.
func Warnf(template string, args ...interface{}) {
	// 调用logger的Warnf方法，传入模板字符串和参数列表
	// logger.Warnf(template, args...)
	logger.Warnf(template, args...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed.
func Panic(args ...interface{}) {
	// 调用logger的Panic方法，传入args参数
	logger.Panic(args...)
}

// Panicf formats and logs a message at PanicLevel.
func Panicf(template string, args ...interface{}) {
	// 调用logger的Panicf方法，传入模板字符串和参数列表
	logger.Panicf(template, args...)
}
