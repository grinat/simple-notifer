package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Telegram struct {
	cfg Params
	sender Sender
	lastMessageId string
}

type TelegramApiAnswer struct {
	Ok bool `json:"ok"`
	Result TelegramApiResult `json:"result"`
}

type TelegramApiResult struct {
	MessageId int `json:"message_id"`
}

func (t Telegram) Notify(r *http.Request) error {
	// grab token
	result := TelegramApiAnswer{
		Result: TelegramApiResult{},
	}
	token := r.FormValue("token")
	if token == "" {
		return errors.New("token field required")
	}

	// change method to post
	t.sender.Method = "POST"

	// send file
	file, header, err := r.FormFile("file")
	if file != nil {
		if err != nil {
			return err
		}
		defer file.Close()
		tmpFile, err := ioutil.TempFile(os.TempDir(), header.Filename)
		if err != nil {
			return err
		}
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())
		io.Copy(tmpFile, file)
		osFile, err := os.Open(tmpFile.Name())
		if err != nil {
			return err
		}
		err, res := t.sender.SendMultipart(
			"https://api.telegram.org/bot" + token + "/sendDocument?chat_id=" + r.FormValue("chat_id"),
			map[string]io.Reader{
				"document": osFile,
			})
		if err != nil {
			return err
		}
		logBody(res.Body)
		err = decodeBody(res.Body, &result)
		if err != nil {
			return err
		}
		t.lastMessageId = string(result.Result.MessageId)
	}

	// send text
	err, res := t.sender.SendMultipart(
		"https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + r.FormValue("chat_id"),
		map[string]io.Reader{
			"text": strings.NewReader(r.FormValue("message")),
		})
	if err != nil {
		return err
	}
	logBody(res.Body)
	err = decodeBody(res.Body, &result)
	if err != nil {
		return err
	}
	t.lastMessageId = string(result.Result.MessageId)

	return nil
}

func (t Telegram) GetMessageId() string {
	return t.lastMessageId
}

func logBody(body io.ReadCloser)  {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))
}

func decodeBody(_ io.ReadCloser, _ *TelegramApiAnswer) error {
	return nil
}
