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

	tt "github.com/eelcocramer/tamtam/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// addressCmd represents the address command
var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "Gets the local bind address of a running agent",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(cfg.RPCAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)
		address, err := c.LocalAddress(context.Background(), &tt.Empty{})
		if err != nil {
			log.Fatalf("could get the address of the agent: %v", err)
		}
		fmt.Printf("The local address of the agent is: [%s]:%d\n", address.IP, address.Port)
	},
}

func init() {
	rootCmd.AddCommand(addressCmd)
}
