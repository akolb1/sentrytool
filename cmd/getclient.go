package cmd

import (
	"fmt"

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
