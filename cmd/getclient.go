package cmd

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/akolb1/sentrytool/sentryapi"
	"github.com/spf13/viper"
)

// getClient returns Sentry API client, extracting parameters like host and port
// from viper.
//
// If component is specified, it uses Generic sentry protocol, otherwise it uses legacy
// protocol
func getClient() (sentryapi.ClientAPI, error) {
	host := viper.GetString(hostOpt)
	user := viper.GetString(userOpt)
	component := viper.GetString(componentOpt)
	port := viper.GetInt(portOpt)

	var errVal error
	parts := strings.Split(host, ",")
	for _, host := range parts {
		if client, err := getClientForHost(host, port,
			user, component); err == nil {
			return client, nil
		} else {
			errVal = err
		}
	}
	return nil, errVal
}

// getCLientForHost gets a client for a specific host
func getClientForHost(host string, port int, user string,
	component string) (sentryapi.ClientAPI, error) {
	// Allow host/port setup as host:port
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		if len(parts) == 2 {
			// host:port is specified
			host = parts[0]
			if portVal, err := strconv.Atoi(parts[1]); err == nil {
				port = portVal
			} else {
				return nil, fmt.Errorf("invalid port %s", parts[1])
			}
		}
	}

	if component == "" {
		return sentryapi.GetClient(sentryapi.PolicyProtocol,
			host, port, component, user)
	}
	return sentryapi.GetClient(sentryapi.GenericPolicyProtocol,
		host, port, component, user)
}

// isValidRole returns true iff role is valid
// Roles are validated against Sentry database, so validation involves a Thrift call.
func isValidRole(client sentryapi.ClientAPI, roleName string) (bool, error) {
	// Get existing roles
	roles, _, err := client.ListRoleByGroup("")
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role == roleName {
			return true, nil
		}
	}
	return false, nil
}

// toError converts ApiError t regular error if java stack was requested.
// All other errors are returned as is.
func toAPIError(err error) error {
	if !viper.GetBool(jstackOpt) {
		return err
	}
	if apiErr, ok := err.(*sentryapi.APIError); ok {
		stack := apiErr.StackTrace
		if stack == "" {
			return apiErr.Err
		}
		return fmt.Errorf("%v\nServer Stacktrace:\n%s", err, stack)
	}
	return err
}
