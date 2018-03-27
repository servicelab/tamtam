/*
Copyright 2018, Eelco Cramer and the TamTam contributors

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
	"github.com/eelcocramer/tamtam/config"

	"github.com/CrowdSurge/banner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
	"os"
)

var cfgFile string

func init() {
	cfg = config.Manager{}
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfg.File, "config", "", "Config file (default is $HOME/.tamtam.yaml)")
	rootCmd.PersistentFlags().BoolVar(&cfg.Verbose, "verbose", false, "turn on verbose logging")
	rootCmd.PersistentFlags().BoolVar(&cfg.Trace, "trace", false, "turn on trace logging")
	rootCmd.PersistentFlags().StringVar(&cfg.RPCAddr, "rpc", "localhost:6262", "The RPC address the agent binds to or other commands query a running agent on")
	rootCmd.PersistentFlags().BoolVar(&cfg.NoBanner, "nobanner", false, "disables printing the banner")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("trace", rootCmd.PersistentFlags().Lookup("trace"))
	viper.BindPFlag("rpc", rootCmd.PersistentFlags().Lookup("rpc"))
	viper.BindPFlag("nobanner", rootCmd.PersistentFlags().Lookup("nobanner"))
	viper.SetDefault("verbose", false)
	viper.SetDefault("trace", false)
	viper.SetDefault("nobanner", false)
	viper.SetDefault("rpc", "localhost:6262")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfg.File != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfg.File)
	}
	viper.SetConfigName(".tamtam")         // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME")) // adding home directory as first search path
	viper.AutomaticEnv()                   // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tamtam",
	Short: "Gossip network",
	Long:  `Gossip network agent and RPC interface to interact with the network.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if !cfg.NoBanner {
		fmt.Fprintln(os.Stderr, banner.PrintS("tamtam"))
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var cfg config.Manager
