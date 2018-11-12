package main

import (
	"github.com/caarlos0/env"
	"log"
	"net/http"
	"os"
)

const VERSION = "0.2.0"

type Params struct {
	Proxy         string `env:"NOTIFIER_PROXY"  envDefault:"socks5://grinat:1122334455@95.216.138.74:1080"`
	Port          string `env:"NOTIFIER_PORT"   envDefault:"9191"`
	HerokuPort    string `env:"PORT"`
	HerokuAppName string `env:"HEROKU_APP_NAME"`
}

type Notifier interface {
	Notify(r *http.Request) error
	GetMessageId() string
}

type NotifyMessage struct {
	service     string
	recipient   string
	token       string
	message     string
	file       *os.File
}

func main() {
	log.Println("App version", VERSION)
	cfg := Params{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	// replace port with heroku dynamic dyno port
	if cfg.HerokuAppName != "" && cfg.HerokuPort != "" {
		log.Println("Port changed to heroku port")
		cfg.Port = cfg.HerokuPort
	}

	log.Println("Server started at port", cfg.Port)

	err = http.ListenAndServe(":" + cfg.Port, Router{cfg:cfg})
	if err != nil {
		panic(err)
	}
}
