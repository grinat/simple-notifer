package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type Sender struct {
	Method string
	Proxy string
	Headers map[string] string
}

func (sender * Sender) SendMultipart(endpoint string, formData map[string]io.Reader) (err error, res *http.Response) {
	client := &http.Client{}
	if sender.Proxy != "" {
		proxyUrl, err := url.Parse(sender.Proxy)
		if err != nil {
			return err, res
		}
		client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		}
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range formData {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// add file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return err, res
			}
		} else {
			// add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return err, res
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, res
		}

	}
	w.Close()

	// create req
	req, err := http.NewRequest(sender.Method, endpoint, &b)
	if err != nil {
		return err, res
	}

	// set headers
	for key, value := range sender.Headers {
		req.Header.Set(key, value)
	}

	// set content type
	req.Header.Set("Content-Type", w.FormDataContentType())

	// submit request
	res, err = client.Do(req)
	if err != nil {
		return err, res
	}

	// check response
	if res.StatusCode != http.StatusOK {
		return err, res
	}
	return nil, res
}
