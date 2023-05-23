package network

const (
	Packet = iota
	Rpc
)

type PacketBase struct {
	PacketID      int
	PacketType    int
	PacketChannel int
}
