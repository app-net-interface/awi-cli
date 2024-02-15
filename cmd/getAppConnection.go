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
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/spf13/cobra"
	awi "github.com/app-net-interface/awi-grpc/pb"
)

// getAppCmd represents the get AppConnection command
var getAppCmd = &cobra.Command{
	Use:   "app-connection",
	Short: "get Application connection",
	RunE:  getApp,
}

func getApp(cmd *cobra.Command, _ []string) error {
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	id := cmd.Flag(idFlag).Value.String()

	// Set up a connection to the server.
	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer connClose(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Infof("sending get AppConnection request")
	cc := awi.NewAppConnectionControllerClient(conn)
	response, err := cc.GetAppConnection(ctx, &awi.GetAppConnectionRequest{ConnectionId: id})
	if err != nil {
		return fmt.Errorf("could not get connection: %v", err)
	}

	b, err := protojson.Marshal(response.AppConnection)
	if err != nil {
		return err
	}
	var obj map[string]interface{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return err
	}
	d, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(d))
	return nil
}

func init() {
	getCmd.AddCommand(getAppCmd)
}
