## Building basic telemetry collector for IOSXR:

Intent of this excercise is to show how to build a very very simple
telemetry collector for IOSXR that can be used with different
encodings. We will be using Go to develop the collector.

This is meant for beginners with the limited knowledge of golang and protobuf to get started with building grpc server or
client. We will start with grpc server which can act as collector for
router to dialout to with json encoding and add other encodings to it.

Table of contents
=================
<!--ts-->
   * [Install go, protoc, grpc](#install-go-protoc-grpc)
      * [Install Go](#install-go)
      * [Install protobuf/protoc](#install-protoc)
      * [Install protoc-gen-go](#install-protoc-gen-go)
      * [Install grpc](#install-grpc-go)
      * [Install elasticsearch go client](#install-elasticsearch)
   * [Get Telemtry Proto for Dialout Services](#get-telemtry-proto-for-dialout-services)
      * [Generate go binding]( #generate-the-go-binding-for-this-proto)
   * [Grpc Server Code](#grpc-server-code)
   * [GPB Encoding](#gpb-encoding)
<!--te-->

If you are familier with Go, skip section for installing go, protobuf,
protoc, grpc and move to collector part.

### Install go, protoc, grpc:
Default GOROOT is picked as /usr/local and GOPATH as ~/go, atleast
in latest releases, if you use standard methods to install go, it
should be available under /usr/local/go. For this excercise, to
isolate everything I am installing everything needed into new
directory $WORKDIR. If you already have go installed then skip this
section.

#### Install Go:
 bash
 ```
  $ export WORKDIR=~/collector
  $ export GOPATH=$WORKDIR
  $ export GOROOT=$WORKDIR/go
  $ export PATH=$GOROOT/bin:$GOPATH/bin:$PATH

  $ mkdir $WORKDIR
  $ cd $WORKDIR
  $ wget https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz
  $ tar xvfz go1.11.5.linux-amd64.tar.gz
```
#### Install protoc:
```
  $ wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0rc2/protoc-3.7.0-rc-2-linux-x86_64.zip
  $ unzip protoc-3.7.0-rc-2-linux-x86_64.zip
```
#### Install protoc-gen-go:
`  $ go get -u github.com/golang/protobuf/protoc-gen-go`

#### Install grpc go:
`  $ go get -u google.golang.org/grpc`

#### Install elasticsearch go client:
   This is only needed if you want to push the data to elasticsearch  
`  $ go get github.com/elastic/go-elasticsearch`

Almost ready to write the collector, just make sure "go", "protoc" and
"protoc-gen-go" are comming from place you intend them to be used from
(you can use "which" command to check the path).

### Get Telemtry Proto for Dialout Services:
Now we need to get the proto which defines the service that the server is going to provide. You can get it from,
https://github.com/cisco-ie/bigmuddy-network-telemetry-pipeline/blob/master/vendor/github.com/cisco/bigmuddy-network-telemetry-proto/staging/mdt_grpc_dialout/mdt_grpc_dialout.proto

```
 $ mkdir $GOPATH/src/mdt_grpc_dialout
 $ vi src/mdt_grpc_dialout/mdt_grpc_dialout.proto  //copy the raw content from above link
 $ more src/mdt_grpc_dialout/mdt_grpc_dialout.proto
syntax = "proto3";

// Package implements gRPC Model Driven Telemetry service
package mdt_dialout;

// gRPCMdtDialout defines service used for client-side streaming pushing MdtDialoutArgs.
service gRPCMdtDialout {
    rpc MdtDialout(stream MdtDialoutArgs) returns(stream MdtDialoutArgs) {};
}

// MdtDialoutArgs is the content pushed to the server
message MdtDialoutArgs {
     int64 ReqId = 1;
     // data carries the payload content.
     bytes data = 2;
     string errors = 3;
}
 $
```
#### Generate the go binding for this proto:
` $ protoc -I$GOPATH/src/mdt_grpc_dialout --go_out=plugins=grpc:src/mdt_grpc_dialout mdt_grpc_dialout.proto`

You should have mdt_grpc_dialout.pb.go newly generated go file in same
direcotry as proto.

### GRPC Server code:
Now we will need to provide the function/rpc defined by gRPCMdtDialout
service. We will do it as part of simple server implementation that
implements this service.

Create new go file with following content, we will walk through
content below.
```
 $ mkdir $GOPATH/src/telemetry_dialout_collector
 $ vi src/telemetry_dialout_collector/telemetry_dialout_collector.go
 $ more src/telemetry_dialout_collector/telemetry_dialout_collector.go
package main

import (
        "flag"
        "fmt"
        "io"
        "time"
        "net"
        "bytes"
        "encoding/json"
        "google.golang.org/grpc"
        "mdt_grpc_dialout"
)

var (
        port          = flag.Int("port", 57400, "The server port")
)

type gRPCMdtDialoutServer struct{}

func (s *gRPCMdtDialoutServer) MdtDialout(stream mdt_dialout.GRPCMdtDialout_MdtDialoutServer) error {
        var numMsgs = 0
        var prettyJSON bytes.Buffer

        for {
                reply, err := stream.Recv()
                if err == io.EOF {
                        fmt.Printf("MdtDialout: Got EOF\n\n")
                        return err
                }
                if err != nil {
                        fmt.Printf("MdtDialout: Stream Recv got error %v", err)
                        return err
                }
                numMsgs++
                t := time.Now()
                fmt.Println(t.Format(time.RFC3339Nano), "message count: ", numMsgs, "bytes ", len(reply.Data))
                err = json.Indent(&prettyJSON, reply.Data, "", "\t")
                if err != nil {
                       fmt.Println("JSON parse error: ", err)
                } else {
                       fmt.Println(string(prettyJSON.Bytes()))
                }
        }
}

func main() {
        var lis net.Listener
        var err error

        flag.Parse()
        var opts []grpc.ServerOption

        grpcServer := grpc.NewServer(opts...)
        s := gRPCMdtDialoutServer{}
        mdt_dialout.RegisterGRPCMdtDialoutServer(grpcServer, &s)

        lis, err = net.Listen("tcp", fmt.Sprintf(":%d", *port))
        if err != nil {
                fmt.Printf("Failed to open listen port %v", err)
                return
        }

        fmt.Printf("mdtDialout server: lport(%v) \n", *port)

        grpcServer.Serve(lis)
        if err != nil {
                fmt.Printf("Server stopped: %v", err)
        }
}
```
###### Going through sections of this code,
1) This is "main" package, that means main is present in this package
2) import all the packages we are using in this is file (there will be
more later(for credentials, tls etc services), for now we only need these)
3) arguments that can be passed to this server, just port number to
listen on for now.
4) gRPCMdtDialoutServer struct that implements generated
gRPCMdtDialoutServer interface
5) MdtDialout rpc implementation for our service
   a) Use mdt_dialout.GRPCMdtDialout_MdtDialoutServer's Recv method to
   read data from client.
   b) pretty print received data using json package
6) main function
   a) parse the input arguments
   b) create server instance
   c) attach the service to grpc server
   d) create a listener and start listening


