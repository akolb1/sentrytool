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

type ProtocolType int

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

// SentryClientAPI is a generic Sentry client interface
type SentryClientAPI interface {
	Close()
	CreateRole(name string) error
	RemoveRole(name string) error
	ListRoleByGroup(group string) ([]string, error)
}

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
