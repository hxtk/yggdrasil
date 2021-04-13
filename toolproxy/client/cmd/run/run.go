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
package run

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hxtk/yggdrasil/common/config/tlsconfig"
	"github.com/hxtk/yggdrasil/toolproxy/client/pkg/rpc"
)

func NewCmdRun() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run a command on the remote host.",
		Long: `Execute a command
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			tlsConfig, err := tlsconfig.FromViper(viper.GetViper())
			if err != nil {
				log.WithError(err).Fatal("Error reading TLS Config")
			}
			client := rpc.New(":8080", tlsConfig)
			client.Run(context.Background(), args)
		},
	}
}
