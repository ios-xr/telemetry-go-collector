package telemetry_decode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"plugin"
	"strconv"
	"strings"
	"unsafe"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"

	"github.com/ios-xr/telemetry-go-collector/telemetry"
)

const ProtocRawDecode string = "protoc --decode_raw "
const ProtocCommandString string = "protoc --decode=Telemetry "
const tmpFileName = "telemetry-msg-*.dat"

// /////////////////////////////////////////////////////////////////////
// /////     O U T P U T   M E S S A G E   H A N D L E R         ///////
// /////////////////////////////////////////////////////////////////////
type MdtOut struct {
	OutFile    string
	Encoding   string
	Decode_raw bool
	DontClean  bool
	ProtoFile  string
	PluginDir  string
	PluginFile string
	DataChan   <-chan []byte
	oFile      *os.File
	esClient   *elasticsearch.Client
}

// message handler
//  1. if encoding is json, pretty print to out file
//  2. if decode_raw is set,
//     a) write the message to tmp file, execute protoc command with --decode_raw option
//     protoc --decode_raw < tmpfile
//     b) print to output file
//     protoc is expected to be present in the PATH, minimum version 3.0.0
//  3. if proto file specified in arguments,
//     a) write the message to tmp file, execute protoc with --decode option
//     protoc --decode=Telemetry <proto> < tmpfile
//     b) write the output to out file
//  4. if none of above, unmarshal message using telemetry.proto
//     a) if self-describing-gpb message, write to outfile
//     b) if gpb(compact),
//     i) using encoding_path find the plugin for the proto
//     ii) if found, decode key and content of all rows using exported symbols from plugin
//     write telemetry header and rows to out file
//     iii) if not found, write the raw content to out file
func (o *MdtOut) MdtOutLoop() {
	var err error

	tmpFile, commandString := o.mdtPrepareDecoding()
	if tmpFile != nil {
		if !o.DontClean {
			defer os.Remove(tmpFile.Name())
		}
		defer tmpFile.Close()
	}
	if o.oFile != nil {
		defer o.oFile.Close()
		fmt.Println("Out file:", o.oFile.Name())
	}

	for {
		data, ok := <-o.DataChan

		if !ok {
			//channel might have been closed
			fmt.Println("Done with output loop..")
			break
		}
		if o.Encoding == "json" {
			o.mdtDumpJsonMessage(data)
		} else if o.Decode_raw || (len(o.ProtoFile) != 0) {
			// use protoc to decode
			/* Write to tmp file and run protoc command to decode */
			_, err = tmpFile.Write(data)
			out, err := exec.Command("sh", "-c", commandString).CombinedOutput()
			if err != nil {
				fmt.Println("Protoc error", err, out)
				fmt.Println("Make sure protoc version in the $PATH is atleast 3.3.0")
			} else {
				_, err := o.oFile.WriteString(string(out))
				if err != nil {
					fmt.Println(err)
				}
				tmpFile.Truncate(0)
				tmpFile.Seek(0, 0)
			}
		} else {
			telem := &telemetry.Telemetry{}

			err = proto.Unmarshal(data, telem)
			if err != nil {
				fmt.Println("Failed to unmarshal:", err)
			}
			if telem.GetDataGpb() != nil {
				//this is gpb message
				o.mdtDumpGPBMessage(telem)
			} else {
				o.mdtDumpKVGPBMessage(telem)
			}
		}
	}
}

func (o *MdtOut) MdtOutSetEncoding(encoding string) {
	o.Encoding = encoding
}

