package vsftp

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/vikingo-project/vsat/utils"
)

// Command represents a Command interface to a ftp command
type Command interface {
	IsExtend() bool
	RequireParam() bool
	RequireAuth() bool
	Execute(*Conn, string)
}

type commandMap map[string]Command

var (
	commands = commandMap{
		"CCC":  commandCcc{},
		"CONF": commandConf{},
		"ENC":  commandEnc{},
		"EPRT": commandEprt{},
		"EPSV": commandEpsv{},
		"FEAT": commandFeat{},
		"LIST": commandList{},
		"LPRT": commandLprt{},
		"NLST": commandNlst{},
		"MDTM": commandMdtm{},
		"MIC":  commandMic{},
		"OPTS": commandOpts{},
		"PBSZ": commandPbsz{},
		"PORT": commandPort{},
		"PROT": commandProt{},
		"REST": commandRest{},
		"SIZE": commandSize{},
		"STRU": commandStru{},
	}
)

type commandOpts struct{}

func (cmd commandOpts) IsExtend() bool {
	return false
}

func (cmd commandOpts) RequireParam() bool {
	return false
}

func (cmd commandOpts) RequireAuth() bool {
	return false
}

func (cmd commandOpts) Execute(conn *Conn, param string) {
	parts := strings.Fields(param)
	if len(parts) != 2 {
		conn.writeMessage(550, "Unknow params")
		return
	}
	if strings.ToUpper(parts[0]) != "UTF8" {
		conn.writeMessage(550, "Unknow params")
		return
	}

	if strings.ToUpper(parts[1]) == "ON" {
		conn.writeMessage(200, "UTF8 mode enabled")
	} else {
		conn.writeMessage(550, "Unsupported non-utf8 mode")
	}
}

type commandFeat struct{}

func (cmd commandFeat) IsExtend() bool {
	return false
}

func (cmd commandFeat) RequireParam() bool {
	return false
}

func (cmd commandFeat) RequireAuth() bool {
	return false
}

var (
	feats    = "Extensions supported:\n%s"
	featCmds = " UTF8\n"
)

func init() {
	for k, v := range commands {
		if v.IsExtend() {
			featCmds = featCmds + " " + k + "\n"
		}
	}
}

func (cmd commandFeat) Execute(conn *Conn, param string) {
	conn.writeMessageMultiline(211, conn.server.feats)
}

// commandEprt responds to the EPRT FTP command. It allows the client to
// request an active data socket with more options than the original PORT
// command. It mainly adds ipv6 support.
type commandEprt struct{}

func (cmd commandEprt) IsExtend() bool {
	return true
}

func (cmd commandEprt) RequireParam() bool {
	return true
}

func (cmd commandEprt) RequireAuth() bool {
	return true
}

func (cmd commandEprt) Execute(conn *Conn, param string) {
	delim := string(param[0:1])
	parts := strings.Split(param, delim)
	addressFamily, err := strconv.Atoi(parts[1])
	if err != nil {
		conn.writeMessage(522, "Network protocol not supported, use (1,2)")
		return
	}
	host := parts[2]
	port, err := strconv.Atoi(parts[3])
	if err != nil {
		conn.writeMessage(522, "Network protocol not supported, use (1,2)")
		return
	}
	if addressFamily != 1 && addressFamily != 2 {
		conn.writeMessage(522, "Network protocol not supported, use (1,2)")
		return
	}
	socket, err := newActiveSocket(host, port, conn.session)
	if err != nil {
		conn.writeMessage(425, "Data connection failed")
		return
	}
	conn.dataConn = socket
	conn.writeMessage(200, "Connection established ("+strconv.Itoa(port)+")")
}

// commandLprt responds to the LPRT FTP command. It allows the client to
// request an active data socket with more options than the original PORT
// command.  FTP Operation Over Big Address Records.
type commandLprt struct{}

func (cmd commandLprt) IsExtend() bool {
	return true
}

