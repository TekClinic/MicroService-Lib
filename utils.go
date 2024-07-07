package ms

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"go.uber.org/zap"
)

const (
	envRunningMode = "MS_MODE"
	envLogMode     = "MS_LOG_MODE"
	envSecureConn  = "MS_SECURE_CONN"

	productionMode   = "prod"
	developmentMode  = "dev"
	debugLoggingMode = "debug"
)

// IsProduction returns true if the microservice is running in production mode.
func IsProduction() bool {
	return GetOptionalEnv(envRunningMode, developmentMode) == productionMode
}

// HasDebugLogging returns true if the microservice logs are set to debug level.
func HasDebugLogging() bool {
	value, set := os.LookupEnv(envLogMode)
	if set {
		return value == debugLoggingMode
	}
	return !IsProduction()
}

// HasSecureConnection returns true if the microservice is running in secure mode.
func HasSecureConnection() bool {
	secure, err := strconv.ParseBool(GetOptionalEnv(envSecureConn, "false"))
	if err != nil {
		zap.L().DPanic(fmt.Sprintf("%s environment variable is not a boolean", envSecureConn))
	}
	return secure
}

// GetRequiredEnv retrieves the value of the environment variable named by the key.
// If the variable is not present in the environment, an error is returned.
func GetRequiredEnv(key string) (string, error) {
	value, set := os.LookupEnv(key)
	if !set {
		return "", fmt.Errorf("%s environment variable is missing", key)
	}
	return value, nil
}

// GetOptionalEnv retrieves the value of the environment variable named by the key.
// If the variable is not present in the environment, the def value is returned.
func GetOptionalEnv(key string, def string) string {
	value, set := os.LookupEnv(key)
	if set {
		return value
	}
	return def
}

// GetGRPCServerOptions returns a list of options for a GRPC server.
func GetGRPCServerOptions() []grpc.ServerOption {
	return getGRPCServerLoggerOptions()
}

// GetGRPCClientOptions returns a list of options for a GRPC client.
func GetGRPCClientOptions() []grpc.DialOption {
	transportCredentials := insecure.NewCredentials()
	if HasSecureConnection() {
		transportCredentials = credentials.NewTLS(&tls.Config{InsecureSkipVerify: false, MinVersion: tls.VersionTLS12})
	}

	return append(getGRPCClientLoggerOptions(), grpc.WithTransportCredentials(transportCredentials))
}
