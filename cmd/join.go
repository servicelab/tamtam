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

	tt "github.com/eelcocramer/tamtam/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join [address:port]",
	Short: "joins a node to the network",
	Long: `Joins a remote node to the network. The bind address should
be either a IPv4 or IPv6 address with a port. For example:

* IPv4 - 127.0.0.1:9999
* IPv6 - [::1]:9999`,
	Args: validateHostPortArg,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(cfg.RPCAddr, grpc.WithInsecure())
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
		_, err = c.Join(context.Background(), &tt.NodeAddress{IP: host, Port: uint32(port)})
		if err != nil {
			log.Fatalf("could not join the host: %v", err)
		}
		fmt.Printf("Host was joined.\n")
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
