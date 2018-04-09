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
	"strconv"

	"github.com/rs/zerolog/log"
	tt "github.com/servicelab/tamtam/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// heartbeatCmd represents the heartbeat command
var heartbeatCmd = &cobra.Command{
	Use:   "heartbeat [millis]",
	Short: "Sets the heartbeat in millis of the agent",
	Long:  "Heartbeat should be at least 10 milliseconds",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires the heartbeat time in milliseconds")
		}

		i, err := strconv.Atoi(args[0])
		if err != nil || i < 10 || i > 1000 {
			return errors.New("milliseconds should be a number larger than 10")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(viper.GetString("rpc"), grpc.WithInsecure())
		if err != nil {
			log.Fatal().Msgf("did not connect to RPC server: %v", err)
		}
		defer conn.Close()
		c := tt.NewTamTamClient(conn)
		millis, _ := strconv.Atoi(args[0])

		// Set the new heartbeat on the server
		_, err = c.SetHeartbeat(context.Background(), &tt.Heartbeat{Millis: int32(millis)})
		if err != nil {
			log.Fatal().Msgf("could not set the new heartbeat: %v", err)
		}
		fmt.Printf("Set the heartbeat of the agent to %d milliseconds\n", millis)
	},
}

func init() {
	configureCmd.AddCommand(heartbeatCmd)
}
