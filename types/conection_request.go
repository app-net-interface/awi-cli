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

import "strings"

const (
	ConditionLabel = "condition"
	OrCondition    = "OR"
	AndCondition   = "AND"
)

type ConnectionRequest struct {
	ID           string                         `mapstructure:"id"`
	Name         string                         `mapstructure:"name"`
	Type         string                         `mapstructure:"type"`
	Source       ConnectionRequestSource        `mapstructure:"source"`
	Destinations []ConnectionRequestDestination `mapstructure:"destinations"`
}

type Metadata struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
}

type ConnectionRequestSource struct {
	ID                string   `mapstructure:"id"`
	Metadata          Metadata `mapstructure:"metadata"`
	SiteID            string   `mapstructure:"site_id"`
	FeatureTemplateID string   `mapstructure:"feature_template_id"`
	Type              string   `mapstructure:"type"`
	Provider          string   `mapstructure:"provider"`
	DefaultAccess     string   `mapstructure:"default_access"`
	Networks          []string `mapstructure:"networks"`
}

type ConnectionRequestDestination struct {
	DbID                   string                 `mapstructure:"db_id"`
	Metadata               Metadata               `mapstructure:"metadata"`
	SiteID                 string                 `mapstructure:"site_id"`
	ID                     string                 `mapstructure:"id"`
	Type                   string                 `mapstructure:"type"`
	Provider               string                 `mapstructure:"provider"`
	RequestedConnectionSLA RequestedConnectionSLA `mapstructure:"requested_connection_sla"`
}

type RequestedConnectionSLA struct {
	Type      string `mapstructure:"type"`
	Bandwidth int    `mapstructure:"bandwidth"`
	Jitter    int    `mapstructure:"jitter"`
	Loss      int    `mapstructure:"loss"`
	Latency   int    `mapstructure:"latency"`
}

func (r *ConnectionRequest) FindDestinationID(destination ConnectionRequestDestination) string {
	if r == nil {
		return ""
	}
	for _, requestDestination := range r.Destinations {
		if requestDestination.equals(destination) {
			return requestDestination.DbID
		}
	}
	return ""
}

func (r *ConnectionRequest) FindDestination(destinationID string) *ConnectionRequestDestination {
	if r == nil {
		return nil
	}
	if split := strings.Split(destinationID, ":"); len(split) > 1 {
		destinationID = split[1]
	}
	for _, requestDestination := range r.Destinations {
		if requestDestination.DbID == destinationID {
			return &requestDestination
		}
	}
	return nil
}

func (d ConnectionRequestDestination) equals(destination ConnectionRequestDestination) bool {
	return d.Metadata.Name == destination.Metadata.Name &&
		d.Type == destination.Type &&
		d.ID == destination.ID &&
		d.Provider == destination.Provider &&
		d.SiteID == destination.SiteID
}
