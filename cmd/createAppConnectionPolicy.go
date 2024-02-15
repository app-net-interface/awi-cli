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

// createAppPolicyCmd represents the createAppPolicy command
var createAppPolicyCmd = &cobra.Command{
	Use:   "app-connection-policy",
	Short: "Create Application Connection Policy",
	RunE:  createAppPolicy,
}

func createAppPolicy(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	connectionConfigPath := cmd.Flag(connectionConfigFlag).Value.String()

	// Set up a connection to the server.
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)

	conf, err := getAppConnectionConfigGRPC(connectionConfigPath)
	if err != nil {
		return fmt.Errorf("could not initialize connection config: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if conf == nil {
		return fmt.Errorf("specify App Connection Policy in the config")
	}
	logger.Infof("sending create AppConnection Policy request")
	cc := awi.NewAppConnectionControllerClient(conn)
	response, err := cc.CreateAppConnectionPolicy(ctx, &awi.CreateAppConnectionPolicyRequest{AppConnection: conf})
	if err != nil {
		return fmt.Errorf("could not create app connection policy: %v", err)
	}
	fmt.Printf("ID: %s\n", response.GetId())
	fmt.Printf("Status: %s\n", response.GetStatus())
	return nil
}

func init() {
	createCmd.AddCommand(createAppPolicyCmd)
}
