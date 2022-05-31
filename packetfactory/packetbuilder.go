package packetfactory

type PacketBuilder struct {
	buffer      []byte
	stack       []Protocol
	DebugBuffer []byte
}

func (pb *PacketBuilder) Init(bufferSize int, protocols ...Protocol) {
	pb.buffer = make([]byte, bufferSize, bufferSize)
	pb.stack = protocols

	pb.DebugBuffer = pb.buffer //for debug
}

func (pb *PacketBuilder) Build() []byte {
	//top most packet starts from the middle, so we could grow data to the right
	//and lower layers headers to the left without re-allocating/copying memory
	windowStart := len(pb.buffer) / 2
	windowEnd := windowStart

	for _, v := range pb.stack {
		windowStart, windowEnd = v.Encode(pb.buffer, windowStart, windowEnd)
	}

	return pb.buffer[windowStart:windowEnd]
}

func CheckSum(data []byte) uint16 {
	var sum uint32
	var i int
	remainingLen := len(data)

	for remainingLen > 1 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
		i += 2
		remainingLen -= 2
	}

	if remainingLen > 0 {
		sum += uint32(data[i])
	}
	sum += sum >> 16

	return uint16(^sum)
}
