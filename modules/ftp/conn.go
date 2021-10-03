// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vsftp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/files"
	"github.com/vikingo-project/vsat/models"
)

type Conn struct {
	conn          net.Conn
	controlReader *bufio.Reader
	controlWriter *bufio.Writer
	dataConn      DataSocket
	vfs           VirtualFS

	server    *Server
	session   string
	curDir    string
	reqUser   string
	user      string
	closed    bool
	EventsAPI *events.EventsAPI
}

// LoginUser returns the login user name if login
func (conn *Conn) LoginUser() string {
	return conn.user
}

// IsLogin returns if user has login
func (conn *Conn) IsLogin() bool {
	return len(conn.user) > 0
}

// PassivePort returns the port which could be used by passive mode.
func (conn *Conn) PassivePort() int {
	if len(conn.server.PassivePorts) > 0 {
		portRange := strings.Split(conn.server.PassivePorts, "-")

		if len(portRange) != 2 {
			log.Println("empty port")
			return 0
		}

		minPort, _ := strconv.Atoi(strings.TrimSpace(portRange[0]))
		maxPort, _ := strconv.Atoi(strings.TrimSpace(portRange[1]))

		return minPort + mrand.Intn(maxPort-minPort)
	}
	// let system automatically chose one port
	return 0
}

// Close will manually close this connection, even if the client isn't ready.
func (conn *Conn) Close() {
	conn.conn.Close()
	conn.closed = true
	conn.reqUser = ""
	conn.user = ""
	if conn.dataConn != nil {
		conn.dataConn.Close()
		conn.dataConn = nil
	}
}

func (conn *Conn) parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return strings.ToUpper(params[0]), strings.TrimSpace(params[1])
}

func (conn *Conn) receiveLine(line string) {
	defer func() {
		if e := recover(); e != nil {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "Handler crashed with error: %v", e)

			for i := 1; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				} else {
					fmt.Fprintf(&buf, "\n")
				}
				fmt.Fprintf(&buf, "%v:%v", file, line)
			}
			log.Println(conn.session, buf.String())
		}
	}()

	cmd, arg := conn.parseLine(line)
	log.Println(conn.session, cmd, arg)
	switch cmd {
	case "USER":
		conn.user = arg
		conn.reqUser = arg
		conn.writeMessage(331, "User name ok, password required")
	case "PASS":
		conn.EventsAPI.PushEvent(models.Event{Name: "auth", Data: map[string]interface{}{"user": conn.user}})
		conn.writeMessage(230, "Password ok, continue")
	case "QUIT":
		conn.writeMessage(221, "Bye")
		conn.Close()
	case "NOOP":
		conn.writeMessage(200, "OK")
	case "FEAT":
		conn.writeMessageMultiline(211, "Features:\r\n EPRT\n EPSV\n PASV\n REST\n SIZE\n TVFS")
	case "ALLO":
		conn.writeMessage(202, "Obsolete")

	case "MODE":
		if strings.ToUpper(arg) == "S" {
			conn.writeMessage(200, "OK")
		} else {
			conn.writeMessage(504, "MODE is an obsolete command")
		}

	case "PASV":
		listenIP := conn.conn.LocalAddr().(*net.TCPAddr).IP.String()
		socket, err := conn.newPassiveSocket()
		if err != nil {
			conn.writeMessage(425, "Data connection failed")
			return
		}

		p1 := socket.Port() / 256
		p2 := socket.Port() - (p1 * 256)

		quads := strings.Split(listenIP, ".")
		target := fmt.Sprintf("(%s,%s,%s,%s,%d,%d)", quads[0], quads[1], quads[2], quads[3], p1, p2)
		conn.writeMessage(227, "Entering Passive Mode "+target)

	case "RNFR":
		conn.writeMessage(350, "Requested file action pending further information.")

	case "MKD", "XMKD":
		path := conn.buildPath(arg)
		conn.vfs.MkDir(path)
		conn.writeMessage(257, "Directory created")

	case "CDUP", "XCUP":
		cmd = "CWD"
		arg = ".."
		fallthrough
	case "CWD", "XCWD":
		path := conn.buildPath(arg)
		node, err := conn.vfs.Stat(path)
		if err != nil {
			conn.writeMessage(550, fmt.Sprint("Directory change to ", path, " failed: ", err))
			return
		}
		if !node.IsDir() {
			conn.writeMessage(550, fmt.Sprint("Directory change to ", path, " is a file"))
			return
		}

		conn.curDir = path
		if err == nil {
			conn.writeMessage(250, "OK")
		} else {
			conn.writeMessage(550, fmt.Sprint("Directory change to ", path, " failed: ", err))
		}

	case "DELE", "RNTO", "RMD", "XRMD":
		path := conn.buildPath(arg)
		conn.vfs.RemoveFile(path)
		conn.writeMessage(250, "OK")
	case "PWD", "XPWD":
		conn.writeMessage(257, "\""+conn.curDir+"\" is the current directory")
	case "HELP":
		conn.EventsAPI.PushEvent(models.Event{Session: conn.session, Name: "Help", Data: map[string]interface{}{
			"command": "HELP",
		}})
		conn.writeMessage(214, "Do you realy need a help?")
	case "STAT":
		path := conn.buildPath(parseListParam(arg))
		if path == "" {
			path = conn.curDir
		}
		conn.EventsAPI.PushEvent(models.Event{Session: conn.session, Name: "Dir info", Data: map[string]interface{}{
			"command": "STAT",
			"path":    path,
		}})
		conn.writeMessage(213, "OK") // todo

	case "LIST":
		path := conn.buildPath(parseListParam(arg))
		if path == "" {
			path = conn.curDir
		}
		conn.EventsAPI.PushEvent(models.Event{Session: conn.session, Name: "List dir", Data: map[string]interface{}{
			"command": "LIST",
			"path":    path,
		}})
		info, err := conn.vfs.Stat(path)
		if err != nil {
			conn.writeMessage(550, err.Error())
			return
		}

		var files []FileInfo
		if info.IsDir() {
			err = conn.vfs.ListDir(path, func(f FileInfo) error {
				files = append(files, f)
				return nil
			})

			if err != nil {
				conn.writeMessage(550, err.Error())
				return
			}

		} else {
			files = append(files, info)
		}
		conn.writeMessage(150, "Opening ASCII mode data connection for file list")
		conn.sendOutofbandData(formatFiles(files))

	case "SYST":
		conn.writeMessage(215, "UNIX Type: L8")
	case "RETR":
		path := conn.buildPath(arg)
		if conn.vfs.Exists(path) {
			data := conn.vfs.ReadFile(path)
			conn.writeMessage(150, fmt.Sprintf("Data transfer starting %d bytes", len(data)))
			err := conn.sendOutofBandDataWriter(data)
			if err != nil {
				conn.writeMessage(551, "Error reading file")
			}
		} else {
			conn.writeMessage(551, "not found")
		}

	case "APPE", "STOR":
		targetPath := conn.buildPath(arg)
		conn.writeMessage(150, "Data transfer starting")
		buf := new(bytes.Buffer)
		buf.ReadFrom(conn.dataConn)

		file := files.PrepareFile(files.GetFilenameFromPath(targetPath), buf.Bytes())
		conn.EventsAPI.PushEvent(models.Event{Session: conn.session, Name: "Upload file", Data: map[string]interface{}{
			"path": targetPath,
			"file": file,
		}})

		conn.vfs.PutFile(targetPath, file)
		conn.writeMessage(226, "OK")

	case "TYPE":
		if strings.ToUpper(arg) == "A" {
			conn.writeMessage(200, "Type set to ASCII")
		} else if strings.ToUpper(arg) == "I" {
			conn.writeMessage(200, "Type set to binary")
		} else {
			conn.writeMessage(500, "Invalid type")
		}
	case "REST":
		conn.writeMessage(350, fmt.Sprint("Start transfer from 0"))

	default:
		conn.writeMessage(500, "Command not found")
		// conn.writeMessage(550, "Action not taken")
	}

	/*
		if cmdObj.RequireParam() && arg == "" {
			conn.writeMessage(553, "action aborted, required param missing")
		} else if cmdObj.RequireAuth() && conn.user == "" {
			conn.writeMessage(530, "not logged in")
		} else {
			cmdObj.Execute(conn, param)
		}
	*/
}

