package main

import (
       "os"
       "os/exec"
       "io/ioutil"
       "log"
       "fmt"
       "plugin"
       "bytes"
       "encoding/json"
       "unsafe"
       "strings"

       "github.com/golang/protobuf/jsonpb"
       "github.com/golang/protobuf/proto"

       "github.com/ios-xr/telemetry-go-collector/telemetry"
       MdtDialin "github.com/ios-xr/telemetry-go-collector/mdt_grpc_dialin"

)

const ProtocRawDecode string     = "protoc --decode_raw "
const ProtocCommandString string = "protoc --decode=Telemetry "

///////////////////////////////////////////////////////////////////////
///////     O U T P U T   M E S S A G E   H A N D L E R         ///////
///////////////////////////////////////////////////////////////////////
// message handler
// 1) if encoding is json, pretty print to out file
// 2) if decode_raw is set,
//    a) write the message to tmp file, execute protoc command with --decode_raw option
//       protoc --decode_raw < tmpfile
//    b) print to output file
//    protoc is expected to be present in the PATH, minimum version 3.0.0
// 3) if proto file specified in arguments,
//    a) write the message to tmp file, execute protoc with --decode option
//       protoc --decode=Telemetry <proto> < tmpfile
//    b) write the output to out file
// 4) if none of above, unmarshal message using telemetry.proto
//    a) if self-describing-gpb message, write to outfile
//    b) if gpb(compact),
//       i) using encoding_path find the plugin for the proto
//       ii) if found, decode key and content of all rows using exported symbols from plugin
//           write telemetry header and rows to out file
//       iii) if not found, write the raw content to out file
//
func mdtOutLoop(dataChan <-chan *MdtDialin.CreateSubsReply, encoding int64) {
     var oFile *os.File
     var tmpfile *os.File
     var commandString string
     var prettyJSON bytes.Buffer
     var err error

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

     for {
         msg, ok := <-dataChan

         if !ok {
            //channel might have been closed
            fmt.Println("Done with output loop..")
            break
         }
         if encoding == 4 {
                err = json.Indent(&prettyJSON, msg.Data, "", "\t")
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
                 _, err = tmpfile.Write(msg.Data)
                 out, err := exec.Command("sh", "-c", commandString).CombinedOutput()
                 if err != nil {
                     fmt.Println("Protoc error", err, out)
                     fmt.Println("Make sure protoc version in the $PATH is atleast 3.3.0")
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

                 err = proto.Unmarshal(msg.Data, telem)
                 if (err != nil) {
                     fmt.Println("Failed to unmarshal:", err)
                 }
                 if telem.GetDataGpb() != nil {
                     //this is gpb message
                     mdtDumpGPBMessage(telem, oFile)
                 } else {
                   j, _ :=  json.MarshalIndent(telem, "", "  ")
                   _, err = oFile.WriteString(string(j))
                 }
         }
     }
}

// plugin info
type gpbPluginInfo struct {
     plug          *plugin.Plugin
     decodedKeys     proto.Message
     decodedContent  proto.Message
}

var pluginTbl map[string]*gpbPluginInfo

// Message type (including header and rows) used for serialisation
type msgToSerialise struct {
     //Source    string
     Telemetry *json.RawMessage
     Rows      []*rowToSerialise `json:"Rows,omitempty"`
}

//
// Message row type used for serialisation
type rowToSerialise struct {
     Timestamp uint64
     Keys      *json.RawMessage
     Content   *json.RawMessage
}

// try to find plugin to decode the gpb content
func mdtDumpGPBMessage(copy *telemetry.Telemetry, oFile *os.File) {
     var err error
     var s msgToSerialise

     if pluginTbl == nil {
        pluginTbl = make(map[string]*gpbPluginInfo)
     }
     gpbPlugin := mdtGetPlugin(copy.EncodingPath)

     if gpbPlugin == nil {
        j, _ :=  json.MarshalIndent(copy, "", "  ")
        _, err = oFile.WriteString(string(j))
        if err != nil {
           fmt.Println("Error writing the output", err)
        }
        return
     }

     marshaller := &jsonpb.Marshaler{
                   EmitDefaults:       true,
                   OrigName:           true}


     for _, row := range copy.GetDataGpb().GetRow() {
         err = proto.Unmarshal(row.Keys, gpbPlugin.decodedKeys)
         if (err != nil) {
            fmt.Println("plugin unmarshal failed", err)
            return
         }

         err = proto.Unmarshal(row.Content, gpbPlugin.decodedContent)
         if (err != nil) {
            fmt.Println("plugin unmarshal failed", err)
            return
         }

         var keys      json.RawMessage
         var content   json.RawMessage

         decodedContentJSON, err := marshaller.MarshalToString(gpbPlugin.decodedContent)
         if err != nil {
             fmt.Println(err)
         } else {
            content = json.RawMessage(decodedContentJSON)
         }

         decodedKeysJSON, err := marshaller.MarshalToString(gpbPlugin.decodedKeys)
         if err != nil {
             fmt.Println(err)
         } else {
             keys = json.RawMessage(decodedKeysJSON)
         }

         s.Rows = append(s.Rows, &rowToSerialise{row.Timestamp, &keys, &content})
     }

     copy.DataGpb = nil
     telemetryJSON, err := marshaller.MarshalToString(copy)
     if err != nil {
        return
     }
     telemetryJSONRaw := json.RawMessage(telemetryJSON)
     s.Telemetry = &telemetryJSONRaw

     b, err := json.Marshal(s)
     if err != nil {
        fmt.Errorf("Marshalling collected content, [%+v][%+v]",
                          s, err)
     }

     var out bytes.Buffer
     json.Indent(&out, b, "", "    ")
     //fmt.Println(out.String())
     _, err = oFile.WriteString(out.String())
     if err != nil {
        fmt.Println("Error writing the output", err)
     }
}

type Plug struct {
     Path    string
     err     string
     _       chan struct{}
     Symbols map[string]interface{}
}

//lookup plugin and exported symbols from plugin.
func mdtGetPlugin(encodingPath string) *gpbPluginInfo {
     var decodedKeys     proto.Message
     var decodedContent  proto.Message

     p, ok := pluginTbl[encodingPath]
     if !ok {
        pluginFileName := mdtGetPluginFileName(encodingPath)
        if len(*pluginDir) > 0 && !strings.HasSuffix(*pluginDir, "/") {
           *pluginDir = *pluginDir + "/"
        }
        plug, err := plugin.Open(*pluginDir + pluginFileName)
        if (err != nil) {
           fmt.Println("plugin open failed", err)
           return nil
        }
        //fmt.Printf("%+v", plug)
        pl := (*Plug)(unsafe.Pointer(plug))
        for name, pointer := range pl.Symbols {
            if strings.HasSuffix(name, "_KEYS") {
               decodedKeys, _ = pointer.(proto.Message)
            } else {
               decodedContent, _ = pointer.(proto.Message)
            }
        }

        p = &gpbPluginInfo{plug, decodedKeys, decodedContent}
        pluginTbl[encodingPath] = p
     }
     return p
}

func mdtGetPluginFileName(encodingPath string) string {
     str := strings.ToLower(encodingPath)
     str = strings.Replace(str, "-", "_", -1)
     return strings.Replace(str, ":", "/", 1) + "/plugin/plugin.so"
}
