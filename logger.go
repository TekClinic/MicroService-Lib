package ms

import (
	"context"

	"github.com/alexlast/bunzap"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() { //nolint:gochecknoinits // required for logger initialization
	// Initialize logger
	logger := zap.Must(zap.NewProduction())
	if HasDebugLogging() {
		logger = zap.Must(zap.NewDevelopment())
	}
	zap.ReplaceGlobals(logger)
}

func getLoggingOptions() []logging.Option {
	return []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithDisableLoggingFields(logging.SystemTag[0]),
	}
}

// getGRPCServerLoggerOptions returns a list of grpc.ServerOption with logging interceptor.
func getGRPCServerLoggerOptions() []grpc.ServerOption {
	opts := getLoggingOptions()
	interceptor := interceptorLogger(zap.L())

	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptor, opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptor, opts...),
		),
	}
}

// getGRPCClientLoggerOptions returns a list of grpc.DialOption with logging interceptor.
func getGRPCClientLoggerOptions() []grpc.DialOption {
	opts := getLoggingOptions()
	interceptor := interceptorLogger(zap.L())

	return []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(interceptor, opts...),
		),
		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(interceptor, opts...),
		),
	}
}

// interceptorLogger adapts zap logger to interceptor logger.
func interceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2) //nolint:gomnd // 2 fields per key-value pair

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			if key == "grpc.start_time" {
				continue
			}
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			logger.DPanic("unknown level", zap.String("level", string(rune(lvl))))
		}
	})
}

// GetDBQueryHook returns a bun.QueryHook that logs queries.
func GetDBQueryHook() bun.QueryHook {
	return bunzap.NewQueryHook(bunzap.QueryHookOptions{
		Logger: zap.L(),
	})
}
