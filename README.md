```
 _                      _
| |_   __ _  _ __ ___  | |_   __ _  _ __ ___
| __| / _` || '_ ` _ \ | __| / _` || '_ ` _ \
| |_ | (_| || | | | | || |_ | (_| || | | | | |
 \__| \__,_||_| |_| |_| \__| \__,_||_| |_| |_|
```

[![pipeline status](https://gitlab.com/eelco/tamtam/badges/master/pipeline.svg)](https://gitlab.com/eelco/tamtam/pipelines)
[![Go Report Card](https://goreportcard.com/badge/github.com/servicelab/tamtam)](https://goreportcard.com/report/github.com/servicelab/tamtam)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/39f89f19b7c7405ab9f6e9a15de19cb5)](https://www.codacy.com/app/eelcocramer/tamtam?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=servicelab/tamtam&amp;utm_campaign=Badge_Grade)
[![codebeat badge](https://codebeat.co/badges/6c1251ca-1aec-4ca3-92af-3cbf1afd2e65)](https://codebeat.co/projects/github-com-servicelab-tamtam-master)

## Introduction

TamTam is a utility that provides a command line and [gRPC](https://grpc.io) interface to [Smudge](https://github.com/clockworksoul/smudge).

TamTam allows you to create a decentralized cluster on top of an IP network (both IPv4 and IPv6 are supported). It has a similar but simplified function as [Serf](https://www.serf.io).

TamTam uses the [Smudge](https://github.com/clockworksoul/smudge) library that provides group member discovery, status dissemination and failure detection to the cluster. Smudge implements the SWIM protocol [(paper)](https://pdfs.semanticscholar.org/8712/3307869ac84fc16122043a4a313604bd948f.pdf) [(explained 1)](https://asafdav2.github.io/2017/swim-protocol/) [(explained 2)](https://prakhar.me/articles/swim/).

## Getting the binary

Binaries are automatically being build by the CI/CD pipeline on each release. The `master` branch always has the latest release. When the build pipeline passed successfully you can download the binaries via the links below:

* [MacOS amd64](https://gitlab.com/eelco/tamtam/-/jobs/artifacts/master/download?job=macos)
* [Windows amd64](https://gitlab.com/eelco/tamtam/-/jobs/artifacts/master/download?job=windows)
* [Linux amd64](https://gitlab.com/eelco/tamtam/-/jobs/artifacts/master/download?job=linux-amd64)
* [Linux i386](https://gitlab.com/eelco/tamtam/-/jobs/artifacts/master/download?job=linux-386)
* [Linux arm32v7](https://gitlab.com/eelco/tamtam/-/jobs/artifacts/master/download?job=linux-arm)
* [Linux arm64v8](https://gitlab.com/eelco/tamtam/-/jobs/artifacts/master/download?job=linux-arm64)

You can also build your binaries from source.

## Command line interface

TamTam has the following options and commands:

### Options

```
      --config string   Config file (default is $HOME/.tamtam.yaml)
  -h, --help            help for tamtam
      --nobanner        disables printing the banner
      --rpc string      The RPC address the agent binds to or other commands query a running agent on (default "localhost:6262")
      --trace           turn on trace logging
      --verbose         turn on verbose logging
```

### Commands
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

## Usage example

In the example below we will run 2 agents that connect to each other. We will demonstrate how to connect to the broadcast stream and how to broadcast a message in the network.

Open a terminal window and issue the following command to start a TamTam agent:

```bash
./tamtam agent
```

You should now see logging indicating that TamTam is running and a warning that it detected a new node (namely itself). Now lets start another agent running on another port. Open another terminal window and execute the following command:

```bash
./tamtam --rpc localhost:6363 agent -p 9998
```

When the second agent is also running you can ask one of the agents to join the other and after that to start listening to broadcast messages. Open another terminal and issue the following commands:

```bash
./tamtam join 127.0.0.1:9998
./tamtam stream
```

The first command should give you a message that the remote host was joined. The second command is waiting for broadcast messages to display. Lets send a broadcast message... Open another terminal and execute the following command:

```bash
./tamtam --rpc localhost:6363 broadcast "Hello world!"
```

The broadcast message should now appear `stream` process at the 3rd terminal window.

The next step is to bind the messages to your own process of script. This can be done by using the gRPC interface or by binding the broadcast stream to a script. Take for example the following `handler.sh` bash script:

```bash
#!/bin/bash
while read msg; do
    printf "${msg}\n"
done
```

You can bind this script to TamTam using the following command:

```bash
./tamtam stream | ./handler.sh
```

### Demo

[![asciicast](https://asciinema.org/a/1XiYR17W8MZsEI0IUM0oBZWJO.png)](https://asciinema.org/a/1XiYR17W8MZsEI0IUM0oBZWJO?speed=2&theme=solarized-dark)

## Configuration

TamTam can be configured using a configuration file, environment variables or command line. Note that command line takes precedence over the environment variables which in their turn take precedence over the values from the configuration file.

The default location for the configuration file is `$HOME/.tamtam.yml` but you can also specify the location via the command line option `--config`. Currently only a few global options can be configured using the environment or configuration file.

Example of the configuration file:

```yml
nobanner: false
verbose: false
trace: false
rpc: "localhost:6262"
```

The following environment variables are supported:

```
Variable        | Description
--------------- | -----------------------------------------------------------------------
TAMTAM_NOBANNER | When set TamTam won't print out the banner on startup
TAMTAM_VERBOSE  | Set to enable verbose logging
TAMTAM_TRACE    | Set to enable trace logging
TAMTAM_RPC      | The host and port to bind the RPC interface to. Default: localhost:6262
```
## gRPC interaface

The gRPC protocol is described in [the protocol definition](service/service.proto). To build an RPC client that connects to the TamTam [agent](docs/tamtam_agent.md) please refer to the gRPC documentation for you programming language or check the source code of TamTam for golang examples. All TamTam commands, with the exception of `agent` and `gendoc` use a gRPC client to connect to the agent.

## Building from source

If you feel lucky you can try `go get -u github.com/servicelab/tamtam`

Otherwise, you can setup the build environment by executing the following script:

```bash
cd $GOPATH
git clone https://github.com/servicelab/tamtam src/github.com/servicelab/tamtam
cd src/github.com/servicelab/tamtam
make # builds TamTam for your current architecture
```

You will find the TamTam binary the current folder.

