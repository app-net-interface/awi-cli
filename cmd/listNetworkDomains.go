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
	"github.com/app-net-interface/awi-infra-guard/grpc/go/infrapb"
	"github.com/app-net-interface/awi-cli/prettyprint"
)

// listNetworkDomainsCmd represents the listNetworkDomainsCmd command
var listNetworkDomainsCmd = &cobra.Command{
	Use:   "network-domains",
	Short: "List all VPCs and VRFs",
	RunE:  listNetworkDomains,
}

type networkDomain struct {
	Type     string
	ID       string
	Name     string
	Provider string
}

func listNetworkDomains(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	printFormat := cmd.Flag(outputFlag).Value.String()

	var vpns []*awi.VPN

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	c := awi.NewCloudClient(conn)
	// list all VPNs
	vpnsList, err := c.ListVPNs(ctx, &awi.ListVPNRequest{})
	if err != nil {
		return err
	}
	vpns = vpnsList.VPNs

	conn2, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn2)
	infraC := infrapb.NewCloudProviderServiceClient(conn2)

	// list All VPCs
	awsvpcs, err := infraC.ListVPC(ctx, &infrapb.ListVPCRequest{
		Provider: "aws",
	})
	if err != nil {
		return err
	}
	gcpvpcs, err := infraC.ListVPC(ctx, &infrapb.ListVPCRequest{
		Provider: "gcp",
	})
	if err != nil {
		return err
	}

	networkDomains := make([]networkDomain, 0, len(vpns)+len(awsvpcs.Vpcs)+len(gcpvpcs.Vpcs))
	for _, vpn := range vpns {
		networkDomains = append(networkDomains, networkDomain{
			Type:     "VRF",
			Provider: "Cisco-SDWAN-vManage", // TODO should be based on info from VPN response message
			ID:       vpn.SegmentID,
			Name:     vpn.SegmentName,
		})
	}
	for _, vpc := range append(awsvpcs.Vpcs, gcpvpcs.Vpcs...) {
		networkDomains = append(networkDomains, networkDomain{
			Type:     "VPC",
			Provider: vpc.Provider,
			ID:       vpc.Id,
			Name:     vpc.Name,
		})
	}

	prettyprint.PrintData(networkDomains, []prettyprint.Display{
		{Name: "Type", Display: "TYPE"},
		{Name: "Provider", Display: "PROVIDER"},
		{Name: "Name", Display: "NAME"},
		{Name: "ID", Display: "ID"},
	}, printFormat)

	return nil
}

func init() {
	listCmd.AddCommand(listNetworkDomainsCmd)
}
