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
type tMPGenericProtocolFactory struct {
}

func (p *tMPGenericProtocolFactory) GetProtocol(t thrift.TTransport) thrift.TProtocol {
	protocol := thrift.NewTBinaryProtocolTransport(t)
	return thrift.NewTMultiplexedProtocol(protocol, sentryGenericProtocol)
}

type genericSentryClient struct {
	component string
	userName  string
	transport thrift.TTransport
	client    *sentry_generic_policy_service.SentryGenericPolicyServiceClient
}

func (c *genericSentryClient) Close() {
	c.transport.Close()
}

func getGenericClient(host string, port int, component string, user string) (*genericSentryClient, error) {
	socket, err := thrift.NewTSocket(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	var transport thrift.TTransport = thrift.NewTBufferedTransport(socket, 1024)
	protocolFactory := &tMPGenericProtocolFactory{}
	client := sentry_generic_policy_service.NewSentryGenericPolicyServiceClientFactory(transport,
		protocolFactory)
	if err := transport.Open(); err != nil {
		return nil, err
	}
	return &genericSentryClient{
		userName:  user,
		transport: transport,
		client:    client,
		component: component,
	}, nil
}

func (c *genericSentryClient) CreateRole(name string) error {
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

func (c *genericSentryClient) RemoveRole(name string) error {
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

func (c *genericSentryClient) ListRoleByGroup(group string) ([]string,
	[]*Role, error) {
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
		return nil, nil, fmt.Errorf("failed to list Sentry roles: %s", err)
	}

	if result.GetStatus().Value != 0 {
		return nil, nil, fmt.Errorf("%s\n%s", result.GetStatus().Message,
			*result.GetStatus().Stack)
	}
	roleNames := make([]string, 0, 8)
	roles := make([]*Role, 0, 8)

	// Collect results
	for role := range result.Roles {
		roleNames = append(roleNames, role.RoleName)
		groupMap := role.Groups
		groups := []string{}
		// Get list of groups
		for group := range groupMap {
			groups = append(groups, group)
		}
		roles = append(roles,
			&Role{Name: role.RoleName, Groups: groups})

	}
	return roleNames, roles, nil
}

func (c *genericSentryClient) AddGroupsToRole(role string, groups []string) error {
	arg := sentry_generic_policy_service.NewTAlterSentryRoleAddGroupsRequest()
	arg.RequestorUserName = c.userName
	arg.RoleName = role
	arg.Component = c.component

	groupsMap := make(map[string]bool)
	for _, group := range groups {
		groupsMap[group] = true
	}
	arg.Groups = groupsMap
	result, err := c.client.AlterSentryRoleAddGroups(arg)
	if err != nil {
		return fmt.Errorf("failed to add groups: %s", err)
	}
	if result.GetStatus().Value != 0 {
		return fmt.Errorf("%s", result.GetStatus().Message)
	}

	return nil
}

func (c *genericSentryClient) RemoveGroupsFromRole(role string, groups []string) error {
	arg := sentry_generic_policy_service.NewTAlterSentryRoleDeleteGroupsRequest()
	arg.RequestorUserName = c.userName
	arg.RoleName = role
	arg.Component = c.component

	groupsMap := make(map[string]bool)
	for _, group := range groups {
		groupsMap[group] = true
	}
	arg.Groups = groupsMap
	result, err := c.client.AlterSentryRoleDeleteGroups(arg)
	if err != nil {
		return fmt.Errorf("failed to remove groups: %s", err)
	}
	if result.GetStatus().Value != 0 {
		return fmt.Errorf("%s", result.GetStatus().Message)
	}

	return nil
}
