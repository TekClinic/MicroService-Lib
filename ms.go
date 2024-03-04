package ms

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"os"
	"strings"
)

const rolesSeparator = "."
const rolesKey = "roles"

// Service defines basic parameters of a microservice for a client.
type Service struct {
	Host string
	Port string
}

// BaseServiceServer contains a port on which the server listens for GRPC requests and provides JWT token verification.
type BaseServiceServer struct {
	Port       string
	authIssuer string
	verifier   *oidc.IDTokenVerifier
}

// Claims contains claims from a token and allows roles-checking for authorization.
type Claims struct {
	PreferredUsername string                         `json:"preferred_username"`
	ResourceAccess    map[string]map[string][]string `json:"resource_access"`
	RealmAccess       map[string][]string            `json:"realm_access"`
}

// GetRequiredEnv retrieves the value of the environment variable named by the key.
// If the variable is not present in the environment an error is returned.
func GetRequiredEnv(key string) (string, error) {
	value, set := os.LookupEnv(key)
	if !set {
		return "", errors.New(key + " environment variable is missing")
	}
	return value, nil
}

// GetOptionalEnv retrieves the value of the environment variable named by the key.
// If the variable is not present in the environment the def values is returned.
func GetOptionalEnv(key string, def string) string {
	value, set := os.LookupEnv(key)
	if set {
		return value
	}
	return def
}

// FetchServiceParameters returns a Service object that describes service named serviceName.
// Host and Port of the Service are retrieved from environment variables MS_{serviceName}_HOST and MS_{serviceName}_PORT
func FetchServiceParameters(serviceName string) (*Service, error) {
	host, err := GetRequiredEnv(fmt.Sprintf("MS_%s_HOST", strings.ToUpper(serviceName)))
	if err != nil {
		return nil, err
	}

	port := GetOptionalEnv(fmt.Sprintf("MS_%s_PORT", strings.ToUpper(serviceName)), "9090")
	return &Service{Host: host, Port: port}, nil
}

// CreateBaseServiceServer initiates BaseServiceServer with parameters from environment variables.
// GRPC_PORT is used to define Port of the server. By default, 9090.
// AUTH_ISSUER is an url to auth provider
func CreateBaseServiceServer() (*BaseServiceServer, error) {
	authIssuer, err := GetRequiredEnv("AUTH_ISSUER")
	if err != nil {
		return nil, err
	}

	return &BaseServiceServer{
		Port:       GetOptionalEnv("GRPC_PORT", "9090"),
		authIssuer: authIssuer,
	}, nil
}

// VerifyToken verifies the provided token via auth provider associated with this BaseServiceServer.
// Returns token claim if it's valid, otherwise returns an error.
// May return an error if provider is misconfigured.
func (server BaseServiceServer) VerifyToken(ctx context.Context, rawToken string) (*Claims, error) {
	if server.verifier == nil {
		provider, err := oidc.NewProvider(context.Background(), server.authIssuer)
		if err != nil {
			return nil, err
		}
		server.verifier = provider.Verifier(&oidc.Config{
			ClientID: "account",
		})
	}

	idToken, err := server.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	var claims Claims
	err = idToken.Claims(&claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

// HasRole checks if these Claims are authorized with a given role.
// There are two types of roles - realm-wise and client-wise.
// realm-wise roles are associated with a user for any client in the realm. Its syntax is just a <role>
// client-wise roles are associated with a user only for specific client. Its syntax is <client>.<role>
func (claims Claims) HasRole(role string) bool {
	parts := strings.SplitN(role, rolesSeparator, 2)

	var roles []string
	var targetRole string
	var exists bool

	if len(parts) == 1 {
		roles, exists = claims.RealmAccess[rolesKey]
		targetRole = role
	} else {
		clientClaims := claims.ResourceAccess[parts[0]]
		roles, exists = clientClaims[rolesKey]
		targetRole = parts[1]
	}

	if !exists {
		return false
	}

	for _, r := range roles {
		if r == targetRole {
			return true
		}
	}
	return false
}
