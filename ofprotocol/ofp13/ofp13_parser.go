package ofp13

import (
	"encoding/binary"
	"errors"
	"net"
)

func Parse(packet []byte) (msg OFMessage) {
	switch packet[1] {
	case OFPT_HELLO:
		msg = new(OfpHello)
		msg.Parse(packet)
	case OFPT_FEATURES_REPLY:
		msg = NewOfpFeaturesReply()
		msg.Parse(packet)
	case OFPT_ECHO_REQUEST:
		msg = NewOfpEchoRequest()
		msg.Parse(packet)
	case OFPT_ECHO_REPLY:
		msg = NewOfpEchoReply()
		msg.Parse(packet)
	case OFPT_PACKET_IN:
		msg = NewOfpPacketIn()
		msg.Parse(packet)
	case OFPT_MULTIPART_REPLY:
		msg = NewOfpMultipartReply()
		msg.Parse(packet)
	default:
	}
	return msg
}

var xid uint32 = 0

func nextXid() uint32 {
	tmp := xid
	xid += 1
	return tmp
}

/*****************************************************/
/* OfpHeader                                         */
/*****************************************************/

/// create OfpHeader instance.
func NewOfpHeader(t uint8) OfpHeader {
	// 4 means ofp version 1.3
	h := OfpHeader{4, t, 8, nextXid()}
	return h
}

/// Serialize OfpHeader and return it as slice of byte.
func (h *OfpHeader) Serialize() []byte {
	packet := make([]byte, 8)
	packet[0] = h.Version
	packet[1] = h.Type
	binary.BigEndian.PutUint16(packet[2:], h.Length)
	binary.BigEndian.PutUint32(packet[4:], h.Xid)
	return packet
}

/// Parse packet data and set value to OfpHeader instance.
func (h *OfpHeader) Parse(packet []byte) {
	h.Version = packet[0]
	h.Type = packet[1]
	h.Length = binary.BigEndian.Uint16(packet[2:])
	h.Xid = binary.BigEndian.Uint32(packet[4:])
}

/// Return OfpHeader's size.
func (h *OfpHeader) Size() int {
	return 8
}

/*****************************************************/
/* Echo Message                                      */
/*****************************************************/
func NewOfpEchoRequest() *OfpHeader {
	echo := NewOfpHeader(OFPT_ECHO_REQUEST)
	return &echo
}

func NewOfpEchoReply() *OfpHeader {
	echo := NewOfpHeader(OFPT_ECHO_REPLY)
	return &echo
}

/*****************************************************/
/* OfpHelloElemHeader                                */
/*****************************************************/
func NewOfpHelloElemHeader() *OfpHelloElemHeader {
	e := new(OfpHelloElemHeader)
	e.Length = 8
	return e
}

func (h *OfpHelloElemHeader) Serialize() []byte {
	packet := make([]byte, 8)
	binary.BigEndian.PutUint16(packet[0:], h.Type)
	binary.BigEndian.PutUint16(packet[2:], h.Length)

	return packet
}

func (h *OfpHelloElemHeader) Parse(packet []byte) {
	h.Type = binary.BigEndian.Uint16(packet[0:])
	h.Length = binary.BigEndian.Uint16(packet[2:])
}

func (h *OfpHelloElemHeader) Size() int {
	return 8
}

/*****************************************************/
/* OfpHello                                          */
/*****************************************************/
func NewOfpHello() *OfpHello {
	hello := new(OfpHello)
	hello.Header = NewOfpHeader(OFPT_HELLO)
	hello.Elements = make([]OfpHelloElemHeader, 0)
	return hello
}

///
///
func (m *OfpHello) Serialize() []byte {
	packet := make([]byte, m.Size())
	// header
	h_packet := m.Header.Serialize()
	// append header
	copy(packet[0:], h_packet)

	// serialize hello body
	index := len(h_packet)
	e_packet := make([]byte, 0)
	for _, elem := range m.Elements {
		e_packet = elem.Serialize()
		copy(packet[index:], elem.Serialize())
		index += len(e_packet)
	}

	return packet
}

func (m *OfpHello) Parse(packet []byte) {
	m.Header.Parse(packet[0:])
	index := 8

	for index < len(packet) {
		e := NewOfpHelloElemHeader()
		e.Parse(packet[index:])
		index += e.Size()
		// m.Elements = append(m.Elements, e)
	}
	return
}

func (m *OfpHello) Size() int {
	size := m.Header.Size()
	for _, e := range m.Elements {
		size += e.Size()
	}
	return size
}

/*****************************************************/
/* OfpSwitchConfig                                   */
/*****************************************************/

/*****************************************************/
/* OfpTableMod                                       */
/*****************************************************/

/*****************************************************/
/* OfpPortStatus                                     */
/*****************************************************/

/*****************************************************/
/* OfpPortMod                                        */
/*****************************************************/

/*****************************************************/
/* OfpFeaturesRequest                                */
/*****************************************************/
func NewOfpFeaturesRequest() *OfpHeader {
	m := NewOfpHeader(OFPT_FEATURES_REQUEST)
	return &m
}

/*****************************************************/
/* OfpSwitchFeatures                                 */
/*****************************************************/
func NewOfpFeaturesReply() *OfpSwitchFeatures {
	m := new(OfpSwitchFeatures)
	m.Header = NewOfpHeader(OFPT_FEATURES_REPLY)
	return m
}

func (m *OfpSwitchFeatures) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()
	copy(packet[0:], h_packet)
	index := m.Header.Size()
	binary.BigEndian.PutUint64(packet[index:8], m.DatapathId)
	index += 8
	binary.BigEndian.PutUint32(packet[index:4], m.NBuffers)
	index += 4
	packet[index] = m.NTables
	index += 1
	packet[index] = m.AuxiliaryId
	index += 1
	packet[index] = m.Pad[0]
	index += 1
	packet[index] = m.Pad[1]
	index += 1
	binary.BigEndian.PutUint32(packet[index:4], m.Capabilities)
	index += 4
	binary.BigEndian.PutUint32(packet[index:4], m.Reserved)

	return packet
}

func (m *OfpSwitchFeatures) Parse(packet []byte) {
	m.Header.Parse(packet)
	index := m.Header.Size()
	m.DatapathId = binary.BigEndian.Uint64(packet[index:])
	index += 8
	m.NBuffers = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.NTables = packet[index]
	index += 1
	m.AuxiliaryId = packet[index]
	index += 1
	m.Pad[0] = packet[index]
	index += 1
	m.Pad[1] = packet[index]
	index += 1
	m.Capabilities = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Reserved = binary.BigEndian.Uint32(packet[index:])
	index += 4
}

func (m *OfpSwitchFeatures) Size() int {
	return m.Header.Size() + 24
}

/*****************************************************/
/* OfpFlowMod                                        */
/*****************************************************/
func NewOfpFlowMod(
	cookie uint64,
	cookieMask uint64,
	tableId uint8,
	command uint8,
	idleTimeout uint16,
	hardTimeout uint16,
	priority uint16,
	bufferId uint32,
	outPort uint32,
	outGroup uint32,
	flags uint16,
) *OfpFlowMod {
	m := new(OfpFlowMod)
	m.Header = NewOfpHeader(OFPT_FLOW_MOD)
	m.Cookie = cookie
	m.CookieMask = cookieMask
	m.TableId = tableId
	m.Command = command
	m.IdleTimeout = idleTimeout
	m.HardTimeout = hardTimeout
	m.Priority = priority
	m.BufferId = bufferId
	m.OutPort = outPort
	m.OutGroup = outGroup
	m.Flags = flags
	m.Match = NewOfpMatch()
	return m
}

func (m *OfpFlowMod) Serialize() []byte {
	packet := make([]byte, m.Size())
	m.Header.Length = uint16(m.Size())
	h_packet := m.Header.Serialize()
	copy(packet[0:], h_packet)
	index := m.Header.Size()

	binary.BigEndian.PutUint64(packet[index:], m.Cookie)
	index += 8
	binary.BigEndian.PutUint64(packet[index:], m.CookieMask)
	index += 8
	packet[index] = m.TableId
	index++
	packet[index] = m.Command
	index++
	binary.BigEndian.PutUint16(packet[index:], m.IdleTimeout)
	index += 2
	binary.BigEndian.PutUint16(packet[index:], m.HardTimeout)
	index += 2
	binary.BigEndian.PutUint16(packet[index:], m.Priority)
	index += 2
	binary.BigEndian.PutUint32(packet[index:], m.BufferId)
	index += 4
	binary.BigEndian.PutUint32(packet[index:], m.OutPort)
	index += 4
	binary.BigEndian.PutUint32(packet[index:], m.OutGroup)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], m.Flags)
	index += 2
	packet[index] = 0x00
	index++
	packet[index] = 0x00
	index++

	m_packet := m.Match.Serialize()
	copy(packet[index:], m_packet)
	//index += m.Match.Size()
	index += len(m_packet)

	for _, inst := range m.Instructions {
		copy(packet[index:], inst.Serialize())
		index += inst.Size()
	}

	return packet
}

func (m *OfpFlowMod) Parse(packet []byte) {
	// not implement
}

func (m *OfpFlowMod) Size() int {
	size := m.Header.Size() + 40 + m.Match.Size()
	for _, inst := range m.Instructions {
		size += inst.Size()
	}

	return size
}

func (m *OfpFlowMod) AppendMatchField(mf OxmField) {
	m.Match.Append(mf)
}

func (m *OfpFlowMod) AppendInstruction(i OfpInstruction) {
	m.Instructions = append(m.Instructions, i)
}

/*****************************************************/
/* OfpBucket                                         */
/*****************************************************/
func NewOfpBucket(weight uint16, watchPort uint32, watchGroup uint32) *OfpBucket {
	bucket := new(OfpBucket)
	bucket.Weight = weight
	bucket.WatchPort = watchPort
	bucket.WatchGroup = watchGroup
	bucket.Actions = make([]OfpAction, 0)

	return bucket
}

func (b *OfpBucket) Serialize() []byte {
	packet := make([]byte, b.Size())
	index := 0

	b.Length = (uint16)(b.Size())
	binary.BigEndian.PutUint16(packet[index:], b.Length)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], b.Weight)
	index += 2

	binary.BigEndian.PutUint32(packet[index:], b.WatchPort)
	index += 4

	binary.BigEndian.PutUint32(packet[index:], b.WatchGroup)
	index += 8

	for _, a := range b.Actions {
		a_packet := a.Serialize()
		copy(packet[index:], a_packet)
		index += a.Size()
	}

	return packet
}

func (b *OfpBucket) Parse() {
}

func (b *OfpBucket) Size() int {
	size := 16
	for _, a := range b.Actions {
		size += a.Size()
	}

	return size
}

func (b *OfpBucket) Append(action OfpAction) {
	b.Actions = append(b.Actions, action)
}

/*****************************************************/
/* OfpGroupMod                                       */
/*****************************************************/
func NewOfpGroupMod(command uint16, t uint8, id uint32) *OfpGroupMod {
	header := NewOfpHeader(OFPT_GROUP_MOD)
	m := new(OfpGroupMod)
	m.Header = header
	m.Command = command
	m.Type = t
	m.GroupId = id
	m.Buckets = make([]*OfpBucket, 0)
	return m
}

func (m *OfpGroupMod) Serialize() []byte {
	packet := make([]byte, m.Size())
	m.Header.Length = (uint16)(m.Size())
	h_packet := m.Header.Serialize()

	index := 0
	copy(packet[index:], h_packet)
	index += m.Header.Size()

	binary.BigEndian.PutUint16(packet[index:], m.Command)
	index += 2

	packet[index] = m.Type
	index += 2

	binary.BigEndian.PutUint32(packet[index:], m.GroupId)
	index += 4

	for _, b := range m.Buckets {
		b_bucket := b.Serialize()
		copy(packet[index:], b_bucket)
		index += b.Size()
	}

	return packet
}

func (m *OfpGroupMod) Parse(packet []byte) {
}

func (m *OfpGroupMod) Size() int {
	size := m.Header.Size() + 8
	for _, b := range m.Buckets {
		size += b.Size()
	}
	return size
}

func (m *OfpGroupMod) Append(bucket *OfpBucket) {
	m.Buckets = append(m.Buckets, bucket)
}

/*****************************************************/
/* OfpPacketOut                                      */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpMeterBandHeader                                */
/*****************************************************/
// TODO: implement
func NewOfpMeterBandHeader(t uint16, rate uint32, burstSize uint32) OfpMeterBandHeader {
	m := OfpMeterBandHeader{}
	m.Type = t
	m.Length = 16
	m.Rate = rate
	m.BurstSize = burstSize
	return m
}

