package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func LogBody(res *http.Response)  {
	requestDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Println("Cant dump response", err)
	} else {
		log.Println(string(requestDump))
	}
}
