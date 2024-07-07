package ms

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	sf "github.com/sa-/slicefunk"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	defaultGRPCPort  = "9090"
	envAuthIssuerURL = "AUTH_ISSUER"
	envGRPCPort      = "GRPC_PORT"

	authAuditor = "account"

	rolesSeparator = "."
	rolesKey       = "roles"
)

// BaseServiceServer contains a port on which the server listens for GRPC requests and provides JWT token verification.
type BaseServiceServer interface {
	// VerifyToken verifies the provided token.
	// Returns token claim if it's valid, otherwise returns an error.
	VerifyToken(ctx context.Context, rawToken string) (Claims, error)
	// GetPort returns a port associated with the server
	GetPort() string
}

// baseServiceServer provides a basic implementation of BaseServiceServer.
type baseServiceServer struct {
	port       string
	authIssuer string
	verifier   *oidc.IDTokenVerifier
}

// Claims contains claims from a token and allows roles-checking for authorization.
type Claims interface {
	// HasRole checks if these Claims are authorized with a given role
	HasRole(role string) bool
	// GetRoles returns all roles associated with these claims
	GetRoles() sets.Set[string]
}

// claims provides a basic implementation of Claims.
type claims struct {
	roles sets.Set[string]
}

// CreateBaseServiceServer initiates baseServiceServer with parameters from environment variables.
// GRPC_PORT is used to define Port of the server. By default, 9090.
// AUTH_ISSUER is an url to auth provider.
func CreateBaseServiceServer() (BaseServiceServer, error) {
	authIssuer, err := GetRequiredEnv(envAuthIssuerURL)
	if err != nil {
		return nil, err
	}

	return &baseServiceServer{
		port:       GetOptionalEnv(envGRPCPort, defaultGRPCPort),
		authIssuer: authIssuer,
	}, nil
}

// retrieveRoles translates token-claims to a set of roles.
// There are two types of roles - realm-wise and client-wise.
// Realm-wise roles are associated with a user for any client in the realm. Its syntax is just a <role>
// client-wise roles are associated with a user only for specific client. Its syntax is <client>.<role>.
func retrieveRoles(token *oidc.IDToken) (sets.Set[string], error) {
	var tokenClaims struct {
		ResourceAccess map[string]map[string][]string `json:"resource_access"`
		RealmAccess    map[string][]string            `json:"realm_access"`
		Roles          []string                       `json:"roles"`
	}
	err := token.Claims(&tokenClaims)
	if err != nil {
		return nil, err
	}

	// Claim defined in standards
	roles := sets.New(tokenClaims.Roles...)

	// KeyCloak specific claims
	if realmRoles, exists := tokenClaims.RealmAccess[rolesKey]; exists {
		sets.Insert(roles, realmRoles...)
	}

	// Per KeyCloak Client claims
	for client, clientClaims := range tokenClaims.ResourceAccess {
		if clientRoles, exists := clientClaims[rolesKey]; exists {
			formattedRoles := sf.Map(clientRoles, func(clientRole string) string {
				return client + rolesSeparator + clientRole
			})
			sets.Insert(roles, formattedRoles...)
		}
	}
	return roles, nil
}

// VerifyToken implements BaseServiceServer.VerifyToken using auth provider associated with this baseServiceServer.
// Returns token claim if it's valid, otherwise returns an error.
// May return an error if the provider is misconfigured.
func (server baseServiceServer) VerifyToken(ctx context.Context, rawToken string) (Claims, error) {
	if server.verifier == nil {
		provider, err := oidc.NewProvider(context.Background(), server.authIssuer)
		if err != nil {
			return nil, err
		}
		server.verifier = provider.Verifier(&oidc.Config{
			ClientID: authAuditor,
		})
	}

	idToken, err := server.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	roles, err := retrieveRoles(idToken)
	if err != nil {
		return nil, err
	}

	return &claims{roles: roles}, nil
}

// GetPort implements BaseServiceServer.GetPort.
func (server baseServiceServer) GetPort() string {
	return server.port
}

// HasRole implements Claims.HasRole.
// There are two types of roles - realm-wise and client-wise.
// Realm-wise roles are associated with a user for any client in the realm. Its syntax is just a <role>
// Client-wise roles are associated with a user only for specific client. Its syntax is <client>.<role>.
func (claims claims) HasRole(role string) bool {
	return claims.roles.Has(role)
}

// GetRoles implements Claims.GetRoles.
// There are two types of roles - realm-wise and client-wise.
// Realm-wise roles are associated with a user for any client in the realm. Its syntax is just a <role>
// Client-wise roles are associated with a user only for specific client. Its syntax is <client>.<role>.
func (claims claims) GetRoles() sets.Set[string] {
	return claims.roles.Clone()
}
