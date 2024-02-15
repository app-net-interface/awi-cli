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
)

// deleteConnectionCmd represents the deleteConnection command
var deleteConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "DeleteConnection VPC from VPN or other VPC",
	RunE:  deleteConnection,
}

func deleteConnection(cmd *cobra.Command, connections []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	for _, connectionID := range connections {
		disconnectRequest := &awi.DisconnectRequest{
			ConnectionId: connectionID,
		}
		c := awi.NewConnectionControllerClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		response, err := c.Disconnect(ctx, disconnectRequest)
		if err != nil {
			return err
		}
		fmt.Printf("Response: %s\n", response.String())
		fmt.Printf("Status: %v\n", response.Status.String())
	}

	return nil
}

func init() {
	deleteCmd.AddCommand(deleteConnectionCmd)
}