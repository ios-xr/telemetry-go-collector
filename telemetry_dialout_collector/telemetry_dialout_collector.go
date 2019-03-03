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
        "encoding/json"
        "path/filepath"
        "google.golang.org/grpc"
        "google.golang.org/grpc/peer"
        "github.com/ios-xr/telemetry-go-collector/mdt_grpc_dialout"
)

var (
        port          = flag.Int("port", 57400, "The server port")
        encoding      = flag.String("encoding", "json", "expected encoding of msg")
        raw_decode    = flag.Bool("decode_raw", false, "Use protoc --decode_raw")
        protoFile     = flag.String("proto", "telemetry.proto", "proto file to use for decode")
        dontClean     = flag.Bool("dont_clean", false, "Don't remove tmp files on exit")
)

const ProtocRawDecode string = "protoc --decode_raw "
const ProtocCommandString string = "protoc --decode=Telemetry "

type gRPCMdtDialoutServer struct{}

func (s *gRPCMdtDialoutServer) MdtDialout(stream mdt_dialout.GRPCMdtDialout_MdtDialoutServer) error {
        var numMsgs = 0
        var prettyJSON bytes.Buffer
        var commandString string

        peer, ok := peer.FromContext(stream.Context())
        if ok {
           fmt.Printf("Session connected from %s\n", peer.Addr.String())
        }

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
                   if *raw_decode {
                     commandString = ProtocRawDecode + "<" + tmpfile.Name()
                   } else {
                     commandString = ProtocCommandString + *protoFile + "<" + tmpfile.Name()
                   }
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

        if !*dontClean {
           // install signal handler for cleaning up tmp files
           sigs := make(chan os.Signal, 1)
           signal.Notify(sigs, os.Interrupt)
           go func() {
              <- sigs
              //cleanup()
              files, _ := filepath.Glob("/tmp/telemetry-msg-*")
              for _, f := range files {
                  if err := os.Remove(f); err != nil {
                     fmt.Printf("Failed to remove tmp file %s\n",f)
                  }
              }
              os.Exit(0)
           }()
        }

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

