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
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	fastping "github.com/tatsushid/go-fastping"

	"github.com/hxtk/yggdrasil/common/config/tlsconfig"
	"github.com/hxtk/yggdrasil/common/grpc/server"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
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

		s := server.New(tlsConfig)

		p := fastping.NewPinger()
		p.MaxRTT = time.Minute

		for _, v := range args {
			err := p.AddIP(v)
			if err != nil {
				log.WithError(err).Fatal("Failed to add IP.")
			}
		}
		attempts := promauto.NewCounter(prometheus.CounterOpts{
			Name: "pinger_ping_attempts_total",
			Help: "The total number of ICMP echo requests attempted per IP",
		})
		attempts.Inc()

		buckets := prometheus.LinearBuckets(0, 5, 100)
		buckets = append(buckets, prometheus.LinearBuckets(500, 100, 5)...)
		buckets = append(buckets, prometheus.LinearBuckets(1000, 1000, 5)...)
		pings := promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "pinger_ping_rtts",
			Buckets: buckets,
		}, []string{"peer"})

		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			pings.WithLabelValues(addr.String()).Observe(float64(rtt.Milliseconds()))
		}
		p.OnIdle = func() {
			attempts.Inc()
		}

		addr := viper.GetString("http.addr")
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.WithError(err).WithField("addr", addr).Fatal("Cannot open http listener.")
		}

		if viper.GetBool("tls.enabled") {
			lis = tls.NewListener(lis, tlsConfig)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			log.WithField("addr", addr).Info("Serving HTTP.")
			err := http.Serve(lis, s)
			if err != nil {
				log.WithError(err).Fatal("http listener returned error.")
			}
			wg.Done()
			p.Stop()
		}()

		wg.Add(1)
		go func() {
			p.RunLoop()
			<-p.Done()
			if err := p.Err(); err != nil {
				log.WithError(err).Error("Ping failed.")
			}
			wg.Done()
		}()

		wg.Wait()
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.server.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
