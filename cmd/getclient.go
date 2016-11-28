package cmd

import (
	"github.com/akolb1/sentrytool/sentryapi"
	"github.com/spf13/viper"
)

// getClient returns Sentry API client, extracting parameters like host and port
// from viper.
//
// If component is specified, it uses Generic sentry protocol, otherwise it uses legacy
// protocol
func getClient() (sentryapi.SentryClientAPI, error) {
	host := viper.Get(hostOpt).(string)
	port := viper.Get(portOpt).(int)
	user := viper.Get(userOpt).(string)
	component := viper.Get(componentOpt).(string)

	if component == "" {
		return sentryapi.GetClient(sentryapi.PolicyProtocol,
			host, port, component, user)
	}
	return sentryapi.GetClient(sentryapi.GenericPolicyProtocol,
		host, port, component, user)
}
