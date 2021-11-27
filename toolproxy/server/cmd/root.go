/*
Copyright © 2021 Peter Sanders

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hxtk/yggdrasil/common/config/postgres"
	"github.com/hxtk/yggdrasil/common/config/tlsconfig"
	"github.com/hxtk/yggdrasil/common/grpc/server"
	"github.com/hxtk/yggdrasil/toolproxy/server/pkg/rpc"
)

var cfgFiles []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tool-server",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tlsConfig, err := tlsconfig.FromViper(viper.GetViper())
		if err != nil {
			log.WithError(err).Fatal("Error reading TLS Config")
		}

		s := server.New()

		db, err := postgres.FromViper(viper.GetViper())
		if err != nil {
			log.WithError(err).Fatal("Error opening database.")
		}
		rpcServer := rpc.New(db)
		s.Register(rpcServer)
		log.Info("Registration complete.")

		addr := viper.GetString("grpc.addr")
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.WithError(err).WithField("addr", addr).Fatal("Cannot open grpc listener.")
		}

		if viper.GetBool("tls.enabled") {
			lis = tls.NewListener(lis, tlsConfig)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			log.Info("gRPC server started.")
			err := s.ServeGRPC(lis)
			if err != nil {
				log.WithError(err).Fatal("gRPC listener returned error.")
			}
			wg.Done()
		}()

		lis, err = net.Listen("tcp", viper.GetString("http.addr"))
		if err != nil {
			log.WithError(err).WithField("addr", addr).Fatal("Cannot open http listener.")
		}

		if viper.GetBool("tls.enabled") {
			lis = tls.NewListener(lis, tlsConfig)
		}

		wg.Add(1)
		go func() {
			log.Info("HTTP server started.")
			err := http.Serve(lis, s)
			if err != nil {
				log.WithError(err).Fatal("http listener returned error.")
			}
			wg.Done()
		}()

		wg.Wait()
		log.Info("Servers shut down.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&cfgFiles, "config", "c", cfgFiles, "Configuration file. If specified more than once, subsequent files will be merged into the first.")

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(cfgFiles) > 0 {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFiles[0])
		fmt.Printf("Config file specified in arguments: %s\n", cfgFiles[0])
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.WithError(err).Fatal("Could not get working directory.")
		}
		log.WithField("cwd", cwd).Println("Retrieved working directory.")

		// Search config in home directory with name "toolproxy.yaml" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(cwd)
		viper.AddConfigPath("/etc/toolproxy")
		viper.SetConfigName("toolproxy")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetConfigType("yaml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if len(cfgFiles) > 0 {
		for _, f := range cfgFiles[1:] {
			viper.SetConfigFile(f)
			if err := viper.MergeInConfig(); err == nil {
				fmt.Println("Using config file:", viper.ConfigFileUsed())
			}
		}
	}

	log.Println(viper.AllKeys())
}
