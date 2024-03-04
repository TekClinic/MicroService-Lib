package ms

import (
	"fmt"
	"strings"
)

const envMSHostPattern = "MS_%s_HOST"
const envMSPortPattern = "MS_%s_PORT"

// Service defines basic parameters of a microservice for a client.
type Service struct {
	Host string
	Port string
}

// GetAddr provides address of the service
func (s Service) GetAddr() string {
	return s.Host + ":" + s.Port
}

// FetchServiceParameters returns a Service object that describes service named serviceName.
// Host and Port of the Service are retrieved from environment variables MS_{serviceName}_HOST and MS_{serviceName}_PORT
func FetchServiceParameters(serviceName string) (*Service, error) {
	host, err := GetRequiredEnv(fmt.Sprintf(envMSHostPattern, strings.ToUpper(serviceName)))
	if err != nil {
		return nil, err
	}

	port := GetOptionalEnv(fmt.Sprintf(envMSPortPattern, strings.ToUpper(serviceName)), defaultGRPCPort)
	return &Service{Host: host, Port: port}, nil
}
