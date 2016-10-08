gofc
==================

[![Build Status](https://travis-ci.org/Kmotiko/gofc.svg?branch=master)](https://travis-ci.org/Kmotiko/gofc)


## What is this?

OpenFlow Controller written in golang.
Now, only support OpenFlow 1.3.

## How to use

To run OpenFlow Controller with gofc, you have to follow the steps outlined bellow.

 * Implement Handler
 * Regist and Run Controller


### Implement Handler
First, Implement handler interface.
For example, if you want to handle PacketIn message, you should implement HandlePacketIn like bellow.

```
import (
	"github.com/Kmotiko/gofc"
	"github.com/Kmotiko/gofc/ofprotocol/ofp13"
}

type SampleController struct {
	// add any paramter used in controller.
}

func (c *SampleController) HandlePacketIn(msg *ofp13.OfpPacketIn, dp *gofc.Datapath) {
	// create flow mod
	fm := ofp13.NewOfpFlowMod(
		0,
		0,
		0,
		ofp13.OFPFC_ADD,
		0,
		0,
		0,
		ofp13.OFP_NO_BUFFER,
		ofp13.OFPP_ANY,
		ofp13.OFPG_ANY,
		ofp13.OFPFF_SEND_FLOW_REM)

	// TODO: set match field and instructions.

	// send FlowMod message
	dp.Send(fm)
}
```

Defenition of gofc interfaces is written in ofp13_handler.go.



### Regist and Run Controller

If you complete to implement of handlers, you can regist that to gofc and start it.

```
func main() {
	// regist app
	ofc := NewSampleController()
	gofc.GetAppManager().RegistApplication(ofc)

	// start Controller
	gofc.ServerLoop()
}
```

## OpenFlow Messages Support Status

### Messages

 - [x] Hello
 - [ ] SwitchConfig
 - [x] EchoRequest
 - [x] EchoReply
 - [x] SwitchFeatures
 - [x] PacketIn
 - [ ] TableMod
 - [ ] PortMod
 - [x] FlowMod(Partial Support)
 - [x] GroupMod
 - [x] MeterMod
 - [ ] PacketOut
 - [ ] FlowRemoved
 - [ ] OfpErroMsg
 - [ ] OfpExperimenterMsg
 - [x] OfpMultipartRequest
 - [x] OfpMultipartReply


### Match Field

 - [x] In Port
 - [x] In Phy Port
 - [x] Metadata
 - [x] Eth Dst
 - [x] Eth Src
 - [x] Eth Type
 - [x] Vlan ID
 - [x] Vlan PCP
 - [x] IP Dscp
 - [x] IP ECN
 - [x] IP Proto
 - [x] IPv4 Src
 - [x] IPv4 Dst
 - [x] TCP Src
 - [x] TCP Dst
 - [x] UDP Src
 - [x] UDP Dst
 - [x] SCTP Src
 - [x] SCTP Dst
 - [x] ICMPv4 Type
 - [x] ICMPv4 Code
 - [x] Arp Op
 - [x] Arp SPA
 - [x] Arp TPA
 - [x] Arp SHA
 - [x] Arp THA
 - [x] IPv6 Src
 - [x] IPv6 Dst
 - [x] IPv6 FLabel
 - [x] IPv6 Type
 - [x] IPv6 Code
 - [x] IPv6 ND Target
 - [x] IPv6 ND SLL
 - [x] IPv6 ND TLL
 - [x] MPLS Label
 - [x] MPLS TC
 - [x] MPLS BOS
 - [x] PBB ISID
 - [x] TUNNEL ID
 - [x] IPv6 ExtraHeader

### Instructions

 - [ ] Go to table
 - [ ] Write Metadata
 - [x] Action
 - [ ] Meter

### Actions

 - [x] Output
 - [x] Copy TTL OUT
 - [x] Copy TTL IN
 - [x] Set MPLS TTL
 - [x] Dec  MPLS TTL
 - [x] PUSH VLAN
 - [x] PUSH MPLS
 - [x] PUSH PBB
 - [x] POP VLAN
 - [x] POP MPLS
 - [x] POP PBB
 - [x] Group
 - [x] Set NW TTL
 - [x] Dec NW TTL
 - [x] Set Field
 - [ ] Experimenter Header

### MeterBand

 - [x] MeterBandDrop
 - [ ] MeterBandDscpRemark
 - [ ] MeterBandExperimenter

### Multipart Message

 - [x] DescStats
 - [x] FlowStats
 - [x] AggregateStats
 - [ ] TableStats
 - [ ] PortStats
 - [ ] QueueStats
 - [ ] GroupStats
 - [ ] GroupDescStats
 - [ ] GroupFeaturesStats
 - [ ] MeterStats
 - [ ] MeterConfigStats
 - [ ] MeterFeaturesStats
 - [ ] TableFeaturesStats
 - [ ] PortDescStats
 - [ ] ExperimenterStats

## License

This software is distributed with MIT license.

```
Copyright (c) 2016 Kmotiko
```
