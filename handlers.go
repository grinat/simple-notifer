package main

import (
	"errors"
	"net/http"
)

func HandleHelp(w http.ResponseWriter, cfg Params)  {
	host := "localhost:" + cfg.Port + "/send-message"
	if cfg.HerokuAppName != "" {
		host = "https://" + cfg.HerokuAppName + ".herokuapp.com/send-message"
	}
	helpMsg := `
Supported services:
telegram

For send notify exec:
curl -i -X POST -F "service=telegram" -F "message=Hello" -F "recipient=your_chat_id" -F "token=your_token" `
	helpMsg += host
	helpMsg += `
    
For send file add:
-F "file=@path/to/file"

Or send get response(in what mode file not supported):
`
	helpMsg += host + "?service=telegram&message=Hello&recipient=your_chat_id&token=your_token"
	helpMsg += `

Return error message with status code 500 or notify id with status code 201
`
	w.Write([]byte(helpMsg))
}

func HandleSendMessage(w http.ResponseWriter, r *http.Request, cfg Params)  {
	msg := NotifyMessage{}
	switch r.Method {
	case "POST":
		err := MultipartToNotifyMessage(r, &msg)
		if err != nil {
			onError(w, err)
			return
		}
	case "GET":
		err := UrlQueryToNotifyMessage(r, &msg)
		if err != nil {
			onError(w, err)
			return
		}
	}

	var notifier Notifier

	// create sender instance
	sender := Sender{}

	// set proxy if exist
	if cfg.Proxy != "" {
		sender.Proxy = cfg.Proxy
	}

	switch r.FormValue("service") {
	case "telegram":
		notifier = &Telegram{
			cfg: cfg,
			sender: sender,
			msg: msg,
		}
	default:
		onError(w, errors.New("Unknown service"))
		return
	}

	err := notifier.Notify(r)
	if err != nil {
		onError(w, err)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(notifier.GetMessageId()))
}


func onError (w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

