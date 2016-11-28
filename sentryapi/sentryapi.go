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
	"github.com/akolb1/sentrytool/sentryapi/thrift/sentry_policy_service"
)

// TMPProtocolFactory is a multiplexing protocol factory
type tMPProtocolFactory struct {
}

func (p *tMPProtocolFactory) GetProtocol(t thrift.TTransport) thrift.TProtocol {
	protocol := thrift.NewTBinaryProtocolTransport(t)
	return thrift.NewTMultiplexedProtocol(protocol, sentryProtocol)
}

type sentryClient struct {
	userName  string
	transport thrift.TTransport
	client    *sentry_policy_service.SentryPolicyServiceClient
}

func (c *sentryClient) Close() {
	c.transport.Close()
}

func getHiveClient(host string, port int, user string) (*sentryClient, error) {
	socket, err := thrift.NewTSocket(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	var transport thrift.TTransport = thrift.NewTBufferedTransport(socket, 1024)
	protocolFactory := &tMPProtocolFactory{}
	client := sentry_policy_service.NewSentryPolicyServiceClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		return nil, err
	}
	return &sentryClient{userName: user, transport: transport, client: client}, nil
}

func (c *sentryClient) CreateRole(name string) error {
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

func (c *sentryClient) RemoveRole(name string) error {
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

func (c *sentryClient) ListRoleByGroup(group string) ([]string, error) {
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

func (c *sentryClient) AddGroupsToRole(role string, groups []string) error {
	arg := sentry_policy_service.NewTAlterSentryRoleAddGroupsRequest()
	arg.RequestorUserName = c.userName
	arg.RoleName = role
	groupsMap := make(map[*sentry_policy_service.TSentryGroup]bool)
	for _, group := range groups {
		tGroup := sentry_policy_service.TSentryGroup{GroupName: group}
		groupsMap[&tGroup] = true
	}
	arg.Groups = groupsMap
	result, err := c.client.AlterSentryRoleAddGroups(arg)
	// fmt.Println(result)
	if err != nil {
		return fmt.Errorf("failed to add groups: %s", err)
	}
	if result.GetStatus().GetValue() != 0 {
		return fmt.Errorf("%s", result.GetStatus().Message)
	}

	return nil
}

func (c *sentryClient) RemoveGroupsFromRole(role string, groups []string) error {
	arg := sentry_policy_service.NewTAlterSentryRoleDeleteGroupsRequest()
	arg.RequestorUserName = c.userName
	arg.RoleName = role
	groupsMap := make(map[*sentry_policy_service.TSentryGroup]bool)
	for _, group := range groups {
		tGroup := sentry_policy_service.TSentryGroup{GroupName: group}
		groupsMap[&tGroup] = true
	}
	arg.Groups = groupsMap
	result, err := c.client.AlterSentryRoleDeleteGroups(arg)
	// fmt.Println(result)
	if err != nil {
		return fmt.Errorf("failed to remove groups: %s", err)
	}
	if result.GetStatus().GetValue() != 0 {
		return fmt.Errorf("%s", result.GetStatus().Message)
	}

	return nil
}