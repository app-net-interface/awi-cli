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
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

const (
	usernameNameFlag        = "username"
	environmentVariableFlag = "VMANAGE_PASSWORD"
)

// generateTokenCmd represents the listVpcTag command
var generateTokenCmd = &cobra.Command{
	Use:   "generate-token",
	Short: "Update config with generated tokens",
	Long:  fmt.Sprintf("Update config with generated tokens. It tries to take password from %v environment variable. If password is not found it will ask user to provide it", environmentVariableFlag),
	RunE:  generateToken,
}

func generateToken(cmd *cobra.Command, _ []string) error {
	password := os.Getenv(environmentVariableFlag)
	if password == "" {
		fmt.Println("Password:")
		passwordByte, err := term.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("could not get password: %v", err)
		}
		password = string(passwordByte)
	}
	if err := initConfig(cmd.Flag(configFlag).Value.String()); err != nil {
		return fmt.Errorf("could not initialize config: %v", err)
	}
	username := cmd.Flag(usernameNameFlag).Value.String()

	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("could not create cookie jar: %v", err)
	}
	client, err := getClientWithJar(jar)
	if err != nil {
		return fmt.Errorf("could not get client: %v", err)
	}

	ctx := context.Background()
	if err := client.Login(ctx, username, password); err != nil {
		return fmt.Errorf("could not login: %v", err)
	}

	vManageURL := viper.GetString(urlFlag)
	u, err := url.Parse(vManageURL)
	if err != nil {
		return fmt.Errorf("could not parse url: %v", u)
	}
	viper.Set(tokenFlag, client.GetToken())
	viper.Set(sessionIDFlag, jar.Cookies(u)[0].Value)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("could not write config: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(generateTokenCmd)
	generateTokenCmd.Flags().StringP(usernameNameFlag, "u", "", "Username")
	_ = generateTokenCmd.MarkFlagRequired(usernameNameFlag)
}
