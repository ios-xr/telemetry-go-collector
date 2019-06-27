package main

import (
        "fmt"
        "io"
        "net"
        "bytes"
        "encoding/binary"

        "github.com/ios-xr/telemetry-go-collector/telemetry_decode"
)

///////////////////////////////////
//      TCP Msg Header
// ----------------------------
//|      MsgType  | MsgEncap   |
// ----------------------------
//|    HdrVersion |  Flags     |
// ----------------------------
//|          Msg Length        |
// ----------------------------
///////////////////////////////////

type encapSTHdrMsgType uint16
const (
      ENC_ST_HDR_MSG_TYPE_UNSED encapSTHdrMsgType = iota
      ENC_ST_HDR_MSG_TYPE_TELEMETRY_DATA
      ENC_ST_HDR_MSG_TYPE_HEARTBEAT
)

type encapSTHdrMsgEncap uint16
const (
      ENC_ST_HDR_MSG_ENCAP_UNSED encapSTHdrMsgEncap = iota
      ENC_ST_HDR_MSG_ENCAP_GPB
      ENC_ST_HDR_MSG_ENCAP_JSON
)

type tcpMsgHdr struct {
     MsgType       encapSTHdrMsgType
     MsgEncap      encapSTHdrMsgEncap
     MsgHdrVersion uint16
     Msgflag       uint16
     Msglen        uint32
}

type tcpSession struct {
     conn          *net.TCPConn
     // tcp msg header
     hdr           []byte
}

func mdtGetEncodeStr(enc encapSTHdrMsgEncap) string {
     switch (enc) {
     case ENC_ST_HDR_MSG_ENCAP_GPB:
         return "gpb"
     case ENC_ST_HDR_MSG_ENCAP_JSON:
         return "json"
     default:
         return "Unknown"
     }
}

func (s *tcpSession) handleConnection() {
     var hdr tcpMsgHdr
     var buf []byte

     dataChan := make(chan []byte, 10000)
     defer close(dataChan)
     o := &telemetry_decode.MdtOut{
                        OutFile:     *outFileName,
                        Encoding:    *encoding,
                        Decode_raw:  *decode_raw,
                        DontClean:   *dontClean,
                        ProtoFile:   *protoFile,
                        PluginDir:   *pluginDir,
                        DataChan:     dataChan,
     }

     go o.MdtOutLoop()

     for {
         // read header for tcp message.
         _, err := io.ReadFull(s.conn, s.hdr)
         if err != nil {
             if err == io.EOF {
                fmt.Printf(".")
                return
             } else {
                fmt.Println("Read error : ", err)
                continue       // should return? allows router to reconnect
             }
         }
         hdrbuf := bytes.NewReader(s.hdr)
         err = binary.Read(hdrbuf, binary.BigEndian, &hdr)
         //fmt.Printf("Received message len: %v encode %v\n", hdr.Msglen, hdr.MsgEncap)
         buf = make([]byte, hdr.Msglen)

         // read rest of the tcp message using length from header.
         _, err = io.ReadFull(s.conn, buf)
         if err != nil {
            fmt.Println(err)
            continue
         }

         // set the encoding from header
         o.MdtOutSetEncoding(mdtGetEncodeStr(hdr.MsgEncap))
         // write to the data channel
         dataChan <- buf
     }
}

func mdtTcpServer(tcpPort string) error {
     var err error
     var hdr tcpMsgHdr

     ServerAddr, err := net.ResolveTCPAddr("tcp", tcpPort)
     if err != nil {
         panic(err)
     }

     // now listen at selected port.
     listener, err := net.ListenTCP("tcp", ServerAddr)
     if err != nil {
         panic(err)
     }
     defer listener.Close()

     fmt.Println("TCP server listening at ", tcpPort)
     for {
         serverConn, err := listener.AcceptTCP()
         if err != nil {
             panic(err)
         }
         defer serverConn.Close()
         fmt.Printf("Session connected from %s\n", serverConn.RemoteAddr())

         s := new(tcpSession)
         s.conn  = serverConn
         s.hdr = make([]byte, binary.Size(hdr))

         go s.handleConnection()
     }

     return nil
}
