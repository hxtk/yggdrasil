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
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hxtk/yggdrasil/common/config/flags"
	"github.com/hxtk/yggdrasil/common/grpc/server"
	"github.com/hxtk/yggdrasil/toolproxy/server/pkg/rpc"
)

var connectionFlags *flags.ConnectionFlags
var cfgFile string

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
		tlsConfig, err := getTLSConfig()
		if err != nil {
			log.WithError(err).Fatal("Error reading TLS Config")
		}

		var s *server.Server
		if viper.GetBool("tls.enabled") {
			s = server.New(tlsConfig)
		} else {
			s = server.New(nil)
		}

		rpcServer := rpc.New(getPostgresDSN())
		s.Register(rpcServer)

		addr := viper.GetString("addr")
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.WithError(err).WithField("addr", addr).Fatal("Cannot open listener.")
		}
		if viper.GetBool("tls.enabled") {
			lis = tls.NewListener(lis, tlsConfig)
		}

		s.GRPCServer.Serve(lis)
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
	connectionFlags = flags.NewConnectionFlags(flags.SERVER_ROLE)
	connectionFlags.AddFlags(rootCmd.Flags())

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "Path to a configuration file to use.")

	cobra.OnInitialize(initConfig)
}

func getPostgresDSN() string {
	var parts []string
	if user := viper.GetString("postgres.user"); user != "" {
		parts = append(parts, fmt.Sprintf("user=%s", user))
	}

	if password := viper.GetString("postgres.password"); password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", password))
	}

	if sslmode := viper.GetString("postgres.sslmode"); sslmode != "" {
		parts = append(parts, fmt.Sprintf("sslmode=%s", sslmode))
	}

	if host := viper.GetString("postgres.host"); host != "" {
		parts = append(parts, fmt.Sprintf("host=%s", host))
	}

	if port := viper.GetString("postgres.port"); port != "" {
		parts = append(parts, fmt.Sprintf("port=%s", port))
	}

	if database := viper.GetString("postgres.database"); database != "" {
		parts = append(parts, fmt.Sprintf("database=%s", database))
	}

	if sslcert := viper.GetString("postgres.sslcert"); sslcert != "" {
		parts = append(parts, fmt.Sprintf("sslcert=%s", sslcert))
	}

	if sslkey := viper.GetString("postgres.sslkey"); sslkey != "" {
		parts = append(parts, fmt.Sprintf("sslkey=%s", sslkey))
	}

	if sslrootcert := viper.GetString("postgres.sslrootcert"); sslrootcert != "" {
		parts = append(parts, fmt.Sprintf("sslrootcert=%s", sslrootcert))
	}

	return strings.Join(parts, " ")
}

func getTLSConfig() (*tls.Config, error) {
	config := &tls.Config{}

	certFile := viper.GetString("tls.certificate")
	keyFile := viper.GetString("tls.key")

	log.WithFields(log.Fields{"cert": certFile, "key": keyFile}).Println("Loading certificates.")
	certificate, err := tls.LoadX509KeyPair(
		certFile,
		keyFile,
	)
	if err != nil {
		return nil, err
	}
	config.Certificates = []tls.Certificate{certificate}

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

	return config, nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		fmt.Printf("Config file specified in arguments: %s\n", cfgFile)
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
		viper.SetConfigName("toolproxy")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetConfigType("yaml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	log.Println(viper.ConfigFileUsed())
	log.Println(viper.AllKeys())
}
