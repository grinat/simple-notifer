package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UrlQueryToNotifyMessage(r *http.Request, msg *NotifyMessage) error {
	msg.service = r.URL.Query().Get("service")
	msg.recipient = r.URL.Query().Get("recipient")
	msg.message = r.URL.Query().Get("message")
	msg.token   = r.URL.Query().Get("token")
	return nil
}

func MultipartToNotifyMessage(r *http.Request, msg *NotifyMessage) error {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}
	msg.service = r.FormValue("service")
	msg.recipient = r.FormValue("recipient")
	msg.message = r.FormValue("message")
	msg.token   = r.FormValue("token")
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
		msg.file = osFile
	}
	return nil
}