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
	"os"
	"time"

	"github.com/mitchellh/mapstructure"

	awi "github.com/app-net-interface/awi-grpc/pb"
	"github.com/app-net-interface/catalyst-sdwan-app-client/vmanage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	useProxyFlag         = globalsFlag + ".use_proxy"
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
	controllerURL := viper.GetString(urlFlag)
	retries := viper.GetInt(longPollRetriesFlag)
	retriesInterval := viper.GetDuration(retriesIntervalFlag)
	client, err := vmanage.NewClient(controllerURL, httpclient, logger, retriesInterval, retries)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// apiPrefixInterceptor attaches the API prefix for the Envoy Proxy
// to distinguish that calls performed by the AWI-CLI should be forwarded
// to backend services rather than front-end.
//
// If the awi-cli calls backend services directly, the prefix should be
// empty - in such case apiPrefixInterceptor will not add anything.
func apiPrefixInterceptor(prefix string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		modifiedMethod := method
		if prefix != "" {
			modifiedMethod = prefix + modifiedMethod
		}

		return invoker(ctx, modifiedMethod, req, reply, cc, opts...)
	}
}

func getGRPCClient() (*grpc.ClientConn, error) {
	address := viper.GetString(urlFlag)
	proxyEnabled := viper.GetBool(useProxyFlag)
	logger.Debugf("connecting to %s", address)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	if proxyEnabled {
		dialOptions = append(
			dialOptions,
			grpc.WithUnaryInterceptor(apiPrefixInterceptor("/grpc")),
		)
	}

	conn, err := grpc.DialContext(
		ctx,
		address,
		dialOptions...,
	)

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
