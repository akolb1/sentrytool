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

import "fmt"

// ProtocolType is enum describing available Apache Sentry protocols. Currently Sentry supports
// two protocols: old protocol and generic protocol.
type ProtocolType int

// PolicyProtocol is the legacy Sentry protocol
// GenericPolicyProtocol is the generic Sentry protocol
const (
	PolicyProtocol ProtocolType = iota
	GenericPolicyProtocol
)

const (
	sentryProtocol        = "SentryPolicyService"
	sentryGenericProtocol = "SentryGenericPolicyService"
)

func (pt ProtocolType) String() string {
	if pt == PolicyProtocol {
		return sentryProtocol
	}
	if pt == GenericPolicyProtocol {
		return sentryGenericProtocol
	}
	return "unknownProtocol"
}

// Role is a representation of Sentry role. Each role has a name and a
// list of groups associated with the role.
type Role struct {
	Name   string
	Groups []string
}

// Privilege is the Sentry privilege representation. It comboines
// Generic model and legacy Hive model
type Privilege struct {
	Scope       string
	Server      string
	Database    string
	Table       string
	Column      string
	URI         string
	Action      string
	Service     string
	GrantOption bool
}

// SentryClientAPI is a generic Apache Sentry client interface.
type ClientAPI interface {
	// Close closes the client connection
	Close()
	// CreateRole creates a role with given name
	//   name - role name
	CreateRole(name string) error
	// RemoveRole removes role with given name
	//   name - role name
	RemoveRole(name string) error
	// ListRoleByGroup returns list of role names for a given group or all
	// roles if group is nil
	//   group - group name
	ListRoleByGroup(group string) ([]string, []*Role, error)
	// AddGroupsToRole adds specified groups to the role
	//   role - role name
	//   groups - list of group names to add
	AddGroupsToRole(role string, groups []string) error
	// RemoveGroupsFromRole removes specified groups from the role
	//   role - role name
	//   groups - list of group names to remove
	RemoveGroupsFromRole(role string, groups []string) error
	// GrantPrivilege grants privilege to the role
	GrantPrivilege(role string, priv *Privilege) error
	// RevokePrivilege revokes privilege from the role
	RevokePrivilege(role string, priv *Privilege) error
}

// GetClient returns a Sentry client implementation
//   protocol - legacy or generic
//   host - server host
//   port - server port
//   component - Sentry component for generic protocol
//   user - Sentry authorization user
func GetClient(protocol ProtocolType, host string, port int,
	component string, user string) (ClientAPI, error) {
	switch protocol {
	case PolicyProtocol:
		return getHiveClient(host, port, user)
	case GenericPolicyProtocol:
		return getGenericClient(host, port, component, user)
	default:
		return nil, fmt.Errorf("invalid protocol %s", protocol.String())
	}
}
