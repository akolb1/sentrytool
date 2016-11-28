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

// SentryClientAPI is a generic Apache Sentry client interface.
type SentryClientAPI interface {
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
	ListRoleByGroup(group string) ([]string, error)
	// AddGroupsToRole adds specified groups to a role
	//   role - role name
	//   groups - list of group names to add
	AddGroupsToRole(role string, groups []string) error
}

// GetClient returns a Sentry client implementation
//   protocol - legacy or generic
//   host - server host
//   port - server port
//   component - Sentry component for generic protocol
//   user - Sentry authorization user
func GetClient(protocol ProtocolType, host string, port int,
	component string, user string) (SentryClientAPI, error) {
	switch protocol {
	case PolicyProtocol:
		return getHiveClient(host, port, user)
	case GenericPolicyProtocol:
		return getGenericClient(host, port, component, user)
	default:
		return nil, fmt.Errorf("invalid protocol %s", protocol.String())
	}
}
