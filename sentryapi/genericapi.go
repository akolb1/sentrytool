// Copyright © 2016 Alex Kolbasov
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sentryapi

import (
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/akolb1/sentrytool/sentryapi/thrift/sentry_generic_policy_service"
)

// TMPGenericProtocolFactory is a multiplexing protocol factory
type TMPGenericProtocolFactory struct {
}

func (p *TMPGenericProtocolFactory) GetProtocol(t thrift.TTransport) thrift.TProtocol {
	protocol := thrift.NewTBinaryProtocolTransport(t)
	return thrift.NewTMultiplexedProtocol(protocol, sentryGenericProtocol)
}

type GenericSentryClient struct {
	component string
	userName  string
	transport thrift.TTransport
	client    *sentry_generic_policy_service.SentryGenericPolicyServiceClient
}

func (c *GenericSentryClient) Close() {
	c.transport.Close()
}

func getGenericClient(host string, port int, component string, user string) (*GenericSentryClient, error) {
	socket, err := thrift.NewTSocket(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	var transport thrift.TTransport = thrift.NewTBufferedTransport(socket, 1024)
	protocolFactory := &TMPGenericProtocolFactory{}
	client := sentry_generic_policy_service.NewSentryGenericPolicyServiceClientFactory(transport,
		protocolFactory)
	if err := transport.Open(); err != nil {
		return nil, err
	}
	return &GenericSentryClient{
		userName:  user,
		transport: transport,
		client:    client,
		component: component,
	}, nil
}

func (c *GenericSentryClient) CreateRole(name string) error {
	arg := sentry_generic_policy_service.NewTCreateSentryRoleRequest()
	arg.RequestorUserName = c.userName
	arg.Component = c.component
	arg.RoleName = name
	result, err := c.client.CreateSentryRole(arg)
	if err != nil {
		return fmt.Errorf("failed to create Sentry role %s: %s", name, err)
	}

	if result.GetStatus().Value != 0 {
		return fmt.Errorf("%s\n%s", result.GetStatus().Message,
			*result.GetStatus().Stack)
	}
	return nil
}

func (c *GenericSentryClient) RemoveRole(name string) error {
	arg := sentry_generic_policy_service.NewTDropSentryRoleRequest()
	arg.RequestorUserName = c.userName
	arg.Component = c.component
	arg.RoleName = name
	result, err := c.client.DropSentryRole(arg)
	if err != nil {
		return fmt.Errorf("failed to remove Sentry role %s: %s", name, err)
	}
	if result.GetStatus().Value != 0 {
		return fmt.Errorf("%s\n%s", result.GetStatus().Message,
			*result.GetStatus().Stack)
	}
	return nil
}

func (c *GenericSentryClient) ListRoleByGroup(group string) ([]string, error) {
	arg := sentry_generic_policy_service.NewTListSentryRolesRequest()
	if group == "" {
		arg.GroupName = nil
	} else {
		arg.GroupName = &group
	}

	arg.RequestorUserName = c.userName
	arg.Component = c.component
	result, err := c.client.ListSentryRolesByGroup(arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list Sentry roles: %s", err)
	}

	if result.GetStatus().Value != 0 {
		return nil, fmt.Errorf("%s\n%s", result.GetStatus().Message,
			*result.GetStatus().Stack)
	}
	roles := make([]string, 0, 8)
	for role := range result.Roles {
		roles = append(roles, role.RoleName)
	}
	return roles, nil
}
