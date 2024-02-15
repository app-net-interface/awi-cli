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
	"github.com/app-net-interface/awi-infra-guard/grpc/go/infrapb"
	"github.com/app-net-interface/awi-cli/prettyprint"
)

const (
	cloudFlag     = "cloud"
	regionFlag    = "region"
	accountIDFlag = "account-id"
	untaggedFlag  = "untagged"
)

// listVPCCmd represents the listVPC command
var listVPCCmd = &cobra.Command{
	Use:   "vpc",
	Short: "List VPCs",
	RunE:  listVPC,
}

func listVPC(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	c := infrapb.NewCloudProviderServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	in := &infrapb.ListVPCRequest{
		Provider:  strings.ToUpper(cmd.Flag(cloudFlag).Value.String()),
		Region:    cmd.Flag(regionFlag).Value.String(),
		AccountId: cmd.Flag(accountIDFlag).Value.String(),
	}
	vpcs, err := c.ListVPC(ctx, in)
	if err != nil {
		return err
	}

	printFormat := cmd.Flag(outputFlag).Value.String()
	prettyprint.PrintData(vpcs.Vpcs, []prettyprint.Display{
		{Name: "Name", Display: "NAME"},
		{Name: "Region", Display: "REGION"},
		{Name: "Id", Display: "ID"},
		{Name: "AccountId", Display: "ACCOUNT_ID"},
	}, printFormat)

	return nil
}

func init() {
	listCmd.AddCommand(listVPCCmd)
	listVPCCmd.Flags().String(cloudFlag, "", "Cloud")
	_ = listVPCCmd.MarkFlagRequired(cloudFlag)
	listVPCCmd.Flags().String(regionFlag, "", "Cloud region")
	listVPCCmd.Flags().String(accountIDFlag, "", "ID of the account")
	listVPCCmd.Flags().String(untaggedFlag, "", "untagged")
}
