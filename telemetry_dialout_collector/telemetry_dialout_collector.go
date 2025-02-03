package main

import (
        "os"
        "os/signal"
        "flag"
        "fmt"
        "io"
        "net"
        "strconv"
        "path/filepath"

        "google.golang.org/grpc"
        "google.golang.org/grpc/peer"
        "google.golang.org/grpc/credentials"

        "github.com/ios-xr/telemetry-go-collector/mdt_grpc_dialout"
        "github.com/ios-xr/telemetry-go-collector/telemetry_decode"
)

var usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])

    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr, "Examples:\n")
    fmt.Fprintf(os.Stderr, "GRPC Server                            : %s -port <> -encoding gpb\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "GRPC with TLS                          : %s -port <> -encoding gpb -cert <> -key <>\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "TCP Server                             : %s -port <> -transport tcp\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "GRPC use protoc to decode              : %s -port <> -encoding gpb -proto cdp_neighbor.proto\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "GRPC use protoc to decode without proto: %s -port <> -encoding gpb -decode_raw\n", os.Args[0])
}
var (
        port         = flag.Int("port", 57400, "The server port to listen on")
        encoding     = flag.String("encoding", "json",
                                   "expected encoding, Options: json,self-describing-gpb,gpb needed only for grpc")
        decode_raw   = flag.Bool("decode_raw", false, "Use protoc --decode_raw")
        protoFile    = flag.String("proto", "", "proto file to use for decode")
        transport    = flag.String("transport", "grpc", "transport to use, grpc, tcp or udp")
        dontClean    = flag.Bool("dont_clean", false, "Don't remove tmp files on exit")
        outFileName  = flag.String("out", "dump_*.txt", "output file to write to")
        pluginDir    = flag.String("plugin_dir", "", "absolute path to directory for proto plugins")
        pluginFile    = flag.String("plugin", "", "plugin file, used to lookup gpb symbol for decode")
        certFile     = flag.String("cert","","TLS cert file")
        keyFile      = flag.String("key","","TLS key file")
)

const tmpFileName                = "telemetry-msg-*.dat"

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
             files, _ := filepath.Glob(os.TempDir() + "/" + tmpFileName)
             for _, f := range files {
                 if err := os.Remove(f); err != nil {
                     fmt.Printf("Failed to remove tmp file %s\n",f)
                 }
             }
             os.Exit(0)
         }()
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

     if *certFile != "" && *keyFile != "" {
         fmt.Printf("Enabled TLS, cert: %v key: %v\n", *certFile, *keyFile)
         creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
         if err != nil {
             fmt.Printf("Failed to generate credentials %v", err)
             return
         }
         opts = []grpc.ServerOption{grpc.Creds(creds)}
     }

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
     peer, ok := peer.FromContext(stream.Context())
     if ok {
         fmt.Printf("Session connected from %s\n", peer.Addr.String())
     }

     dataChan := make(chan []byte, 10000)
     defer close(dataChan)
     o := &telemetry_decode.MdtOut{
                        OutFile:     *outFileName,
                        Encoding:    *encoding,
                        Decode_raw:  *decode_raw,
                        DontClean:   *dontClean,
                        ProtoFile:   *protoFile,
                        PluginDir:   *pluginDir,
                        PluginFile:  *pluginFile,
                        DataChan:     dataChan,
     }
     // handler for decoding the data, reads data from dataChan
     go o.MdtOutLoop()

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

         dataChan <- reply.Data
     }

     return nil
}
