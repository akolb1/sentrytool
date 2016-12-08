package cmd

import (
	"github.com/akolb1/sentrytool/sentryapi"
	"github.com/spf13/viper"
	"strconv"
	"fmt"
)

// getClient returns Sentry API client, extracting parameters like host and port
// from viper.
//
// If component is specified, it uses Generic sentry protocol, otherwise it uses legacy
// protocol
func getClient() (sentryapi.ClientAPI, error) {
	host := viper.Get(hostOpt).(string)
	port := viper.Get(portOpt).(string)
	user := viper.Get(userOpt).(string)
	component := viper.Get(componentOpt).(string)

	portVal, err := strconv.Atoi(port)
	if (err != nil) {
		return nil, fmt.Errorf("invalid port: %v", err)
	}

	if component == "" {
		return sentryapi.GetClient(sentryapi.PolicyProtocol,
			host, portVal, component, user)
	}
	return sentryapi.GetClient(sentryapi.GenericPolicyProtocol,
		host, portVal, component, user)
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
