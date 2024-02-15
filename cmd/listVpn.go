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

// listVPNCmd represents the listVPN command
var listVPNCmd = &cobra.Command{
	Use:   "vpn",
	Short: "List VPNs",
	RunE:  listVPN,
}

func listVPN(cmd *cobra.Command, _ []string) error {
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
	in := &awi.ListVPNRequest{
		Provider: strings.ToUpper(cmd.Flag(cloudFlag).Value.String()),
	}
	vpns, err := c.ListVPNs(ctx, in)
	if err != nil {
		return err
	}
	printFormat := cmd.Flag(outputFlag).Value.String()
	prettyprint.PrintData(vpns.VPNs, []prettyprint.Display{
		{Name: "SegmentID", Display: "SEGMENT_ID"},
		{Name: "SegmentName", Display: "SEGMENT_NAME"},
		{Name: "ID", Display: "ID"},
	}, printFormat)

	return nil
}

func init() {
	listCmd.AddCommand(listVPNCmd)
	listVPNCmd.Flags().String(cloudFlag, "", "Cloud")
}
