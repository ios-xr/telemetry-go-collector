package main

import (
       "bytes"
       "encoding/json"
       "flag"
       "fmt"
       "io"
       "io/ioutil"
       "log"
       "os"
       "os/exec"
       "os/signal"
       "strings"
       "path/filepath"

       "golang.org/x/net/context"
       "google.golang.org/grpc"
        "github.com/golang/protobuf/proto"

       MdtDialin "github.com/ios-xr/telemetry-go-collector/mdt_grpc_dialin"
       "github.com/ios-xr/telemetry-go-collector/telemetry"
)

const tmpFileName                = "telemetry-msg-*.dat"
const ProtocRawDecode string     = "protoc --decode_raw "
const ProtocCommandString string = "protoc --decode=Telemetry "

const NotConfigured = 0xffff

var telemetryEncoding = map[string]int64{
    "gpb":                 2,
    "self-describing-gpb": 3,
    "json":                4,
}

var usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])

    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr, "Examples:\n")
    fmt.Fprintf(os.Stderr, "Subscribe: %s -server <ip:port> -subscription <> -encoding self-describing-gpb -username <> -password <>\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "Get proto for yang path:   %s -server <ip:port> -oper get-proto -yang <yang model or xpath> -out <filename> -username <> -password <>\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "Subscribe, use protoc to decode:   %s -server <ip:port> -subscription <> -encoding gpb -username <> -password <> -proto cdp_neighbor.proto\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "Subscribe, use protoc to decode without proto: %s %s -server <ip:port> -subscription <> -encoding gpb -decode_raw\n", os.Args[0])
}

var (
        serverAddr   = flag.String("server", "", "The server address, host:port")
        operation    = flag.String("oper", "subscribe", "Operation: subscribe, get-proto")
        subIds       = flag.String("subscription", "", "Subscription name to subscribe to")
        encoding     = flag.String("encoding", "json",
                                   "encoding to use, Options: json,self-describing-gpb,gpb")
        qos          = flag.Uint("qos", NotConfigured, "Qos to use for the session")
        yangPath     = flag.String("yang_path", "", "Yang path for get-proto")
        outFile      = flag.String("out", "", "output file to write to")
        username     = flag.String("username", "",
                                   "Username for the client connection")
        password     = flag.String("password", "",
                                   "Password for the client connection")
        decode_raw   = flag.Bool("decode_raw", false, "Use protoc --decode_raw")
        protoFile    = flag.String("proto", "", "proto file to use for decode")
        dontClean    = flag.Bool("dont_clean", false, "Don't remove tmp files on exit")
)

func main() {
     flag.Usage = usage
     flag.Parse()
     var opts []grpc.DialOption
     var cred passCredential

     if !*dontClean {
           // install signal handler for cleaning up tmp files
           sigs := make(chan os.Signal, 1)
           signal.Notify(sigs, os.Interrupt)
           go func() {
              <- sigs
              //cleanup()
              files, _ := filepath.Glob("/tmp/" + tmpFileName)
              for _, f := range files {
                  if err := os.Remove(f); err != nil {
                     fmt.Printf("Failed to remove tmp file %s\n",f)
                  }
              }
              os.Exit(0)
           }()
     }

     // TODO: add TLS support
     opts = append(opts, grpc.WithInsecure())
     opts = append(opts, grpc.WithPerRPCCredentials(cred))

     conn, err := grpc.Dial(*serverAddr, opts...)
     if err != nil {
        log.Fatalf("fail to dial: %v", err)
     }
     defer conn.Close()

     configOperClient := MdtDialin.NewGRPCConfigOperClient(conn)

     reqId := int64(os.Getpid())
     telemetryEncode, ok := telemetryEncoding[*encoding]
     if !ok {
        log.Fatalf("Not supported encoding: %s", *encoding)
     }
     telemetrySubIdstr := *subIds
     telemetryQos := (uint32)(*qos)

     if strings.EqualFold(*operation, "subscribe") {
        subidstrings := strings.Split(telemetrySubIdstr, "#")

        var marking *MdtDialin.QOSMarking
        if telemetryQos != NotConfigured {
           marking = &MdtDialin.QOSMarking{Marking: telemetryQos}
        }

        createSubsArgs := MdtDialin.CreateSubsArgs{
                          ReqId:         reqId,
                          Encode:        telemetryEncode,
                          Subscriptions: subidstrings,
                          Qos:           marking}

        mdtSubscribe(configOperClient, &createSubsArgs)
     } else if strings.EqualFold(*operation, "get-proto") {
        if len(*yangPath) > 0 {
           getProtoArgs := MdtDialin.GetProtoFileArgs{ReqId: reqId, YangPath: *yangPath}
           mdtGetProto(configOperClient, &getProtoArgs)
        } else {
           fmt.Println("No yang path specified!")
        }
     } else {
        fmt.Println("Unsupported operation!")
     }
}

