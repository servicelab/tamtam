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
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/clockworksoul/smudge"
	tt "github.com/eelcocramer/tamtam/service"
	util "github.com/eelcocramer/tamtam/util"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	bind string
	port int
	hbm  int
)

type server struct{}

// Join adds remote node to the network
func (s *server) Join(ctx context.Context, in *tt.NodeAddress) (*tt.Response, error) {
	node, err := smudge.CreateNodeByIP(net.ParseIP(in.IP), uint16(in.Port))
	if err == nil {
		node, err = smudge.AddNode(node)
		if err == nil {
			log.Printf("Added node: %s\n", node.Address())
			return &tt.Response{Code: tt.Response_OK}, nil
		}
	}
	return &tt.Response{Code: tt.Response_ERROR}, err
}

// Leave removes a remote node from the network
func (s *server) Leave(ctx context.Context, in *tt.NodeAddress) (*tt.Response, error) {
	node, err := smudge.CreateNodeByIP(net.ParseIP(in.IP), uint16(in.Port))
	if err == nil {
		node, err = smudge.RemoveNode(node)
		if err == nil {
			log.Printf("Removed node: %s\n", node.Address())
			return &tt.Response{Code: tt.Response_OK}, nil
		}
	}
	return &tt.Response{Code: tt.Response_ERROR}, err
}

// Broadcast sends an arbitrary data message to other
// healthy nodes in the network.
func (s *server) Broadcast(ctx context.Context, in *tt.Message) (*tt.Response, error) {
	err := smudge.BroadcastBytes(in.Bytes)
	if err != nil {
		return &tt.Response{Code: tt.Response_ERROR}, err
	}
	log.Printf("Send broadcast to the network\n")
	return &tt.Response{Code: tt.Response_OK}, nil
}