func (m *OfpMeterBandHeader) Serialize() []byte {
	packet := make([]byte, m.Size())
	index := 0
	binary.BigEndian.PutUint16(packet[index:], m.Type)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.Length)
	index += 2

	binary.BigEndian.PutUint32(packet[index:], m.Rate)
	index += 4

	binary.BigEndian.PutUint32(packet[index:], m.BurstSize)
	index += 4

	return packet
}

func (m *OfpMeterBandHeader) Parse(packet []byte) {
}

func (m *OfpMeterBandHeader) Size() int {
	return 12
}

/*****************************************************/
/* OfpMeterBandDrop                                  */
/*****************************************************/
// TODO: implement
func NewOfpMeterBandDrop(rate uint32, burstSize uint32) *OfpMeterBandDrop {
	m := new(OfpMeterBandDrop)
	m.Header = NewOfpMeterBandHeader(OFPMBT_DROP, rate, burstSize)
	return m
}

func (m *OfpMeterBandDrop) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()
	index := 0
	copy(packet[index:], h_packet)

	return packet
}

func (m *OfpMeterBandDrop) Parse(packet []byte) {
}

func (m *OfpMeterBandDrop) Size() int {
	return m.Header.Size() + 4
}

func (m *OfpMeterBandDrop) MeterBandType() uint16 {
	return m.Header.Type
}

/*****************************************************/
/* OfpMeterBandDscpRemark                            */
/*****************************************************/
// TODO: implement
func NewOfpMeterBandDscpRemark(rate uint32, burstSize uint32, precLevel uint8) *OfpMeterBandDscpRemark {
	m := new(OfpMeterBandDscpRemark)
	m.Header = NewOfpMeterBandHeader(OFPMBT_DSCP_REMARK, rate, burstSize)
	m.PrecLevel = precLevel
	return m
}

func (m *OfpMeterBandDscpRemark) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()
	index := 0
	copy(packet[index:], h_packet)

	packet[index] = m.PrecLevel

	return packet
}

func (m *OfpMeterBandDscpRemark) Parse(packet []byte) {
}

func (m *OfpMeterBandDscpRemark) Size() int {
	return m.Header.Size() + 4
}

func (m *OfpMeterBandDscpRemark) MeterBandType() uint16 {
	return m.Header.Type
}

/*****************************************************/
/* OfpMeterBandExperimenter                          */
/*****************************************************/
// TODO: implement
func NewOfpMeterBandExperimenter(rate uint32, burstSize uint32, experimenter uint32) *OfpMeterBandExperimenter {
	m := new(OfpMeterBandExperimenter)
	m.Header = NewOfpMeterBandHeader(OFPMBT_EXPERIMENTER, rate, burstSize)
	m.Experimenter = experimenter
	return m
}

func (m *OfpMeterBandExperimenter) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()
	index := 0
	copy(packet[index:], h_packet)
	index += m.Header.Size()

	binary.BigEndian.PutUint32(packet[index:], m.Experimenter)

	return packet
}

func (m *OfpMeterBandExperimenter) Parse(packet []byte) {
}

func (m *OfpMeterBandExperimenter) Size() int {
	return m.Header.Size() + 4
}

func (m *OfpMeterBandExperimenter) MeterBandType() uint16 {
	return m.Header.Type
}

/*****************************************************/
/* OfpMeterMod                                       */
/*****************************************************/
func NewOfpMeterMod(command uint16, flags uint16, id uint32) *OfpMeterMod {
	m := new(OfpMeterMod)
	m.Header = NewOfpHeader(OFPT_METER_MOD)
	m.Header.Length = 16
	m.Command = command
	m.Flags = flags
	m.MeterId = id
	m.Bands = make([]OfpMeterBand, 0)
	return m
}

func (m *OfpMeterMod) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()

	index := 0
	copy(packet[index:], h_packet)
	index += m.Header.Size()

	binary.BigEndian.PutUint16(packet[index:], m.Command)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.Flags)
	index += 2

	binary.BigEndian.PutUint32(packet[index:], m.MeterId)
	index += 4

	for _, b := range m.Bands {
		b_packet := b.Serialize()
		copy(packet[index:], b_packet)
		index += b.Size()
	}

	return packet
}

func (m *OfpMeterMod) Parse(packet []byte) {
}

func (m *OfpMeterMod) Size() int {
	size := m.Header.Size() + 8
	for _, b := range m.Bands {
		size += b.Size()
	}
	return size
}

func (m *OfpMeterMod) AppendMeterBand(mb OfpMeterBand) {
	m.Bands = append(m.Bands, mb)
	m.Header.Length += (uint16)(mb.Size())
}

/*****************************************************/
/* OfpPacketIn                                       */
/*****************************************************/
func NewOfpPacketIn() *OfpPacketIn {
	m := new(OfpPacketIn)
	m.Header = NewOfpHeader(OFPT_PACKET_IN)
	m.Header.Type = OFPT_PACKET_IN
	m.Match = NewOfpMatch()
	return m
}

func (m *OfpPacketIn) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()
	copy(packet[0:], h_packet)
	index := m.Header.Size()

	binary.BigEndian.PutUint32(packet[index:4], m.BufferId)
	index += 4
	binary.BigEndian.PutUint16(packet[index:2], m.TotalLen)
	index += 2
	packet[index] = m.Reason
	index++
	packet[index] = m.TableId
	index++

	m_packet := m.Match.Serialize()
	copy(packet[index:], m_packet)

	return packet
}

func (m *OfpPacketIn) Parse(packet []byte) {
	m.Header.Parse(packet)
	index := m.Header.Size()

	m.BufferId = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.TotalLen = binary.BigEndian.Uint16(packet[index:])
	index += 2
	m.Reason = packet[index]
	index++
	m.TableId = packet[index]
	index++
	m.Cookie = binary.BigEndian.Uint64(packet[index:])
	index += 8

	// parse match field
	m.Match.Parse(packet[index:])
}

func (m *OfpPacketIn) Size() int {
	return m.Header.Size() + 16 + m.Match.Size() + 2 + len(m.Data)
}

/*****************************************************/
/* OfpFlowRemoved                                    */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpMatch                                          */
/*****************************************************/
/*
 in_port		OFPXMT_OFB_IN_PORT
 in_phy_port	OFPXMT_OFB_IN_PHY_PORT
 metadata		OFPXMT_OFB_METADATA
 eth_dst		OFPXMT_OFB_ETH_DST
 eth_src		OFPXMT_OFB_ETH_SRC
 eth_type		OFPXMT_OFB_ETH_TYPE
 vlan_vid		OFPXMT_OFB_VLAN_VID
 vlan_pcp		OFPXMT_OFB_VLAN_PCP
 ip_dscp		OFPXMT_OFB_IP_DSCP
 ip_ecn			OFPXMT_OFB_IP_ECN
 ip_proto		OFPXMT_OFB_IP_PROTO
 ipv4_src		OFPXMT_OFB_IPV4_SRC
 ipv4_dst		OFPXMT_OFB_IPV4_DST
 tcp_src		OFPXMT_OFB_TCP_SRC
 tcp_dst		OFPXMT_OFB_TCP_DST
 udp_src		OFPXMT_OFB_UDP_SRC
 udp_dst		OFPXMT_OFB_UDP_DST
 sctp_src		OFPXMT_OFB_SCTP_SRC
 sctp_dst		OFPXMT_OFB_SCTP_DST
 icmpv4_typ		OFPXMT_OFB_ICMPV4_TYPE
 icmpv4_code	OFPXMT_OFB_ICMPV4_CODE
 arp_op			OFPXMT_OFB_ARP_OP
 arp_spa		OFPXMT_OFB_ARP_SPA
 arp_tpa		OFPXMT_OFB_ARP_TPA
 arp_sha		OFPXMT_OFB_ARP_SHA
 arp_tha		OFPXMT_OFB_ARP_THA
 ipv6_src		OFPXMT_OFB_IPV6_SRC
 ipv6_dst		OFPXMT_OFB_IPV6_DST
 ipv6_flabel	OFPXMT_OFB_IPV6_FLABEL
 icmpv6_type	OFPXMT_OFB_ICMPV6_TYPE
 icmpv6_code	OFPXMT_OFB_ICMPV6_CODE
 ipv6_nd_target	OFPXMT_OFB_IPV6_ND_TARGET
 ipv6_nd_sll	OFPXMT_OFB_IPV6_ND_SLL
 ipv6_nd_tll	OFPXMT_OFB_IPV6_ND_TLL
 mpls_label		OFPXMT_OFB_MPLS_LABEL
 mpls_tc		OFPXMT_OFB_MPLS_TC
 mpls_bos		OFPXMT_OFB_MPLS_BOS
 pbb_isid		OFPXMT_OFB_PBB_ISID
 tunnel_id		OFPXMT_OFB_TUNNEL_ID
 ipv6_exthdr	OFPXMT_OFB_IPV6_EXTHDR
*/
// func NewOfpMatch() *OfpMatch {
// 	m := new(OfpMatch)
// 	m.OxmFields = make([]OxmField, 0)
// 	return m
// }

//func NewOfpMatch(fields []OxmField) *OfpMatch {
func NewOfpMatch() *OfpMatch {
	m := new(OfpMatch)
	m.Type = OFPMT_OXM
	m.OxmFields = make([]OxmField, 0)
	return m
}

func (m *OfpMatch) Serialize() []byte {
	// TODO: set Size
	m.Length = 4
	for _, e := range m.OxmFields {
		m.Length += uint16(e.Size())
	}
	packet := make([]byte, m.Size())
	index := 0
	binary.BigEndian.PutUint16(packet[index:], m.Type)
	index += 2
	binary.BigEndian.PutUint16(packet[index:], m.Length)
	index += 2
	for _, e := range m.OxmFields {
		mf_packet := e.Serialize()
		copy(packet[index:], mf_packet)
		index += e.Size()
	}
	return packet
}

func (m *OfpMatch) Parse(packet []byte) {
	index := 0
	m.Type = binary.BigEndian.Uint16(packet[index:])
	index += 2
	m.Length = binary.BigEndian.Uint16(packet[index:])
	index += 2

	for index < (int(m.Length) - 4) {
		mf := parseOxmField(packet[index:])
		m.OxmFields = append(m.OxmFields, mf)
		index += mf.Size()
	}
}

func (m *OfpMatch) Size() int {
	size := 4
	for _, e := range m.OxmFields {
		size += e.Size()
	}
	size += (8 - (size % 8))
	return size
}

func (m *OfpMatch) Append(f OxmField) {
	m.OxmFields = append(m.OxmFields, f)
}

/*
 * TODO: implements OxmField
 */

