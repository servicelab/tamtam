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
	"fmt"
	"log"
	"net"
	"strconv"

	tt "github.com/servicelab/tamtam/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// leaveCmd represents the leave command
var leaveCmd = &cobra.Command{
	Use:   "leave [address:port]",
	Short: "Removes a node from the network",
	Long: `Removes a remote node from the network. The remote node address should
be either a IPv4 or IPv6 address with a port. For example:

* IPv4 - 127.0.0.1:9000
* IPv6 - [::1]:9000`,
	Args: validateHostPortArg,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(viper.GetString("rpc"), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)
		host, p, err := net.SplitHostPort(args[0])
		if err == nil {
			port, err = strconv.Atoi(p)
		}
		if err != nil {
			log.Fatal("Should not happen")
		}

		// Add the host
		_, err = c.Leave(context.Background(), &tt.NodeAddress{IP: host, Port: uint32(port)})
		if err != nil {
			log.Fatalf("could not leave the host: %v", err)
		}
		fmt.Printf("Host was left.\n")
	},
}

func init() {
	rootCmd.AddCommand(leaveCmd)
}