// Stream creates a stream with arbitrary broadcast messages
// received from other nodes in the network.
func (s *server) Stream(in *tt.Empty, stream tt.TamTam_StreamServer) error {
	ch := make(chan []byte)
	ctx := stream.Context()
	util.AddBroadcastChannel(ctx, ch)
	defer util.RemoveBroadcastChannel(ctx)
	defer log.Printf("Broadcast listener went away\n")
	for {
		select {
		case v := <-ch:
			if err := stream.Send(&tt.Message{Bytes: v}); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// Monitor creates a stream with network status changes.
func (s *server) Monitor(in *tt.Empty, stream tt.TamTam_MonitorServer) error {
	ch := make(chan *tt.Node)
	ctx := stream.Context()
	util.AddMonitorChannel(ctx, ch)
	defer util.RemoveMonitorChannel(ctx)
	defer log.Printf("Monitor listener went away\n")
	for {
		select {
		case v := <-ch:
			if err := stream.Send(v); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// NodeList creates a list containing information about
// other nodes in the network.
func (s *server) Nodes(ctx context.Context, in *tt.Empty) (*tt.NodeList, error) {
	all := smudge.AllNodes()
	nodes := make([]*tt.Node, 0, len(all))
	for _, n := range all {
		nodes = append(nodes, util.SmudgeToTamTamNode(n))
	}
	return &tt.NodeList{Nodes: nodes}, nil
}

// Healthy creates a list containing information about
// all healthy nodes in the network
func (s *server) Healthy(ctx context.Context, in *tt.Empty) (*tt.NodeList, error) {
	healthy := smudge.HealthyNodes()
	nodes := make([]*tt.Node, 0, len(healthy))
	for _, n := range healthy {
		nodes = append(nodes, util.SmudgeToTamTamNode(n))
	}
	return &tt.NodeList{Nodes: nodes}, nil
}

// SetHeartbeat changes the frequency this nodes sends out
// heartbeat messages to the network.
func (s *server) SetHeartbeat(ctx context.Context, in *tt.Heartbeat) (*tt.Response, error) {
	if in.Millis < 10 && in.Millis > 1000 {
		return &tt.Response{Code: tt.Response_ERROR}, errors.New("heartbeat should be between 10 and 1000")
	}
	smudge.SetHeartbeatMillis(int(in.Millis))
	log.Printf("Heartbeat has been changed to %d\n", in.Millis)
	return &tt.Response{Code: tt.Response_OK}, nil
}

// SetLogThreshold changes the log level the network agent output logs to the console
func (s *server) SetLogThreshold(ctx context.Context, in *tt.LogLevel) (*tt.Response, error) {
	switch in.Level {
	case tt.LogLevel_ALL:
		smudge.SetLogThreshold(smudge.LogAll)
	case tt.LogLevel_TRACE:
		smudge.SetLogThreshold(smudge.LogTrace)
	case tt.LogLevel_DEBUG:
		smudge.SetLogThreshold(smudge.LogDebug)
	case tt.LogLevel_INFO:
		smudge.SetLogThreshold(smudge.LogInfo)
	case tt.LogLevel_WARN:
		smudge.SetLogThreshold(smudge.LogWarn)
	case tt.LogLevel_ERROR:
		smudge.SetLogThreshold(smudge.LogError)
	case tt.LogLevel_FATAL:
		smudge.SetLogThreshold(smudge.LogFatal)
	case tt.LogLevel_OFF:
		smudge.SetLogThreshold(smudge.LogOff)
	default:
		return nil, errors.New("unknown loglevel")
	}
	return &tt.Response{Code: tt.Response_OK}, nil
}

// Ping sends a ping message to another node in the network.
func (s *server) Ping(ctx context.Context, in *tt.NodeAddress) (*tt.Response, error) {
	ip := net.ParseIP(in.IP)
	if ip == nil {
		return nil, errors.New("not a valid ip address to ping")
	}
	n, err := smudge.CreateNodeByIP(net.ParseIP(in.IP), uint16(in.Port))
	if err != nil {
		return nil, err
	}
	err = smudge.PingNode(n)
	if err != nil {
		return nil, err
	}
	return &tt.Response{Code: tt.Response_OK}, nil
}

// LocalAddress returns the address of the local network agent.
func (s *server) LocalAddress(ctx context.Context, in *tt.Empty) (*tt.NodeAddress, error) {
	ip := smudge.GetListenIP()
	port := smudge.GetListenPort()
	return &tt.NodeAddress{IP: ip.String(), Port: uint32(port)}, nil
}

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Starts a TamTam agent.",
	Long: `Starts a TamTam agent that keeps a connection to the gossip network
and listens to RPC command on the RCP interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", cfg.RPCAddr)
		if err != nil {
			log.Fatalf("error starting the rpc server, failed to listen: %v", err)
		}
		s := grpc.NewServer()
		tt.RegisterTamTamServer(s, &server{})
		reflection.Register(s)
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Failed to start the rpc server: %v", err)
			}
		}()

		// configure smudge
		ip := net.ParseIP(bind)
		smudge.SetListenIP(ip)
		smudge.SetLogThreshold(smudge.LogWarn)
		if cfg.Verbose {
			smudge.SetLogThreshold(smudge.LogDebug)
		}
		if cfg.Trace {
			smudge.SetLogThreshold(smudge.LogTrace)
		}
		smudge.SetListenPort(port)
		smudge.SetHeartbeatMillis(hbm)

		if ip.To4() == nil {
			smudge.SetMaxBroadcastBytes(512)
		}
		smudge.SetMulticastEnabled(false)
		smudge.SetClusterName("tamtam")

		go func() {
			log.Printf("Listening for gRPC at %s\n", cfg.RPCAddr)
			if ip.To4() != nil {
				log.Printf("Listening for gossip at %s:%d\n", bind, port)
			} else {
				log.Printf("Listening for gossip at [%s]:%d\n", bind, port)
			}
			log.Println("Running agent, press CRTL-C to abort...")
			smudge.Begin()
		}()
		// Handle SIGINT and SIGTERM.
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-quit)
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().IntVarP(&port, "port", "p", smudge.GetListenPort(), "list port for the gossip network")
	agentCmd.Flags().StringVarP(&bind, "bind", "b", "127.0.0.1", "listen address for the gossip network")
	agentCmd.Flags().IntVar(&hbm, "heartbeat", smudge.GetHeartbeatMillis(), "heartbeat used within the gossip network")
}
