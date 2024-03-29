//go:build !full
// +build !full

package main

import (
	"embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vikingo-project/vsat/ctrl"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/shared"
	"github.com/vikingo-project/vsat/update"
	"github.com/vikingo-project/vsat/utils"
	"gopkg.in/yaml.v2"
)

//go:embed frontend/dist
var _static embed.FS

func main() {
	db.Init()
	initAuth()
	manager.Start()

	// handle signals ctrl+c, etc
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// database should be closed before exit
		db.Close()
		os.Exit(0)
	}()

	go update.CheckNewVersion()

	ctrl.Assets = _static
	ctrlServer := ctrl.NewCtrlServer()
	err := ctrlServer.Run()
	if err != nil {
		panic(err)
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
