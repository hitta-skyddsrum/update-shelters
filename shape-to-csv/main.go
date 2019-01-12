// +build !lambda

package main

import (
  "flag"
  "path/filepath"
)

func main() {
  flag.Parse()

  filePath := flag.Args()[0]

  ExportShapeToCSV(filePath, filepath.Base(filePath))
}

