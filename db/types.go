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

package db

import "github.com/app-net-interface/awi-cli/types"

type ACL struct {
	ID              string
	Provider        string
	SecurityGroupID string
	Source          *ACLSource
	Destination     *ACLDestination
}

type ACLSource struct {
	VPCID     string
	IPs       []string
	Labels    types.Label
	Condition string
}

type ACLDestination struct {
	VPCID       string
	InstanceIDs []string
	Labels      types.Label
	Condition   string
}
