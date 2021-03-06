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
// Attributes:
//   Name - Role name
//   Groups - list of groups for the role
type Role struct {
	Name   string
	Groups []string
}

// Privilege is the Sentry privilege representation. It comboines
// Generic model and legacy Hive model
type Privilege struct {
	Scope            string
	Server           string
	Database         string
	Table            string
	Column           string
	URI              string
	Action           string
	Service          string
	GrantOption      bool
	UnsetGrantOption bool // True is grant option is unset
}

// ClientAPI is a generic Apache Sentry client interface.
// The API is the same for the Generic or Hive model.
type ClientAPI interface {
	// Close closes the client connection
	Close()
	// CreateRole creates a role with given name
	//   roleName - role name
	CreateRole(roleName string) error
	// RemoveRole removes role with given name
	//   roleName - role name
	RemoveRole(roleName string) error
	// ListRoleByGroup returns list of role names for a given group or all
	// roles if group is nil
	//   groupName - group name
	ListRoleByGroup(groupName string) ([]string, []*Role, error)
	// AddGroupsToRole adds specified groups to the role
	//   roleName - role name
	//   groups - list of group names to add
	AddGroupsToRole(roleName string, groups []string) error
	// RemoveGroupsFromRole removes specified groups from the role
	//   roleName - role name
	//   groups - list of group names to remove
	RemoveGroupsFromRole(roleName string, groups []string) error
	// GrantPrivilege grants privilege to the role
	//   roleName - role name
	//   priv - privilege to grant
	GrantPrivilege(roleName string, priv *Privilege) error
	// RevokePrivilege revokes privilege from the role
	//  roleName - role name
	//  priv - privilege to revoke
	RevokePrivilege(roleName string, priv *Privilege) error
	// ListPrivilegesByRole returns a list of privileges for the role
	// role.
	// If template is not NULL, only return privileges matching template
	ListPrivilegesByRole(roleName string, template *Privilege) ([]*Privilege, error)
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

// APIError is an extension to error type that also contains source stack trace.
type APIError struct {
	Err        error
	StackTrace string
}

func (err *APIError) Error() string {
	return err.Err.Error()
}

// newApiError returns an initialized instance of ApiError.
func newAPIError(err error, stackP *string) *APIError {
	var stack string
	if stackP != nil {
		stack = *stackP
	}
	return &APIError{
		Err:        err,
		StackTrace: stack,
	}
}