func parseOxmField(packet []byte) OxmField {
	header := binary.BigEndian.Uint32(packet[0:])
	switch oxmField(header) {
	case OFPXMT_OFB_IN_PORT:
		mf := NewOxmInPort(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IN_PHY_PORT:
		mf := NewOxmInPhyPort(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_METADATA:
		mf := NewOxmMetadata(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ETH_DST:
		mf, err := NewOxmEthDst("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ETH_SRC:
		mf, err := NewOxmEthSrc("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ETH_TYPE:
		mf := NewOxmEthType(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_VLAN_VID:
		mf := NewOxmVlanVid(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_VLAN_PCP:
		mf := NewOxmVlanPcp(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IP_DSCP:
		mf := NewOxmIpDscp(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IP_ECN:
		mf := NewOxmIpEcn(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IP_PROTO:
		mf := NewOxmIpProto(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV4_SRC:
		mf, err := NewOxmIpv4Src("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV4_DST:
		mf, err := NewOxmIpv4Dst("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_TCP_SRC:
		mf := NewOxmTcpSrc(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_TCP_DST:
		mf := NewOxmTcpDst(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_UDP_SRC:
		mf := NewOxmUdpSrc(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_UDP_DST:
		mf := NewOxmUdpDst(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_SCTP_SRC:
		mf := NewOxmSctpSrc(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_SCTP_DST:
		mf := NewOxmSctpDst(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ICMPV4_TYPE:
		mf := NewOxmIcmpType(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ICMPV4_CODE:
		mf := NewOxmIcmpCode(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ARP_OP:
		mf := NewOxmArpOp(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ARP_SPA:
		mf, err := NewOxmArpSpa("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ARP_TPA:
		mf, err := NewOxmArpTpa("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ARP_SHA:
		mf, err := NewOxmArpSha("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ARP_THA:
		mf, err := NewOxmArpTha("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV6_SRC:
		mf, err := NewOxmIpv6Src("::")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV6_DST:
		mf, err := NewOxmIpv6Dst("::")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV6_FLABEL:
		mf := NewOxmIpv6FLabel(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ICMPV6_TYPE:
		mf := NewOxmIcmpv6Type(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_ICMPV6_CODE:
		mf := NewOxmIcmpv6Code(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV6_ND_TARGET:
		mf, err := NewOxmIpv6NdTarget("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV6_ND_SLL:
		mf, err := NewOxmIpv6NdSll("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV6_ND_TLL:
		mf, err := NewOxmIpv6NdTll("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_MPLS_LABEL:
		mf := NewOxmMplsLabel(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_MPLS_TC:
		mf := NewOxmMplsTc(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_MPLS_BOS:
		mf := NewOxmMplsBos(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_PBB_ISID:
		mf := NewOxmPbbIsid([3]uint8{0, 0, 0})
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_TUNNEL_ID:
		mf := NewOxmTunnelId(0)
		mf.Parse(packet)
		return mf
	case OFPXMT_OFB_IPV6_EXTHDR:
		mf := NewOxmIpv6ExtHeader(0)
		mf.Parse(packet)
		return mf
	default:
		return nil
	}
}

/*
 *
 */
func NewOxmInPort(port uint32) *OxmInPort {
	// create tlv header
	header := OXM_OF_IN_PORT

	// create OxmField
	field := OxmInPort{header, port}

	return &field
}

// Serialize
func (m *OxmInPort) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	// serialize header
	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	// serialize value
	binary.BigEndian.PutUint32(packet[index:], m.Value)

	return packet
}

// Parse
func (m *OxmInPort) Parse(packet []byte) {
	index := 0
	// parse header
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	// parse value
	m.Value = binary.BigEndian.Uint32(packet[index:])
}

// OxmClass
func (m *OxmInPort) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

// OxmField
func (m *OxmInPort) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

// OxmHasMask
func (m *OxmInPort) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

// Length
func (m *OxmInPort) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmInPort) Size() int {
	return int(m.Length() + 4)
}

func NewOxmInPhyPort(port uint32) *OxmInPhyPort {
	// create tlv header
	header := OXM_OF_IN_PHY_PORT

	// create OxmField
	field := OxmInPhyPort{header, port}

	return &field
}

func (m *OxmInPhyPort) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	binary.BigEndian.PutUint32(packet[index:], m.Value)

	return packet
}

func (m *OxmInPhyPort) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint32(packet[index:])
}

func (m *OxmInPhyPort) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmInPhyPort) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmInPhyPort) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmInPhyPort) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmInPhyPort) Size() int {
	return int(m.Length() + 4)
}

func NewOxmMetadata(metadata uint64) *OxmMetadata {
	// create tlv header
	header := OXM_OF_METADATA

	// create OxmField
	field := OxmMetadata{header, metadata, 0}

	return &field
}

func NewOxmMetadataW(metadata uint64, mask uint64) *OxmMetadata {
	// create tlv header
	header := OXM_OF_METADATA_W

	// create field value
	field := OxmMetadata{header, metadata, mask}

	return &field
}

func (m *OxmMetadata) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	binary.BigEndian.PutUint64(packet[index:], m.Value)
	index += 8
	if oxmHasMask(m.TlvHeader) == 1 {
		binary.BigEndian.PutUint64(packet[index:], m.Mask)
	}

	return packet
}

func (m *OxmMetadata) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint64(packet[index:])
	index += 8
	if oxmHasMask(m.TlvHeader) == 1 {
		m.Mask = binary.BigEndian.Uint64(packet[index:])
	}
}

func (m *OxmMetadata) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmMetadata) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmMetadata) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmMetadata) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmMetadata) Size() int {
	return int(m.Length() + 4)
}

func NewOxmEthDst(hwAddr string) (*OxmEth, error) {
	return NewOxmEth(OXM_OF_ETH_DST, hwAddr)
}
func NewOxmEthDstW(hwAddr string, mask string) (*OxmEth, error) {
	return NewOxmEthW(OXM_OF_ETH_DST_W, hwAddr, mask)
}
func NewOxmEthSrc(hwAddr string) (*OxmEth, error) {
	return NewOxmEth(OXM_OF_ETH_SRC, hwAddr)
}
func NewOxmEthSrcW(hwAddr string, mask string) (*OxmEth, error) {
	return NewOxmEthW(OXM_OF_ETH_SRC_W, hwAddr, mask)
}

func NewOxmEth(header uint32, hwAddr string) (*OxmEth, error) {
	// convert str to uint
	value, err := net.ParseMAC(hwAddr)
	if err != nil {
		return nil, err
	}

	// create field value
	field := OxmEth{header, value, nil}

	return &field, nil
}

func NewOxmEthW(header uint32, hwAddr string, mask string) (*OxmEth, error) {
	// convert str to uint
	value, err := net.ParseMAC(hwAddr)
	if err != nil {
		return nil, err
	}
	maskAddr, merr := net.ParseMAC(mask)
	if merr != nil {
		return nil, merr
	}

	// create field value
	field := OxmEth{header, value, maskAddr}

	return &field, nil
}

func (m *OxmEth) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 6; i++ {
		packet[index] = m.Value[i]
		index++
	}

	if oxmHasMask(m.TlvHeader) == 1 {
		for i := 0; i < 6; i++ {
			packet[index] = m.Mask[0]
			index++
		}
	}

	return packet
}

func (m *OxmEth) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := []byte{packet[index], packet[index+1], packet[index+2],
		packet[index+3], packet[index+4], packet[index+5]}
	m.Value = addr
	index += 6
	if oxmHasMask(m.TlvHeader) == 1 {
		mask := []byte{packet[index], packet[index+1], packet[index+2],
			packet[index+3], packet[index+4], packet[index+5]}
		m.Mask = mask
	}
}

func (m *OxmEth) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmEth) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmEth) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmEth) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmEth) Size() int {
	return int(m.Length() + 4)
}

func NewOxmEthType(ethType uint16) *OxmEthType {
	// create tlv header
	header := OXM_OF_ETH_TYPE

	// create field value
	field := OxmEthType{header, ethType}

	return &field
}

func (m *OxmEthType) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], m.Value)
	return packet
}

func (m *OxmEthType) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint16(packet[index:])
}

func (m *OxmEthType) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmEthType) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmEthType) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmEthType) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmEthType) Size() int {
	return int(m.Length() + 4)
}

func NewOxmVlanVid(vid uint16) *OxmVlanVid {
	// create tlv header
	header := OXM_OF_VLAN_VID

	// create field value
	field := OxmVlanVid{header, vid, 0}

	return &field
}

func NewOxmVlanVidW(vid uint16, mask uint16) *OxmVlanVid {
	// create tlv header
	header := OXM_OF_VLAN_VID_W

	// create field value
	field := OxmVlanVid{header, vid, mask}

	return &field
}

func (m *OxmVlanVid) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], m.Value)
	index += 2
	if oxmHasMask(m.TlvHeader) == 1 {
		binary.BigEndian.PutUint16(packet[index:], m.Mask)
	}
	return packet
}

func (m *OxmVlanVid) Parse(packet []byte) {
	index := 0

	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint16(packet[index:])
	index += 2

	if oxmHasMask(m.TlvHeader) == 1 {
		m.Mask = binary.BigEndian.Uint16(packet[index:])
	}
}

func (m *OxmVlanVid) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmVlanVid) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmVlanVid) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmVlanVid) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmVlanVid) Size() int {
	return int(m.Length() + 4)
}

func NewOxmVlanPcp(pcp uint8) *OxmVlanPcp {
	// create tlv header
	header := OXM_OF_VLAN_PCP

	// create field value
	field := OxmVlanPcp{header, pcp}

	return &field
}

func (m *OxmVlanPcp) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	packet[index] = m.Value

	return packet
}

func (m *OxmVlanPcp) Parse(packet []byte) {
	index := 0

	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = packet[index]
}

func (m *OxmVlanPcp) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmVlanPcp) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmVlanPcp) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmVlanPcp) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmVlanPcp) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpDscp(dscp uint8) *OxmIpDscp {
	// create tlv header
	header := OXM_OF_IP_DSCP

	// create field value
	field := OxmIpDscp{header, dscp}

	return &field
}

func (m *OxmIpDscp) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	packet[index] = m.Value

	return packet
}

func (m *OxmIpDscp) Parse(packet []byte) {
	index := 0

	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = packet[index]
}

func (m *OxmIpDscp) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpDscp) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpDscp) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpDscp) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpDscp) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpEcn(ecn uint8) *OxmIpEcn {
	// create tlv header
	header := OXM_OF_IP_ECN

	// create field value
	field := OxmIpEcn{header, ecn}

	return &field
}

func (m *OxmIpEcn) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	packet[index] = m.Value

	return packet
}

func (m *OxmIpEcn) Parse(packet []byte) {
	index := 0

	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = packet[index]
}

func (m *OxmIpEcn) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpEcn) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpEcn) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpEcn) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpEcn) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpProto(proto uint8) *OxmIpProto {
	// create tlv header
	header := OXM_OF_IP_PROTO

	// create field value
	field := OxmIpProto{header, proto}

	return &field
}

func (m *OxmIpProto) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	packet[index] = m.Value

	return packet
}

func (m *OxmIpProto) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = packet[index]
}

func (m *OxmIpProto) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpProto) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpProto) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpProto) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpProto) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpv4Src(addr string) (*OxmIpv4, error) {
	return NewOxmIpv4(OXM_OF_IPV4_SRC, addr)
}
func NewOxmIpv4SrcW(addr string, mask int) (*OxmIpv4, error) {
	return NewOxmIpv4W(OXM_OF_IPV4_SRC_W, addr, mask)
}
func NewOxmIpv4Dst(addr string) (*OxmIpv4, error) {
	return NewOxmIpv4(OXM_OF_IPV4_DST, addr)
}
func NewOxmIpv4DstW(addr string, mask int) (*OxmIpv4, error) {
	return NewOxmIpv4W(OXM_OF_IPV4_DST_W, addr, mask)
}

func NewOxmIpv4(header uint32, addr string) (*OxmIpv4, error) {
	// parse string as IPAddr
	v4addr := net.ParseIP(addr)
	if v4addr == nil {
		return nil, errors.New("failed to parse IPv4 address.")
	}

	// create field value
	field := OxmIpv4{header, v4addr, nil}

	return &field, nil
}

func NewOxmIpv4W(header uint32, addr string, mask int) (*OxmIpv4, error) {
	// parse string as IPAddr
	v4addr := net.ParseIP(addr)
	if v4addr == nil {
		return nil, errors.New("failed to parse IPv4 address.")
	}
	ipmask := net.CIDRMask(mask, 32)

	// create field value
	field := OxmIpv4{header, v4addr, ipmask}

	return &field, nil
}

func (m *OxmIpv4) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 4; i++ {
		packet[index] = m.Value[12+i]
		index++
	}

	if oxmHasMask(m.TlvHeader) == 1 {
		for i := 0; i < 4; i++ {
			packet[index] = m.Mask[i]
			index++
		}
	}

	return packet
}

func (m *OxmIpv4) Parse(packet []byte) {
	index := 0

	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := make([]byte, 4)
	for i := 0; i < 4; i++ {
		addr[i] = packet[index]
		index++
	}
	m.Value = addr

	if oxmHasMask(m.TlvHeader) == 1 {
		mask := make([]byte, 4)
		for i := 0; i < 4; i++ {
			mask[i] = packet[index]
			index++
		}
		m.Mask = mask
	}
}

func (m *OxmIpv4) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpv4) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpv4) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpv4) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpv4) Size() int {
	return int(m.Length() + 4)
}

func NewOxmTcpSrc(port uint16) *OxmTcp {
	return NewOxmTcp(OXM_OF_TCP_SRC, port)
}
func NewOxmTcpDst(port uint16) *OxmTcp {
	return NewOxmTcp(OXM_OF_TCP_DST, port)
}

func NewOxmTcp(header uint32, port uint16) *OxmTcp {
	// create field value
	field := OxmTcp{header, port}
	return &field
}

func (m *OxmTcp) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], m.Value)

	return packet
}

func (m *OxmTcp) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = binary.BigEndian.Uint16(packet[index:])
}

