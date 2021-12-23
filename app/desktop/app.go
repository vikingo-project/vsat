package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asticode/go-astikit"
	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/jbowes/whatsnew"

	"github.com/vikingo-project/vsat/ctrl"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/shared"
	"github.com/vikingo-project/vsat/utils"
	"gopkg.in/yaml.v2"
)

var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

var (
	fs = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	w  *astilectron.Window
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	log.Println("handle", m.Name)
	return
}

const htmlAbout = `Welcome on <b>Vikingo Satellite</b>`

func main() {
	db.Init()
	initAuth()

	manager := manager.NewManager()
	manager.Start()

	// handle signals ctrl+c, etc
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		// database should be closed validly
		db.Close()
		os.Exit(0)
	}()

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

	ctrlServer := ctrl.NewCtrlServer()
	go ctrlServer.Run(manager)

	if err := bootstrap.Run(bootstrap.Options{
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug: true,

		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								log.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
								return
							}
							log.Printf("About modal has been displayed and payload is %s!\n", s)
						}); err != nil {
							log.Println(fmt.Errorf("sending about event failed: %w", err))
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
					log.Println(fmt.Errorf("sending check.out.menu event failed: %w", err))
				}
			}()
			return nil
		},
		Windows: []*bootstrap.Window{{
			Homepage:       "http://127.0.0.1:1025",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#f0f4f5"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(800),
				Width:           astikit.IntPtr(800),
				MinWidth:        astikit.IntPtr(400),
				MinHeight:       astikit.IntPtr(400),
			},
		}},
	}); err != nil {
		log.Fatal(fmt.Errorf("error: %w", err))
	}

}

func init() {
	_db := flag.String("db", "storage.db", "Path to database. If the parameter is not set storage.db will be used")
	_config := flag.String("c", "", "Path to config.")
	_token := flag.String("token", "", "New access token")
	_ctrl_listen_addr := flag.String("l", "0.0.0.0:1025", "Control interface listen address")
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
	shared.Config.Token = *_token
	shared.Config.Listen = *_ctrl_listen_addr
	shared.DesktopMode = "true"
}

func initAuth() {
	if shared.Config.Token == "" {
		var auth models.Auth
		db.GetConnection().Model(&auth).First(&auth)
		// if the token is not in DB generate a new one
		if auth.Token == "" {
			token := utils.EasyHash(true)
			db.GetConnection().Save(&models.Auth{Token: token})
			shared.Config.Token = token
		} else {
			shared.Config.Token = auth.Token
		}
	}

	fmt.Printf(`
**********************************************************
***** Access token is: %s ******
**********************************************************
`, shared.Config.Token)
}
