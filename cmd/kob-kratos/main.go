package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"kob-kratos/internal/conf"

	"github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/natefinch/lumberjack"
	"google.golang.org/protobuf/encoding/protojson"

	_ "go.uber.org/automaxprocs"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")

	// 增加这段代码
	json.MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true, // 默认值不忽略
		UseProtoNames:   true, // 使用proto name返回http字段
	}
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func newLogger(c *conf.Log) (log.Logger, error) {
	var (
		logger log.Logger
		core   zapcore.Core
		err    error
	)
	switch c.GetMode() {
	case "file":
		core, err = getFileZapCore(c)
		if err != nil {
			return nil, err
		}
	case "console":
		core, err = getConsoleZapCore(c)
		if err != nil {
			return nil, err
		}
	}

	logger = zap.NewLogger(uberzap.New(core, uberzap.AddCaller(), uberzap.AddCallerSkip(2)))

	return logger, nil
}

// getZapLevel 转换日志级别
func getZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func getFileZapCore(c *conf.Log) (zapcore.Core, error) {
	// 1. 创建目录
	if err := os.MkdirAll(c.GetPath(), 0o755); err != nil {
		return nil, err
	}

	// 2. 解析级别文件映射（示例配置）
	levelMap := map[string]string{
		"info":  "app.log",   // info级别日志
		"error": "err.log",   // error级别日志
		"debug": "debug.log", // 新增debug级别日志
	}

	encoderConfig := uberzap.NewProductionEncoderConfig()
	omitZapEncoderConfig(&encoderConfig)
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 3. 动态创建核心
	var cores []zapcore.Core

	for levelStr, filename := range levelMap {
		// 获取当前级别
		level := getZapLevel(levelStr)

		// 创建对应文件核心
		core := newFileCore(
			filepath.Join(c.GetPath(), filename),
			encoder,
			level, // 该文件的最小记录级别
			c.MaxSize,
			c.KeepDays,
			c.MaxBackups,
			c.Compress,
		)
		cores = append(cores, core)
	}

	return zapcore.NewTee(cores...), nil
}

// 通用文件核心创建
func newFileCore(
	filename string,
	encoder zapcore.Encoder,
	minLevel zapcore.Level,
	maxSize, keepDays, maxBackups int32,
	compress bool,
) zapcore.Core {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    int(maxSize),    // 日志文件最大大小(MB)
		MaxBackups: int(maxBackups), // 保留旧日志文件数量
		MaxAge:     int(keepDays),   // 保留天数
		Compress:   compress,        // 是否压缩旧日志
	}

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(lumberjackLogger),
		uberzap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == minLevel
		}),
	)
}

func getConsoleZapCore(c *conf.Log) (zapcore.Core, error) {
	encoderConfig := uberzap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	omitZapEncoderConfig(&encoderConfig)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		uberzap.NewAtomicLevelAt(getZapLevel(c.GetLevel())),
	), nil
}

func omitZapEncoderConfig(conf *zapcore.EncoderConfig) {
	// conf.TimeKey = zapcore.OmitKey
	// conf.LevelKey = zapcore.OmitKey
	// conf.NameKey = zapcore.OmitKey
	conf.CallerKey = zapcore.OmitKey
	conf.FunctionKey = zapcore.OmitKey
	// conf.MessageKey = zapcore.OmitKey
	conf.StacktraceKey = zapcore.OmitKey
}