func (m *OxmTcp) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmTcp) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmTcp) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmTcp) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmTcp) Size() int {
	return int(m.Length() + 4)
}

func NewOxmUdpSrc(port uint16) *OxmUdp {
	return NewOxmUdp(OXM_OF_UDP_SRC, port)
}
func NewOxmUdpDst(port uint16) *OxmUdp {
	return NewOxmUdp(OXM_OF_UDP_DST, port)
}

func NewOxmUdp(header uint32, port uint16) *OxmUdp {
	// create field value
	field := OxmUdp{header, port}
	return &field
}

func (m *OxmUdp) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], m.Value)

	return packet
}

func (m *OxmUdp) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = binary.BigEndian.Uint16(packet[index:])
}

func (m *OxmUdp) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmUdp) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmUdp) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmUdp) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmUdp) Size() int {
	return int(m.Length() + 4)
}

func NewOxmSctpSrc(port uint16) *OxmSctp {
	return NewOxmSctp(OXM_OF_SCTP_SRC, port)
}
func NewOxmSctpDst(port uint16) *OxmSctp {
	return NewOxmSctp(OXM_OF_SCTP_DST, port)
}

func NewOxmSctp(header uint32, port uint16) *OxmSctp {
	// create field value
	field := OxmSctp{header, port}
	return &field
}

func (m *OxmSctp) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], m.Value)

	return packet
}

func (m *OxmSctp) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = binary.BigEndian.Uint16(packet[index:])
}

func (m *OxmSctp) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmSctp) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmSctp) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmSctp) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmSctp) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIcmpType(value uint8) *OxmIcmpType {
	// create tlv header
	header := OXM_OF_ICMPV4_TYPE

	// create field value
	field := OxmIcmpType{header, value}

	return &field
}

func (m *OxmIcmpType) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	packet[index] = m.Value

	return packet
}

func (m *OxmIcmpType) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = packet[index]
}

func (m *OxmIcmpType) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIcmpType) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIcmpType) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIcmpType) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIcmpType) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIcmpCode(value uint8) *OxmIcmpCode {
	// create tlv header
	header := OXM_OF_ICMPV4_CODE

	// create field value
	field := OxmIcmpCode{header, value}

	return &field
}

func (m *OxmIcmpCode) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	packet[index] = m.Value

	return packet
}

func (m *OxmIcmpCode) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = packet[index]
}

func (m *OxmIcmpCode) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIcmpCode) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIcmpCode) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIcmpCode) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIcmpCode) Size() int {
	return int(m.Length() + 4)
}

func NewOxmArpOp(op uint16) *OxmArpOp {
	// create tlv header
	header := OXM_OF_ARP_OP

	// create field value
	field := OxmArpOp{header, op}

	return &field
}

func (m *OxmArpOp) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], m.Value)

	return packet
}

func (m *OxmArpOp) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4
	m.Value = binary.BigEndian.Uint16(packet[index:])
}

func (m *OxmArpOp) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmArpOp) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmArpOp) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmArpOp) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmArpOp) Size() int {
	return int(m.Length() + 4)
}

func NewOxmArpSpa(addr string) (*OxmArpPa, error) {
	return NewOxmArpPa(OXM_OF_ARP_SPA, addr)
}
func NewOxmArpSpaW(addr string, mask int) (*OxmArpPa, error) {
	return NewOxmArpPaW(OXM_OF_ARP_SPA_W, addr, mask)
}
func NewOxmArpTpa(addr string) (*OxmArpPa, error) {
	return NewOxmArpPa(OXM_OF_ARP_TPA, addr)
}
func NewOxmArpTpaW(addr string, mask int) (*OxmArpPa, error) {
	return NewOxmArpPaW(OXM_OF_ARP_TPA_W, addr, mask)
}

func NewOxmArpPa(header uint32, addr string) (*OxmArpPa, error) {
	// parse addr
	v4addr := net.ParseIP(addr)
	if v4addr == nil {
		return nil, errors.New("failed to parse IPv4 address.")
	}

	// create field value
	field := OxmArpPa{header, v4addr, nil}
	return &field, nil
}

func NewOxmArpPaW(header uint32, addr string, mask int) (*OxmArpPa, error) {
	// parse addr
	v4addr := net.ParseIP(addr)
	if v4addr == nil {
		return nil, errors.New("failed to parse IPv4 address.")
	}
	ipmask := net.CIDRMask(mask, 32)

	// create field value
	field := OxmArpPa{header, v4addr, ipmask}
	return &field, nil
}

func (m *OxmArpPa) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 4; i++ {
		packet[index] = m.Value[12+i]
		index++
	}

	if oxmHasMask(m.TlvHeader) == 1 {
		for i := 0; i < 4; i++ {
			packet[index] = m.Mask[i]
			index++
		}
	}

	return packet
}

func (m *OxmArpPa) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := make([]byte, 4)
	for i := 0; i < 4; i++ {
		addr[i] = packet[index]
		index++
	}
	m.Value = addr

	if oxmHasMask(m.TlvHeader) == 1 {
		mask := make([]byte, 4)
		for i := 0; i < 4; i++ {
			mask[i] = packet[index]
			index++
		}
		m.Mask = mask
	}
}

func (m *OxmArpPa) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmArpPa) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmArpPa) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmArpPa) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmArpPa) Size() int {
	return int(m.Length() + 4)
}

func NewOxmArpSha(hwAddr string) (*OxmArpHa, error) {
	header := OXM_OF_ARP_SHA
	return NewOxmArpHa(header, hwAddr)
}
func NewOxmArpTha(hwAddr string) (*OxmArpHa, error) {
	header := OXM_OF_ARP_THA
	return NewOxmArpHa(header, hwAddr)
}

func NewOxmArpHa(header uint32, hwAddr string) (*OxmArpHa, error) {
	// create field value
	value, err := net.ParseMAC(hwAddr)
	if err != nil {
		return nil, err
	}

	field := OxmArpHa{header, value}
	return &field, nil
}

func (m *OxmArpHa) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 6; i++ {
		packet[index] = m.Value[i]
		index++
	}

	return packet
}

func (m *OxmArpHa) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := make([]byte, 6)
	for i := 0; i < 6; i++ {
		addr[i] = packet[index]
		index++
	}
	m.Value = addr
}

func (m *OxmArpHa) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmArpHa) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmArpHa) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmArpHa) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmArpHa) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpv6Src(addr string) (*OxmIpv6, error) {
	// create tlv header
	header := OXM_OF_IPV6_SRC

	return NewOxmIpv6(header, addr)
}
func NewOxmIpv6SrcW(addr string, mask int) (*OxmIpv6, error) {
	// create tlv header
	header := OXM_OF_IPV6_SRC_W

	return NewOxmIpv6W(header, addr, mask)
}
func NewOxmIpv6Dst(addr string) (*OxmIpv6, error) {
	// create tlv header
	header := OXM_OF_IPV6_DST

	return NewOxmIpv6(header, addr)
}
func NewOxmIpv6DstW(addr string, mask int) (*OxmIpv6, error) {
	// create tlv header
	header := OXM_OF_IPV6_DST_W

	return NewOxmIpv6W(header, addr, mask)
}

func NewOxmIpv6(header uint32, addr string) (*OxmIpv6, error) {
	// create field value
	v6addr := net.ParseIP(addr)
	if v6addr == nil {
		return nil, errors.New("failed to parse IPv6 address.")
	}

	field := OxmIpv6{header, v6addr, nil}
	return &field, nil
}
func NewOxmIpv6W(header uint32, addr string, mask int) (*OxmIpv6, error) {
	// create field value
	v6addr := net.ParseIP(addr)
	if v6addr == nil {
		return nil, errors.New("failed to parse IPv6 address.")
	}
	ipmask := net.CIDRMask(mask, 128)

	field := OxmIpv6{header, v6addr, ipmask}
	return &field, nil
}

func (m *OxmIpv6) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 16; i++ {
		packet[index] = m.Value[i]
		index++
	}

	if oxmHasMask(m.TlvHeader) == 1 {
		for i := 0; i < 16; i++ {
			packet[index] = m.Mask[i]
			index++
		}
	}

	return packet
}

func (m *OxmIpv6) Parse(packet []byte) {
	index := 0

	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := make([]byte, 16)
	for i := 0; i < 16; i++ {
		addr[i] = packet[index]
		index++
	}
	m.Value = addr

	if oxmHasMask(m.TlvHeader) == 1 {
		mask := make([]byte, 16)
		for i := 0; i < 16; i++ {
			mask[i] = packet[index]
			index++
		}
		m.Mask = mask
	}
}

func (m *OxmIpv6) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpv6) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpv6) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpv6) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpv6) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpv6FLabel(label uint32) *OxmIpv6FLabel {
	// create tlv header
	header := OXM_OF_IPV6_FLABEL

	// create field value
	field := OxmIpv6FLabel{header, label, 0}

	return &field
}

func NewOxmIpv6FLabelW(label uint32, mask uint32) *OxmIpv6FLabel {
	// create tlv header
	header := OXM_OF_IPV6_FLABEL_W

	// create field value
	field := OxmIpv6FLabel{header, label, mask}

	return &field
}

func (m *OxmIpv6FLabel) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint32(packet[index:], m.Value)
	index += 4

	if oxmHasMask(m.TlvHeader) == 1 {
		binary.BigEndian.PutUint32(packet[index:], m.Mask)
	}

	return packet
}

func (m *OxmIpv6FLabel) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint32(packet[index:])
	index += 4
	if oxmHasMask(m.TlvHeader) == 1 {
		m.Mask = binary.BigEndian.Uint32(packet[index:])
	}
}

func (m *OxmIpv6FLabel) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpv6FLabel) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpv6FLabel) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpv6FLabel) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpv6FLabel) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIcmpv6Type(value uint8) *OxmIcmpv6Type {
	// create tlv header
	header := OXM_OF_ICMPV6_TYPE

	// create field value
	field := OxmIcmpv6Type{header, value}

	return &field
}

func (m *OxmIcmpv6Type) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	packet[index] = m.Value

	return packet
}

func (m *OxmIcmpv6Type) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = packet[index]
}

func (m *OxmIcmpv6Type) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIcmpv6Type) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIcmpv6Type) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIcmpv6Type) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIcmpv6Type) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIcmpv6Code(value uint8) *OxmIcmpv6Code {
	// create tlv header
	header := OXM_OF_ICMPV6_CODE

	// create field value
	field := OxmIcmpv6Code{header, value}

	return &field
}

func (m *OxmIcmpv6Code) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	packet[index] = m.Value

	return packet
}

func (m *OxmIcmpv6Code) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = packet[index]
}

func (m *OxmIcmpv6Code) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIcmpv6Code) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIcmpv6Code) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIcmpv6Code) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIcmpv6Code) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpv6NdTarget(addr string) (*OxmIpv6NdTarget, error) {
	// create tlv header
	header := OXM_OF_IPV6_ND_TARGET

	v6addr := net.ParseIP(addr)
	if v6addr == nil {
		return nil, errors.New("failed to parse IPv6 address.")
	}

	// create field value
	field := OxmIpv6NdTarget{header, v6addr}

	return &field, nil
}

func (m *OxmIpv6NdTarget) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 16; i++ {
		packet[index] = m.Value[i]
		index++
	}

	return packet
}

func (m *OxmIpv6NdTarget) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := make([]byte, 16)
	for i := 0; i < 16; i++ {
		addr[i] = packet[index]
		index++
	}
	m.Value = addr
}

func (m *OxmIpv6NdTarget) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpv6NdTarget) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpv6NdTarget) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpv6NdTarget) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpv6NdTarget) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpv6NdSll(hwAddr string) (*OxmIpv6NdSll, error) {
	// create tlv header
	header := OXM_OF_IPV6_ND_SLL

	// create field value
	value, err := net.ParseMAC(hwAddr)
	if err != nil {
		return nil, err
	}

	field := OxmIpv6NdSll{header, value}

	return &field, nil
}

func (m *OxmIpv6NdSll) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 6; i++ {
		packet[index] = m.Value[i]
		index++
	}

	return packet
}

func (m *OxmIpv6NdSll) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := make([]byte, 6)
	for i := 0; i < 6; i++ {
		addr[i] = packet[index]
		index++
	}
	m.Value = addr
}

