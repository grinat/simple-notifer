package main

import (
	"errors"
	"github.com/caarlos0/env"
	"log"
	"net/http"
)

const VERSION = "0.1.0"

type Params struct {
	Proxy       string `env:"NOTIFIER_PROXY"  envDefault:""`
	Port        string `env:"NOTIFIER_PORT"   envDefault:"9191"`
	HerokuPort  string `env:"PORT"`
	HerokuDyno  string `env:"DYNO"`
}
var cfg Params

type Notifier interface {
	Notify(r *http.Request) error
	GetMessageId() string
}

func main() {
	log.Println("App version", VERSION)
	cfg = Params{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	port := cfg.Port

	// replace port with heroku dynamic dyno port
	if cfg.HerokuDyno != "" && cfg.HerokuPort != "" {
		port = cfg.HerokuPort
	}

	log.Println("Server started at port ", port)

	http.HandleFunc("/send-message", sendMessage)

	err = http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		onError(w, err)
		return
	}

	var notifier Notifier
	sender := Sender{}

	// set proxy if exist
	if cfg.Proxy != "" {
		sender.Proxy = cfg.Proxy
	}

	switch r.FormValue("service") {
	case "telegram":
		notifier = Telegram{
			cfg: cfg,
			sender: sender,
		}
	default:
		onError(w, errors.New("Unknown service"))
		return
	}

	err = notifier.Notify(r)
	if err != nil {
		onError(w, err)
		return
	}

	w.Write([]byte(notifier.GetMessageId()))
}


func onError (w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}