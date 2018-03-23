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
	"bufio"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	tt "github.com/eelcocramer/tamtam/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	encoding  string
	encodings = []string{"string", "base64", "bytes"}
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Streams broadcast events that are received from other nodes in the network",
	Long: `Streams broadcast messages that are received from other nodes. The encoding type
of the messages can be specified on the command line using the 'encoding' flag. The following
encodings are valid:

* string [default]
* base64
* bytes
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if exists, _ := inArray(encoding, encodings); exists {
			return nil
		}

		return errors.New("invalid value for encoding")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var stdout bufio.Writer
		conn, err := grpc.Dial(cfg.RPCAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)
		stream, err := c.Stream(context.Background(), &tt.Empty{})
		if err != nil {
			log.Fatalf("%v.Stream(_) = _, %v", c, err)
		}

		waitc := make(chan struct{})
		if encoding == "bytes" {
			stdout := bufio.NewWriter(os.Stdout)
			defer stdout.Flush()
		}

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
				switch encoding {
				case "string":
					fmt.Printf("%s\n", string(in.Bytes))
				case "base64":
					fmt.Printf("%s\n", b64.StdEncoding.EncodeToString(in.Bytes))
				case "bytes":
					stdout.Write(in.Bytes)
				}
			}
		}()
		<-waitc
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.Flags().StringVar(&encoding, "encoding", "string", "Encoding to output broadcast messages with")
}