func (m *OxmIpv6NdSll) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpv6NdSll) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpv6NdSll) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpv6NdSll) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpv6NdSll) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpv6NdTll(hwAddr string) (*OxmIpv6NdTll, error) {
	// create tlv header
	header := OXM_OF_IPV6_ND_TLL

	// create field value
	value, err := net.ParseMAC(hwAddr)
	if err != nil {
		return nil, err
	}

	field := OxmIpv6NdTll{header, value}

	return &field, nil
}

func (m *OxmIpv6NdTll) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 6; i++ {
		packet[index] = m.Value[i]
		index++
	}

	return packet
}

func (m *OxmIpv6NdTll) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	addr := make([]byte, 6)
	for i := 0; i < 6; i++ {
		addr[i] = packet[index]
		index++
	}
	m.Value = addr
}

func (m *OxmIpv6NdTll) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpv6NdTll) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpv6NdTll) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpv6NdTll) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpv6NdTll) Size() int {
	return int(m.Length() + 4)
}

func NewOxmMplsLabel(label uint32) *OxmMplsLabel {
	// create tlv header
	header := OXM_OF_MPLS_LABEL

	// create field value
	field := OxmMplsLabel{header, label}

	return &field
}

func (m *OxmMplsLabel) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4
	binary.BigEndian.PutUint32(packet[index:], m.Value)

	return packet
}

func (m *OxmMplsLabel) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint32(packet[index:])
}

func (m *OxmMplsLabel) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmMplsLabel) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmMplsLabel) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmMplsLabel) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmMplsLabel) Size() int {
	return int(m.Length() + 4)
}

func NewOxmMplsTc(tc uint8) *OxmMplsTc {
	// create tlv header
	header := OXM_OF_MPLS_TC

	// create field value
	field := OxmMplsTc{header, tc}

	return &field
}

func (m *OxmMplsTc) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	packet[index] = m.Value

	return packet
}

func (m *OxmMplsTc) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = packet[index]
}

func (m *OxmMplsTc) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmMplsTc) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmMplsTc) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmMplsTc) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmMplsTc) Size() int {
	return int(m.Length() + 4)
}

func NewOxmMplsBos(bos uint8) *OxmMplsBos {
	// create tlv header
	header := OXM_OF_MPLS_BOS

	// create field value
	field := OxmMplsBos{header, bos}

	return &field
}

func (m *OxmMplsBos) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	packet[index] = m.Value

	return packet
}

func (m *OxmMplsBos) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = packet[index]
}

func (m *OxmMplsBos) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmMplsBos) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmMplsBos) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmMplsBos) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmMplsBos) Size() int {
	return int(m.Length() + 4)
}

func NewOxmPbbIsid(isid [3]uint8) *OxmPbbIsid {
	// create tlv header
	header := OXM_OF_PBB_ISID

	// create field value
	field := OxmPbbIsid{header, isid, [3]uint8{0, 0, 0}}

	return &field
}

func NewOxmPbbIsidW(isid [3]uint8, mask [3]uint8) *OxmPbbIsid {
	// create tlv header
	header := OXM_OF_PBB_ISID_W

	// create field value
	field := OxmPbbIsid{header, isid, mask}

	return &field
}

func (m *OxmPbbIsid) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	for i := 0; i < 3; i++ {
		packet[index] = m.Value[i]
		index++
	}

	if oxmHasMask(m.TlvHeader) == 1 {
		for i := 0; i < 3; i++ {
			packet[index] = m.Mask[i]
			index++
		}
	}

	return packet
}

func (m *OxmPbbIsid) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	for i := 0; i < 3; i++ {
		m.Value[i] = packet[index]
		index++
	}

	if oxmHasMask(m.TlvHeader) == 1 {
		for i := 0; i < 3; i++ {
			m.Mask[i] = packet[index]
			index++
		}
	}
}

func (m *OxmPbbIsid) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmPbbIsid) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmPbbIsid) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmPbbIsid) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmPbbIsid) Size() int {
	return int(m.Length() + 4)
}

func NewOxmTunnelId(id uint64) *OxmTunnelId {
	// create tlv header
	header := OXM_OF_TUNNEL_ID

	// create field value
	field := OxmTunnelId{header, id, 0}

	return &field
}

func NewOxmTunnelIdW(id uint64, mask uint64) *OxmTunnelId {
	// create tlv header
	header := OXM_OF_TUNNEL_ID_W

	// create field value
	field := OxmTunnelId{header, id, mask}

	return &field
}

func (m *OxmTunnelId) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	binary.BigEndian.PutUint64(packet[index:], m.Value)
	index += 8

	if oxmHasMask(m.TlvHeader) == 1 {
		binary.BigEndian.PutUint64(packet[index:], m.Mask)
	}

	return packet
}

func (m *OxmTunnelId) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint64(packet[index:])
	index += 8

	if oxmHasMask(m.TlvHeader) == 1 {
		m.Mask = binary.BigEndian.Uint64(packet[index:])
	}
}

func (m *OxmTunnelId) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmTunnelId) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmTunnelId) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmTunnelId) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmTunnelId) Size() int {
	return int(m.Length() + 4)
}

func NewOxmIpv6ExtHeader(value uint16) *OxmIpv6ExtHeader {
	// create tlv header
	header := OXM_OF_IPV6_EXTHDR

	// create field value
	field := OxmIpv6ExtHeader{header, value, 0}

	return &field
}

func NewOxmIpv6ExtHeaderW(value uint16, mask uint16) *OxmIpv6ExtHeader {
	// create tlv header
	header := OXM_OF_IPV6_EXTHDR_W

	// create field value
	field := OxmIpv6ExtHeader{header, value, mask}

	return &field
}

func (m *OxmIpv6ExtHeader) Serialize() []byte {
	index := 0
	packet := make([]byte, m.Size())

	binary.BigEndian.PutUint32(packet[index:], m.TlvHeader)
	index += 4

	binary.BigEndian.PutUint16(packet[index:], m.Value)
	index += 2

	if oxmHasMask(m.TlvHeader) == 1 {
		binary.BigEndian.PutUint16(packet[index:], m.Value)
	}

	return packet
}

func (m *OxmIpv6ExtHeader) Parse(packet []byte) {
	index := 0
	m.TlvHeader = binary.BigEndian.Uint32(packet[index:])
	index += 4

	m.Value = binary.BigEndian.Uint16(packet[index:])

	if oxmHasMask(m.TlvHeader) == 1 {
		m.Mask = binary.BigEndian.Uint16(packet[index:])
	}
}

func (m *OxmIpv6ExtHeader) OxmClass() uint32 {
	return oxmClass(m.TlvHeader)
}

func (m *OxmIpv6ExtHeader) OxmField() uint32 {
	return oxmField(m.TlvHeader)
}

func (m *OxmIpv6ExtHeader) OxmHasMask() uint32 {
	return oxmHasMask(m.TlvHeader)
}

func (m *OxmIpv6ExtHeader) Length() uint32 {
	return oxmLength(m.TlvHeader)
}

func (m *OxmIpv6ExtHeader) Size() int {
	return int(m.Length() + 4)
}

/*****************************************************/
/* OfpInstruction                                    */
/*****************************************************/
func NewOfpInstructionHeader(t uint16) OfpInstructionHeader {
	header := OfpInstructionHeader{t, 4}
	return header
}

func (h *OfpInstructionHeader) Serialize() []byte {
	packet := make([]byte, h.Size())
	index := 0
	binary.BigEndian.PutUint16(packet[index:], h.Type)
	index += 2
	binary.BigEndian.PutUint16(packet[index:], h.Length)
	return packet
}

func (h *OfpInstructionHeader) Parse(packet []byte) {
	index := 0
	h.Type = binary.BigEndian.Uint16(packet[index:])
	index += 2
	h.Length = binary.BigEndian.Uint16(packet[index:])
}

func (i OfpInstructionHeader) Size() int {
	return 4
}

/*
 * OfpInstructionGotoTable
 */
func NewOfpInstructionGotoTable(id uint8) *OfpInstructionGotoTable {
	header := NewOfpInstructionHeader(OFPIT_GOTO_TABLE)
	i := new(OfpInstructionGotoTable)
	i.Header = header
	i.TableId = id
	header.Length = uint16(i.Size())
	return i
}

func (i *OfpInstructionGotoTable) Serialize() []byte {
	packet := make([]byte, i.Size())
	index := 0
	h_packet := i.Header.Serialize()
	copy(packet[0:], h_packet)
	index += i.Header.Size()
	packet[index] = i.TableId
	index += 1
	return packet
}

func (i *OfpInstructionGotoTable) Parse(packet []byte) {
	header := NewOfpInstructionHeader(OFPIT_GOTO_TABLE)
	header.Parse(packet)
	index := header.Size()
	i.Header = header
	i.TableId = packet[index]
}

func (i *OfpInstructionGotoTable) Size() int {
	return i.Header.Size() + 4
}

func (i *OfpInstructionGotoTable) InstructionType() uint16 {
	return OFPIT_GOTO_TABLE
}

/*
 * OfpInstructionWriteMetadata
 */
func NewOfpInstructionWriteMetadata(metadata uint64, mask uint64) *OfpInstructionWriteMetadata {
	i := new(OfpInstructionWriteMetadata)
	header := NewOfpInstructionHeader(OFPIT_WRITE_METADATA)
	i.Header = header
	i.Metadata = metadata
	i.MetadataMask = mask

	return i
}

func (i *OfpInstructionWriteMetadata) Serialize() []byte {
	packet := make([]byte, i.Header.Size())
	index := 0
	h_packet := i.Header.Serialize()
	copy(packet[index:], h_packet)
	index += i.Header.Size()
	index += 4
	binary.BigEndian.PutUint64(packet[index:], i.Metadata)
	index += 8
	binary.BigEndian.PutUint64(packet[index:], i.MetadataMask)
	return packet
}

func (i *OfpInstructionWriteMetadata) Parse(packet []byte) {
	header := NewOfpInstructionHeader(OFPIT_WRITE_METADATA)
	index := 0
	header.Parse(packet)
	index += i.Header.Size()
	index += 4
	i.Metadata = binary.BigEndian.Uint64(packet[index:])
	index += 4
	i.MetadataMask = binary.BigEndian.Uint64(packet[index:])
}

func (i *OfpInstructionWriteMetadata) Size() int {
	return 24
}

func (i *OfpInstructionWriteMetadata) InstructionType() uint16 {
	return OFPIT_WRITE_METADATA
}

/*
 * OfpInstructionActions
 */
func NewOfpInstructionActions(
	t uint16) *OfpInstructionActions {
	// TODO:check t is one of following actions.
	// WRITE_ACTION
	// APPLY_ACTION
	// CLEAR_ACTION
	i := new(OfpInstructionActions)
	header := NewOfpInstructionHeader(t)
	i.Header = header
	i.Actions = make([]OfpAction, 0)
	return i
}

func (i *OfpInstructionActions) Serialize() []byte {
	packet := make([]byte, i.Size())
	index := 0
	// set actual length
	i.Header.Length = uint16(i.Size())
	h_packet := i.Header.Serialize()
	copy(packet[index:], h_packet)
	index += i.Header.Size()

	// Padding
	index += 4

	// Actions
	for _, a := range i.Actions {
		a_packet := a.Serialize()
		copy(packet[index:], a_packet)
		index += a.Size()
	}
	return packet
}

