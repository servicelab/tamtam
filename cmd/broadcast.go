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
	"fmt"
	"io"
	"log"
	"os"

	tt "github.com/eelcocramer/tamtam/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	base64 bool
)

// broadcastCmd represents the broadcast command
var broadcastCmd = &cobra.Command{
	Use:   "broadcast [message]",
	Short: "Broadcasts a message to other nodes in the network",
	Long: `Sends a broadcast message to other nodes in the gossip network.
The message is send in bytes and may be specified either as a string, a base64 encoded
byte array, that will be converted to bytes before sending, or the input will be taken from
stdin, in which case the message will only be send after an EOF character is received.`,
	Run: func(cmd *cobra.Command, args []string) {
		var bytes []byte
		if len(args) == 0 {
			// no args getting data from stdin
			r := bufio.NewReader(os.Stdin)
			bytes = make([]byte, 0, 2048)
			l := 0
			for {
				n, err := r.Read(bytes[l : cap(bytes)-l])
				l = l + n
				bytes = bytes[:l]
				if n == 0 {
					if err == io.EOF {
						break
					}
					log.Fatal("Message to large")
				}
				if l > 1024 {
					log.Fatal("Message to large")
				}
			}
		} else {
			bytes = []byte(args[0])
		}

		// decoding the base64 string
		if base64 {
			data, err := b64.StdEncoding.DecodeString(string(bytes))
			if err != nil {
				log.Fatalf("error decoding the base64 string: %v", err)
			}
			bytes = data
		}

		conn, err := grpc.Dial(viper.GetString("rpc"), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)

		// Send message
		_, err = c.Broadcast(context.Background(), &tt.Message{Bytes: bytes})
		if err != nil {
			log.Fatalf("could not send message: %v", err)
		}
		fmt.Printf("broadcast message was send: %s\n", string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)
	broadcastCmd.Flags().BoolVar(&base64, "base64", false, "string argument is base64 encoded")
}
