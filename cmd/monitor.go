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
	"io"
	"log"
	"os"
	"text/tabwriter"

	tt "github.com/eelcocramer/tamtam/service"
	"github.com/eelcocramer/tamtam/util"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitors the network status",
	Long:  `Monitor logs network changes that occur in the gossip network`,
	Run: func(cmd *cobra.Command, args []string) {
		w := new(tabwriter.Writer)
		// Format in tab-separated columns with a tab stop of 8.
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "NODE\tAGE\tEMIT #\tMILLIS\tSTATUS")

		conn, err := grpc.Dial(cfg.RPCAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)
		stream, err := c.Monitor(context.Background(), &tt.Empty{})
		if err != nil {
			log.Fatalf("%v.Stream(_) = _, %v", c, err)
		}

		waitc := make(chan struct{})
		go func() {
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					// read done.
					close(waitc)
					return
				}
				if err != nil {
					log.Fatalf("Failed to receive a broadcast: %v", err)
				}
				fmt.Fprintf(w, "[%s]:%d \t%d\t%d\t%d\t%s\t\n", in.Address.IP,
					in.Address.Port, in.Age, in.EmitCounter, in.PingMillis,
					util.StatusToString(in.Status))
				w.Flush()
			}
		}()
		<-waitc
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
