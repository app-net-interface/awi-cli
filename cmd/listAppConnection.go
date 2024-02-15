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

// listAppConnectionCmd represents the listAppConnection command
var listAppConnectionCmd = &cobra.Command{
	Use:   "app-connection",
	Short: "List App Connection",
	RunE:  listAppConnection,
}

func listAppConnection(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	c := awi.NewAppConnectionControllerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connections, err := c.ListConnectedApps(ctx, &awi.ListAppConnectionsRequest{})
	if err != nil {
		return err
	}

	printFormat := cmd.Flag(outputFlag).Value.String()
	prettyprint.PrintConvertedData(connections.AppConnections, convertAppConnectionRequest(connections.AppConnections), []prettyprint.Display{
		{Name: "ID", Display: "ID"},
		{Name: "Name", Display: "NAME"},
		{Name: "NetworkDomainConnectionName", Display: "NETWORK_DOMAIN_CONNECTION_NAME"},
		{Name: "CreationTimestamp", Display: "CREATE_TIME"},
		{Name: "ModificationTimestamp", Display: "MOD_TIME"},
		{Name: "Status", Display: "STATUS"},
	}, printFormat)

	return nil
}

type appReqDisplay struct {
	ID                          string
	Name                        string
	NetworkDomainConnectionName string
	Status                      string
	ModificationTimestamp       string
	CreationTimestamp           string
}

func convertAppConnectionRequest(crs []*awi.AppConnectionInformation) []appReqDisplay {
	displays := make([]appReqDisplay, 0, len(crs))
	for _, request := range crs {
		display := appReqDisplay{
			ID:                          request.GetId(),
			Name:                        request.GetAppConnectionConfig().GetMetadata().GetName(),
			NetworkDomainConnectionName: request.GetNetworkDomainConnectionName(),
			Status:                      awi.Status_name[int32(request.Status)],
			ModificationTimestamp:       request.GetAppConnectionConfig().GetMetadata().GetModificationTimestamp(),
			CreationTimestamp:           request.GetAppConnectionConfig().GetMetadata().GetCreationTimestamp(),
		}
		displays = append(displays, display)
	}
	return displays
}

func init() {
	listCmd.AddCommand(listAppConnectionCmd)
}
