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

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	awi "github.com/app-net-interface/awi-grpc/pb"

	"github.com/app-net-interface/awi-cli/prettyprint"
)

// listAccessPolicyCmd represents the listAccessPolicy command
var listAccessPolicyCmd = &cobra.Command{
	Use:   "access-policy",
	Short: "List AccessPolicies",
	RunE:  listAccessPolicy,
}

func listAccessPolicy(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	c := awi.NewSecurityPolicyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	AccessPolicies, err := c.ListAccessPolicies(ctx, &awi.AccessPolicyListRequest{})
	if err != nil {
		return err
	}

	printFormat := cmd.Flag(outputFlag).Value.String()
	prettyprint.PrintConvertedData(AccessPolicies.GetAccessPolicies(), convertAccessPolicies(AccessPolicies.GetAccessPolicies()), []prettyprint.Display{
		{Name: "Name", Display: "NAME"},
	}, printFormat)

	return nil
}

type AccessPolicyDisplay struct {
	Name string
}

func convertAccessPolicies(crs []*awi.Security_AccessPolicy) []AccessPolicyDisplay {
	displays := make([]AccessPolicyDisplay, 0, len(crs))
	for _, sla := range crs {
		display := AccessPolicyDisplay{
			Name: sla.GetMetadata().GetName(),
		}
		displays = append(displays, display)
	}
	return displays
}

func init() {
	listCmd.AddCommand(listAccessPolicyCmd)
}
