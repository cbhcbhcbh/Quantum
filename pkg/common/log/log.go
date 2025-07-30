package log

import (
	"io"
	"os"

	"github.com/cbhcbhcbh/Quantum/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HttpLog struct {
	*zap.Logger
}

type GrpcLog struct {
	*zap.Logger
}

func init() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
}

func NewHttpLog(config *config.Config) (HttpLog, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	logger := zap.New(core).With(zap.String("proto", "http"))

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Writer(os.Stderr)

	return HttpLog{logger}, nil
}

func NewGrpcLog(config *config.Config) (GrpcLog, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.ErrorLevel,
	)
	logger := zap.New(core).With(zap.String("proto", "grpc"))

	return GrpcLog{logger}, nil
}