#### Compile the server code:
 from $GOPATH directory,
 
` $ go build -i -v -o bin/telemetry_dialout_collector telemetry_dialout_collector`

Ready to use the collector,
```
 $ telemetry_dialout_collector -h
Usage of telemetry_dialout_collector:
  -port int
        The server port (default 57400)
 $
```
###### Output example:
```
 $ /ws/adithyas-sjc/collector/bin/telemetry_dialout_collector -port 57500
mdtDialout server: lport(57500) 
2019-02-15T09:09:07.167004137-08:00 message count:  1 bytes  783
{
        "node_id_str": "adithyas-1",
        "subscription_id_str": "sw-ver",
        "encoding_path": "Cisco-IOS-XR-spirit-install-instmgr-oper:software-install/version",
        "collection_id": "33",
        "collection_start_time": "1550250518925",
        "msg_timestamp": "1550250518929",
        "data_json": [
           {
                "timestamp": "1550250518928",
                "keys": [],
                "content": {
                        "package": [
                                {
                                        "name": "IOS-XR",
                                        "version": "6.5.3.02I",
                                        "built-by": "satbhatt",
                                        "built-on": "Mon Feb 11 14:46:45 PST 2019",
                                        "build-host": "sjc-ads-1163",
                                        "workspace": "/nobackup/satbhatt/mdt-r65x"
                                }
                        ],
                        "location": "/opt/cisco/XR/packages/",
                        "copyright-info": "Cisco IOS XR Software, Version 6.5.3.02I\nCopyright (c) 2013-2019 by Cisco Systems, Inc.",
                        "hardware-info": "cisco IOS-XRv 9000 () processor",
                        "system-uptime": "System uptime is 3 days 17 hours 23 minutes"
                }
           }
        ],
        "collection_end_time": "1550250518929"
}

^C
 $
```
### GPB encoding:
To support Self-Describing-GPB or Compact GPB encoding, there is
atleast one more step that needs to happen in the collector to decode
the data according to proto defination for the message.

