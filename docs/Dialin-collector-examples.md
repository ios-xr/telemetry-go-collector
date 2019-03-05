### Subscribe rpc
```
 $ telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab -encoding self-describing-gpb -qos 10
mdtSubscribe: Dialin ReqId 2752 sub_idstr cdp-neighbor
node_id_str: "adithyas-1"
subscription_id_str: "cdp-neighbor"
encoding_path: "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail"
collection_id: 37
collection_start_time: 1551807027918
msg_timestamp: 1551807027918
data_gpbkv {
  timestamp: 1551807027925
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
        uint32_value: 150
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
          string_value: " 7.0.1.122I"
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
data_gpbkv {
  timestamp: 1551807027925
  fields {
    name: "keys"
    fields {
      name: "node-name"
      string_value: "0/0/CPU0"
    }
    fields {
      name: "interface-name"
      string_value: "GigabitEthernet0/0/0/2"
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
        string_value: "GigabitEthernet0/0/0/2"
      }
      fields {
        name: "device-id"
        string_value: "adithyas-2"
      }
      fields {
        name: "port-id"
        string_value: "GigabitEthernet0/0/0/2"
      }
      fields {
        name: "header-version"
        uint32_value: 2
      }
      fields {
        name: "hold-time"
        uint32_value: 158
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
                string_value: "2.2.2.2"
              }
            }
          }
        }
        fields {
          name: "version"
          string_value: " 7.0.1.122I"
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
data_gpbkv {
  timestamp: 1551807027925
  fields {
    name: "keys"
    fields {
      name: "node-name"
      string_value: "0/0/CPU0"
    }
    fields {
      name: "interface-name"
      string_value: "GigabitEthernet0/0/0/1"
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
        string_value: "GigabitEthernet0/0/0/1"
      }
      fields {
        name: "device-id"
        string_value: "adithyas-2"
      }
      fields {
        name: "port-id"
        string_value: "GigabitEthernet0/0/0/1"
      }
      fields {
        name: "header-version"
        uint32_value: 2
      }
      fields {
        name: "hold-time"
        uint32_value: 154
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
                string_value: "5.0.0.2"
              }
            }
          }
        }
        fields {
          name: "version"
          string_value: " 7.0.1.122I"
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
collection_end_time: 1551807027926

^C $
 $ telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab -encoding gpb -qos 10 -proto telemetry.proto
mdtSubscribe: Dialin ReqId 2912 sub_idstr cdp-neighbor
node_id_str: "adithyas-1"
subscription_id_str: "cdp-neighbor"
encoding_path: "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail"
collection_id: 39
collection_start_time: 1551807125694
msg_timestamp: 1551807125694
data_gpb {
  row {
    timestamp: 1551807125700
    keys: "\n\0100/0/CPU0\022\026GigabitEthernet0/0/0/0\032\nadithyas-2"
    content: "\222\003\256\001\n\026GigabitEthernet0/0/0/0\022\nadithyas-2\032\026GigabitEthernet0/0/0/0 \002(\254\0012\001R:\022cisco IOS-XRv 9000BT\n(\n\021\n\017\n\004ipv4\022\0074.0.0.2\n\023\n\021\n\004ipv6\032\t2002::1:2\022\013 7.0.1.122I(\0002\rcdp-dplx-none:\nadithyas-2"
  }
  row {
    timestamp: 1551807125700
    keys: "\n\0100/0/CPU0\022\026GigabitEthernet0/0/0/2\032\nadithyas-2"
    content: "\222\003\230\001\n\026GigabitEthernet0/0/0/2\022\nadithyas-2\032\026GigabitEthernet0/0/0/2 \002(x2\001R:\022cisco IOS-XRv 9000B?\n\023\n\021\n\017\n\004ipv4\022\0072.2.2.2\022\013 7.0.1.122I(\0002\rcdp-dplx-none:\nadithyas-2"
  }
  row {
    timestamp: 1551807125700
    keys: "\n\0100/0/CPU0\022\026GigabitEthernet0/0/0/1\032\nadithyas-2"
    content: "\222\003\231\001\n\026GigabitEthernet0/0/0/1\022\nadithyas-2\032\026GigabitEthernet0/0/0/1 \002(\260\0012\001R:\022cisco IOS-XRv 9000B?\n\023\n\021\n\017\n\004ipv4\022\0075.0.0.2\022\013 7.0.1.122I(\0002\rcdp-dplx-none:\nadithyas-2"
  }
}
collection_end_time: 1551807125701

^C $
 $ telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab -encoding gpb -qos 10 roto cdp_neighbor_compact.proto 
mdtSubscribe: Dialin ReqId 2948 sub_idstr cdp-neighbor
node_id_str: "adithyas-1"
subscription_id_str: "cdp-neighbor"
encoding_path: "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail"
collection_id: 40
collection_start_time: 1551807148466
msg_timestamp: 1551807148466
data_gpb {
  row {
    timestamp: 1551807148472
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
        hold_time: 150
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
          version: " 7.0.1.122I"
          duplex: "cdp-dplx-none"
          system_name: "adithyas-2"
        }
      }
    }
  }
  row {
    timestamp: 1551807148472
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
          version: " 7.0.1.122I"
          duplex: "cdp-dplx-none"
          system_name: "adithyas-2"
        }
      }
    }
  }
  row {
    timestamp: 1551807148472
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
          version: " 7.0.1.122I"
          duplex: "cdp-dplx-none"
          system_name: "adithyas-2"
        }
      }
    }
  }
}
collection_end_time: 1551807148481

^C $
 $ telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab -encoding gpb -qos 10 -decode_raw
mdtSubscribe: Dialin ReqId 3000 sub_idstr cdp-neighbor
1: "adithyas-1"
3: "cdp-neighbor"
6: "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail"
7: "2019-01-09"
8: 41
9: 1551807181158
10: 1551807181158
13: 1551807181165
12 {
  1 {
    1: 1551807181164
    10 {
      1: "0/0/CPU0"
      2: "GigabitEthernet0/0/0/0"
      3: "adithyas-2"
    }
    11 {
      50 {
        1: "GigabitEthernet0/0/0/0"
        2: "adithyas-2"
        3: "GigabitEthernet0/0/0/0"
        4: 2
        5: 177
        6: "R"
        7: "cisco IOS-XRv 9000"
        8 {
          1 {
            1 {
              1 {
                1: "ipv4"
                2: "4.0.0.2"
              }
            }
            1 {
              1 {
                1: "ipv6"
                3: "2002::1:2"
              }
            }
          }
          2: " 7.0.1.122I"
          5: 0
          6: "cdp-dplx-none"
          7: "adithyas-2"
        }
      }
    }
  }
  1 {
    1: 1551807181164
    10 {
      1: "0/0/CPU0"
      2: "GigabitEthernet0/0/0/2"
      3: "adithyas-2"
    }
    11 {
      50 {
        1: "GigabitEthernet0/0/0/2"
        2: "adithyas-2"
        3: "GigabitEthernet0/0/0/2"
        4: 2
        5: 124
        6: "R"
        7: "cisco IOS-XRv 9000"
        8 {
          1 {
            1 {
              1 {
                1: "ipv4"
                2: "2.2.2.2"
              }
            }
          }
          2: " 7.0.1.122I"
          5: 0
          6: "cdp-dplx-none"
          7: "adithyas-2"
        }
      }
    }
  }
  1 {
    1: 1551807181164
    10 {
      1: "0/0/CPU0"
      2: "GigabitEthernet0/0/0/1"
      3: "adithyas-2"
    }
    11 {
      50 {
        1: "GigabitEthernet0/0/0/1"
        2: "adithyas-2"
        3: "GigabitEthernet0/0/0/1"
        4: 2
        5: 121
        6: "R"
        7: "cisco IOS-XRv 9000"
        8 {
          1 {
            1 {
              1 {
                1: "ipv4"
                2: "5.0.0.2"
              }
            }
          }
          2: " 7.0.1.122I"
          5: 0
          6: "cdp-dplx-none"
          7: "adithyas-2"
        }
      }
    }
  }
}
^C $
 $ /ws/adithyas-sjc/collector/bin/telemetry_dialin_collector -server "192.168.122.157:57500" -subscription cdp-neighbor -oper subscribe -username root -password lab                        
mdtSubscribe: Dialin ReqId 3063 sub_idstr cdp-neighbor
{
        "node_id_str": "adithyas-1",
        "subscription_id_str": "cdp-neighbor",
        "encoding_path": "Cisco-IOS-XR-cdp-oper:cdp/nodes/node/neighbors/details/detail",
        "collection_id": "42",
        "collection_start_time": "1551807230354",
        "msg_timestamp": "1551807230359",
        "data_json": [
                {
                        "timestamp": "1551807230358",
                        "keys": [
                                {
                                        "node-name": "0/0/CPU0"
                                },
                                {
                                        "interface-name": "GigabitEthernet0/0/0/0"
                                },
                                {
                                        "device-id": "adithyas-2"
                                }
                        ],
                        "content": {
                                "cdp-neighbor": [
                                        {
                                                "receiving-interface-name": "GigabitEthernet0/0/0/0",
                                                "device-id": "adithyas-2",
                                                "port-id": "GigabitEthernet0/0/0/0",
                                                "header-version": 2,
                                                "hold-time": 128,
                                                "capabilities": "R",
                                                "platform": "cisco IOS-XRv 9000",
                                                "detail": {
                                                        "network-addresses": {
                                                                "cdp-addr-entry": [
                                                                        {
                                                                                "address": {
                                                                                        "address-type": "ipv4",
                                                                                        "ipv4-address": "4.0.0.2"
                                                                                }
                                                                        },
                                                                        {
                                                                                "address": {
                                                                                        "address-type": "ipv6",
                                                                                        "ipv6-address": "2002::1:2"
                                                                                }
                                                                        }
                                                                ]
                                                        },
                                                        "version": " 7.0.1.122I",
                                                        "native-vlan": 0,
                                                        "duplex": "cdp-dplx-none",
                                                        "system-name": "adithyas-2"
                                                }
                                        }
                                ]
                        }
                },
                {
                        "timestamp": "1551807230359",
                        "keys": [
                                {
                                        "node-name": "0/0/CPU0"
                                },
                                {
                                        "interface-name": "GigabitEthernet0/0/0/2"
                                },
                                {
                                        "device-id": "adithyas-2"
                                }
                        ],
                        "content": {
                                "cdp-neighbor": [
                                        {
                                                "receiving-interface-name": "GigabitEthernet0/0/0/2",
                                                "device-id": "adithyas-2",
                                                "port-id": "GigabitEthernet0/0/0/2",
                                                "header-version": 2,
                                                "hold-time": 135,
                                                "capabilities": "R",
                                                "platform": "cisco IOS-XRv 9000",
                                                "detail": {
                                                        "network-addresses": {
                                                                "cdp-addr-entry": [
                                                                        {
                                                                                "address": {
                                                                                        "address-type": "ipv4",
                                                                                        "ipv4-address": "2.2.2.2"
                                                                                }
                                                                        }
                                                                ]
                                                        },
                                                        "version": " 7.0.1.122I",
                                                        "native-vlan": 0,
                                                        "duplex": "cdp-dplx-none",
                                                        "system-name": "adithyas-2"
                                                }
                                        }
                                ]
                        }
                },
                {
                        "timestamp": "1551807230359",
                        "keys": [
                                {
                                        "node-name": "0/0/CPU0"
                                },
                                {
                                        "interface-name": "GigabitEthernet0/0/0/1"
                                },
                                {
                                        "device-id": "adithyas-2"
                                }
                        ],
                        "content": {
                                "cdp-neighbor": [
                                        {
                                                "receiving-interface-name": "GigabitEthernet0/0/0/1",
                                                "device-id": "adithyas-2",
                                                "port-id": "GigabitEthernet0/0/0/1",
                                                "header-version": 2,
                                                "hold-time": 131,
                                                "capabilities": "R",
                                                "platform": "cisco IOS-XRv 9000",
                                                "detail": {
                                                        "network-addresses": {
                                                                "cdp-addr-entry": [
                                                                        {
                                                                                "address": {
                                                                                        "address-type": "ipv4",
                                                                                        "ipv4-address": "5.0.0.2"
                                                                                }
                                                                        }
                                                                ]
                                                        },
                                                        "version": " 7.0.1.122I",
                                                        "native-vlan": 0,
                                                        "duplex": "cdp-dplx-none",
                                                        "system-name": "adithyas-2"
                                                }
                                        }
                                ]
                        }
                }
        ],
        "collection_end_time": "1551807230359"
}^C $
```