// createSubs rpc to subscribe
func mdtSubscribe(client MdtDialin.GRPCConfigOperClient, args *MdtDialin.CreateSubsArgs) {
     var oFile *os.File
     var tmpfile *os.File
     var commandString string
     var prettyJSON bytes.Buffer
     var err error

     fmt.Printf("mdtSubscribe: Dialin %d subscription %s\n", args.ReqId, args.Subscriptions)

     oFile = os.Stdout
     if len(*outFile) != 0 {
        oFile, _ = os.Create(*outFile)
        defer oFile.Close()
     }

     if *decode_raw || (len(*protoFile) != 0) {
         tmpfile, err = ioutil.TempFile("", tmpFileName)
         if (err != nil) {
             log.Fatal("Failed to create tmp file for writing", err)
         }
         defer os.Remove(tmpfile.Name())
         defer tmpfile.Close()

         if *decode_raw {
             commandString = ProtocRawDecode + "<" + tmpfile.Name()
         } else {
             commandString = ProtocCommandString + *protoFile + "<" + tmpfile.Name()
         }
     }

     stream, err := client.CreateSubs(context.Background(), args)
     if err != nil {
        log.Fatalf("mdtSubscribe: ReqId %d, %v", args.ReqId, err)
     }

     for {
         reply, err := stream.Recv()
         if err == io.EOF {
            fmt.Printf("Subscribe: Got EOF\n\n")
            break
         }
         if err != nil {
            log.Fatalf("Subscribe: ReqId %d, %v", args.ReqId, err)
         }

         if len(reply.Data) == 0 {
            if len(reply.Errors) != 0 {
               fmt.Printf("Subscribe: Received ReqId %d, error:\n%s\n", args.ReqId, reply.Errors)
               break
            }
         } else {
             if args.Encode == 4 {
                err = json.Indent(&prettyJSON, reply.Data, "", "\t")
                if err != nil {
                   fmt.Println("JSON parse error: ", err)
                } else {
                   _, err := oFile.WriteString(string(prettyJSON.Bytes()))
                   if err != nil {
                      fmt.Println(err)
                   }
                }
                continue
             }
             if *decode_raw || (len(*protoFile) != 0) {
                 // use protoc to decode
                 /* Write to tmp file and run protoc command to decode */
                 _, err = tmpfile.Write(reply.Data)
                 out, err := exec.Command("sh", "-c", commandString).CombinedOutput()
                 if err != nil {
                     fmt.Println("Protoc error", err, out)
                 } else {
                     _, err := oFile.WriteString(string(out))
                     if err != nil {
                         fmt.Println(err)
                     }
                     tmpfile.Truncate(0)
                     tmpfile.Seek(0,0)
                 }
             } else {
                 telem := &telemetry.Telemetry{}
                 err = proto.Unmarshal(reply.Data, telem)
                 if (err != nil) {
                     fmt.Println("Failed to unmarshal:", err)
                 }
                 j, _ :=  json.MarshalIndent(telem, "", "  ")
                 _, err = oFile.WriteString(string(j))
             }
         }
     }

}

// Get Proto request
func mdtGetProto(client MdtDialin.GRPCConfigOperClient, args *MdtDialin.GetProtoFileArgs) int64 {
     var oFile *os.File

     stream, err := client.GetProtoFile(context.Background(), args)
     if err != nil {
        log.Fatalf("GetProto: ReqId %d, %v", args.ReqId, err)
        return 0
     }

     oFile = os.Stdout
     if len(*outFile) != 0 {
        oFile, err = os.Create(*outFile)
        defer oFile.Close()
     }

     for {
         reply, err := stream.Recv()
         if err == io.EOF {
            break
         }
         if err != nil {
            log.Fatalf("GetProto: ReqId %d, %v", args.ReqId, err)
            return 0
         }

         if len(reply.Errors) != 0 {
            fmt.Printf("GetProto: ReqId %d, received error: %s\n", args.ReqId, reply.Errors)
            return 0
         } else if reply.ReqId != args.ReqId {
            fmt.Printf("GetProto: mismatch sent ReqID %d, Received ReqId %d\n",
                                         args.ReqId, reply.ReqId)
            return 0
         } else {
            if len(reply.ProtoContent) == 0 {
               fmt.Printf("GetProto: Received ReqId %d \n", reply.ReqId)
            } else {
               _, err := oFile.WriteString(reply.ProtoContent)
               if err != nil {
                  fmt.Println(err)
               }
            }
         }
     }

     return 0
}


type passCredential int
func (passCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
     return map[string]string{
                "username": *username,
                "password": *password,
            }, nil
}

func (passCredential) RequireTransportSecurity() bool {
     return false
}