func (cmd commandLprt) RequireParam() bool {
	return true
}

func (cmd commandLprt) RequireAuth() bool {
	return true
}

func (cmd commandLprt) Execute(conn *Conn, param string) {
	// No tests for this code yet

	parts := strings.Split(param, ",")

	addressFamily, err := strconv.Atoi(parts[0])
	if err != nil {
		conn.writeMessage(522, "Network protocol not supported, use 4")
		return
	}
	if addressFamily != 4 {
		conn.writeMessage(522, "Network protocol not supported, use 4")
		return
	}

	addressLength, err := strconv.Atoi(parts[1])
	if err != nil {
		conn.writeMessage(522, "Network protocol not supported, use 4")
		return
	}
	if addressLength != 4 {
		conn.writeMessage(522, "Network IP length not supported, use 4")
		return
	}

	host := strings.Join(parts[2:2+addressLength], ".")

	portLength, err := strconv.Atoi(parts[2+addressLength])
	if err != nil {
		conn.writeMessage(522, "Network protocol not supported, use 4")
		return
	}
	portAddress := parts[3+addressLength : 3+addressLength+portLength]

	// Convert string[] to byte[]
	portBytes := make([]byte, portLength)
	for i := range portAddress {
		p, _ := strconv.Atoi(portAddress[i])
		portBytes[i] = byte(p)
	}

	// convert the bytes to an int
	port := int(binary.BigEndian.Uint16(portBytes))

	// if the existing connection is on the same host/port don't reconnect
	if conn.dataConn.Host() == host && conn.dataConn.Port() == port {
		return
	}

	socket, err := newActiveSocket(host, port, conn.session)
	if err != nil {
		conn.writeMessage(425, "Data connection failed")
		return
	}
	conn.dataConn = socket
	conn.writeMessage(200, "Connection established ("+strconv.Itoa(port)+")")
}

// commandEpsv responds to the EPSV FTP command. It allows the client to
// request a passive data socket with more options than the original PASV
// command. It mainly adds ipv6 support, although we don't support that yet.
type commandEpsv struct{}

func (cmd commandEpsv) IsExtend() bool {
	return true
}

func (cmd commandEpsv) RequireParam() bool {
	return false
}

func (cmd commandEpsv) RequireAuth() bool {
	return true
}

func (cmd commandEpsv) Execute(conn *Conn, param string) {
	socket, err := conn.newPassiveSocket()
	if err != nil {
		log.Println(conn.session, "%s\n", err)
		conn.writeMessage(425, "Data connection failed")
		return
	}

	msg := fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", socket.Port())
	conn.writeMessage(229, msg)
}

// commandList responds to the LIST FTP command. It allows the client to retreive
// a detailed listing of the contents of a directory.
type commandList struct{}

func (cmd commandList) IsExtend() bool {
	return false
}

func (cmd commandList) RequireParam() bool {
	return false
}

func (cmd commandList) RequireAuth() bool {
	return true
}

func (cmd commandList) Execute(conn *Conn, param string) {

}

func parseListParam(param string) (path string) {
	if len(param) == 0 {
		path = param
	} else {
		fields := strings.Fields(param)
		i := 0
		for _, field := range fields {
			if !strings.HasPrefix(field, "-") {
				break
			}
			i = strings.LastIndex(param, " "+field) + len(field) + 1
		}
		path = strings.TrimLeft(param[i:], " ") //Get all the path even with space inside
	}
	return path
}

// commandNlst responds to the NLST FTP command. It allows the client to
// retreive a list of filenames in the current directory.
type commandNlst struct{}

func (cmd commandNlst) IsExtend() bool {
	return false
}

func (cmd commandNlst) RequireParam() bool {
	return false
}

func (cmd commandNlst) RequireAuth() bool {
	return true
}