Proto for GPB encoding can be obtained from,
https://github.com/cisco/bigmuddy-network-telemetry-proto/blob/master/staging/telemetry.proto

There might be updated versions of this proto released but it will
always be backward compatible, so we should be fine using this proto
for decoding the data.
Simplest option to decode the data is to pass in received message as
input to protoc. To do this, we will save the message to tmp file and
provide it as input to protoc just to show how it can be done.

Code:
This is updated code for the server, new code added is for handling
the gpb encoding.
```
 $ more $GOPATH/src/telemetry_dialout_collector/telemetry_dialout_collector.go
package main

import (
        "os"
        "os/exec"
        "flag"
        "fmt"
        "io"
        "io/ioutil"
        "log"
        "time"
        "net"
        "bytes"
        "encoding/json"
        "google.golang.org/grpc"
        "mdt_grpc_dialout"
)

var (
        port          = flag.Int("port", 57400, "The server port")
        encoding      = flag.String("encoding", "json", "expected encoding of msg")
        protoFile     = flag.String("proto", "telemetry.proto", "proto file to use for decode")
)

const ProtocCommandString string = "protoc --decode=Telemetry "

type gRPCMdtDialoutServer struct{}

func (s *gRPCMdtDialoutServer) MdtDialout(stream mdt_dialout.GRPCMdtDialout_MdtDialoutServer) error {
        var numMsgs = 0
        var prettyJSON bytes.Buffer
        var commandString string

        tmpfile, err := ioutil.TempFile("", "telemetry-msg-")
        if (err != nil) {
           log.Fatal("Failed to create tmp file for writing", err)
        }
        defer os.Remove(tmpfile.Name())
        for {
                reply, err := stream.Recv()
                if err == io.EOF {
                        fmt.Printf("MdtDialout: Got EOF\n\n")
                        return err
                }
                if err != nil {
                        fmt.Printf("MdtDialout: Stream Recv got error %v", err)
                        return err
                }
                numMsgs++
                t := time.Now()
                fmt.Println(t.Format(time.RFC3339Nano), "message count: ", numMsgs, "bytes ", len(reply.Data))

                if (*encoding == "json") {
                   err = json.Indent(&prettyJSON, reply.Data, "", "\t")
                   if err != nil {
                       fmt.Println("JSON parse error: ", err)
                   } else {
                       fmt.Println(string(prettyJSON.Bytes()))
                   }
                } else if (*encoding == "self-describing-gpb") {
                   // Decode the data using protoc
                   _, err = tmpfile.Write(reply.Data)
                   commandString = ProtocCommandString + *protoFile + "<" + tmpfile.Name()
                   out, err := exec.Command("sh", "-c", commandString).Output()
                   fmt.Printf("%s\n", out)
                   if err == nil {
                      tmpfile.Truncate(0)
                      tmpfile.Seek(0,0)
                   }

                }
        }
        tmpfile.Close()
        return nil
}

func main() {
        var lis net.Listener
        var err error

        flag.Parse()
        var opts []grpc.ServerOption

        grpcServer := grpc.NewServer(opts...)
        s := gRPCMdtDialoutServer{}
        mdt_dialout.RegisterGRPCMdtDialoutServer(grpcServer, &s)

        lis, err = net.Listen("tcp", fmt.Sprintf(":%d", *port))
        if err != nil {
                fmt.Printf("Failed to open listen port %v", err)
                return
        }

        fmt.Printf("mdtDialout server: lport(%v) \n", *port)

        grpcServer.Serve(lis)
        if err != nil {
                fmt.Printf("Server stopped: %v", err)
        }
}
```

