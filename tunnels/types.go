package tunnels

type AuthReqMsg struct {
	Token string // tunnel token
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
