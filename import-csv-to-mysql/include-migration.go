// +build never-build-only-generate

package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	f, _ := ioutil.ReadFile("./schema.sql")
	out, _ := os.Create("schema.go")
	out.Write([]byte("package main \n\nconst (\n"))
	out.Write([]byte("schema = `"))
	out.Write([]byte(strings.Replace(string(f), "`", "\"", -1)))
	out.Write([]byte("`\n"))
	out.Write([]byte(")\n"))
}
