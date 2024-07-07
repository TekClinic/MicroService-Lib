# MicroService-Lib

The `MicroService-Lib` package provides essential utilities for managing microservice configurations,
including environment-based service parameters, logging, gRPC server/client setup, and token verification.

This package simplifies common tasks in microservice architectures such as service discovery,
configuration, and secure communication.

## Configurable Environment Variables

This application uses several environment variables to configure various aspects of 
the microservices, authentication, and logging. 

Below is a comprehensive list of these environment variables and their descriptions.

### Client-Side Microservice Configuration

These environment variables are used to discover and connect to other microservices.

- **`MS_{serviceName}_HOST`**
    - **Description:** The hostname of the microservice specified by `{serviceName}`.
    - **Usage:** Defines the host for a given microservice.
    - **Example:** For a service named `auth`, the variable would be `MS_AUTH_HOST`.

- **`MS_{serviceName}_PORT`**
    - **Description:** The port number on which the microservice specified by `{serviceName}` listens for GRPC requests.
    - **Default:** `"9090"`
    - **Usage:** Defines the port for a given microservice.
    - **Example:** For a service named `auth`, the variable would be `MS_AUTH_PORT`.

- **`MS_SECURE_CONN`**
  - **Description:** Indicates whether the microservice should use a secure connection to other services.
  - **Options:** `"true"` (use secure connection) or `"false"` (use insecure connection).
  - **Default:** `"false"`
  - **Usage:** Configures the transport security for GRPC connections.

### Server-Side Configuration

- **`GRPC_PORT`**
    - **Description:** The port on which the GRPC server listens to.
    - **Default:** `"9090"`
    - **Usage:** Configures the listening port of the GRPC server.

- **`AUTH_ISSUER`**
    - **Description:** The URL of the authentication provider.
    - **Usage:** Specifies the issuer for JWT token verification.
    - **Required:** Yes

### Logging Configuration

- **`MS_LOG_MODE`**
    - **Description:** Sets the logging mode for the microservice.
    - **Options:** `"debug"` (enables debug logging) or any other value (disables debug logging).
    - **Default:** `"dev"` mode if not specified.
    - **Usage:** Controls the verbosity of the logs.

### Deployment Mode Configuration

- **`MS_MODE`**
    - **Description:** Specifies the running mode of the microservice.
    - **Options:** `"dev"` (development mode) or `"prod"` (production mode).
    - **Default:** `"dev"`
    - **Usage:** Determines the mode in which the microservice operates, affecting logging and other settings.
