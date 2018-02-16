package main

import (
  "io/ioutil"
)

func StoreShapefile(data []byte) string {
  storagePath := "/go/src/update-shelters/skyddsrum.zip"
  err := ioutil.WriteFile(storagePath, data, 0644)
  if err != nil {
    panic(err)
  }

  return storagePath
}