func (i *OfpInstructionActions) Parse(packet []byte) {
	index := 0
	i.Header.Parse(packet[index:])
	index += i.Header.Size()

	// Pad
	index += 4

	// for index < len(packet) {
	for index < (int)(i.Header.Length) {
		a_type := binary.BigEndian.Uint16(packet[index:])
		switch a_type {
		//TODO:implement
		case OFPAT_OUTPUT:
			action := NewOfpActionOutput(0, 0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_COPY_TTL_OUT:
			action := NewOfpActionCopyTtlOut()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_COPY_TTL_IN:
			action := NewOfpActionCopyTtlIn()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_SET_MPLS_TTL:
			action := NewOfpActionSetMplsTtl(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_DEC_MPLS_TTL:
			action := NewOfpActionDecMplsTtl()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_PUSH_VLAN:
			action := NewOfpActionPushVlan()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_POP_VLAN:
			action := NewOfpActionPopVlan(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_PUSH_MPLS:
			action := NewOfpActionPushMpls()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_POP_MPLS:
			action := NewOfpActionPopMpls(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_SET_QUEUE:
			action := NewOfpActionSetQueue(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_GROUP:
			action := NewOfpActionGroup(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_SET_NW_TTL:
			action := NewOfpActionSetNwTtl(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_DEC_NW_TTL:
			action := NewOfpActionDecNwTtl()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_SET_FIELD:
			action := newEmptyOfpActionSetField()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_PUSH_PBB:
			action := NewOfpActionPushPbb()
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_POP_PBB:
			action := NewOfpActionPopPbb(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		case OFPAT_EXPERIMENTER:
			action := NewOfpActionExperimenter(0)
			action.Parse(packet[index:])
			i.Append(action)
			index += action.Size()
		default:
		}
	}
}

func (i *OfpInstructionActions) Size() int {
	size := i.Header.Size() + 4
	for _, a := range i.Actions {
		size += a.Size()
	}
	return size
}

func (i *OfpInstructionActions) InstructionType() uint16 {
	return i.Header.Type
}

func (i *OfpInstructionActions) Append(a OfpAction) {
	i.Actions = append(i.Actions, a)
}

/*****************************************************/
/* OfpAction                                         */
/*****************************************************/
func NewOfpActionHeader(t uint16, length uint16) OfpActionHeader {
	header := OfpActionHeader{t, length}
	return header
}

func (h *OfpActionHeader) Serialize() []byte {
	packet := make([]byte, h.Size())
	binary.BigEndian.PutUint16(packet[0:], h.Type)
	binary.BigEndian.PutUint16(packet[2:], h.Length)

	return packet
}

func (h *OfpActionHeader) Parse(packet []byte) {
	h.Type = binary.BigEndian.Uint16(packet[0:])
	h.Length = binary.BigEndian.Uint16(packet[2:])
}

func (h OfpActionHeader) Size() int {
	return 4
}

/*
 * OfpActionOutput
 */
func NewOfpActionOutput(port uint32, max_len uint16) *OfpActionOutput {
	header := NewOfpActionHeader(OFPAT_OUTPUT, 16)
	action := new(OfpActionOutput)
	action.ActionHeader = header
	action.Port = port
	action.MaxLen = max_len
	return action
}

// func NewOfpActionOutput(port uint32) {
// 	h := NewOfpActionOutput(OFPAT_OUTPUT, 16)
// 	h.Port = port
// 	h.MaxLen = OFPCML_MAX
// }

func (a *OfpActionOutput) Serialize() []byte {
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	index := 0
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	binary.BigEndian.PutUint32(packet[index:], a.Port)
	index += 4
	binary.BigEndian.PutUint16(packet[index:], a.MaxLen)

	return packet
}

func (a *OfpActionOutput) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet[index:])
	index += a.ActionHeader.Size()

	a.Port = binary.BigEndian.Uint32(packet[index:])
	index += 4
	a.MaxLen = binary.BigEndian.Uint16(packet[index:])
}

func (a *OfpActionOutput) Size() int {
	return 16
}

func (a *OfpActionOutput) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionCopyTtlOut
 */
func NewOfpActionCopyTtlOut() *OfpActionCopyTtlOut {
	action := new(OfpActionCopyTtlOut)
	header := NewOfpActionHeader(OFPAT_COPY_TTL_OUT, 8)
	action.ActionHeader = header
	return action
}

func (a *OfpActionCopyTtlOut) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()

	return packet
}

func (a *OfpActionCopyTtlOut) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
}

func (a *OfpActionCopyTtlOut) Size() int {
	return 8
}

func (a *OfpActionCopyTtlOut) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionCopyTtlIn
 */
// TODO: implement
func NewOfpActionCopyTtlIn() *OfpActionCopyTtlIn {
	action := new(OfpActionCopyTtlIn)
	header := NewOfpActionHeader(OFPAT_COPY_TTL_IN, 8)
	action.ActionHeader = header
	return action
}

func (a *OfpActionCopyTtlIn) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()

	return packet
}

func (a *OfpActionCopyTtlIn) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
}

func (a *OfpActionCopyTtlIn) Size() int {
	return 8
}

func (a *OfpActionCopyTtlIn) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionSetMplsTtl
 */
func NewOfpActionSetMplsTtl(ttl uint8) *OfpActionSetMplsTtl {
	action := new(OfpActionSetMplsTtl)
	header := NewOfpActionHeader(OFPAT_SET_MPLS_TTL, 8)
	action.ActionHeader = header
	action.MplsTtl = ttl
	return action
}

func (a *OfpActionSetMplsTtl) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	packet[index] = a.MplsTtl

	return packet
}

func (a *OfpActionSetMplsTtl) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
	a.MplsTtl = packet[index]
}

func (a *OfpActionSetMplsTtl) Size() int {
	return 8
}

func (a *OfpActionSetMplsTtl) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionDecMplsTtl
 */
func NewOfpActionDecMplsTtl() *OfpActionDecMplsTtl {
	action := new(OfpActionDecMplsTtl)
	header := NewOfpActionHeader(OFPAT_DEC_MPLS_TTL, 8)
	action.ActionHeader = header
	return action
}

func (a *OfpActionDecMplsTtl) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()

	return packet
}

func (a *OfpActionDecMplsTtl) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
}

func (a *OfpActionDecMplsTtl) Size() int {
	return 8
}

func (a *OfpActionDecMplsTtl) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionPush
 */

func NewOfpActionPushVlan() *OfpActionPush {
	action := new(OfpActionPush)
	header := NewOfpActionHeader(OFPAT_PUSH_VLAN, 8)
	action.ActionHeader = header
	action.EtherType = 0x8100

	return action
}

func NewOfpActionPushMpls() *OfpActionPush {
	action := new(OfpActionPush)
	header := NewOfpActionHeader(OFPAT_PUSH_MPLS, 8)
	action.ActionHeader = header
	action.EtherType = 0x8847 //MPLS UNI CAST

	return action
}

func NewOfpActionPushPbb() *OfpActionPush {
	action := new(OfpActionPush)
	header := NewOfpActionHeader(OFPAT_PUSH_PBB, 8)
	action.ActionHeader = header
	action.EtherType = 0x88e7

	return action
}

func NewOfpActionPush(actionType uint16, etherType uint16) *OfpActionPush {
	action := new(OfpActionPush)
	header := NewOfpActionHeader(actionType, 8)
	action.ActionHeader = header
	action.EtherType = etherType

	return action
}

func (a *OfpActionPush) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	binary.BigEndian.PutUint16(packet[index:], a.EtherType)

	return packet
}

func (a *OfpActionPush) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
	a.EtherType = binary.BigEndian.Uint16(packet[index:])
}

func (a *OfpActionPush) Size() int {
	return 8
}

func (a *OfpActionPush) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionPop
 */
func NewOfpActionPopVlan(etherType uint16) *OfpActionPop {
	action := new(OfpActionPop)
	header := NewOfpActionHeader(OFPAT_POP_VLAN, 8)
	action.ActionHeader = header
	action.EtherType = etherType

	return action
}

func NewOfpActionPopMpls(etherType uint16) *OfpActionPop {
	action := new(OfpActionPop)
	header := NewOfpActionHeader(OFPAT_POP_MPLS, 8)
	action.ActionHeader = header
	action.EtherType = etherType

	return action
}

func NewOfpActionPopPbb(etherType uint16) *OfpActionPop {
	action := new(OfpActionPop)
	header := NewOfpActionHeader(OFPAT_POP_PBB, 8)
	action.ActionHeader = header
	action.EtherType = etherType

	return action
}

func (a *OfpActionPop) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	binary.BigEndian.PutUint16(packet[index:], a.EtherType)

	return packet
}

func (a *OfpActionPop) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
	a.EtherType = binary.BigEndian.Uint16(packet[index:])
}

func (a *OfpActionPop) Size() int {
	return 8
}

func (a *OfpActionPop) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionGroup
 */
func NewOfpActionGroup(id uint32) *OfpActionGroup {
	action := new(OfpActionGroup)
	header := NewOfpActionHeader(OFPAT_GROUP, 8)
	action.ActionHeader = header
	action.GroupId = id

	return action
}

func (a *OfpActionGroup) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	binary.BigEndian.PutUint32(packet[index:], a.GroupId)

	return packet
}

func (a *OfpActionGroup) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
	a.GroupId = binary.BigEndian.Uint32(packet[index:])
}

func (a *OfpActionGroup) Size() int {
	return 8
}

func (a *OfpActionGroup) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionSetQueue
 */
func NewOfpActionSetQueue(id uint32) *OfpActionSetQueue {
	action := new(OfpActionSetQueue)
	header := NewOfpActionHeader(OFPAT_SET_QUEUE, 8)
	action.ActionHeader = header
	action.QueueId = id

	return action
}

func (a *OfpActionSetQueue) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	binary.BigEndian.PutUint32(packet[index:], a.QueueId)

	return packet
}

func (a *OfpActionSetQueue) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
	a.QueueId = binary.BigEndian.Uint32(packet[index:])
}

func (a *OfpActionSetQueue) Size() int {
	return 8
}

func (a *OfpActionSetQueue) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionSetNwTtl
 */
func NewOfpActionSetNwTtl(ttl uint8) *OfpActionSetNwTtl {
	action := new(OfpActionSetNwTtl)
	header := NewOfpActionHeader(OFPAT_SET_NW_TTL, 8)
	action.ActionHeader = header
	action.NwTtl = ttl

	return action
}

func (a *OfpActionSetNwTtl) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	packet[index] = a.NwTtl

	return packet
}

func (a *OfpActionSetNwTtl) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
	a.NwTtl = packet[index]
}

func (a *OfpActionSetNwTtl) Size() int {
	return 8
}

func (a *OfpActionSetNwTtl) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionDecNwTtl
 */
func NewOfpActionDecNwTtl() *OfpActionDecNwTtl {
	action := new(OfpActionDecNwTtl)
	header := NewOfpActionHeader(OFPAT_DEC_NW_TTL, 8)
	action.ActionHeader = header
	return action
}

func (a *OfpActionDecNwTtl) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)

	return packet
}

func (a *OfpActionDecNwTtl) Parse(packet []byte) {
	a.ActionHeader.Parse(packet)
}

func (a *OfpActionDecNwTtl) Size() int {
	return 8
}

func (a *OfpActionDecNwTtl) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*
 * OfpActionSetField
 */
// TODO: implements
//
func NewOfpActionSetField(oxm OxmField) *OfpActionSetField {
	a := new(OfpActionSetField)
	length := 4 + oxm.Size()
	length += (8 - (length % 8))
	header := NewOfpActionHeader(OFPAT_SET_FIELD, (uint16)(4+oxm.Size()+2))
	a.ActionHeader = header
	a.Oxm = oxm
	return a
}

func newEmptyOfpActionSetField() *OfpActionSetField {
	return new(OfpActionSetField)
}

func (a *OfpActionSetField) Serialize() []byte {
	index := 0
	packet := make([]byte, a.Size())
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()

	m_packet := a.Oxm.Serialize()
	copy(packet[index:], m_packet)

	return packet
}

