## Functions

### Environment Variables

- **GetRequiredEnv(key string) (string, error)**
    - Retrieves the value of the required environment variable. Returns an error if the variable is not set.

- **GetOptionalEnv(key string, def string) string**
    - Retrieves the value of the optional environment variable. Returns the default value if the variable is not set.

### Microservice Configuration

- **FetchServiceParameters(serviceName string) (\*Service, error)**
    - Returns a `Service` object with `Host` and `Port` values derived from environment
      variables `MS_{serviceName}_HOST` and `MS_{serviceName}_PORT`.

### gRPC Server/Client Options

- **GetGRPCServerOptions() []grpc.ServerOption**
    - Returns options for configuring a gRPC server with logging interceptors.

- **GetGRPCClientOptions() []grpc.DialOption**
    - Returns options for configuring a gRPC client with logging interceptors and transport credentials based on
      environment settings.

- **GetGRPCServerLoggerOptions() []grpc.ServerOption**
    - Returns gRPC server options with logging interceptors.

- **GetGRPCClientLoggerOptions() []grpc.DialOption**
    - Returns gRPC client options with logging interceptors.

### Database Logging

- **GetDBQueryHook() bun.QueryHook**
    - Returns a query hook for logging bun database queries using zap.

### Server Configuration

- **CreateBaseServiceServer() (BaseServiceServer, error)**
    - Creates a `BaseServiceServer` with parameters from environment variables, including `GRPC_PORT` and `AUTH_ISSUER`.

### Token Verification

- **(server baseServiceServer) VerifyToken(ctx context.Context, rawToken string) (Claims, error)**
    - Verifies a JWT token using the associated auth provider. Returns token claims or an error if verification fails.

### Claims

- **(claims claims) HasRole(role string) bool**
    - Checks if the claims contain the specified role.

- **(claims claims) GetRoles() sets.Set[string]**
    - Returns all roles associated with the claims.

### Runtime Environment

- **IsProduction() bool**
    - Returns `true` if the service is running in production mode, determined by the `MS_MODE` environment variable.

- **HasDebugLogging() bool**
    - Returns `true` if the service logging level is set to debug, determined by the `MS_LOG_MODE` environment variable.

- **HasSecureConnection() bool**
    - Returns `true` if the service is configured for secure connections, determined by the `MS_SECURE_CONN` environment
      variable.