## Get Proto example:
```
 $ telemetry_dialin_collector -server "192.168.122.157:57500" -oper get-proto -username root -password lab -yang_path Cisco-IOS-XR-infra-statsd-oper:infra-statistics/interfaces/interface/latest/generic-counters
//                                 Apache License
//                           Version 2.0, January 2004
//                        http://www.apache.org/licenses/
//
//   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
//
//   1. Definitions.
//
//      "License" shall mean the terms and conditions for use, reproduction,
//      and distribution as defined by Sections 1 through 9 of this document.
//
//      "Licensor" shall mean the copyright owner or entity authorized by
//      the copyright owner that is granting the License.
//
//      "Legal Entity" shall mean the union of the acting entity and all
//      other entities that control, are controlled by, or are under common
//      control with that entity. For the purposes of this definition,
//      "control" means (i) the power, direct or indirect, to cause the
//      direction or management of such entity, whether by contract or
//      otherwise, or (ii) ownership of fifty percent (50%) or more of the
//      outstanding shares, or (iii) beneficial ownership of such entity.
//
//      "You" (or "Your") shall mean an individual or Legal Entity
//      exercising permissions granted by this License.
//
//      "Source" form shall mean the preferred form for making modifications,
//      including but not limited to software source code, documentation
//      source, and configuration files.
//
//      "Object" form shall mean any form resulting from mechanical
//      transformation or translation of a Source form, including but
//      not limited to compiled object code, generated documentation,
//      and conversions to other media types.
//
//      "Work" shall mean the work of authorship, whether in Source or
//      Object form, made available under the License, as indicated by a
//      copyright notice that is included in or attached to the work
//      (an example is provided in the Appendix below).
//
//      "Derivative Works" shall mean any work, whether in Source or Object
//      form, that is based on (or derived from) the Work and for which the
//      editorial revisions, annotations, elaborations, or other modifications
//      represent, as a whole, an original work of authorship. For the purposes
//      of this License, Derivative Works shall not include works that remain
//      separable from, or merely link (or bind by name) to the interfaces of,
//      the Work and Derivative Works thereof.
//
//      "Contribution" shall mean any work of authorship, including
//      the original version of the Work and any modifications or additions
//      to that Work or Derivative Works thereof, that is intentionally
//      submitted to Licensor for inclusion in the Work by the copyright owner
//      or by an individual or Legal Entity authorized to submit on behalf of
//      the copyright owner. For the purposes of this definition, "submitted"
//      means any form of electronic, verbal, or written communication sent
//      to the Licensor or its representatives, including but not limited to
//      communication on electronic mailing lists, source code control systems,
//      and issue tracking systems that are managed by, or on behalf of, the
//      Licensor for the purpose of discussing and improving the Work, but
//      excluding communication that is conspicuously marked or otherwise
//      designated in writing by the copyright owner as "Not a Contribution."
//
//      "Contributor" shall mean Licensor and any individual or Legal Entity
//      on behalf of whom a Contribution has been received by Licensor and
//      subsequently incorporated within the Work.
//
//   2. Grant of Copyright License. Subject to the terms and conditions of
//      this License, each Contributor hereby grants to You a perpetual,
//      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
//      copyright license to reproduce, prepare Derivative Works of,
//      publicly display, publicly perform, sublicense, and distribute the
//      Work and such Derivative Works in Source or Object form.
//
//   3. Grant of Patent License. Subject to the terms and conditions of
//      this License, each Contributor hereby grants to You a perpetual,
//      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
//      (except as stated in this section) patent license to make, have made,
//      use, offer to sell, sell, import, and otherwise transfer the Work,
//      where such license applies only to those patent claims licensable
//      by such Contributor that are necessarily infringed by their
//      Contribution(s) alone or by combination of their Contribution(s)
//      with the Work to which such Contribution(s) was submitted. If You
//      institute patent litigation against any entity (including a
//      cross-claim or counterclaim in a lawsuit) alleging that the Work
//      or a Contribution incorporated within the Work constitutes direct
//      or contributory patent infringement, then any patent licenses
//      granted to You under this License for that Work shall terminate
//      as of the date such litigation is filed.
//
//   4. Redistribution. You may reproduce and distribute copies of the
//      Work or Derivative Works thereof in any medium, with or without
//      modifications, and in Source or Object form, provided that You
//      meet the following conditions:
//
//      (a) You must give any other recipients of the Work or
//          Derivative Works a copy of this License; and
//
//      (b) You must cause any modified files to carry prominent notices
//          stating that You changed the files; and
//
//      (c) You must retain, in the Source form of any Derivative Works
//          that You distribute, all copyright, patent, trademark, and
//          attribution notices from the Source form of the Work,
//          excluding those notices that do not pertain to any part of
//          the Derivative Works; and
//
//      (d) If the Work includes a "NOTICE" text file as part of its
//          distribution, then any Derivative Works that You distribute must
//          include a readable copy of the attribution notices contained
//          within such NOTICE file, excluding those notices that do not
//          pertain to any part of the Derivative Works, in at least one
//          of the following places: within a NOTICE text file distributed
//          as part of the Derivative Works; within the Source form or
//          documentation, if provided along with the Derivative Works; or,
//          within a display generated by the Derivative Works, if and
//          wherever such third-party notices normally appear. The contents
//          of the NOTICE file are for informational purposes only and
//          do not modify the License. You may add Your own attribution
//          notices within Derivative Works that You distribute, alongside
//          or as an addendum to the NOTICE text from the Work, provided
//          that such additional attribution notices cannot be construed
//          as modifying the License.
//
//      You may add Your own copyright statement to Your modifications and
//      may provide additional or different license terms and conditions
//      for use, reproduction, or distribution of Your modifications, or
//      for any such Derivative Works as a whole, provided Your use,
//      reproduction, and distribution of the Work otherwise complies with
//      the conditions stated in this License.
//
//   5. Submission of Contributions. Unless You explicitly state otherwise,
//      any Contribution intentionally submitted for inclusion in the Work
//      by You to the Licensor shall be under the terms and conditions of
//      this License, without any additional terms or conditions.
//      Notwithstanding the above, nothing herein shall supersede or modify
//      the terms of any separate license agreement you may have executed
//      with Licensor regarding such Contributions.
//
//   6. Trademarks. This License does not grant permission to use the trade
//      names, trademarks, service marks, or product names of the Licensor,
//      except as required for reasonable and customary use in describing the
//      origin of the Work and reproducing the content of the NOTICE file.
//
//   7. Disclaimer of Warranty. Unless required by applicable law or
//      agreed to in writing, Licensor provides the Work (and each
//      Contributor provides its Contributions) on an "AS IS" BASIS,
//      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
//      implied, including, without limitation, any warranties or conditions
//      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
//      PARTICULAR PURPOSE. You are solely responsible for determining the
//      appropriateness of using or redistributing the Work and assume any
//      risks associated with Your exercise of permissions under this License.
//
//   8. Limitation of Liability. In no event and under no legal theory,
//      whether in tort (including negligence), contract, or otherwise,
//      unless required by applicable law (such as deliberate and grossly
//      negligent acts) or agreed to in writing, shall any Contributor be
//      liable to You for damages, including any direct, indirect, special,
//      incidental, or consequential damages of any character arising as a
//      result of this License or out of the use or inability to use the
//      Work (including but not limited to damages for loss of goodwill,
//      work stoppage, computer failure or malfunction, or any and all
//      other commercial damages or losses), even if such Contributor
//      has been advised of the possibility of such damages.
//
//   9. Accepting Warranty or Additional Liability. While redistributing
//      the Work or Derivative Works thereof, You may choose to offer,
//      and charge a fee for, acceptance of support, warranty, indemnity,
//      or other liability obligations and/or rights consistent with this
//      License. However, in accepting such obligations, You may act only
//      on Your own behalf and on Your sole responsibility, not on behalf
//      of any other Contributor, and only if You agree to indemnify,
//      defend, and hold each Contributor harmless for any liability
//      incurred by, or claims asserted against, such Contributor by reason
//      of your accepting any such warranty or additional liability.
//
//   END OF TERMS AND CONDITIONS
//
//   APPENDIX: How to apply the Apache License to your work.
//
//      To apply the Apache License to your work, attach the following
//      boilerplate notice, with the fields enclosed by brackets "{}"
//      replaced with your own identifying information. (Don't include
//      the brackets!)  The text should be enclosed in the appropriate
//      comment syntax for the file format. We also recommend that a
//      file or class name and description of purpose be included on the
//      same "printed page" as the copyright notice for easier
//      identification within third-party archives.
//
//   Copyright (c) 2017 Cisco
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

syntax = "proto3";

//Path: Cisco-IOS-XR-infra-statsd-oper:infra-statistics/interfaces/interface/latest/generic-counters

package cisco_ios_xr_infra_statsd_oper.infra_statistics.interfaces.interface.latest.generic_counters;

// Generic set of interface counters
message ifstatsbag_generic_KEYS {
    string interface_name = 1;
}

message ifstatsbag_generic {
    // Packets received
    uint64 packets_received = 50;
    // Bytes received
    uint64 bytes_received = 51;
    // Packets sent
    uint64 packets_sent = 52;
    // Bytes sent
    uint64 bytes_sent = 53;
    // Multicast packets received
    uint64 multicast_packets_received = 54;
    // Broadcast packets received
    uint64 broadcast_packets_received = 55;
    // Multicast packets sent
    uint64 multicast_packets_sent = 56;
    // Broadcast packets sent
    uint64 broadcast_packets_sent = 57;
    // Total output drops
    uint32 output_drops = 58;
    // Output queue drops
    uint32 output_queue_drops = 59;
    // Total input drops
    uint32 input_drops = 60;
    // Input queue drops
    uint32 input_queue_drops = 61;
    // Received runt packets
    uint32 runt_packets_received = 62;
    // Received giant packets
    uint32 giant_packets_received = 63;
    // Received throttled packets
    uint32 throttled_packets_received = 64;
    // Received parity packets
    uint32 parity_packets_received = 65;
    // Unknown protocol packets received
    uint32 unknown_protocol_packets_received = 66;
    // Total input errors
    uint32 input_errors = 67;
    // Input CRC errors
    uint32 crc_errors = 68;
    // Input overruns
    uint32 input_overruns = 69;
    // Framing-errors received
    uint32 framing_errors_received = 70;
    // Input ignored packets
    uint32 input_ignored_packets = 71;
    // Input aborts
    uint32 input_aborts = 72;
    // Total output errors
    uint32 output_errors = 73;
    // Output underruns
    uint32 output_underruns = 74;
    // Output buffer failures
    uint32 output_buffer_failures = 75;
    // Output buffers swapped out
    uint32 output_buffers_swapped_out = 76;
    // Applique
    uint32 applique = 77;
    // Number of board resets
    uint32 resets = 78;
    // Carrier transitions
    uint32 carrier_transitions = 79;
    // Availability bit mask
    uint32 availability_flag = 80;
    // Time when counters were last written (in seconds)
    uint32 last_data_time = 81;
    // Number of seconds since last clear counters
    uint32 seconds_since_last_clear_counters = 82;
    // SysUpTime when counters were last reset (in seconds)
    uint32 last_discontinuity_time = 83;
    // Seconds since packet received
    uint32 seconds_since_packet_received = 84;
    // Seconds since packet sent
    uint32 seconds_since_packet_sent = 85;
}

GetProto: Received ReqId 7646 
 $ 

```