func (cmd commandNlst) Execute(conn *Conn, param string) {
	path := conn.buildPath(parseListParam(param))
	info, err := conn.vfs.Stat(path)
	if err != nil {
		conn.writeMessage(550, err.Error())
		return
	}
	if !info.IsDir() {
		conn.writeMessage(550, param+" is not a directory")
		return
	}

	var files []FileInfo
	err = conn.vfs.ListDir(path, func(f FileInfo) error {
		files = append(files, f)
		return nil
	})
	if err != nil {
		conn.writeMessage(550, err.Error())
		return
	}
	conn.writeMessage(150, "Opening ASCII mode data connection for file list")
	conn.sendOutofbandData(listFormatter(files).Short())
}

// commandMdtm responds to the MDTM FTP command. It allows the client to
// retreive the last modified time of a file.
type commandMdtm struct{}

func (cmd commandMdtm) IsExtend() bool {
	return false
}

func (cmd commandMdtm) RequireParam() bool {
	return true
}

func (cmd commandMdtm) RequireAuth() bool {
	return true
}

func (cmd commandMdtm) Execute(conn *Conn, param string) {
	path := conn.buildPath(param)
	stat, err := conn.vfs.Stat(path)
	if err == nil {
		conn.writeMessage(213, stat.ModTime().Format("20060102150405"))
	} else {
		conn.writeMessage(450, "File not available")
	}
}

/*
// commandMkd responds to the MKD FTP command. It allows the client to create
// a new directory
type commandMkd struct{}

func (cmd commandMkd) IsExtend() bool {
	return false
}

func (cmd commandMkd) RequireParam() bool {
	return true
}

func (cmd commandMkd) RequireAuth() bool {
	return true
}

func (cmd commandMkd) Execute(conn *Conn, param string) {
	path := conn.buildPath(param)

	err := conn.storage.MakeDir(path)

	if err == nil {
		conn.writeMessage(257, "Directory created")
	} else {
		conn.writeMessage(550, fmt.Sprint("Action not taken: ", err))
	}
}
*/

// commandPort responds to the PORT FTP command.
//
// The client has opened a listening socket for sending out of band data and
// is requesting that we connect to it
type commandPort struct{}

func (cmd commandPort) IsExtend() bool {
	return false
}

func (cmd commandPort) RequireParam() bool {
	return true
}

func (cmd commandPort) RequireAuth() bool {
	return true
}

func (cmd commandPort) Execute(conn *Conn, param string) {
	nums := strings.Split(param, ",")
	portOne, _ := strconv.Atoi(nums[4])
	portTwo, _ := strconv.Atoi(nums[5])
	port := (portOne * 256) + portTwo
	host := nums[0] + "." + nums[1] + "." + nums[2] + "." + nums[3]
	socket, err := newActiveSocket(host, port, conn.session)
	if err != nil {
		conn.writeMessage(425, "Data connection failed")
		return
	}
	conn.dataConn = socket
	conn.writeMessage(200, "Connection established ("+strconv.Itoa(port)+")")
}

type commandRest struct{}

func (cmd commandRest) IsExtend() bool {
	return false
}

func (cmd commandRest) RequireParam() bool {
	return true
}

func (cmd commandRest) RequireAuth() bool {
	return true
}

func (cmd commandRest) Execute(conn *Conn, param string) {
	conn.writeMessage(350, fmt.Sprint("Start transfer from 0"))
}

type commandCcc struct{}

func (cmd commandCcc) IsExtend() bool {
	return false
}

func (cmd commandCcc) RequireParam() bool {
	return true
}

func (cmd commandCcc) RequireAuth() bool {
	return true
}

func (cmd commandCcc) Execute(conn *Conn, param string) {
	conn.writeMessage(550, "Action not taken")
}

type commandEnc struct{}

func (cmd commandEnc) IsExtend() bool {
	return false
}

func (cmd commandEnc) RequireParam() bool {
	return true
}

func (cmd commandEnc) RequireAuth() bool {
	return true
}

func (cmd commandEnc) Execute(conn *Conn, param string) {
	conn.writeMessage(550, "Action not taken")
}

type commandMic struct{}

