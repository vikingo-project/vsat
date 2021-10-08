package ctrl

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/shared"
)

var wsServer *socketio.Server

func startWS() {
	wsServer = socketio.NewServer(nil)
	wsServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})
	wsServer.OnDisconnect("/", func(s socketio.Conn, _ string) {
		// get out here
		wsServer.LeaveRoom("/", "users", s)
	})

	wsServer.OnEvent("/", "auth", func(s socketio.Conn, token string) {
		var auth models.Auth
		db.GetConnection().Model(&auth).First(&auth)
		if auth.Token == token {
			wsServer.JoinRoom("/", "users", s)
		}

	})

	go func() {
		for msg := range shared.Notifications {
			wsServer.BroadcastToRoom("/", "users", "notifications", msg)
		}
	}()
}
