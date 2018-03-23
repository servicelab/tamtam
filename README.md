```
 _                      _
| |_   __ _  _ __ ___  | |_   __ _  _ __ ___
| __| / _` || '_ ` _ \ | __| / _` || '_ ` _ \
| |_ | (_| || | | | | || |_ | (_| || | | | | |
 \__| \__,_||_| |_| |_| \__| \__,_||_| |_| |_|
```

[![Go Report Card](https://goreportcard.com/badge/github.com/eelcocramer/tamtam)](https://goreportcard.com/report/github.com/eelcocramer/tamtam)

## Introduction

TamTam is a convenience layer that provides a command line and [gRPC](https://grpc.io) interface to [Smudge](https://github.com/clockworksoul/smudge).

TamTam allows you to create a decentralized cluster on top of an IP network (both IPv4 and IPv6 are supported). It has a similar but simplified function as [Serf](https://www.serf.io).

TamTam uses the [Smudge](https://github.com/clockworksoul/smudge) library that provides group member discovery, status dissemination and failure detection to the cluster. Smudge implements the SWIM protocol [(paper)](https://pdfs.semanticscholar.org/8712/3307869ac84fc16122043a4a313604bd948f.pdf) [(explained 1)](https://asafdav2.github.io/2017/swim-protocol/) [(explained 2)](https://prakhar.me/articles/swim/).


### Command line interface

TamTam has the following options and commands:

#### Options

```
      --config string   Config file (default is $HOME/.tamtam.yaml)
  -h, --help            help for tamtam
      --nobanner        disables printing the banner
      --rpc string      The RPC address the agent binds to or other commands query a running agent on (default "localhost:6262")
      --trace           turn on trace logging
      --verbose         turn on verbose logging
```

#### Commands
* [tamtam address](docs/tamtam_address.md)	 - Gets the local bind address of a running agent
* [tamtam agent](docs/tamtam_agent.md)	 - Starts a TamTam agent.
* [tamtam broadcast](docs/tamtam_broadcast.md)	 - Broadcasts a message to other nodes in the network
* [tamtam configure](docs/tamtam_configure.md)	 - Configure a running agent
* [tamtam gendoc](docs/tamtam_gendoc.md)	 - Generates markdown docs for TamTam
* [tamtam healthy](docs/tamtam_healthy.md)	 - Displays the healthy nodes in the network
* [tamtam join](docs/tamtam_join.md)	 - joins a node to the network
* [tamtam leave](docs/tamtam_leave.md)	 - Removes a node from the network
* [tamtam monitor](docs/tamtam_monitor.md)	 - Monitors the network status
* [tamtam nodes](docs/tamtam_nodes.md)	 - Displays a list of known nodes in the network
* [tamtam ping](docs/tamtam_ping.md)	 - Pings a remote node in the network
* [tamtam stream](docs/tamtam_stream.md)	 - Streams broadcast events that are received from other nodes in the network
* [tamtam version](docs/tamtam_version.md)	 - Prints the version number

### gRPC interaface

The gRPC protocol is described in [the protocol definition](service/service.proto). To build an RPC client that connects to the TamTam [agent](docs/tamtam_agent.md) please refer to the gRPC documentation for you programming language or check the source code of TamTam for golang examples. All TamTam commands, with the exception of `agent` and `gendoc` use a gRPC client to connect to the agent.

