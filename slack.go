package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Slack struct {
	cfg Params
	sender Sender
	lastMessageId string
	msg NotifyMessage
}

type SlackApiAnswer struct {
	Ok     bool     `json:"ok"`
	Error  string   `json:"error"`
	Ts     string   `json:"ts"`
}

func (s *Slack) Notify(_ *http.Request) error {
	result := SlackApiAnswer{}

	// check token
	token := s.msg.token
	if token == "" {
		return errors.New("token field required")
	}

	s.sender.Method = "POST"
	s.sender.Headers = map[string]string{
		"Authorization": "Bearer " + token,
	}

	// send text
	err, res := s.sender.SendMultipart(
		"https://slack.com/api/chat.postMessage",
		map[string]io.Reader{
			"channel": strings.NewReader(s.msg.recipient),
			"text": strings.NewReader(s.msg.message),
		})
	if err != nil {
		return err
	}
	LogBody(res)

	err = s.decodeBody(res, &result)
	if err != nil {
		return err
	}

	// slack always answer with status code 200
	// check server response status
	if result.Ok == false {
		return errors.New(result.Error)
	}

	s.lastMessageId = result.Ts

	return nil
}

func (s *Slack) GetMessageId() string {
	return s.lastMessageId
}

func (s *Slack) decodeBody(res *http.Response, model *SlackApiAnswer) error {
	err := json.NewDecoder(res.Body).Decode(&model)
	if err != nil {
		return err
	}
	return nil
}