// writeMessage will send a standard FTP response back to the client.
func (conn *Conn) writeMessage(code int, message string) {
	log.Println(conn.session, code, message)
	line := fmt.Sprintf("%d %s\r\n", code, message)
	_, _ = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
}

// writeMessage will send a standard FTP response back to the client.
func (conn *Conn) writeMessageMultiline(code int, message string) {
	log.Println(conn.session, code, message)
	line := fmt.Sprintf("%d-%s\r\n%d END\r\n", code, message, code)
	_, _ = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
}

// buildPath takes a client supplied path or filename and generates a safe
// absolute path within their account sandbox.
//
//    buildpath("/")
//    => "/"
//    buildpath("one.txt")
//    => "/one.txt"
//    buildpath("/files/two.txt")
//    => "/files/two.txt"
//    buildpath("files/two.txt")
//    => "/files/two.txt"
//    buildpath("/../../../../etc/passwd")
//    => "/etc/passwd"
//
// The driver implementation is responsible for deciding how to treat this path.
// Obviously they MUST NOT just read the path off disk. The probably want to
// prefix the path with something to scope the users access to a sandbox.
func (conn *Conn) buildPath(filename string) (fullPath string) {
	if len(filename) > 0 && filename[0:1] == "/" {
		fullPath = filepath.Clean(filename)
	} else if len(filename) > 0 && filename != "-a" {
		fullPath = filepath.Clean(conn.curDir + "/" + filename)
	} else {
		fullPath = filepath.Clean(conn.curDir)
	}
	fullPath = strings.Replace(fullPath, "//", "/", -1)
	fullPath = strings.Replace(fullPath, string(filepath.Separator), "/", -1)
	return
}

// sendOutofbandData will send a string to the client via the currently open
// data socket. Assumes the socket is open and ready to be used.
func (conn *Conn) sendOutofbandData(data []byte) {
	bytes := len(data)
	if conn.dataConn != nil {
		_, _ = conn.dataConn.Write(data)
		conn.dataConn.Close()
		conn.dataConn = nil
	}
	message := "Closing data connection, sent " + strconv.Itoa(bytes) + " bytes"
	conn.writeMessage(226, message)
}

func (conn *Conn) sendOutofBandDataWriter(data []byte) error {
	bytes, err := io.Copy(conn.dataConn, bytes.NewReader(data))
	if err != nil {
		conn.dataConn.Close()
		conn.dataConn = nil
		return err
	}
	message := "Closing data connection, sent " + strconv.Itoa(int(bytes)) + " bytes"
	conn.writeMessage(226, message)
	conn.dataConn.Close()
	conn.dataConn = nil

	return nil
}