func (cmd commandMic) IsExtend() bool {
	return false
}

func (cmd commandMic) RequireParam() bool {
	return true
}

func (cmd commandMic) RequireAuth() bool {
	return true
}

func (cmd commandMic) Execute(conn *Conn, param string) {
	conn.writeMessage(550, "Action not taken")
}

type commandPbsz struct{}

func (cmd commandPbsz) IsExtend() bool {
	return false
}

func (cmd commandPbsz) RequireParam() bool {
	return true
}

func (cmd commandPbsz) RequireAuth() bool {
	return false
}

func (cmd commandPbsz) Execute(conn *Conn, param string) {
	conn.writeMessage(550, "Action not taken")
}

type commandProt struct{}

func (cmd commandProt) IsExtend() bool {
	return false
}

func (cmd commandProt) RequireParam() bool {
	return true
}

func (cmd commandProt) RequireAuth() bool {
	return false
}

func (cmd commandProt) Execute(conn *Conn, param string) {
	conn.writeMessage(550, "Action not taken")
}

type commandConf struct{}

func (cmd commandConf) IsExtend() bool {
	return false
}

func (cmd commandConf) RequireParam() bool {
	return true
}

func (cmd commandConf) RequireAuth() bool {
	return true
}

func (cmd commandConf) Execute(conn *Conn, param string) {
	conn.writeMessage(550, "Action not taken")
}

// commandSize responds to the SIZE FTP command. It returns the size of the
// requested path in bytes.
type commandSize struct{}

func (cmd commandSize) IsExtend() bool {
	return false
}

func (cmd commandSize) RequireParam() bool {
	return true
}

func (cmd commandSize) RequireAuth() bool {
	return true
}

func (cmd commandSize) Execute(conn *Conn, param string) {
	path := conn.buildPath(param)
	stat, err := conn.vfs.Stat(path)
	if err != nil {
		utils.PrintDebug(conn.session, "Size: error(%s)\n", err)
		conn.writeMessage(450, fmt.Sprint("path", path, "not found"))
	} else {
		conn.writeMessage(213, strconv.Itoa(int(stat.Size())))
	}
}

/*
// commandStor responds to the STOR FTP command. It allows the user to upload a
// new file.
type commandStor struct{}

func (cmd commandStor) IsExtend() bool {
	return false
}

func (cmd commandStor) RequireParam() bool {
	return true
}

func (cmd commandStor) RequireAuth() bool {
	return true
}

func (cmd commandStor) Execute(conn *Conn, param string) {
	targetPath := conn.buildPath(param)
	conn.writeMessage(150, "Data transfer starting")

	defer func() {
		conn.appendData = false
	}()

	conn.server.notifiers.BeforePutFile(conn, targetPath)
	size, err := conn.storage.PutFile(targetPath, conn.dataConn, conn.appendData)
	conn.server.notifiers.AfterFilePut(conn, targetPath, size, err)
	if err == nil {
		msg := fmt.Sprintf("OK, received %d bytes", size)
		conn.writeMessage(226, msg)
	} else {
		conn.writeMessage(450, fmt.Sprint("error during transfer: ", err))
	}
}
*/

// commandStru responds to the STRU FTP command.
//
// like the MODE and TYPE commands, stru[cture] dates back to a time when the
// FTP protocol was more aware of the content of the files it was transferring,
// and would sometimes be expected to translate things like EOL markers on the
// fly.
//
// These days files are sent unmodified, and F(ile) mode is the only one we
// really need to support.
type commandStru struct{}

func (cmd commandStru) IsExtend() bool {
	return false
}

func (cmd commandStru) RequireParam() bool {
	return true
}

func (cmd commandStru) RequireAuth() bool {
	return true
}

func (cmd commandStru) Execute(conn *Conn, param string) {
	if strings.ToUpper(param) == "F" {
		conn.writeMessage(200, "OK")
	} else {
		conn.writeMessage(504, "STRU is an obsolete command")
	}
}
