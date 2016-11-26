package sentryapi

import (
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/akolb1/sentrytool/sentryapi/thrift/sentry_policy_service"
)

const (
	sentry_protocol = "SentryPolicyService"
)

// TMPProtocolFactory is a multiplexing protocol factory
type TMPProtocolFactory struct {
}

func (p *TMPProtocolFactory) GetProtocol(t thrift.TTransport) thrift.TProtocol {
	protocol := thrift.NewTBinaryProtocolTransport(t)
	return thrift.NewTMultiplexedProtocol(protocol, sentry_protocol)
}

type SentryClient struct {
	userName  string
	transport thrift.TTransport
	client    *sentry_policy_service.SentryPolicyServiceClient
}

func (c *SentryClient) Close() {
	c.transport.Close()
}

func GetClient(host string, port int, user string) (*SentryClient, error) {
	socket, err := thrift.NewTSocket(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	var transport thrift.TTransport = thrift.NewTBufferedTransport(socket, 1024)
	protocolFactory := &TMPProtocolFactory{}
	client := sentry_policy_service.NewSentryPolicyServiceClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		return nil, err
	}
	return &SentryClient{userName: user, transport: transport, client: client}, nil
}

func (c *SentryClient) CreateRole(name string) error {
	arg := sentry_policy_service.NewTCreateSentryRoleRequest()
	arg.RequestorUserName = c.userName
	arg.RoleName = name
	result, err := c.client.CreateSentryRole(arg)
	if err != nil {
		return fmt.Errorf("failed to create Sentry role %s: %s", name, err)
	}

	if result.GetStatus().GetValue() != 0 {
		return fmt.Errorf("%s", result.GetStatus().Message)
	}
	return nil
}

func (c *SentryClient) RemoveRole(name string) error {
	arg := sentry_policy_service.NewTDropSentryRoleRequest()
	arg.RequestorUserName = c.userName
	arg.RoleName = name
	result, err := c.client.DropSentryRole(arg)
	if err != nil {
		return fmt.Errorf("failed to remove Sentry role %s: %s", name, err)
	}
	if result.GetStatus().GetValue() != 0 {
		return fmt.Errorf("%s", result.GetStatus().Message)
	}
	return nil
}

func (c *SentryClient) ListRoleByGroup(group string) ([]string, error) {
	arg := sentry_policy_service.NewTListSentryRolesRequest()
	if group == "" {
		arg.GroupName = nil
	} else {
		arg.GroupName = &group
	}

	arg.RequestorUserName = c.userName
	result, err := c.client.ListSentryRolesByGroup(arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list Sentry roles: %s", err)
	}

	if result.GetStatus().GetValue() != 0 {
		return nil, fmt.Errorf("%s", result.GetStatus().Message)
	}
	roles := make([]string, 0, 8)
	for role := range result.Roles {
		roles = append(roles, role.RoleName)
	}
	return roles, nil
}
