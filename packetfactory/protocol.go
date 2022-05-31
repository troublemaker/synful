package packetfactory

type Protocol interface {
	Encode([]byte, int, int) (int, int)
}
