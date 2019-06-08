package main

import (
	"io/ioutil"
)

func StoreShapefile(name string, data []byte) string {
	storagePath := "./" + name + ".zip"
	err := ioutil.WriteFile(storagePath, data, 0644)
	if err != nil {
		panic(err)
	}

	return storagePath
}