// json walk and dump
func (o *MdtOut) mdtDumpJsonMessage(copy []byte) {
	var prettyJSON bytes.Buffer

	if o.esClient != nil {
		m := make(map[string]interface{})
		err := json.Unmarshal(copy, &m)
		if err != nil {
			log.Fatal(err)
		}

		for i, row := range m["data_json"].([]interface{}) {
			j, _ := json.Marshal(row)
			o.elasticSearchOutput(string(j),
				m["encoding_path"].(string),
				m["node_id_str"].(string),
				m["collection_id"].(string),
				i)
		}
	} else {
		err := json.Indent(&prettyJSON, copy, "", "\t")
		if err != nil {
			fmt.Println("JSON parse error: ", err)
		} else {
			_, err = o.oFile.WriteString(string(prettyJSON.Bytes()))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// kvgpb walk and dump
func (o *MdtOut) mdtDumpKVGPBMessage(copy *telemetry.Telemetry) {

	if o.esClient != nil {
		for i, row := range copy.GetDataGpbkv() {
			j, _ := json.Marshal(row)
			o.elasticSearchOutput(string(j),
				copy.GetEncodingPath(),
				copy.GetNodeIdStr(),
				strconv.FormatUint(copy.GetCollectionId(), 10),
				i)
		}
	} else {
		j, _ := json.MarshalIndent(copy, "", "  ")
		_, err := o.oFile.WriteString(string(j))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// plugin info
type gpbPluginInfo struct {
	plug           *plugin.Plugin
	decodedKeys    proto.Message
	decodedContent proto.Message
}

var pluginTbl map[string]*gpbPluginInfo

// Message type (including header and rows) used for serialisation
type msgToSerialise struct {
	//Source    string
	Telemetry *json.RawMessage
	Rows      []*rowToSerialise `json:"Rows,omitempty"`
}

// Message row type used for serialisation
type rowToSerialise struct {
	Timestamp uint64
	Keys      *json.RawMessage
	Content   *json.RawMessage
}

// try to find plugin to decode the gpb content
func (o *MdtOut) mdtDumpGPBMessage(copy *telemetry.Telemetry) {
	var err error
	var s msgToSerialise

	gpbPlugin := mdtGetPlugin(copy.EncodingPath, o.PluginDir, o.PluginFile)

	if gpbPlugin == nil {
		j, _ := json.MarshalIndent(copy, "", "  ")
		_, err = o.oFile.WriteString(string(j))
		if err != nil {
			fmt.Println("Error writing the output", err)
		}
		return
	}

	marshaller := &jsonpb.Marshaler{
		EmitDefaults: true,
		OrigName:     true}

	for i, row := range copy.GetDataGpb().GetRow() {
		err = proto.Unmarshal(row.Keys, gpbPlugin.decodedKeys)
		if err != nil {
			fmt.Println("plugin unmarshal failed", err)
			return
		}

		err = proto.Unmarshal(row.Content, gpbPlugin.decodedContent)
		if err != nil {
			fmt.Println("plugin unmarshal failed", err)
			return
		}

		var keys json.RawMessage
		var content json.RawMessage

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

		if o.esClient != nil {
			b, _ := json.Marshal(&rowToSerialise{row.Timestamp, &keys, &content})
			o.elasticSearchOutput(string(b),
				copy.GetEncodingPath(),
				copy.GetNodeIdStr(),
				strconv.FormatUint(copy.GetCollectionId(), 10),
				i)
		}
	}

	if o.esClient == nil {
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
		_, err = o.oFile.WriteString(out.String())
		if err != nil {
			fmt.Println("Error writing the output", err)
		}
	}

}

type Plug struct {
	Path    string
	err     string
	_       chan struct{}
	Symbols map[string]interface{}
}

// lookup plugin and exported symbols from plugin.
// Two ways to provide plugin to use to the collector,
//  1. plugin_dir, plugin directory where to look for, in this case plugin.so is
//     expected to be present in directory hirarchy that matches the sensor-path
//  2. plugin, plugin file to use, it will be used to lookup sysmbol using sensor-path
func mdtGetPlugin(encodingPath string, pluginDir string, pluginFile string) *gpbPluginInfo {
	var decodedKeys proto.Message
	var decodedContent proto.Message
	var plug *plugin.Plugin
	var err error

	if pluginTbl == nil {
		pluginTbl = make(map[string]*gpbPluginInfo)
	}

	p, ok := pluginTbl[encodingPath]
	if !ok {
		if pluginFile != "" {
			plug, err = plugin.Open(pluginFile)
			if err != nil {
				fmt.Println("plugin open failed", err)
				return nil
			}
			symStr := strings.ToLower(encodingPath)
			symStr = strings.Replace(symStr, "-", "_", -1)
			symStr = strings.Replace(symStr, "/", "_", -1)
			symStr = strings.Replace(symStr, ":", "_", -1)

			symKey, err := plug.Lookup("KEYS_" + symStr)
			if err != nil {
				fmt.Println("plugin symbol not found", err)
				return nil
			}
			symContent, err := plug.Lookup("CONTENT_" + symStr)
			if err != nil {
				fmt.Println("plugin symbol not found", err)
				return nil
			}
			decodedKeys, _ = symKey.(proto.Message)
			decodedContent, _ = symContent.(proto.Message)
		} else {
			pluginFileName := mdtGetPluginFileName(encodingPath)
			if len(pluginDir) > 0 && !strings.HasSuffix(pluginDir, "/") {
				pluginDir = pluginDir + "/"
			}
			plug, err = plugin.Open(pluginDir + pluginFileName)
			if err != nil {
				fmt.Println("plugin open failed", err)
				return nil
			}

			pl := (*Plug)(unsafe.Pointer(plug))
			for name, pointer := range pl.Symbols {
				if strings.HasSuffix(name, "_KEYS") {
					decodedKeys, _ = pointer.(proto.Message)
				} else {
					decodedContent, _ = pointer.(proto.Message)
				}
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

// create tmp and output file
func (o *MdtOut) mdtPrepareDecoding() (*os.File, string) {
	var commandString string
	var err error

	outN := strings.SplitN(o.OutFile, ":", 2)
	if outN[0] == "elasticsearch" {
		o.esClient = elasticSearchClientInit("http://" + outN[1])
		return nil, ""
	}

	// create/open output file
	if len(o.OutFile) != 0 {
		o.oFile, err = ioutil.TempFile(".", o.OutFile)
		if err != nil {
			log.Fatal("Failed to create output file for writing", err)
		}
	} else {
		o.oFile = os.Stdout
	}

	if o.Decode_raw || (len(o.ProtoFile) != 0) {
		// temp file to write message to for decoding
		tmpFile, err := ioutil.TempFile("", tmpFileName)
		if err != nil {
			log.Fatal("Failed to create tmp file for writing", err)
		}

		// proto command to use for decoding gpb message
		if o.Decode_raw {
			commandString = ProtocRawDecode + " < " + tmpFile.Name()
		} else {
			commandString = ProtocCommandString + o.ProtoFile + " < " + tmpFile.Name()
		}
		return tmpFile, commandString
	}

	return nil, ""
}

var replacer = strings.NewReplacer("/", "_", ":", "_")

// elastic search functions
func (o *MdtOut) elasticSearchOutput(data string, index string, node_id string, coll_id string, i int) {
	if o.esClient != nil {

		req := esapi.IndexRequest{
			Index:      strings.ToLower(replacer.Replace(index)),
			DocumentID: fmt.Sprintf("%s.%s.%d", node_id, coll_id, i),
			Body:       strings.NewReader(data),
			Refresh:    "true",
		}

		// Perform the request with the client.
		res, err := req.Do(context.Background(), o.esClient)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document", res.Status())
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and indexed document version.
				log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}

	} else {
		// last resort, dump to stdout
		fmt.Println(data)
	}
}

func elasticSearchClientInit(esServer string) *elasticsearch.Client {
	var r map[string]interface{}

	//esServer := "http://localhost:9200"

	// Initialize a client with the default settings.
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//es, err := elasticsearch.NewDefaultClient()
	cfg := elasticsearch.Config{
		Addresses: []string{
			esServer,
		},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
	return es
}
