//go:build desktop
// +build desktop

package main

import (
	"context"
	"embed"
	"flag"
	"fmt"

	"io/ioutil"
	"log"
	"time"

	"github.com/jbowes/whatsnew"
	"github.com/vikingo-project/vsat/api"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/logger"
	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/shared"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
)

//go:embed frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// NewApp creates a new App application struct
func NewApp() *App {
	ctx := context.WithValue(context.Background(), "start_time", time.Now())
	return &App{
		ctx: ctx,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogTrace(a.ctx, "Startup...")
	go func() {
		// use desktop notification stream instead of ws from server package
		//for data := range shared.WSMessagesChan {
		// check broadcast/dst
		// log.Println("got ws message", data)
		//runtime.EventsEmit(a.ctx, data.Event, data)
		// WSServer.BroadcastToRoom("/", data.Room, data.Event, data)
		//}
	}()
}

func main() {
	shared.DesktopMode = "true"
	fileLogger := logger.Initialize(false)
	fileLogger.Info("Starting")
	db.Init()
	manager.Start()
	defer db.Close()

	go func() {
		// check new version
		for {
			ctx := context.Background()
			fut := whatsnew.Check(ctx, &whatsnew.Options{
				Slug:    "vikingo-project/vsat",
				Version: shared.Version,
			})
			if v, _ := fut.Get(); v != "" {
				fmt.Printf("new version available: %s\n", v)
			}
			time.Sleep(30 * time.Minute)
		}
	}()

	app := NewApp()
	log.Println("app", app)
	err := wails.Run(&options.App{
		Title:     "Vikingo Satellite",
		Width:     1024,
		Height:    800,
		Assets:    assets,
		OnStartup: app.startup,
		// Logger:    fileLogger,
		LogLevel: 1, //logger.INFO,
		Bind: []interface{}{
			app, &api.Instance,
		},
	})

	if err != nil {
		println("Error:", err)
	}
}

func init() {
	_db := flag.String("db", "storage.db", "Path to database. If the parameter is not set storage.db will be used")
	_config := flag.String("c", "", "Path to config.")
	flag.Parse()

	if *_config != "" {
		data, err := ioutil.ReadFile(*_config)
		if err != nil {
			log.Printf("Config file not found. Use default configuration.")
		} else {
			err = yaml.Unmarshal(data, &shared.Config)
			if err != nil {
				log.Fatalf("Failed to parse config file; %s", err)
			}
		}
	}

	shared.Config.DB = *_db
	shared.DesktopMode = "true"
}