func (a *OfpActionSetField) Parse(packet []byte) {
	// TODO: implement OxmField Parser
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()

	// parse tlv header
	tlvheader := binary.BigEndian.Uint32(packet[index:])
	switch oxmField(tlvheader) {
	case OFPXMT_OFB_IN_PORT:
		mf := NewOxmInPort(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IN_PHY_PORT:
		mf := NewOxmInPhyPort(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_METADATA:
		mf := NewOxmMetadata(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ETH_DST:
		mf, err := NewOxmEthDst("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ETH_SRC:
		mf, err := NewOxmEthSrc("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ETH_TYPE:
		mf := NewOxmEthType(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_VLAN_VID:
		mf := NewOxmVlanVid(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_VLAN_PCP:
		mf := NewOxmVlanPcp(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IP_DSCP:
		mf := NewOxmIpDscp(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IP_ECN:
		mf := NewOxmIpEcn(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IP_PROTO:
		mf := NewOxmIpProto(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV4_SRC:
		mf, err := NewOxmIpv4Src("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV4_DST:
		mf, err := NewOxmIpv4Dst("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_TCP_SRC:
		mf := NewOxmTcpSrc(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_TCP_DST:
		mf := NewOxmTcpDst(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_UDP_SRC:
		mf := NewOxmUdpSrc(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_UDP_DST:
		mf := NewOxmUdpDst(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_SCTP_SRC:
		mf := NewOxmSctpSrc(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_SCTP_DST:
		mf := NewOxmSctpDst(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ICMPV4_TYPE:
		mf := NewOxmIcmpType(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ICMPV4_CODE:
		mf := NewOxmIcmpCode(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ARP_OP:
		mf := NewOxmArpOp(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ARP_SPA:
		mf, err := NewOxmArpSpa("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ARP_TPA:
		mf, err := NewOxmArpTpa("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ARP_SHA:
		mf, err := NewOxmArpSha("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ARP_THA:
		mf, err := NewOxmArpTha("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV6_SRC:
		mf, err := NewOxmIpv6Src("::")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV6_DST:
		mf, err := NewOxmIpv6Dst("::")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV6_FLABEL:
		mf := NewOxmIpv6FLabel(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ICMPV6_TYPE:
		mf := NewOxmIcmpv6Type(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_ICMPV6_CODE:
		mf := NewOxmIcmpv6Code(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV6_ND_TARGET:
		mf, err := NewOxmIpv6NdTarget("0.0.0.0")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV6_ND_SLL:
		mf, err := NewOxmIpv6NdSll("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV6_ND_TLL:
		mf, err := NewOxmIpv6NdTll("00:00:00:00:00:00")
		if err != nil {
			// TODO: error handling
		}
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_MPLS_LABEL:
		mf := NewOxmMplsLabel(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_MPLS_TC:
		mf := NewOxmMplsTc(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_MPLS_BOS:
		mf := NewOxmMplsBos(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_PBB_ISID:
		mf := NewOxmPbbIsid([3]uint8{0, 0, 0})
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_TUNNEL_ID:
		mf := NewOxmTunnelId(0)
		mf.Parse(packet)
		a.Oxm = mf
	case OFPXMT_OFB_IPV6_EXTHDR:
		mf := NewOxmIpv6ExtHeader(0)
		mf.Parse(packet)
		a.Oxm = mf
	default:
		//TODO: Error handling
	}

	return
}

func (a *OfpActionSetField) Size() int {
	size := 4 + a.Oxm.Size()
	size += (8 - (size % 8))
	return size
}

func (a *OfpActionSetField) OfpActionType() uint16 {
	return OFPAT_SET_FIELD
}

/*
 * OfpActionExperimenter
 */
func NewOfpActionExperimenter(experimenter uint32) *OfpActionExperimenter {
	action := new(OfpActionExperimenter)
	header := NewOfpActionHeader(OFPAT_EXPERIMENTER, 8)
	action.ActionHeader = header
	action.Experimenter = experimenter

	return action
}

func (a *OfpActionExperimenter) Serialize() []byte {
	index := 0
	packet := make([]byte, 8)
	h_packet := a.ActionHeader.Serialize()
	copy(packet[index:], h_packet)
	index += a.ActionHeader.Size()
	binary.BigEndian.PutUint32(packet[index:], a.Experimenter)

	return packet
}

func (a *OfpActionExperimenter) Parse(packet []byte) {
	index := 0
	a.ActionHeader.Parse(packet)
	index += a.ActionHeader.Size()
	a.Experimenter = binary.BigEndian.Uint32(packet[index:])
}

func (a *OfpActionExperimenter) Size() int {
	return 8
}

func (a *OfpActionExperimenter) OfpActionType() uint16 {
	return a.ActionHeader.Type
}

/*****************************************************/
/* OfpErrorMsg                                       */
/*****************************************************/
func NewOfpErrorMsg() *OfpErrorMsg {
	header := NewOfpHeader(OFPT_ERROR)
	header.Type = OFPT_ERROR
	return nil
}

func (m *OfpErrorMsg) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()
	copy(packet[0:], h_packet)
	index := m.Header.Size()
	binary.BigEndian.PutUint16(packet[index:], m.Type)
	index += 2
	binary.BigEndian.PutUint16(packet[index:], m.Code)
	index += 2
	for _, d := range m.Data {
		packet[index] = d
		index += 1
	}
	return packet
}

func (m *OfpErrorMsg) Parse(packet []byte) {
	m.Header.Parse(packet)
	index := m.Header.Size()
	m.Type = binary.BigEndian.Uint16(packet[index:])
	index += 2
	m.Code = binary.BigEndian.Uint16(packet[index:])
	index += 2
	for int(index) < len(packet) {
		m.Data = append(m.Data, packet[index])
		index += 1
	}
}

func (m *OfpErrorMsg) Size() int {
	return m.Header.Size() + 8 + len(m.Data)
}

/*****************************************************/
/* OfpErrorExperimenterMsg                           */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpMultipartRequest                               */
/*****************************************************/
// TODO: implement
/**
OFPMP_DESC
OFPMP_FLOW
OFPMP_AGGREGATE
OFPMP_TABLE
OFPMP_PORT_STATS
OFPMP_QUEUE
OFPMP_GROUP
OFPMP_GROUP_DESC
OFPMP_GROUP_FEATURES
OFPMP_METER
OFPMP_METER_CONFIG
OFPMP_METER_FEATURES
OFPMP_TABLE_FEATURES
OFPMP_PORT_DESC
OFPMP_EXPERIMENTER
*/

func NewOfpDescStatsRequest(flags uint16) *OfpMultipartRequest {
	m := NewOfpMultipartRequest(OFPMP_DESC, flags)
	return m
}

func NewOfpFlowStatsRequest(
	flags uint16,
	tableId uint8,
	outPort uint32,
	outGroup uint32,
	cookie uint64,
	cookieMask uint64,
	match *OfpMatch) *OfpMultipartRequest {
	m := NewOfpMultipartRequest(OFPMP_FLOW, flags)
	m.Body = newOfpFlowStatsRequestBody(
		tableId,
		outPort,
		outGroup,
		cookie,
		cookieMask,
		match)
	m.Header.Length += (uint16)(m.Body.Size())
	return m
}

func NewOfpAggregateStatsRequest(
	flags uint16,
	tableId uint8,
	outPort uint32,
	outGroup uint32,
	cookie uint64,
	cookieMask uint64,
	match *OfpMatch) *OfpMultipartRequest {
	m := NewOfpMultipartRequest(OFPMP_AGGREGATE, flags)
	m.Body = newOfpAggregateStatsRequestBody(
		tableId,
		outPort,
		outGroup,
		cookie,
		cookieMask,
		match)
	m.Header.Length += (uint16)(m.Body.Size())
	return m
}

func NewOfpTableStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpPortStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpQueueStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpGroupStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpGroupDescStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpGroupFeaturesStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpMeterStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpMeterConfigStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpMeterFeaturesStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpTableFeaturesStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpPortDescStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpExperimenterStatsRequest() *OfpMultipartRequest {
	return nil
}

func NewOfpMultipartRequest(t uint16, flags uint16) *OfpMultipartRequest {
	m := new(OfpMultipartRequest)
	m.Header = NewOfpHeader(OFPT_MULTIPART_REQUEST)
	m.Header.Length = (uint16)(m.Header.Size()) + 8
	m.Type = t
	m.Flags = flags
	return m
}

func (m *OfpMultipartRequest) Serialize() []byte {
	packet := make([]byte, m.Size())
	h_packet := m.Header.Serialize()

	index := 0
	copy(packet[index:], h_packet)
	index += m.Header.Size()

	binary.BigEndian.PutUint16(packet[index:], m.Type)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.Flags)
	index += 6

	if m.Body != nil {
		mp_packet := m.Body.Serialize()
		copy(packet[index:], mp_packet)
		index += m.Body.Size()
	}

	return packet
}

func (m *OfpMultipartRequest) Parse(packet []byte) {
	return
}

func (m *OfpMultipartRequest) Size() int {
	size := m.Header.Size() + 8
	if m.Body != nil {
		size += m.Body.Size()
	}
	return size
}

/*****************************************************/
/* OfpMultipartReply                                 */
/*****************************************************/
// TODO: implement
func NewOfpMultipartReply() *OfpMultipartReply {
	m := new(OfpMultipartReply)
	header := NewOfpHeader(OFPT_MULTIPART_REPLY)
	m.Header = header
	// m.Type = t
	// m.Flags = flags
	m.Body = make([]OfpMultipartBody, 0)
	return m
}

func (m *OfpMultipartReply) Serialize() []byte {
	return nil
}

func (m *OfpMultipartReply) Parse(packet []byte) {
	index := 0
	m.Header.Parse(packet[index:])
	index += m.Header.Size()

	m.Type = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.Flags = binary.BigEndian.Uint16(packet[index:])
	index += 6

	switch m.Type {
	case OFPMP_DESC:
		for (uint16)(index) < m.Header.Length {
			mp := newOfpDescStats()
			mp.Parse(packet[index:])
			m.Append(mp)
			index += mp.Size()
		}
	case OFPMP_FLOW:
		for (uint16)(index) < m.Header.Length {
			mp := newOfpFlowStats()
			mp.Parse(packet[index:])
			m.Append(mp)
			index += mp.Size()
		}
	case OFPMP_AGGREGATE:
		for (uint16)(index) < m.Header.Length {
			mp := newOfpAggregateStats()
			mp.Parse(packet[index:])
			m.Append(mp)
			index += mp.Size()
		}
	case OFPMP_TABLE:
		// TODO: implements
	case OFPMP_PORT_STATS:
		// TODO: implements
	case OFPMP_QUEUE:
		// TODO: implements
	case OFPMP_GROUP:
		// TODO: implements
	case OFPMP_GROUP_DESC:
		// TODO: implements
	case OFPMP_GROUP_FEATURES:
		// TODO: implements
	case OFPMP_METER:
		// TODO: implements
	case OFPMP_METER_CONFIG:
		// TODO: implements
	case OFPMP_METER_FEATURES:
		// TODO: implements
	case OFPMP_TABLE_FEATURES:
		// TODO: implements
	case OFPMP_PORT_DESC:
		// TODO: implements
	case OFPMP_EXPERIMENTER:
		// TODO: implements
	default:
	}

	return
}

func (m *OfpMultipartReply) Size() int {
	size := m.Header.Size() + 8
	for _, mp := range m.Body {
		size += mp.Size()
	}
	return size
}

func (m *OfpMultipartReply) Append(mp OfpMultipartBody) {
	m.Body = append(m.Body, mp)
}

/*****************************************************/
/* OfpDesc                                           */
/*****************************************************/
// TODO: implement
func newOfpDescStats() *OfpDescStats {
	mp := new(OfpDescStats)
	mp.MfrDesc = make([]byte, DESC_STR_LEN)
	mp.HwDesc = make([]byte, DESC_STR_LEN)
	mp.SwDesc = make([]byte, DESC_STR_LEN)
	mp.SerialNum = make([]byte, SERIAL_NUM_LEN)
	mp.DpDesc = make([]byte, DESC_STR_LEN)
	return mp
}

func (mp *OfpDescStats) Serialize() []byte {
	// not implement
	return nil
}

func (mp *OfpDescStats) Parse(packet []byte) {
	index := 0
	copy(mp.MfrDesc, packet[index:(index+DESC_STR_LEN)])
	index += DESC_STR_LEN
	copy(mp.HwDesc, packet[index:(index+DESC_STR_LEN)])
	index += DESC_STR_LEN
	copy(mp.SwDesc, packet[index:(index+DESC_STR_LEN)])
	index += DESC_STR_LEN
	copy(mp.SerialNum, packet[index:(index+SERIAL_NUM_LEN)])
	index += SERIAL_NUM_LEN
	copy(mp.DpDesc, packet[index:(index+DESC_STR_LEN)])
	index += DESC_STR_LEN

	return
}

func (mp *OfpDescStats) Size() int {
	//return len(mp.MfrDesc) + len(mp.HwDesc) +
	//	len(mp.SwDesc) + len(mp.SerialNum) + len(mp.DpDesc)
	return 1056
}

func (mp *OfpDescStats) MPType() uint16 {
	return OFPMP_DESC
}

/*****************************************************/
/* OfpFlowStatsRequest                               */
/*****************************************************/
// TODO: implement
func newOfpFlowStatsRequestBody(
	tableId uint8,
	outPort uint32,
	outGroup uint32,
	cookie uint64,
	cookieMask uint64,
	match *OfpMatch) *OfpFlowStatsRequest {
	req := new(OfpFlowStatsRequest)
	req.TableId = tableId
	req.OutPort = outPort
	req.OutGroup = outGroup
	req.Cookie = cookie
	req.CookieMask = cookieMask
	req.Match = match
	return req
}

func (m *OfpFlowStatsRequest) Serialize() []byte {
	packet := make([]byte, m.Size())
	index := 0

	packet[index] = m.TableId
	index += 4

	binary.BigEndian.PutUint32(packet[index:], m.OutPort)
	index += 4

	binary.BigEndian.PutUint32(packet[index:], m.OutGroup)
	index += 8

	binary.BigEndian.PutUint64(packet[index:], m.Cookie)
	index += 8

	binary.BigEndian.PutUint64(packet[index:], m.CookieMask)
	index += 8

	m_packet := m.Match.Serialize()
	copy(packet[index:], m_packet)

	return packet
}

func (m *OfpFlowStatsRequest) Parse(packet []byte) {
	return
}

func (m *OfpFlowStatsRequest) Size() int {
	size := 32 + m.Match.Size()
	return size
}

func (m *OfpFlowStatsRequest) MPType() uint16 {
	return OFPMP_FLOW
}

/*****************************************************/
/* OfpFlowStats                                      */
/*****************************************************/
// TODO: implement
func newOfpFlowStats() *OfpFlowStats {
	m := new(OfpFlowStats)
	return m
}

func (mp *OfpFlowStats) Serialize() []byte {
	return nil
}

func (mp *OfpFlowStats) Parse(packet []byte) {
	index := 0
	mp.Length = binary.BigEndian.Uint16(packet[index:])
	index += 2

	mp.TableId = packet[index]
	index += 2 // include Padding

	mp.DurationSec = binary.BigEndian.Uint32(packet[index:])
	index += 4

	mp.DurationNSec = binary.BigEndian.Uint32(packet[index:])
	index += 4

	mp.Priority = binary.BigEndian.Uint16(packet[index:])
	index += 2

	mp.IdleTimeout = binary.BigEndian.Uint16(packet[index:])
	index += 2

	mp.HardTimeout = binary.BigEndian.Uint16(packet[index:])
	index += 2

	mp.Flags = binary.BigEndian.Uint16(packet[index:])
	index += 6 // include Padding

	mp.Cookie = binary.BigEndian.Uint64(packet[index:])
	index += 8

	mp.PacketCount = binary.BigEndian.Uint64(packet[index:])
	index += 8

	mp.ByteCount = binary.BigEndian.Uint64(packet[index:])
	index += 8

	mp.Match = NewOfpMatch()
	mp.Match.Parse(packet[index:])
	index += mp.Match.Size()

	for index < (int)(mp.Length) {
		t := binary.BigEndian.Uint16(packet[index:])

		// don't forward index ,
		// because type and length will be parsed in instruction's paraser

		// TODO: implements parse process for other instruction type
		switch t {
		case OFPIT_GOTO_TABLE:
			instruction := NewOfpInstructionGotoTable(0)
			instruction.Parse(packet[index:])
			mp.Instructions = append(mp.Instructions, instruction)
			index += instruction.Size()
		case OFPIT_WRITE_METADATA:
			instruction := NewOfpInstructionWriteMetadata(0, 0)
			instruction.Parse(packet[index:])
			mp.Instructions = append(mp.Instructions, instruction)
			index += instruction.Size()
		case OFPIT_WRITE_ACTIONS:
			instruction := NewOfpInstructionActions(OFPIT_WRITE_ACTIONS)
			instruction.Parse(packet[index:])
			mp.Instructions = append(mp.Instructions, instruction)
			index += instruction.Size()
		case OFPIT_APPLY_ACTIONS:
			instruction := NewOfpInstructionActions(OFPIT_APPLY_ACTIONS)
			instruction.Parse(packet[index:])
			mp.Instructions = append(mp.Instructions, instruction)
			index += instruction.Size()
		case OFPIT_CLEAR_ACTIONS:
			instruction := NewOfpInstructionActions(OFPIT_CLEAR_ACTIONS)
			instruction.Parse(packet[index:])
			mp.Instructions = append(mp.Instructions, instruction)
			index += instruction.Size()
		case OFPIT_METER:
		case OFPIT_EXPERIMENTER:
		default:

		}
	}

	return
}

func (mp *OfpFlowStats) Size() int {
	size := 48 + mp.Match.Size()
	for _, i := range mp.Instructions {
		size += i.Size()
	}
	return size
}

func (mp *OfpFlowStats) MPType() uint16 {
	return OFPMP_FLOW
}

/*****************************************************/
/* OfpAggregateStatsRequest                          */
/*****************************************************/
// TODO: implement
func newOfpAggregateStatsRequestBody(
	tableId uint8,
	outPort uint32,
	outGroup uint32,
	cookie uint64,
	cookieMask uint64,
	match *OfpMatch) *OfpAggregateStatsRequest {
	req := new(OfpAggregateStatsRequest)
	req.TableId = tableId
	req.OutPort = outPort
	req.OutGroup = outGroup
	req.Cookie = cookie
	req.CookieMask = cookieMask
	req.Match = match
	return req
}

func (mp *OfpAggregateStatsRequest) Serialize() []byte {
	packet := make([]byte, mp.Size())
	index := 0

	packet[index] = mp.TableId
	index += 4

	binary.BigEndian.PutUint32(packet[index:], mp.OutPort)
	index += 4

	binary.BigEndian.PutUint32(packet[index:], mp.OutGroup)
	index += 8

	binary.BigEndian.PutUint64(packet[index:], mp.Cookie)
	index += 8

	binary.BigEndian.PutUint64(packet[index:], mp.CookieMask)
	index += 8

	m_packet := mp.Match.Serialize()
	copy(packet[index:], m_packet)

	return packet
}

func (mp *OfpAggregateStatsRequest) Parse(packet []byte) {
	return
}

func (mp *OfpAggregateStatsRequest) Size() int {
	size := 32 + mp.Match.Size()
	return size
}

func (mp *OfpAggregateStatsRequest) MPType() uint16 {
	return OFPMP_AGGREGATE
}

/*****************************************************/
/* OfpAggregateStats                               */
/*****************************************************/
// TODO: implement
func newOfpAggregateStats() *OfpAggregateStats {
	mp := new(OfpAggregateStats)
	return mp
}

func (mp *OfpAggregateStats) Serialize() []byte {
	return nil
}

func (mp *OfpAggregateStats) Parse(packet []byte) {
	index := 0
	mp.PacketCount = binary.BigEndian.Uint64(packet[index:])
	index += 8

	mp.ByteCount = binary.BigEndian.Uint64(packet[index:])
	index += 8

	mp.FlowCount = binary.BigEndian.Uint32(packet[index:])
	return
}

func (mp *OfpAggregateStats) Size() int {
	return 24
}

func (mp *OfpAggregateStats) MPType() uint16 {
	return OFPMP_AGGREGATE
}

/*****************************************************/
/* OfpTableFeaturePropHeader                         */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableFeaturePropInstructions                   */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableFeaturePropNextTables                     */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableFeaturePropActions                        */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableFeaturePropOxm                            */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableFeaturePropExperimenter                   */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableFeatures                                  */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableStatsRequest                              */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpTableStats                                     */
/*****************************************************/
// TODO: implement
func newOfpTableStats() *OfpTableStats {
	return nil
}

func (mp *OfpTableStats) Serialize() []byte {
	return nil
}

func (mp *OfpTableStats) Parse(packet []byte) {
	return
}

func (mp *OfpTableStats) Size() int {
	return 0
}

func (mp *OfpTableStats) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpPortStatsRequest                               */
/*****************************************************/
// TODO: implement
func newOfpPortStatsRequestBody() *OfpPortStatsRequest {
	return nil
}

func (mp *OfpPortStatsRequest) Serialize() []byte {
	return nil
}

func (mp *OfpPortStatsRequest) Parse(packet []byte) {
	return
}

func (mp *OfpPortStatsRequest) Size() int {
	return 0
}

func (mp *OfpPortStatsRequest) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpPortStats                                      */
/*****************************************************/
// TODO: implement
func newOfpPortStats() *OfpPortStats {
	return nil
}

func (mp *OfpPortStats) Serialize() []byte {
	return nil
}

func (mp *OfpPortStats) Parse(packet []byte) {
	return
}

func (mp *OfpPortStats) Size() int {
	return 0
}

func (mp *OfpPortStats) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpQueueStatsRequest                              */
/*****************************************************/
// TODO: implement
func newOfpQueueStatsRequestBody() *OfpQueueStatsRequest {
	return nil
}

func (mp *OfpQueueStatsRequest) Serialize() []byte {
	return nil
}

func (mp *OfpQueueStatsRequest) Parse(packet []byte) {
	return
}

func (mp *OfpQueueStatsRequest) Size() int {
	return 0
}

func (mp *OfpQueueStatsRequest) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpQueueStats                                     */
/*****************************************************/
// TODO: implement
func newOfpQueueStats() *OfpQueueStats {
	return nil
}

func (mp *OfpQueueStats) Serialize() []byte {
	return nil
}

func (mp *OfpQueueStats) Parse(packet []byte) {
	return
}

func (mp *OfpQueueStats) Size() int {
	return 0
}

func (mp *OfpQueueStats) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpGroupStatsRequest                              */
/*****************************************************/
// TODO: implement
func newOfpGroupStatsRequestBody() *OfpGroupStatsRequest {
	return nil
}

func (mp *OfpGroupStatsRequest) Serialize() []byte {
	return nil
}

func (mp *OfpGroupStatsRequest) Parse(packet []byte) {
	return
}

func (mp *OfpGroupStatsRequest) Size() int {
	return 0
}

func (mp *OfpGroupStatsRequest) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpGroupStats                                     */
/*****************************************************/
// TODO: implement
func newOfpGroupStats() *OfpGroupStats {
	return nil
}

func (mp *OfpGroupStats) Serialize() []byte {
	return nil
}

func (mp *OfpGroupStats) Parse(packet []byte) {
	return
}

func (mp *OfpGroupStats) Size() int {
	return 0
}

func (mp *OfpGroupStats) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpGroupFeatures                                  */
/*****************************************************/
// TODO: implement
func newOfpGroupFeaturesStats() *OfpGroupFeaturesStats {
	return nil
}

func (mp *OfpGroupFeaturesStats) Serialize() []byte {
	return nil
}

func (mp *OfpGroupFeaturesStats) Parse(packet []byte) {
	return
}

func (mp *OfpGroupFeaturesStats) Size() int {
	return 0
}

func (mp *OfpGroupFeaturesStats) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpMeterMultipartRequest                          */
/*****************************************************/
// TODO: implement
func newOfpMeterMultipartRequestBody() *OfpMeterMultipartRequest {
	return nil
}

func (mp *OfpMeterMultipartRequest) Serialize() []byte {
	return nil
}

func (mp *OfpMeterMultipartRequest) Parse(packet []byte) {
	return
}

func (mp *OfpMeterMultipartRequest) Size() int {
	return 0
}

func (mp *OfpMeterMultipartRequest) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpMeterStats                                      */
/*****************************************************/
// TODO: implement
func newOfpMeterStats() *OfpMeterStats {
	return nil
}

func (mp *OfpMeterStats) Serialize() []byte {
	return nil
}

func (mp *OfpMeterStats) Parse(packet []byte) {
	return
}

func (mp *OfpMeterStats) Size() int {
	return 0
}

func (mp *OfpMeterStats) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpMeterConfig                                    */
/*****************************************************/
// TODO: implement
func newOfpMeterConfig() *OfpMeterConfig {
	return nil
}

func (mp *OfpMeterConfig) Serialize() []byte {
	return nil
}

func (mp *OfpMeterConfig) Parse(packet []byte) {
	return
}

func (mp *OfpMeterConfig) Size() int {
	return 0
}

func (mp *OfpMeterConfig) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpMeterFeatures                                  */
/*****************************************************/
// TODO: implement
func newOfpMeterFeatures() *OfpMeterFeatures {
	return nil
}

func (mp *OfpMeterFeatures) Serialize() []byte {
	return nil
}

func (mp *OfpMeterFeatures) Parse(packet []byte) {
	return
}

func (mp *OfpMeterFeatures) Size() int {
	return 0
}

func (mp *OfpMeterFeatures) MPType() uint16 {
	return 0
}

/*****************************************************/
/* OfpExperimenterMultipartHeader                    */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpQueuePropHeader                                */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpQueuePropMinRate                               */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpQueuePropMaxRate                               */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpQueuePropExperimenter                          */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpPacketQueue                                    */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpQueueGetConfigRequest                          */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpActionSetQueue                                 */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpQueueStatsRequest                              */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpQueueStats                                     */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpRoleRequest                                    */
/*****************************************************/
// TODO: implement

/*****************************************************/
/* OfpAsyncConfig                                    */
/*****************************************************/
// TODO: implement
