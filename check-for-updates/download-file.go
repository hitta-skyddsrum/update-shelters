package main

import (
	"io/ioutil"
	"net/http"
)

func DownloadFile(url string) []byte {
	rs, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)
	}

	return bodyBytes
}
