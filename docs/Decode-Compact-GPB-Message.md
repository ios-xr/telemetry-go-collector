## Decoding Compact GPB message

IOSXR telemetry supports json, self-describing-gpb and gpb encoding. 
Json is human-readable text format encoded in UTF-8, it can
be consumed without need for any decode. Self-describing-gpb and gpb, both need proto to
decode/unmarshal the message.

Telemetry message defined
in
[telemetry.proto](https://github.com/ios-xr/telemetry-go-collector/blob/master/telemetry/telemetry.proto) can
be used to decode gpb/self-describing-gpb message.

###Self-describing-GPB:
As name suggests self-describing-gpb message carries leaf names and values, its also
reffered to in many documents as key-value GPB. Just telemetry.proto
is needed to decode the message for any sensor-path/model. Repeated
TelemetryField is used to encode keys and content.

###GPB:
Same "Telemetry" message is used to encode GPB message. Keys and
Content include only the values and not the leaf names. To unmarshall
the keys and content part of the message, proto specific to the path
being streamed is needed.
TelemetryGPBTable field is set when GPB encoding is used, this field
contains repeated TelemetryRowGPB, each row containing timestamp,
keys, and content.

Protos are available at
https://github.com/ios-xr/model-driven-telemetry

A Script is available to generate go bindings and plugins for all
these protos.

###Decode logic in the collector:
> if transport is TCP or UDP, 12 byte header in the message has encode type

 ---------------------------------------------
|      MsgType(2 Bytes)  | MsgEncap(2 Bytes)   |
 ---------------------------------------------
|    HdrVersion(2 Bytes) |  Flags(2 Bytes)     |
 ---------------------------------------------
|          Msg Length(4 Bytes)                 |
 ---------------------------------------------

MsgEncap: 1:GPB and 2:JSON
Decode will happen based on the encode type from the header.

> GRPC dialin as well s dailout, collector needs encoding to be passed
> as input argument.

1) If received message is json, write the message to out file as json
pretty print
2) If message is gpb/self-describing-gpb,
   a) if decode_raw or proto is passed as arguments,
      1) write the message to tmp file
      2) use protoc to decode the message
         (minimum version protoc in the PATH has to be 3.3.0)
      3) write the decoded message to out file
   b) Unmarshall the message using Telemetry message
      1) If data_gpb field is not set, this is self-describing-gpb
         message already decoded, write to out file
      2) if compact gpb message,
         a) replace "-" to "_" and ":" to "/" is the encoding_path,
         add it to absolute path passed in --plugin_dir
         b) look for plugin.so under this directory and open it
         c) lookup exported symbols in the plugin
         b) Use exported plugin symbols to unmarshal keys and content
         fields in each of the rows in the message.
         d) write the header and all the rows to out file
