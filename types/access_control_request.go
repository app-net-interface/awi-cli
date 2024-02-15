// Copyright (c) 2023 Cisco Systems, Inc. and its affiliates
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package types

type Label map[string]string
type ProtocolPorts map[string][]string

type ACLData struct {
	Protocols []string
	Ports     []string
}

type AccessControlRequest struct {
	Name                string                   `mapstructure:"name"`
	ClusterConnectionID string                   `mapstructure:"cluster_connection_reference"`
	NewConnection       NewConnection            `mapstructure:"new_connection"`
	Source              AccessControlSource      `mapstructure:"source"`
	Destination         AccessControlDestination `mapstructure:"destination"`
}

type NewConnection struct {
	Create bool   `mapstructure:"create"`
	Access string `mapstructure:"access"`
}

type AccessControlSource struct {
	Metadata Metadata `mapstructure:"metadata"`
	Endpoint Endpoint `mapstructure:"endpoints"`
	Subnet   Subnet   `mapstructure:"subnets"`
	Host     Host     `mapstructure:"hosts"`
}

type AccessControlDestination struct {
	AccessControlSource `mapstructure:",squash"`
	URI                 URI                    `mapstructure:"uri"`
	SLA                 RequestedConnectionSLA `mapstructure:"requested_connection_sla"`
}

type Common struct {
	Metadata      Metadata      `mapstructure:"metadata"`
	ProtocolPorts ProtocolPorts `mapstructure:"proto_and_ports"`
	Access        string        `mapstructure:"access"`
}

type Endpoint struct {
	Common `mapstructure:",squash"`
	Labels Label `mapstructure:"labels"`
}

type Subnet struct {
	Endpoint `mapstructure:",squash"`
	Prefixes []string `mapstructure:"prefix"`
}

type Host struct {
	Common `mapstructure:",squash"`
	IPs    []string `mapstructure:"ip"`
	FQDNs  []string `mapstructure:"fqdn"`
}

type URI struct {
	Common `mapstructure:",squash"`
	URIs   []string `mapstructure:"uri"`
}

func (a *AccessControlSource) GetLabels() Label {
	labels := a.Subnet.Labels
	if labels == nil {
		labels = a.Endpoint.Labels
	}
	return labels
}

func (a *AccessControlSource) GetProtocols() ProtocolPorts {
	protocolPorts := a.Subnet.ProtocolPorts
	if protocolPorts == nil {
		protocolPorts = a.Endpoint.ProtocolPorts
	}
	return protocolPorts
}

func (a *AccessControlSource) GetName() string {
	name := a.Subnet.Metadata.Name
	if name == "" {
		name = a.Endpoint.Metadata.Name
	}
	return name
}

func (a *AccessControlSource) GetType() string { // FIXME
	if a.Endpoint.Metadata.Name != "" {
		return "instance"
	}
	return "subnet"
}
