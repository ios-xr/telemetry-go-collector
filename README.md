## IOS XR Telemetry collector

A Simple "Go" based collector for IOS XR Telemetry that can get you started on trying telemetry. IOS XR supports dialout and dialin modes for connecting to router for streaming, this repo will have both the collectors. These collectors only read the data as its sent from the router and dump on to stdout, you can add on to these scripts to do many fancy things. This is meant for begginers, if you are familier with Go, GRPC and IOSXR, you should be using and modifying/extending [Pipeline](https://github.com/cisco-ie/bigmuddy-network-telemetry-pipeline) for your needs.

### MDT Dialout Collector:
If you want to start from scratch or do not have "Go", Protoc, protoc-gen-go or grpc installed, you can check out [Dialout-collector-howto.md](Dialout-collector-howto.md). It has instructions on how to get started and how to write simple collector.

Assuming "Go" is already installed, following intructions are for getting collector, building it and running it.

##### Install instructions:
`go get -d github.com/adithyasesani/test-collector`

##### Build
`go build -o bin/telemetry_collector src/github.com/adithyasesani/test-collector/telemetry_dialout_collector/telemetry_dialout_collector.go`

##### Run
```
$ ./bin/telemetry_collector -h
Usage of ./bin/telemetry_collector:
  -dont_clean
    	Don't remove tmp files on exit
  -encoding string
    	expected encoding of msg (default "json")
  -port int
    	The server port (default 57400)
  -proto string
    	proto file to use for decode (default "telemetry.proto")
  -raw_decode
    	Use protoc --decode_raw
 $
 ```


![](docs/dialout-build.gif)
