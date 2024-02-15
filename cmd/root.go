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
	"crypto/tls"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	awi "github.com/app-net-interface/awi-grpc/pb"
	"github.com/app-net-interface/catalyst-sdwan-app-client/vmanage"

	"github.com/app-net-interface/awi-cli/types"
)

const (
	configFlag           = "config"
	globalsFlag          = "Globals"
	logFileFlag          = globalsFlag + ".log_file"
	logLevelFlag         = globalsFlag + ".log_level"
	dbFileName           = globalsFlag + ".db_name"
	controllersFlag      = "controllers.sdwan"
	sessionIDFlag        = controllersFlag + ".session_id"
	tokenFlag            = controllersFlag + ".token"
	urlFlag              = globalsFlag + ".grpc_url"
	secureConnectionFlag = controllersFlag + ".secure_connection"
	longPollRetriesFlag  = controllersFlag + ".controller_connection_retries"
	retriesIntervalFlag  = controllersFlag + ".retries_interval"
	labelsSeparator      = ","
	keyValueSepartor     = "="
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awi",
	Short: "CLI for connecting networking and application resources through Cisco Catalyst WAN",
}

var logger = log.New()

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP(configFlag, "c", "config.yaml", "Configuration file in YAML format")
}

func getClient() (vmanage.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create cookie jar: %v", err)
	}
	vManageURL := viper.GetString(urlFlag)
	u, err := url.Parse(vManageURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %v", err)
	}
	sessionID := viper.GetString(sessionIDFlag)
	if err != nil {
		return nil, fmt.Errorf("error while reading credentials: %v", err)
	}
	cookie := &http.Cookie{Name: "JSESSIONID", Value: sessionID}
	jar.SetCookies(u, []*http.Cookie{cookie})
	client, err := getClientWithJar(jar)
	if err != nil {
		return nil, err
	}
	client.SetToken(viper.GetString(tokenFlag))
	return client, nil
}

func getClientWithJar(jar *cookiejar.Jar) (vmanage.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: !viper.GetBool(secureConnectionFlag),
		MinVersion:         tls.VersionTLS12,
	}
	httpclient := &http.Client{
		Transport: transport,
		Jar:       jar,
	}
	vManageURL := viper.GetString(urlFlag)
	retries := viper.GetInt(longPollRetriesFlag)
	retriesInterval := viper.GetDuration(retriesIntervalFlag)
	client, err := vmanage.NewClient(vManageURL, httpclient, logger, retriesInterval, retries)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getGRPCClient() (*grpc.ClientConn, error) {
	address := viper.GetString(urlFlag)
	logger.Debugf("connecting to %s", address)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		//conn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		//if err != nil {
		return nil, fmt.Errorf("could not connect to grpc server: %v", err)
		//}
	}
	logger.Debugf("connected")
	return conn, nil
}

func initConfig(configFilePath string) error {
	viper.AutomaticEnv()
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		// Return userError if the config file doesn't exist or has unsupported config type.
		if _, match := err.(viper.UnsupportedConfigError); match || errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("unsupported Config or File doesn't exist")
		}
		return err
	}

	if err := initLogger(); err != nil {
		return err
	}

	logger.Debugf("Using config file: %s", viper.ConfigFileUsed())

	return nil
}

func initConnectionConfig(configFilePath string) ([]types.ConnectionRequest, []types.AccessControlRequest, error) {
	if err := loadConfig(configFilePath); err != nil {
		return nil, nil, err
	}
	logger.Infof("Using connection config file: %s", viper.ConfigFileUsed())
	var requests []types.ConnectionRequest
	if err := viper.UnmarshalKey(connectionRequestFlag, &requests); err != nil {
		return nil, nil, fmt.Errorf("could not read connections: %v", err)
	}
	var acls []types.AccessControlRequest
	if err := viper.UnmarshalKey(accessRequestFlag, &acls); err != nil {
		return nil, nil, fmt.Errorf("could not read access controls: %v", err)
	}
	return requests, acls, nil
}

func getConnectionConfigGRPC(configFilePath string) (*awi.ConnectionRequest, error) {
	if err := loadConfig(configFilePath); err != nil {
		return nil, err
	}
	logger.Infof("Using connection config file: %s", viper.ConfigFileUsed())
	request := &awi.ConnectionRequest{}

	if err := viper.UnmarshalKey(specFlag, &request.Spec,
		func(config *mapstructure.DecoderConfig) { config.ErrorUnused = true }); err != nil {
		return nil, fmt.Errorf("could not read connection spec: %v", err)
	}
	if err := viper.UnmarshalKey(metadataFlag, &request.Metadata); err != nil {
		return nil, fmt.Errorf("could not read connection metadata: %v", err)
	}
	return request, nil
}

func getAppConnectionConfigGRPC(configFilePath string) (*awi.AppConnection, error) {
	if err := loadConfig(configFilePath); err != nil {
		return nil, err
	}
	logger.Infof("Using connection config file: %s", viper.ConfigFileUsed())
	var acl *awi.AppConnection
	err := viper.UnmarshalKey(accessRequestFlag, &acl, func(config *mapstructure.DecoderConfig) {
		config.ErrorUnused = true
	})
	if err != nil {
		return nil, fmt.Errorf("could not read app connection config: %v", err)
	}
	if acl == nil {
		err := viper.UnmarshalKey(fmt.Sprintf(specFlag+"."+accessRequestFlag), &acl, func(config *mapstructure.DecoderConfig) {
			config.ErrorUnused = true
		})
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %v", err)
		}
		if acl == nil {
			return nil, fmt.Errorf("wrong configuration")
		}
	}
	return acl, nil
}

func getNetworkSLAConfigGRPC(configFilePath string) (*awi.NetworkSLA, error) {
	if err := loadConfig(configFilePath); err != nil {
		return nil, err
	}
	logger.Infof("Using networkSLA config file: %s", viper.ConfigFileUsed())
	var request *awi.NetworkSLA
	if err := viper.UnmarshalKey(networkSLAFlag, &request); err != nil {
		return nil, fmt.Errorf("could not read networkSLA config: %v", err)
	}
	return request, nil
}

func getAccessControlConfigGRPC(configFilePath string) (*awi.Security_AccessPolicy, error) {
	if err := loadConfig(configFilePath); err != nil {
		return nil, err
	}
	logger.Infof("Using Access Policy config file: %s", viper.ConfigFileUsed())
	var request *awi.Security_AccessPolicy
	if err := viper.UnmarshalKey(specFlag, &request, func(config *mapstructure.DecoderConfig) {
		config.ErrorUnused = true
	}); err != nil {
		return nil, fmt.Errorf("could not read access policy config: %v", err)
	}
	return request, nil
}

func loadConfig(configFilePath string) error {
	viper.SetConfigFile(configFilePath)
	if err := viper.MergeInConfig(); err != nil {
		return err
	}
	return nil
}

func initLogger() error {
	logLevel := viper.GetString(logLevelFlag)
	logger = log.New()
	if logLevel != "" {
		var err error
		logger.Level, err = log.ParseLevel(logLevel)
		if err != nil {
			return err
		}
	} else {
		logger.Level = log.InfoLevel
	}
	logFile := viper.GetString(logFileFlag)
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		logger.Out = file
	} else {
		logger.Out = os.Stdout
	}
	return nil
}

func connClose(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		fmt.Printf("Error during gRPC connection conenction close: %v", err)
	}
}
