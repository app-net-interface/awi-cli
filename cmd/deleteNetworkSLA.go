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

// deleteNetworkSLACmd represents the deleteNetworkSLA command
var deleteNetworkSLACmd = &cobra.Command{
	Use:   "network-sla",
	Short: "delete NetworkSLA",
	RunE:  deleteNetworkSLA,
}

func deleteNetworkSLA(cmd *cobra.Command, networkSLAs []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)
	for _, sla := range networkSLAs {
		deleteRequest := &awi.NetworkSLADeleteRequest{
			Name: sla,
		}
		c := awi.NewNetworkSLAServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		response, err := c.DeleteNetworkSLA(ctx, deleteRequest)
		if err != nil {
			cancel()
			return err
		}
		if response.String() == "" {
			fmt.Println("Response: successfully deleted NetworkSLA")
		} else {
			fmt.Printf("Response: %s\n", response.String())
		}
		cancel()
	}
	return nil
}

func init() {
	deleteCmd.AddCommand(deleteNetworkSLACmd)
}
