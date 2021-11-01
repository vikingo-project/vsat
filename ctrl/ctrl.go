package ctrl

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/shared"
	"github.com/vikingo-project/vsat/utils"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}
func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: root, Fallback: "index.html"}
	return &binaryFileSystem{
		fs,
	}
}

type Ctrl struct {
}

func NewCtrlServer() *Ctrl {
	return &Ctrl{}
}

func (c *Ctrl) Run(mgr *manager.Manager) error {
	gin.SetMode(gin.ReleaseMode)
	_, certErr := os.Stat("./vsat.crt")
	_, keyErr := os.Stat("./vsat.key")
	if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
		utils.PrintDebug("SSL cert/key not found, start generaing a new pair")
		cert, key, err := utils.GenerateCertAndKey(utils.EasyHash(false) + "vikingo.satellite")
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile("./vsat.crt", []byte(cert), 0666)
		ioutil.WriteFile("./vsat.key", []byte(key), 0666)
	}

	// start socketio server
	startWS()
	go func() {
		if err := wsServer.Serve(); err != nil {
			log.Printf("socketio listen error: %s\n", err)
		}
	}()
	defer wsServer.Close()

	router := gin.New()
	router.Use(prepare(mgr))
	router.GET("/socket.io/*any", gin.WrapH(wsServer))
	router.POST("/socket.io/*any", gin.WrapH(wsServer))
	router.GET("/api/files/download/:hash/", hDownloadFile) // file download should work without auth ;)

	authorized := router.Group("/")
	authorized.Use(auth())
	{
		api := authorized.Group("/api")
		{
			api.GET("/sql/", sql)
			api.GET("/about/", about)
			api.GET("/ping/", ping)
			api.GET("/networks/", networks)

			// service handlers
			api.GET("/services/", services)
			api.POST("/services/create/", httpCreateService)
			api.POST("/services/update/", httpUpdateService)
			api.POST("/services/remove/", httpRemoveService)
			api.POST("/services/start/", httpStartService)
			api.POST("/services/stop/", httpStopService)

			// tunnels handlers
			api.GET("/tunnels/", httpTunnels)
			api.POST("/tunnels/create/", httpCreateTunnel)
			api.POST("/tunnels/update/", httpUpdateTunnel)
			api.POST("/tunnels/remove/", httpRemoveTunnel)
			api.POST("/tunnels/start/", httpStartTunnel)
			api.POST("/tunnels/stop/", httpStopTunnel)

			// sessions and events
			api.GET("/sessions/", hSessions)
			api.GET("/sessions/view/:hash/", hEvents)
			api.POST("/sessions/remove/", hRemoveSession)

			// modules
			api.GET("/modules/", hModules)

			// Files and certs
			api.GET("/files/", hFiles)
			api.POST("/files/remove/", hRemoveFile)
			api.POST("/files/upload/", hUploadFiles)
			api.GET("/files/types/", hTypes)
			api.GET("/certs/", hCerts)
		}
	}

	// serve static files
	router.Use(static.Serve("/", BinaryFileSystem("dist")))

	// dev mode enables extended logging; ctrl server uses HTTP instead of HTTPS
	if !utils.IsDevMode() {
		return router.RunTLS(shared.Config.Listen, "./vsat.crt", "./vsat.key")
	}
	return router.Run(shared.Config.Listen)

}

func prepare(mgr *manager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("mgr", mgr)
	}
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// DO NOT CHECK AUTH IF DEBUG ENABLED AND CLIENT FROM LOCALHOST
		if utils.IsDevMode() && c.ClientIP() == "::1" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		splitToken := strings.Split(token, "Bearer")
		if len(splitToken) == 2 {
			tokenFromUser := strings.TrimSpace(splitToken[1])
			if shared.Config.Token == tokenFromUser {
				// auth ok
				c.Next()
				return
			}
		}
		c.JSON(200, gin.H{"error": "auth required"})
		c.Abort()
	}
}

func getManager(c *gin.Context) *manager.Manager {
	m, _ := c.Get("mgr")
	return m.(*manager.Manager)
}