#### Build:
` $ go build -i -v -o bin/telemetry_dialout_collector telemetry_dialout_collector`

When running this collector for gpb encoding, make sure protoc in
your path is atleast 3.0.0 version (you can check using protoc
--version) and telemetry.proto is present in the directory where you
are running the colletor.

##### Self-Describing-GPB example:
```
 $ /ws/adithyas-sjc/collector/bin/telemetry_dialout_collector -port 57500 -encoding self-describing-gpb
mdtDialout server: lport(57500) 
2019-02-18T22:33:09.662131641-08:00 message count:  1 bytes  1810
node_id_str: "adithyas-1"
subscription_id_str: "cdp-neighbor"
encoding_path: "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail"
collection_id: 23
collection_start_time: 1550557984617
msg_timestamp: 1550557984617
data_gpbkv {
  timestamp: 1550557984628
  fields {
    name: "keys"
    fields {
      name: "node-name"
      string_value: "0/0/CPU0"
    }
    fields {
      name: "interface-name"
      string_value: "GigabitEthernet0/0/0/0"
    }
    fields {
      name: "device-id"
      string_value: "adithyas-2"
    }
  }
  fields {
    name: "content"
    fields {
      name: "cdp-neighbor"
      fields {
        name: "receiving-interface-name"
        string_value: "GigabitEthernet0/0/0/0"
      }
      fields {
        name: "device-id"
        string_value: "adithyas-2"
      }
      fields {
        name: "port-id"
        string_value: "GigabitEthernet0/0/0/0"
      }
      fields {
        name: "header-version"
        uint32_value: 2
      }
      fields {
        name: "hold-time"
        uint32_value: 174
      }
      fields {
        name: "capabilities"
        string_value: "R"
      }
      fields {
        name: "platform"
        string_value: "cisco IOS-XRv 9000"
      }
      fields {
        name: "detail"
        fields {
          name: "network-addresses"
          fields {
            name: "cdp-addr-entry"
            fields {
              name: "address"
              fields {
                name: "address-type"
                string_value: "ipv4"
              }
              fields {
                name: "ipv4-address"
                string_value: "4.0.0.2"
              }
            }
          }
          fields {
            name: "cdp-addr-entry"
            fields {
              name: "address"
              fields {
                name: "address-type"
                string_value: "ipv6"
              }
              fields {
                name: "ipv6-address"
                string_value: "2002::1:2"
              }
            }
          }
        }
        fields {
          name: "version"
          string_value: " 6.5.3.02I"
        }
        fields {
          name: "native-vlan"
          uint32_value: 0
        }
        fields {
          name: "duplex"
          string_value: "cdp-dplx-none"
        }
        fields {
          name: "system-name"
          string_value: "adithyas-2"
        }
      }
    }
  }
}
collection_end_time: 1550557984634

^C
 $
```
##### Compact GPB message with same collector:
```
 $ ./bin/telemetry_dialout_collector -port 57500 -encoding self-describing-gpb
mdtDialout server: lport(57500) 
2019-02-18T22:36:55.258671652-08:00 message count:  1 bytes  799
node_id_str: "adithyas-1"
subscription_id_str: "cdp-neighbor"
encoding_path: "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail"
collection_id: 25
collection_start_time: 1550558210191
msg_timestamp: 1550558210191
data_gpb {
  row {
    timestamp: 1550558210206
    keys: "\n\0100/0/CPU0\022\026GigabitEthernet0/0/0/0\032\nadithyas-2"
    content: "\222\003\255\001\n\026GigabitEthernet0/0/0/0\022\nadithyas-2\032\026GigabitEthernet0/0/0/0 \002(\201\0012\001R:\022cisco IOS-XRv 9000BS\n(\n\021\n\017\n\004ipv4\022\0074.0.0.2\n\023\n\021\n\004ipv6\032\t2002::1:2\022\n 6.5.3.02I(\0002\rcdp-dplx-none:\nadithyas-2"
  }
}
collection_end_time: 1550558210211

^C
 $
```
As you can see, message encoded with gpb can also be decoded using
same proto as self-describing-gpb except for key and content part. We
will need proto specific to encoding_path to decode keys and content.

