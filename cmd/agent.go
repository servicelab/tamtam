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
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/clockworksoul/smudge"
	"github.com/rs/zerolog/log"
	tt "github.com/servicelab/tamtam/service"
	util "github.com/servicelab/tamtam/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	bind              string
	port              int
	hbm               int
	multicast         bool
	multicastAddress  string
	multicastPort     int
	multicastInterval int
	clusterName       string
	ipv6              bool
	max               int
)

type server struct{}

// Join adds remote node to the network
func (s *server) Join(ctx context.Context, in *tt.NodeAddress) (*tt.Response, error) {
	node, err := smudge.CreateNodeByIP(net.ParseIP(in.IP), uint16(in.Port))
	if err == nil {
		node, err = smudge.AddNode(node)
		if err == nil {
			log.Info().Msgf("Added node: %s", node.Address())
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
			log.Info().Msgf("Removed node: %s", node.Address())
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
	log.Info().Msg("Send broadcast to the network")
	return &tt.Response{Code: tt.Response_OK}, nil
}

// Stream creates a stream with arbitrary broadcast messages
// received from other nodes in the network.
func (s *server) Stream(in *tt.Empty, stream tt.TamTam_StreamServer) error {
	ch := make(chan []byte)
	ctx := stream.Context()
	util.AddBroadcastChannel(ctx, ch)
	defer util.RemoveBroadcastChannel(ctx)
	defer log.Info().Msg("Broadcast listener went away")
	for {
		select {
		case v := <-ch:
			log.Debug().Msgf("Streaming %d bytes to subscriber", len(v))
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
	defer log.Info().Msg("Monitor listener went away")
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
	log.Info().Msgf("Heartbeat has been changed to %d", in.Millis)
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

// Get preferred outbound ip of this machine
func getOutboundIP(ipv6 bool) (string, error) {
	var (
		conn net.Conn
		err  error
	)

	if ipv6 {
		conn, err = net.Dial("udp", "2001:4860:4860::8888")
	} else {
		conn, err = net.Dial("udp", "8.8.8.8:80")
	}
	if err != nil {
		return "", err
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Starts a TamTam agent.",
	Long: `Starts a TamTam agent that keeps a connection to the gossip network
and listens to RPC command on the RCP interface.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if ipv6 {
			if bind != "0.0.0.0" {
				return errors.New("--ipv6 flag can only be used with the default bind address")
			}
			bind = "::"
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", viper.GetString("rpc"))
		if err != nil {
			log.Fatal().Msgf("error starting the rpc server, failed to listen: %v", err)
		}
		s := grpc.NewServer()
		tt.RegisterTamTamServer(s, &server{})
		reflection.Register(s)
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatal().Msgf("Failed to start the rpc server: %v", err)
			}
		}()

		// configure smudge
		if bind == "0.0.0.0" {
			bind, err = getOutboundIP(ipv6)
			if err != nil {
				log.Fatal().Msgf("Failed to determine bind address: %v", err)
			}
		}
		ip := net.ParseIP(bind)
		if ip == nil {
			log.Fatal().Msgf("Failed to parse ip address: %s", bind)
		}
		smudge.SetListenIP(ip)
		smudge.SetLogger(util.SmudgeLogger{})
		smudge.SetLogThreshold(smudge.LogWarn)
		if viper.GetBool("verbose") {
			smudge.SetLogThreshold(smudge.LogDebug)
		}
		if viper.GetBool("trace") {
			smudge.SetLogThreshold(smudge.LogTrace)
		}
		smudge.SetListenPort(port)
		smudge.SetHeartbeatMillis(hbm)

		if max == 0 {
			if ipv6 {
				max = 512
			} else {
				max = 256
			}
		} else {
			log.Info().Msgf("Setting the maximum size for a broadcast message to %d bytes", max)
		}
		smudge.SetMaxBroadcastBytes(max)

		smudge.SetMulticastEnabled(multicast)
		smudge.SetMulticastPort(multicastPort)
		log.Debug().Msgf("multicast interval = %d", multicastInterval)
		smudge.SetMulticastAnnounceIntervalSeconds(multicastInterval)
		if multicastAddress != "" {
			smudge.SetMulticastAddress(multicastAddress)
		}
		smudge.SetClusterName(clusterName)

		go func() {
			log.Info().Msgf("Listening for gRPC at %s", viper.GetString("rpc"))
			if ip.To4() != nil || strings.HasPrefix(bind, "[") {
				log.Info().Msgf("Listening for gossip at %s:%d", bind, port)
			} else {
				log.Info().Msgf("Listening for gossip at [%s]:%d", bind, port)
			}
			log.Info().Msg("Running agent, press CRTL-C to abort...")
			smudge.Begin()
		}()
		// Handle SIGINT and SIGTERM.
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().IntVarP(&port, "port", "p", smudge.GetListenPort(), "list port for the gossip network")
	agentCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "listen address for the gossip network")
	agentCmd.Flags().BoolVarP(&ipv6, "ipv6", "6", false, "alias for -b [::], listens to all IPv6 interfaces")
	agentCmd.Flags().IntVar(&hbm, "heartbeat", smudge.GetHeartbeatMillis(), "heartbeat used within the gossip network")
	agentCmd.Flags().IntVar(&max, "max", 0, "maximum size for broadcast messages, set to 0 means 256 for IPv4 and 512 for IPv6")
	agentCmd.Flags().BoolVar(&multicast, "multicast", false, "enable multicast node discovery")
	agentCmd.Flags().StringVar(&clusterName, "clustername", "tamtam", "name for the multicast cluster")
	agentCmd.Flags().StringVar(&multicastAddress, "multicast-address", "", "address for multicast discovery messages defaults to 224.0.0.0 or [ff02::1]")
	agentCmd.Flags().IntVar(&multicastInterval, "multicast-interval", 0, "seconds between mutlicast disovery messages")
	agentCmd.Flags().IntVar(&multicastPort, "multicast-port", 9998, "port to listen for multicast discovery messages")
}
