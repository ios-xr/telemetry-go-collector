package main

import (
        "os"
        "os/exec"
        "os/signal"
        "flag"
        "fmt"
        "io"
        "io/ioutil"
        "log"
        "time"
        "net"
        "bytes"
        "strconv"
        "encoding/json"
        "path/filepath"

        "google.golang.org/grpc"
        "google.golang.org/grpc/peer"
        "github.com/golang/protobuf/proto"

        "mdt_grpc_dialout"
        "telemetry"
)

var usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])

    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr, "Examples:\n")
    fmt.Fprintf(os.Stderr, "GRPC Server:                             %s -port <> -encoding gpb\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "TCP Server:                              %s -port <> -transport tcp\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "GRPC use protoc to decode:               %s -port <> -encoding gpb -proto cdp_neighbor.proto\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "GRPC use protoc to decode without proto: %s -port <> -encoding gpb -decode_raw\n", os.Args[0])
}
var (
        port          = flag.Int("port", 57400, "The server port to listen on")
        encoding      = flag.String("encoding", "json",
                                    "expected encoding, Options: json,self-describing-gpb,gpb, needed only for grpc")
        decode_raw    = flag.Bool("decode_raw", false, "Use protoc --decode_raw")
        protoFile     = flag.String("proto", "", "proto file to use for decode")
        transport     = flag.String("transport", "grpc", "transport to use, grpc, tcp or udp")
        dontClean     = flag.Bool("dont_clean", false, "Don't remove tmp files on exit")
        outFile      = flag.String("out", "", "output file to write to")
)

// output file
var oFile *os.File

const tmpFileName                = "telemetry-msg-*.dat"
const ProtocRawDecode string     = "protoc --decode_raw "
const ProtocCommandString string = "protoc --decode=Telemetry "

func main() {
     flag.Usage = usage
     flag.Parse()

     if !*dontClean {
         // install signal handler for cleaning up tmp files
         sigs := make(chan os.Signal, 1)
         signal.Notify(sigs, os.Interrupt)
         go func() {
             <- sigs
             // cleanup
             files, _ := filepath.Glob("/tmp/" + tmpFileName)
             for _, f := range files {
                 if err := os.Remove(f); err != nil {
                     fmt.Printf("Failed to remove tmp file %s\n",f)
                 }
             }
             os.Exit(0)
         }()
     }

     // write data to stdout unless output file is specified
     oFile = os.Stdout
     if len(*outFile) != 0 {
        oFile, _ = os.Create(*outFile)
        defer oFile.Close()
     }

     if (*transport == "tcp") {
         mdtTcpServer(":" + strconv.Itoa(*port))
     } else if (*transport == "udp") {
         mdtUdpServer(":" + strconv.Itoa(*port))
     } else {
         mdtGrpcServer(":" + strconv.Itoa(*port))
     }
}

// grpc server
func mdtGrpcServer(grpcPort string) {
     var lis net.Listener
     var err error
     var opts []grpc.ServerOption

     grpcServer := grpc.NewServer(opts...)
     s := gRPCMdtDialoutServer{}
     mdt_dialout.RegisterGRPCMdtDialoutServer(grpcServer, &s)

     lis, err = net.Listen("tcp", grpcPort)
     if err != nil {
         fmt.Printf("Failed to open listen port %v", err)
         return
     }

     fmt.Println("GRPC server listening at ", grpcPort)
     grpcServer.Serve(lis)
     if err != nil {
         fmt.Printf("Server stopped: %v", err)
     }
}

type gRPCMdtDialoutServer struct{}

func (s *gRPCMdtDialoutServer) MdtDialout(stream mdt_dialout.GRPCMdtDialout_MdtDialoutServer) error {
     var numMsgs = 0
     var commandString string

     peer, ok := peer.FromContext(stream.Context())
     if ok {
         fmt.Printf("Session connected from %s\n", peer.Addr.String())
     }

     tmpFile, commandString := mdtPrepareDecoding()
     if tmpFile != nil {
         defer tmpFile.Close()
     }
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

         mdtDumpData(reply.Data, *encoding, commandString, tmpFile)
     }

     return nil
}

/////////////////
// handle output
/////////////////
func mdtDumpData(data []byte, enc string, commandString string, tmpFile *os.File) {
     var prettyJSON bytes.Buffer
     var err error

     if (enc == "json") {
        err = json.Indent(&prettyJSON, data, "", "\t")
        if err != nil {
           fmt.Println("JSON parse error: ", err)
        } else {
           _, err = oFile.WriteString(string(prettyJSON.Bytes()))
        }
     } else if ((enc == "self-describing-gpb") ||
                (enc == "gpb")) {
         if (tmpFile == nil) {
             telem := &telemetry.Telemetry{}
             err = proto.Unmarshal(data, telem)
             if (err != nil) {
                 fmt.Println("Failed to unmarshal:", err)
             }
             j, _ :=  json.MarshalIndent(telem, "", "  ")
             _, err = oFile.WriteString(string(j))
         } else {
             // Write to tmp file
             _, err = tmpFile.Write(data)
             if (err != nil) {
                 fmt.Println("Failed to write to tmp file", err)
             }

             // Decode the data using protoc
             out, err := exec.Command("sh", "-c", commandString).CombinedOutput()
             if err == nil {
                 _, err = oFile.WriteString(string(out))
                 tmpFile.Truncate(0)
                 tmpFile.Seek(0,0)
             } else {
                 fmt.Println("Protoc error", err, out)
             }
         }
     } else {
        fmt.Println("Unknown encoding")
     }
}

func mdtPrepareDecoding() (*os.File, string) {
     var commandString string

     if *decode_raw || (len(*protoFile) != 0) {
         // temp file to write message to for decoding
         tmpFile, err := ioutil.TempFile("", tmpFileName)
         if (err != nil) {
             log.Fatal("Failed to create tmp file for writing", err)
         }

         // proto command to use for decoding gpb message
         if *decode_raw {
             commandString = ProtocRawDecode + " < " + tmpFile.Name()
         } else {
             commandString = ProtocCommandString + *protoFile + " < " + tmpFile.Name()
         }
         return tmpFile, commandString
     }
     return nil, ""
}
