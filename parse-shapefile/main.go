package main

import (
  "github.com/jonas-p/go-shp"
  "fmt"
  "flag"
  "os"
)

func main() {
  listFields := flag.Bool("list-fields", false, "List all fields in shapefile")
  showExample := flag.Bool("show-example", false, "Show an example shape from the shapefile")
  flag.Parse()

  shapefile := flag.Args()[0]

  shape, err := shp.OpenZip(shapefile)
  if err != nil {
    panic(err)
  }

  defer shape.Close()

  fields := shape.Fields()

  if *listFields == true {
    fmt.Printf("Fields in %s", shapefile)
    fmt.Println()

    for k := range fields {
      fmt.Print(fields[k])
      fmt.Println()
    }

    os.Exit(0)
  }

  if *showExample == true {
    shape.Next()

    for k, f := range fields {
      val := shape.Attribute(k)
      fmt.Printf("\t%v: %v\n", f, val)
    }

    fmt.Println()

    os.Exit(0)
  }
}
