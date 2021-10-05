package vsroguemysql

import (
	"bufio"
	"bytes"
	"net"

	"github.com/vikingo-project/vsat/utils"
)

const MaxPayloadLen int = 1<<24 - 1

type ClientCapabilities struct {
	LongPassword              bool
	FoundRows                 bool
	LongColumn                bool
	ConnectWithDatabase       bool
	DontAllowDatabaseDotTable bool
	CanCompression            bool
	ODBCClient                bool
	CanUseLoadDataLocal       bool
	IgnoreSpacesBrackets      bool
	Proto41                   bool
	Interactive               bool
	ToSSLAfterHandshake       bool
	IgnoreSigpipes            bool
	KnowAboutTransactions     bool
	Proto41Old                bool
	Proto41Auth               bool
}

type Session struct {
	net.Conn
	// conn     net.Conn
	state  string
	reader *bufio.Reader

	FileData  *bytes.Buffer
	PacketNum int
}

func (sess *Session) SetState(s string) {
	utils.PrintDebug("set state %s", s)
	sess.state = s
}

func (sess *Session) GetState() string {
	return sess.state
}
