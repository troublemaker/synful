package packetfactory

import (
	"encoding/binary"
	"net"
	"runtime"
)

/*
   https://datatracker.ietf.org/doc/html/rfc791

    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Version|  IHL  |Type of Service|          Total Length         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |         Identification        |Flag|      Fragment Offset    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Time to Live |    Protocol   |         Header Checksum       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                       Source Address                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Destination Address                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options                    |    Padding    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

const ip4Version = 4
const ip4MinHeaderLen = 20
const ip4MaxHeaderLen = 60 //todo:check max size

type IP4 struct {
	Version    int
	TOS        int
	ID         int
	Flag       int
	FragOffset int
	TTL        int
	Protocol   int
	Src        []byte
	Dst        []byte
	Options    []byte
}

func (ip4 *IP4) Encode(raw []byte, start int, end int) (int, int) {
	optionsSize := len(ip4.Options)
	ihl := ip4MinHeaderLen + optionsSize

	start = start - ihl
	buffer := raw[start:end]
	ip4.buildHeader(buffer)
	return start, end
}

func (ip4 *IP4) buildHeader(buffer []byte) {

	//calculated fields
	optionsSize := len(ip4.Options)
	ihl := ip4MinHeaderLen + optionsSize

	totalLength := len(buffer)

	buffer[0] = byte(ip4Version<<4 | ihl>>2)
	buffer[1] = byte(ip4.TOS)

	binary.BigEndian.PutUint16(buffer[4:6], uint16(ip4.ID))

	flagsFragOffset := ip4.Flag<<13 | (ip4.FragOffset & 0x1fff)

	//Weird ordering in MacOSX / FreeBSD
	//https://cseweb.ucsd.edu//~braghava/notes/freebsd-sockets.txt
	if runtime.GOOS == "darwin" {
		binary.LittleEndian.PutUint16(buffer[2:4], uint16(totalLength))
		binary.LittleEndian.PutUint16(buffer[6:8], uint16(flagsFragOffset))
	} else {
		binary.BigEndian.PutUint16(buffer[2:4], uint16(totalLength))
		binary.BigEndian.PutUint16(buffer[6:8], uint16(flagsFragOffset))
	}

	buffer[8] = byte(ip4.TTL)
	buffer[9] = byte(ip4.Protocol)

	copy(buffer[12:16], ip4.Src[:net.IPv4len])
	copy(buffer[16:20], ip4.Dst[:net.IPv4len])

	if len(ip4.Options) > 0 {
		copy(buffer[ip4MinHeaderLen:], ip4.Options)
	}

	binary.BigEndian.PutUint16(buffer[10:12], CheckSum(buffer[:ihl]))

}
