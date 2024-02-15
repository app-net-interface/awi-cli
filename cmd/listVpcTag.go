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

	"github.com/spf13/cobra"
	awi "github.com/app-net-interface/awi-grpc/pb"

	"github.com/app-net-interface/awi-cli/prettyprint"
)

const (
	tagNameFlag = "tag-name"
)

// listVPCTagCmd represents the listVpcTag command
var listVPCTagCmd = &cobra.Command{
	Use:   "vpc-tag",
	Short: "List VPCs with tag",
	RunE:  listVpcTag,
}

func listVpcTag(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	c := awi.NewCloudClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	in := &awi.ListVPCTagRequest{
		Provider: strings.ToUpper(cmd.Flag(cloudFlag).Value.String()),
		Region:   cmd.Flag(regionFlag).Value.String(),
		Tag:      cmd.Flag(tagNameFlag).Value.String(),
	}
	vpcs, err := c.ListVPCTags(ctx, in)
	if err != nil {
		return err
	}

	printFormat := cmd.Flag(outputFlag).Value.String()
	prettyprint.PrintData(vpcs.VPCs, []prettyprint.Display{
		{Name: "Name", Display: "NAME"},
		{Name: "Tag", Display: "TAG"},
		{Name: "Region", Display: "REGION"},
		{Name: "ID", Display: "ID"},
		{Name: "AccountName", Display: "ACCOUNT_NAME"},
	}, printFormat)

	return nil
}

func init() {
	listCmd.AddCommand(listVPCTagCmd)
	listVPCTagCmd.Flags().String(cloudFlag, "", "Cloud")
	_ = listVPCTagCmd.MarkFlagRequired(cloudFlag)
	listVPCTagCmd.Flags().String(regionFlag, "", "Cloud region")
	listVPCTagCmd.Flags().String(tagNameFlag, "", "Name of the tag")
}
