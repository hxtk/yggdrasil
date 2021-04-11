/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hxtk/yggdrasil/toolproxy/client/pkg/rpc"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command",
	Long: `Execute a command
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tlsConfig, err := getTLSConfig()
		if err != nil {
			log.WithError(err).Fatal("Error loading TLS Config.")
		}
		client := rpc.New(":8080", tlsConfig)
		client.Run(context.Background(), args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTLSConfig() (*tls.Config, error) {
	config := &tls.Config{}

	if viper.IsSet("tls.certificate") && viper.IsSet("tls.key") {
		certFile := viper.GetString("tls.certificate")
		keyFile := viper.GetString("tls.key")
		log.WithFields(log.Fields{"cert": certFile, "key": keyFile}).Debugln("Loading certificates.")
		certificate, err := tls.LoadX509KeyPair(
			certFile,
			keyFile,
		)
		if err != nil {
			return nil, err
		}
		config.Certificates = []tls.Certificate{certificate}
	}

	config.InsecureSkipVerify = viper.GetBool("tls.insecure")

	if viper.GetString("tls.ca") != "" {
		certBytes, err := ioutil.ReadFile(viper.GetString("tls.ca"))
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(certBytes); !ok {
			return nil, errors.New("error parsing CA")
		}

		config.ClientCAs = certPool
		config.RootCAs = certPool
	}

	config.ServerName = viper.GetString("tls.hostname")

	return config, nil
}
