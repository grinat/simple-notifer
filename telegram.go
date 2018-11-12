package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

type Telegram struct {
	cfg Params
	sender Sender
	lastMessageId string
	msg NotifyMessage
}

type TelegramApiAnswer struct {
	Ok bool `json:"ok"`
	Result TelegramApiResult `json:"result"`
}

type TelegramApiResult struct {
	MessageId int `json:"message_id"`
}

func (t *Telegram) Notify(_ *http.Request) error {
	result := TelegramApiAnswer{
		Result: TelegramApiResult{},
	}

	// check token
	token := t.msg.token
	if token == "" {
		return errors.New("token field required")
	}

	// change method to post
	t.sender.Method = "POST"

	// send file if exist
	if t.msg.file != nil {
		err, res := t.sender.SendMultipart(
			"https://api.telegram.org/bot" + token + "/sendDocument?chat_id=" + t.msg.recipient,
			map[string]io.Reader{
				"document": t.msg.file,
			})
		if err != nil {
			return err
		}
		logBody(res)
		err = decodeBody(res, &result)
		if err != nil {
			return err
		}
		t.lastMessageId = strconv.Itoa(result.Result.MessageId)
	}

	// send text
	err, res := t.sender.SendMultipart(
		"https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + t.msg.recipient,
		map[string]io.Reader{
			"text": strings.NewReader(t.msg.message),
		})
	if err != nil {
		return err
	}
	logBody(res)
	err = decodeBody(res, &result)
	if err != nil {
		return err
	}
	t.lastMessageId = strconv.Itoa(result.Result.MessageId)
	return nil
}

func (t *Telegram) GetMessageId() string {
	return t.lastMessageId
}

func logBody(res *http.Response)  {
	requestDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		fmt.Println("Cant dump response", err)
	} else {
		fmt.Println(string(requestDump))
	}
}

func decodeBody(res *http.Response, model *TelegramApiAnswer) error {
	err := json.NewDecoder(res.Body).Decode(&model)
	if err != nil {
		return err
	}
	return nil
}
