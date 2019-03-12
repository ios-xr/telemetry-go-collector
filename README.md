## IOS XR Telemetry collector

A Simple "Go" based collector for IOS XR Telemetry that can get you started on trying telemetry. IOS XR supports dialout and dialin modes for connecting to router for streaming, this repo will have both the collectors. These collectors only read the data as its sent from the router and dump on to stdout, you can add on to these scripts to do many fancy things. This is meant for beginners, if you are familiar with Go, GRPC and IOSXR, you should be using and modifying/extending [Pipeline](https://github.com/cisco-ie/bigmuddy-network-telemetry-pipeline) for your needs.

If you want to start from scratch or do not have "Go", Protoc, protoc-gen-go or grpc installed, you can check out [Dialout-collector-howto.md](Dialout-collector-howto.md). It has instructions on how to get started and how to write a simple collector.

Assuming "Go" is already installed, following instructions are for getting collector, building it and running it.

Note:
Dialout collector supports grpc, tcp and udp transports
Dialin Collector supports subscribe and get-proto RPCs to IOSXR device

#### Install instructions:
`go get -d github.com/ios-xr/telemetry-go-collector`

alternately, use git clone to get the collector to $GOPATH/src directory

`git clone github.com/ios-xr/telemetry-go-collector $GOPATH/src`

If git clone is used, change the import of mdt_grpc_dialout in telemetry_dialout_collector.go to make sure, path is correct relative to $GOPATH/src

#### Requirements
Need following to be available to be able to build the collectors, binaries can be used as is on linux
* go
* protoc
* grpc

Install instructions are present in [Dialout-collector-howto.md](Dialout-collector-howto.md)

### MDT Dialout Collector:
##### Build
`go build -o bin/telemetry_dialout_collector github.com/ios-xr/telemetry-go-collector/telemetry_dialout_collector`

prebuilt binary can be used from bin/telemetry_dialout_collector on Linux.

##### Run
```
 $ ./bin/telemetry_dialout_collector -h
Usage: ./bin/telemetry_dialout_collector [options]
  -decode_raw
        Use protoc --decode_raw
  -dont_clean
        Don't remove tmp files on exit
  -encoding string
        expected encoding, Options: json,self-describing-gpb,gpb, needed only for grpc (default "json")
  -out string
        output file to write to
  -port int
        The server port to listen on (default 57400)
  -proto string
        proto file to use for decode
  -transport string
        transport to use, grpc, tcp or udp (default "grpc")
Examples:
GRPC Server:                             ./bin/telemetry_dialout_collector -port <> -encoding gpb
TCP Server:                              ./bin/telemetry_dialout_collector -port <> -transport tcp
GRPC use protoc to decode:               ./bin/telemetry_dialout_collector -port <> -encoding gpb -proto cdp_neighbor.proto
GRPC use protoc to decode without proto: ./bin/telemetry_dialout_collector -port <> -encoding gpb -decode_raw
 $
 ```

![](docs/dialout-build.gif)

### MDT Dialin Collector:
##### Build
`go build -o bin/telemetry_dialin_collector github.com/ios-xr/telemetry-go-collector/telemetry_dialin_collector`

prebuilt binary can be used from bin/telemetry_dialin_collector on Linux.

##### Run
```
 $ ./bin/telemetry_dialin_collector -h
Usage: ./bin/telemetry_dialin_collector [options]
  -decode_raw
        Use protoc --decode_raw
  -dont_clean
        Don't remove tmp files on exit
  -encoding string
        encoding to use, Options: json,self-describing-gpb,gpb (default "json")
  -oper string
        Operation: subscribe, get-proto (default "subscribe")
  -out string
        output file to write to
  -password string
        Password for the client connection
  -proto string
        proto file to use for decode
  -qos uint
        Qos to use for the session (default 65535)
  -server string
        The server address, host:port
  -subscription string
        Subscription name to subscribe to
  -username string
        Username for the client connection
  -yang_path string
        Yang path for get-proto
Examples:
Subscribe: ./bin/telemetry_dialin_collector -server <ip:port> -subscription <> -encoding self-describing-gpb -username <> -password <>
Get proto for yang path:   ./bin/telemetry_dialin_collector -server <ip:port> -oper get-proto -yang <yang model or xpath> -out <filename> -username <> -password <>
Subscribe, use protoc to decode:   ./bin/telemetry_dialin_collector -server <ip:port> -subscription <> -encoding gpb -username <> -password <> -proto cdp_neighbor.proto
Subscribe, use protoc to decode without proto: ./bin/telemetry_dialin_collector %!s(MISSING) -server <ip:port> -subscription <> -encoding gpb -decode_raw
 $
```

### Example usage:
#### Dialout Server:
```
  // default encoding, expects json message to be streamed
  telemetry_dialout_collector -port 57500
  // Uses telemetry.proto from current directory to decode, needs protoc to be present in $PATH
  telemetry_dialout_collector -port 57500 -encoding self-describing-gpb
  // decode gpb message with proto
  telemetry_dialout_collector -port 57500 -encoding gpb -decode_raw
```
#### Dialin client:
```
  telemetry_dialin_collector -server "<router-ip-address>:<grpc-port>" -subscription <subscription-name> -oper subscribe -username <username> -password <passwd> -encoding <> -qos <dscp>
```
###### Subscribe to a subscription configured on the router
```
  telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab -encoding gpb -qos 10 -proto cdp_neighbor_compact.proto 
  telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab -encoding gpb -qos 10 -decode_raw
  telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab
```
###### Get Proto for an oper model (Supported from 6.5.1 IOS XR release)
```
  telemetry_dialin_collector -server "192.168.122.157:57500" -oper get-proto -username root -password lab -yang_path Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail
  telemetry_dialin_collector -server "192.168.122.157:57500" -oper get-proto -username root -password lab -yang_path Cisco-IOS-XR-cdp-oper:cdp -out cdp.proto
  telemetry_dialin_collector -server "192.168.122.157:57500" -oper get-proto -username root -password lab -yang_path Cisco-IOS-XR-*statsd*
```
Sample output messages from dialin collector are
at [docs/Dialin-collector-examples.md](docs/Dialin-collector-examples.md)
