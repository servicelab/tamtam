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
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	tt "github.com/eelcocramer/tamtam/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping [address:port]",
	Short: "Pings a remote node in the network",
	Long: `Pings a remote node in the network. The address should
be either a IPv4 or IPv6 address with a port. For example:

* IPv4 - 127.0.0.1:9999
* IPv6 - [::1]:9999`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires the remote binding address")
		}

		_, _, err := net.SplitHostPort(args[0])
		if err != nil {
			return errors.New("error parsing the remote binding address")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip, ps, err := net.SplitHostPort(args[0])
		if err != nil {
			log.Fatalf("should not happen: %v", err)
		}
		port, err := strconv.Atoi(ps)
		if err != nil {
			log.Fatalf("should not happen: %v", err)
		}
		conn, err := grpc.Dial(cfg.RPCAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)
		_, err = c.Ping(context.Background(), &tt.NodeAddress{IP: ip, Port: uint32(port)})
		if err != nil {
			log.Fatalf("could not ping the node: %v", err)
		}
		fmt.Printf("Send a ping message to node %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
