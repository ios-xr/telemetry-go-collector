package main

import (
        "fmt"
        "io"
        "net"
        "bytes"
        "encoding/binary"
)

///////////////////////////////////
//      Msg Header
// ----------------------------
//|      MsgType  | MsgEncap   |
// ----------------------------
//|    HdrVersion |  Flags     |
// ----------------------------
//|          Msg Length        |
// ----------------------------
///////////////////////////////////

func mdtUdpServer(udpPort string) error {
     var err error
     var hdr tcpMsgHdr

     dataChan := make(chan []byte, 10000)
     defer close(dataChan)
     o := &mdtOut{
                outFile:     *outFileName,
                encoding:    *encoding,
                decode_raw:  *decode_raw,
                protoFile:   *protoFile,
                dataChan:     dataChan,
     }

     go o.mdtOutLoop()

     ServerAddr, err := net.ResolveUDPAddr("udp", udpPort)
     if err != nil {
         panic(err)
     }

     // now listen at selected port.
     ServerConn, err := net.ListenUDP("udp", ServerAddr)
     if err != nil {
         panic(err)
     }
     defer ServerConn.Close()
     fmt.Println("UDP server listening at ", udpPort)

     buf := make([]byte, 64*1024)
     for {
         n, addr, err := ServerConn.ReadFromUDP(buf)
         if (err != nil) || (n == 0) {
             if err == io.EOF {
                fmt.Printf(".")
                return err
             } else {
                fmt.Println("Read error:", err, "from", addr)
                continue
             }
         }

         // msg header is 12 bytes
         hdrbuf := bytes.NewReader(buf[:12])
         err = binary.Read(hdrbuf, binary.BigEndian, &hdr)
         //fmt.Printf("From %s received message len: %v encode %v\n",
         //           addr, hdr.Msglen, hdr.MsgEncap)

         // set the encoding from header
         o.mdtOutSetEncoding(mdtGetEncodeStr(hdr.MsgEncap))
         // write to data channel
         dataChan <- buf[12:n]
     }

     return nil
}
