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

// listNetworkSLACmd represents the listNetworkSLA command
var listNetworkSLACmd = &cobra.Command{
	Use:   "network-sla",
	Short: "List NetworkSLAs",
	RunE:  listNetworkSLA,
}

func listNetworkSLA(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	c := awi.NewNetworkSLAServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	networkSLAs, err := c.ListNetworkSLAs(ctx, &awi.NetworkSLAListReqest{})
	if err != nil {
		return err
	}

	printFormat := cmd.Flag(outputFlag).Value.String()
	prettyprint.PrintConvertedData(networkSLAs.GetNetworkSLAs(), convertNetworksSLAs(networkSLAs.GetNetworkSLAs()), []prettyprint.Display{
		{Name: "Name", Display: "NAME"},
		{Name: "Description", Display: "DESCRIPTION"},
		{Name: "Bandwidth", Display: "BANDWIDTH[Mbps]"},
		{Name: "Jitter", Display: "JITTER[Ms]"},
		{Name: "Latency", Display: "LATENCY[Ms]"},
		{Name: "Loss", Display: "LOSS[%]"},
		{Name: "Priority", Display: "PRIORITY"},
		{Name: "EnforcementRequestType", Display: "ENFORCEMENT_REQUEST_TYPE"},
	}, printFormat)

	return nil
}

type networkSLADisplay struct {
	Name                   string
	Description            string
	Bandwidth              string
	Jitter                 string
	Latency                string
	Loss                   string
	Priority               string
	EnforcementRequestType string
}

func convertNetworksSLAs(crs []*awi.NetworkSLA) []networkSLADisplay {
	displays := make([]networkSLADisplay, 0, len(crs))
	for _, sla := range crs {
		display := networkSLADisplay{
			Name:                   sla.GetMetadata().GetName(),
			Description:            sla.GetMetadata().GetName(),
			Bandwidth:              fmt.Sprintf("%.1f", sla.GetTrafficProfile().GetBandwidth()),
			Jitter:                 fmt.Sprintf("%.2f", sla.GetTrafficProfile().GetJitter()),
			Latency:                fmt.Sprintf("%.2f", sla.GetTrafficProfile().GetLatency()),
			Loss:                   fmt.Sprintf("%.2f", sla.GetTrafficProfile().GetLoss()),
			Priority:               sla.GetPriority(),
			EnforcementRequestType: sla.GetEnforcementRequest().GetType(),
		}
		displays = append(displays, display)
	}
	return displays
}

func init() {
	listCmd.AddCommand(listNetworkSLACmd)
}
