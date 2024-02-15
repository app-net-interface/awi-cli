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

// listConnectionCmd represents the listConnection command
var listConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "List Connections",
	RunE:  listConnection,
}

func listConnection(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	c := awi.NewConnectionControllerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connections, err := c.ListConnections(ctx, &awi.ListConnectionsRequest{})
	if err != nil {
		return err
	}

	printFormat := cmd.Flag(outputFlag).Value.String()
	prettyprint.PrintConvertedData(connections.GetConnections(), convertConnectionRequest(connections.GetConnections()), []prettyprint.Display{
		{Name: "ID", Display: "ID"},
		{Name: "Name", Display: "NAME"},
		{Name: "SourceName", Display: "SRC_NAME"},
		{Name: "SourceType", Display: "SRC_TYPE"},
		{Name: "SourceProvider", Display: "SRC_PROVIDER"},
		{Name: "DestinationName", Display: "DEST_NAME"},
		{Name: "DestinationType", Display: "DEST_TYPE"},
		{Name: "DestinationProvider", Display: "DEST_PROVIDER"},
		//{Name: "DefaultAccess", Display: "DEFAULT_ACCESS"},
		{Name: "CreationTimestamp", Display: "CREATE_TIME"},
		{Name: "ModificationTimestamp", Display: "MOD_TIME"},
		{Name: "Status", Display: "STATUS"},
	}, printFormat)

	return nil
}

type reqDisplay struct {
	ID                  string
	Name                string
	SourceName          string
	SourceType          string
	SourceProvider      string
	DestinationName     string
	DestinationType     string
	DestinationProvider string
	//DefaultAccess         string
	Status                string
	ModificationTimestamp string
	CreationTimestamp     string
}

func convertConnectionRequest(crs []*awi.ConnectionInformation) []reqDisplay {
	displays := make([]reqDisplay, 0, len(crs))
	for _, conn := range crs {
		display := reqDisplay{
			ID:                  conn.GetId(),
			Name:                conn.GetMetadata().GetName(),
			SourceName:          conn.GetSource().GetName(),
			SourceType:          conn.GetSource().GetType(),
			SourceProvider:      conn.GetSource().GetProvider(),
			DestinationName:     conn.GetDestination().GetName(),
			DestinationType:     conn.GetDestination().GetType(),
			DestinationProvider: conn.GetDestination().GetProvider(),
			//DefaultAccess:         conn.GetDestination().GetDefaultAccessControl(),
			Status:                awi.Status_name[int32(conn.GetStatus())],
			ModificationTimestamp: conn.ModificationTimestamp,
			CreationTimestamp:     conn.CreationTimestamp,
		}
		displays = append(displays, display)
	}
	return displays
}

func init() {
	listCmd.AddCommand(listConnectionCmd)
	listConnectionCmd.Flags().String(cloudFlag, "", "Cloud")
}
