package tunnels

const (
	AuthReq uint8 = iota
	AuthRes
	NewConReq
	NewConRes
	Ping
	Pong
	Stat
)

type AuthReqMsg struct {
	Hash        string // tunnel hash
	APIKey      string // user secret
	Type        string // tunnel type: RDP, HTTP, TCP
	Destination string
	Extra       map[string]string
}

type Message struct {
	Level   uint8
	Message string
}

type AuthResMsg struct {
	Message    Message
	PublicAddr string // secret token
}

type NewConReqMsg struct {
	Session string
}

type NewConResMsg struct {
	Session string
}
