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
	"strings"
	"time"
	"github.com/app-net-interface/awi-infra-guard/grpc/go/infrapb"

	"github.com/spf13/cobra"
	"github.com/app-net-interface/awi-cli/prettyprint"
)

const (
	vpcFlag        = "vpc"
	tagFlag        = "labels"
	zoneFlag       = "zone"
	cidrFlag       = "cidr"
	showLabelsFlag = "show-labels"
)

// listSubnetCmd represents the listSubnet command
var listSubnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "List Subnets",
	RunE:  listSubnet,
}

func listSubnet(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	cloud := cmd.Flag(cloudFlag).Value.String()
	printFormat := cmd.Flag(outputFlag).Value.String()
	vpcID := cmd.Flag(vpcFlag).Value.String()
	tag := cmd.Flag(tagFlag).Value.String()
	labels := make(map[string]string)
	if tag != "" {
		tags := strings.Split(tag, labelsSeparator)
		for _, t := range tags {
			keyVal := strings.Split(t, keyValueSepartor)
			if len(keyVal) != 2 {
				return fmt.Errorf("specify tag in format 'key1=value1,key2=value2'")
			}
			labels[keyVal[0]] = keyVal[1]
		}
	}
	zone := cmd.Flag(zoneFlag).Value.String()
	cidr := cmd.Flag(cidrFlag).Value.String()
	showLabels, err := cmd.Flags().GetBool(showLabelsFlag)
	if err != nil {
		showLabels = false
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	c := infrapb.NewCloudProviderServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	in := &infrapb.ListSubnetsRequest{
		Zone:     zone,
		VpcId:    vpcID,
		Provider: cloud,
		Cidr:     cidr,
		Labels:   labels,
	}
	subnets, err := c.ListSubnets(ctx, in)
	if err != nil {
		return err
	}

	displays := []prettyprint.Display{
		{Name: "SubnetId", Display: "SUBNET_ID"},
		{Name: "Name", Display: "SUBNET_NAME"},
		{Name: "VpcId", Display: "VPC_ID"},
		{Name: "Zone", Display: "ZONE"},
		{Name: "CidrBlock", Display: "CIDR_BLOCK"},
	}
	if showLabels {
		displays = append(displays, prettyprint.Display{Name: "Labels", Display: "LABELS"})
	}
	prettyprint.PrintData(subnets.Subnets, displays, printFormat)

	return nil
}

func init() {
	listCmd.AddCommand(listSubnetCmd)
	listSubnetCmd.Flags().String(cloudFlag, "", "Cloud")
	_ = listSubnetCmd.MarkFlagRequired(cloudFlag)
	listSubnetCmd.Flags().String(vpcFlag, "", "VPC ID")
	listSubnetCmd.Flags().String(tagFlag, "", "Labels in key1=value1,key2=value2 format")
	listSubnetCmd.Flags().String(zoneFlag, "", "Availability Zone")
	listSubnetCmd.Flags().String(cidrFlag, "", "CIDR")
	listSubnetCmd.Flags().Bool(showLabelsFlag, false, "Display labels")
}