You can get proto for this path from,
https://github.com/cisco/bigmuddy-network-telemetry-proto/blob/master/staging/cisco_ios_xr_cdp_oper/cdp/nodes/node/neighbors/details/detail/cdp_neighbor.proto

key and content messages have to be updated in telemetry.proto to
change them from bytes to message defination in above proto for the
sensor path we are testing.

Final proto we are going to use for cdp neighbor details path is,
```
 $ more cdp_neighbor_compact.proto 
/* ----------------------------------------------------------------------------
 * telemetry_bis.proto - Telemetry protobuf definitions
 *
 * August 2016
 *
 * Copyright (c) 2016 by Cisco Systems, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ----------------------------------------------------------------------------
 */

syntax = "proto3";
option go_package="telemetry";

// Telemetry message is the outermost payload message used to stream
// telemetry in a Model Driven Telemetry (MDT) system. MDT provides a
// mechanism for an external entity to subscribe to a data set defined in
// a Yang model and receive periodic or event-based updates of the data
// set from an MDT-capable device.
message Telemetry {
  //
  // node_id_str is a string encoded unique node ID of the MDT-capable
  // device producing the message. (node_id_uuid alternative is not currently
  // produced in IOS-XR)
  oneof node_id {
    string node_id_str = 1;
    // bytes node_id_uuid = 2;
  }
  //
  // subscription_id_str is the name of the subscription against which
  // this content is being produced. (subscription_id alternative is not
  //  currently produced in IOS-XR)
  oneof subscription {
    string   subscription_id_str = 3;
    // uint32   subscription_id = 4;
  }
  //
  // sensor_path is not currently produced in IOS-XR
  // string   sensor_path = 5;
  //
  // encoding_path is the Yang path leading to the content in this message.
  // The Yang tree encoded in the content section of this message is rooted
  // at the point described by the encoding_path.
  string   encoding_path = 6;
  //
  // model_version is not currently produced in IOS-XR
  // string   model_version = 7;
  //
  // collection_id identifies messages belonging to a collection round.
  // Multiple message may be generated from a collection round.
  uint64   collection_id = 8;
  //
  // collection_start_time is the time when the collection identified by
  // the collection_id begins - encoded as milliseconds since the epoch.
  // If a single collection is spread over multiple Telemetry Messages,
  // collection_start_time may be encoded in the first Telemetry Message
  // for the collection only.
  uint64   collection_start_time = 9;
  //
  // msg_timestamp is the time when the data encoded in the Telemetry
  // message is generated - encoded as milliseconds since the epoch.
  uint64   msg_timestamp = 10;
  //
  // data_gpbkv contains the payload data if data is being encoded in the
  // self-describing GPB-KV format.
  repeated TelemetryField data_gpbkv = 11;
  //
  // data_gpb contains the payload data if data is being encoded as
  // serialised GPB messages.
  TelemetryGPBTable data_gpb = 12;
  //
  // collection_end_time is the timestamp when the last Telemetry message
  // for a collection has been encoded - encoded as milliseconds since the
  // epoch. If a single collection is spread over multiple Telemetry
  // messages, collection_end_time is encoded in the last Telemetry Message
  // for the collection only.
  uint64 collection_end_time = 13;
  //
  // heartbeat_sequence_number is not currently produced in IOS-XR
  // uint64   heartbeat_sequence_number = 14; // not produced
}

//
// TelemetryField messages are used to export content in the self
// describing GPB KV form. The TelemetryField message is sufficient to
// decode telemetry messages for all models. KV-GPB encoding is very
// similar in concept, to JSON encoding
message TelemetryField {
  //
  // timestamp represents the starting time of the generation of data
  // starting from this key, value pair in this message - encoded as
  // milliseconds since the epoch. It is encoded when different from the
  // msg_timestamp in the containing Telemetry Message. This field can be
  // omitted if the value is the same as a TelemetryField message up the
  // hierarchy within the same Telemetry Message as well.
  uint64         timestamp = 1;
  //
  // name: string encoding of the name in the key, value pair. It is
  // the corresponding YANG element name.
  string         name = 2;
  //
  // value_by_type, if present, for the corresponding YANG element
  // represented by the name field in the same TelemetryField message. The
  // value is encoded to the matching type as defined in the YANG model.
  // YANG models often define new types (derived types) using one or more
  // base types.  The types included in the oneof grouping is sufficient to
  // represent such derived types. Derived types represented as a Yang
  // container are encoded using the nesting primitive defined in this
  // encoding proposal.
  oneof value_by_type {
    bytes          bytes_value = 4;
    string         string_value = 5;
    bool           bool_value = 6;
    uint32         uint32_value = 7;
    uint64         uint64_value = 8;
    sint32         sint32_value = 9;
    sint64         sint64_value = 10;
    double         double_value = 11;
    float          float_value = 12;
  }
  //
  // The Yang model may include nesting (e.g hierarchy of containers). The
  // next level of nesting, if present, is encoded, starting from fields.
  repeated TelemetryField fields = 15;
}

// TelemetryGPBTable contains a repeated number of TelemetryRowGPB,
// each of which represents content from a subtree instance in the
// the YANG model. For example; a TelemetryGPBTable might contain
// the interface statistics of a collection of interfaces.
message TelemetryGPBTable {
  repeated TelemetryRowGPB row = 1;
}

//
// TelemetryRowGPB, in conjunction with the Telemetry encoding_path and
// model_version, unambiguously represents the root of a subtree in
// the YANG model, and content from that subtree encoded in serialised
// GPB messages. For example; a TelemetryRowGPB might contain the
// interface statistics of one interface. Per encoding-path .proto
// messages are required to decode keys/content pairs below.
message TelemetryRowGPB {
  //
  // timestamp at which the data for this instance of the TelemetryRowGPB
  // message was generated by an MDT-capable device - encoded as
  // milliseconds since the epoch.  When included, this is typically
  // different from the msg_timestamp in the containing Telemetry message.
  uint64 timestamp = 1;
  //
  // keys: if the encoding-path includes one or more list elements, and/or
  // ends in a list element, the keys field is a GPB encoded message that
  // contains the sequence of key values for each such list element in the
  // encoding-path traversed starting from the root.  The set of keys
  // unambiguously identifies the instance of data encoded in the
  // TelemetryRowGPB message. Corresponding protobuf message definition will
  // be required to decode the byte stream. The encoding_path field in
  // Telemetry message, together with model_version field should be
  // sufficient to identify the corresponding protobuf message.
  cdp_neighbor_KEYS keys = 10;
  //
  // content: the content field is a GPB encoded message that contains the
  // data for the corresponding encoding-path. A separate decoding pass
  // would be performed by consumer with the content field as a GPB message
  // and the matching .proto used to decode the message. Corresponding
  // protobuf message definition will be required to decode the byte
  // stream. The encoding_path field in Telemetry message, together with
  // model_version field should be sufficient to identify the corresponding
  // protobuf message. The decoded combination of keys (when present) and
  // content, unambiguously represents an instance of the data set, as
  // defined in the Yang model, identified by the encoding-path in the
  // containing Telemetry message.
  cdp_neighbor content = 11;
}

// CDP neighbor info
message cdp_neighbor_KEYS {
    string node_name = 1;
    string interface_name = 2;
    string device_id = 3;
}

message cdp_neighbor {
    // Next neighbor in the list
    repeated cdp_neighbor_item cdp_neighbor = 50;
}

message cdp_neighbor_item {
    // Interface the neighbor entry was received on 
    string receiving_interface_name = 1;
    // Device identifier
    string device_id = 2;
    // Outgoing port identifier
    string port_id = 3;
    // Version number
    uint32 header_version = 4;
    // Remaining hold time
    uint32 hold_time = 5;
    // Capabilities
    string capabilities = 6;
    // Platform type
    string platform = 7;
    // Detailed neighbor info
    cdp_neighbor_detail detail = 8;
}

message in6_addr_td {
    string value = 1;
}

message cdp_l3_addr {
    string address_type = 1;
    // IPv4 address
    string ipv4_address = 2;
    // IPv6 address
    //in6_addr_td ipv6_address = 3;
}

message cdp_addr_entry {
    // Next address entry in list
    repeated cdp_addr_entry_item cdp_addr_entry = 1;
}

message cdp_addr_entry_item {
    // Network layer address
    cdp_l3_addr address = 1;
}

message cdp_prot_hello_entry {
    // Next protocol hello entry in list
    repeated cdp_prot_hello_entry_item cdp_prot_hello_entry = 1;
}

message cdp_prot_hello_entry_item {
    // Protocol Hello msg
    bytes hello_message = 1;
}

message cdp_neighbor_detail {
    // List of network addresses 
    cdp_addr_entry network_addresses = 1;
    // Version TLV
    string version = 2;
    // List of protocol hello entries
    cdp_prot_hello_entry protocol_hello_list = 3;
    // VTP domain
    string vtp_domain = 4;
    // Native VLAN
    uint32 native_vlan = 5;
    // Duplex setting
    string duplex = 6;
    // SysName
    string system_name = 7;
}
 $
```
###### Decoded output example for compact gpb:
```
 $ telemetry_dialout_collector -port 57500 -encoding self-describing-gpb -proto cdp_neighbor_compact.proto
mdtDialout server: lport(57500) 
Session connected from 192.168.122.41:59837

2019-02-26T12:42:48.59473828-08:00 message count:  1 bytes  799
node_id_str: "adithyas-1"
subscription_id_str: "cdp-neighbor"
encoding_path: "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail"
collection_id: 82
collection_start_time: 1551213760962
msg_timestamp: 1551213760962
data_gpb {
  row {
    timestamp: 1551213760972
    keys {
      node_name: "0/0/CPU0"
      interface_name: "GigabitEthernet0/0/0/0"
      device_id: "adithyas-2"
    }
    content {
      cdp_neighbor {
        receiving_interface_name: "GigabitEthernet0/0/0/0"
        device_id: "adithyas-2"
        port_id: "GigabitEthernet0/0/0/0"
        header_version: 2
        hold_time: 149
        capabilities: "R"
        platform: "cisco IOS-XRv 9000"
        detail {
          network_addresses {
            cdp_addr_entry {
              address {
                address_type: "ipv4"
                ipv4_address: "4.0.0.2"
              }
            }
            cdp_addr_entry {
              address {
                address_type: "ipv6"
                3: "2002::1:2"
              }
            }
          }
          version: " 6.5.3.02I"
          duplex: "cdp-dplx-none"
          system_name: "adithyas-2"
        }
      }
    }
  }
  row {
    timestamp: 1551213760972
    keys {
      node_name: "0/0/CPU0"
      interface_name: "GigabitEthernet0/0/0/2"
      device_id: "adithyas-2"
    }
    content {
      cdp_neighbor {
        receiving_interface_name: "GigabitEthernet0/0/0/2"
        device_id: "adithyas-2"
        port_id: "GigabitEthernet0/0/0/2"
        header_version: 2
        hold_time: 157
        capabilities: "R"
        platform: "cisco IOS-XRv 9000"
        detail {
          network_addresses {
            cdp_addr_entry {
              address {
                address_type: "ipv4"
                ipv4_address: "2.2.2.2"
              }
            }
          }
          version: " 6.5.3.02I"
          duplex: "cdp-dplx-none"
          system_name: "adithyas-2"
        }
      }
    }
  }
  row {
    timestamp: 1551213760972
    keys {
      node_name: "0/0/CPU0"
      interface_name: "GigabitEthernet0/0/0/1"
      device_id: "adithyas-2"
    }
    content {
      cdp_neighbor {
        receiving_interface_name: "GigabitEthernet0/0/0/1"
        device_id: "adithyas-2"
        port_id: "GigabitEthernet0/0/0/1"
        header_version: 2
        hold_time: 153
        capabilities: "R"
        platform: "cisco IOS-XRv 9000"
        detail {
          network_addresses {
            cdp_addr_entry {
              address {
                address_type: "ipv4"
                ipv4_address: "5.0.0.2"
              }
            }
          }
          version: " 6.5.3.02I"
          duplex: "cdp-dplx-none"
          system_name: "adithyas-2"
        }
      }
    }
  }
}
collection_end_time: 1551213760974

^C
 $
```
