## Types

### Service

Defines a basic microservice with `Host` and `Port`.

- **GetAddr() string**
    - Returns the address of the service in the format `Host:Port`.

### BaseServiceServer

Interface for a gRPC server that verifies JWT tokens.

- **VerifyToken(ctx context.Context, rawToken string) (Claims, error)**
    - Verifies a JWT token and returns the token claims if valid.

- **GetPort() string**
    - Returns the port the server is listening on.

### Claims

Interface for handling JWT token claims.

- **HasRole(role string) bool**
    - Checks if the claims contain the specified role.

- **GetRoles() sets.Set[string]**
    - Returns all roles associated with the claims.