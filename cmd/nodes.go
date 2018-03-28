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
	"os"
	"text/tabwriter"

	tt "github.com/servicelab/tamtam/service"
	"github.com/servicelab/tamtam/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Displays a list of known nodes in the network",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(viper.GetString("rpc"), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)

		// Send message
		nl, err := c.Nodes(context.Background(), &tt.Empty{})
		if err != nil {
			log.Fatalf("could not get node list: %v", err)
		}
		w := new(tabwriter.Writer)
		// Format in tab-separated columns with a tab stop of 8.
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "NODE\tAGE\tEMIT #\tMILLIS\tSTATUS")
		for _, n := range nl.Nodes {
			fmt.Fprintf(w, "[%s]:%d \t%d\t%d\t%d\t%s\t\n", n.Address.IP,
				n.Address.Port, n.Age, n.EmitCounter, n.PingMillis,
				util.StatusToString(n.Status))
		}
		fmt.Fprintln(w, "")
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(nodesCmd)
}
