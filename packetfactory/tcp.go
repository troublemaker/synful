package packetfactory

import (
	"encoding/binary"
	"net"
	"syscall"
)

/*
	https://datatracker.ietf.org/doc/html/rfc793

    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          Source Port          |       Destination Port        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                        Sequence Number                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Acknowledgment Number                      |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Data |           |U|A|P|R|S|F|                               |
   | Offset| Reserved  |R|C|S|S|Y|I|            Window             |
   |       |           |G|K|H|T|N|N|                               |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           Checksum            |         Urgent Pointer        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options                    |    Padding    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                             data                              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/
const tcpMinHeaderLen = 20
const tcpMaxHeaderLen = 60 //todo:check max size

type TCP struct {
	SrcPort int
	DstPort int
	Seq     int
	Ack     int
	Len     int
	Flag    int
	WS      int
	Urp     int
	Src     []byte
	Dst     []byte
	Options []byte
	Data    []byte
}

func (tcp *TCP) Encode(raw []byte, start int, end int) (int, int) {
	thl := tcpMinHeaderLen + len(tcp.Options)
	dataLen := len(tcp.Data)

	start = start - thl
	end = end + dataLen
	buffer := raw[start:end]

	//build TCP header w/o checksum
	tcp.buildHeader(buffer)

	//Calc pseudo header & TCP checksum
	//1. reserve 12 bytes for the pseudo header
	bufferPD := raw[start-12 : end]

	//2. copy Src & Dst IPs
	copy(bufferPD[0:4], tcp.Src[:net.IPv4len])
	copy(bufferPD[4:8], tcp.Dst[:net.IPv4len])

	//3. reserved byte, protocol & total segment length
	bufferPD[8] = 0
	bufferPD[9] = syscall.IPPROTO_TCP
	binary.BigEndian.PutUint16(bufferPD[10:12], uint16(end-start))

	//4. Calc & Store Checksum
	tcpCheckSum := CheckSum(bufferPD)
	binary.BigEndian.PutUint16(buffer[16:18], tcpCheckSum)

	return start, end
}

func (tcp *TCP) buildHeader(buffer []byte) {

	thl := tcpMinHeaderLen + len(tcp.Options)

	binary.BigEndian.PutUint16(buffer[0:2], uint16(tcp.SrcPort))
	binary.BigEndian.PutUint16(buffer[2:4], uint16(tcp.DstPort))

	binary.BigEndian.PutUint32(buffer[4:8], uint32(tcp.Seq))
	binary.BigEndian.PutUint32(buffer[8:12], uint32(tcp.Ack))

	buffer[12] = byte(thl>>2) << 4 //TCP header size + 3 reserved bits
	buffer[13] = uint8(tcp.Flag)

	binary.BigEndian.PutUint16(buffer[14:16], uint16(tcp.WS))
	binary.BigEndian.PutUint16(buffer[16:18], uint16(0)) //checksum is calculated later
	binary.BigEndian.PutUint16(buffer[18:20], uint16(tcp.Urp))

	if len(tcp.Options) > 0 {
		copy(buffer[thl:], tcp.Options)
	}

}
