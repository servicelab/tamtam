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

	tt "github.com/eelcocramer/tamtam/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	levels = []string{"all", "trace", "debug", "info", "warn", "error", "fatal", "off"}
)

// thresholdCmd represents the threshold command
var thresholdCmd = &cobra.Command{
	Use:   "threshold [threshold]",
	Short: "Sets the threshold of the console output logging of the agent",
	Long: `Sets the threshold of the logging output. The threshold should be 
one of the following values:

* all
* trace
* debug
* info
* warn
* error
* fatal
* off
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a log level to set")
		}

		if exsist, _ := inArray(args[0], levels); !exsist {
			return errors.New("invalid argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("threshold called: " + args[0])
		conn, err := grpc.Dial(cfg.RPCAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)

		var l tt.LogLevel_Level
		switch args[0] {
		case "all":
			l = tt.LogLevel_ALL
		case "trace":
			l = tt.LogLevel_TRACE
		case "debug":
			l = tt.LogLevel_DEBUG
		case "info":
			l = tt.LogLevel_INFO
		case "warn":
			l = tt.LogLevel_WARN
		case "error":
			l = tt.LogLevel_ERROR
		case "fatal":
			l = tt.LogLevel_FATAL
		case "off":
			l = tt.LogLevel_OFF
		default:
			log.Fatalf("Unexpected log level: %s\n", args[0])
		}

		// Set the new heartbeat on the server
		_, err = c.SetLogThreshold(context.Background(), &tt.LogLevel{Level: l})
		if err != nil {
			log.Fatalf("could not set the new log threshold: %v", err)
		}
		fmt.Printf("Set the log treshold of the agent to %s\n", args[0])
	},
}

func init() {
	configureCmd.AddCommand(thresholdCmd)
}
