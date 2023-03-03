package tunnels

const (
	AuthReq uint8 = iota
	AuthRes
	NewConReq
	NewConRes
	Ping
	Pong
)

type AuthReqMsg struct {
	Token       string // tunnel token
	Type        string // tunnel type: RDP, HTTP, TCP
	Destination string
}

type AuthResMsg struct {
	PublicAddr string // secret token
}

type NewConReqMsg struct {
	Session string
}

type NewConResMsg struct {
	Session string
}
