package main

import (
  "github.com/jonas-p/go-shp"
  "fmt"
  "flag"
  "os"
  "encoding/json"
  "io/ioutil"
)

func ListFields(shape *shp.ZipReader) {
  fields := shape.Fields()

  for k := range fields {
    fmt.Print(fields[k])
    fmt.Println()
  }
}

func ShowExample(shape *shp.ZipReader) {
  fields := shape.Fields()
  shape.Next()

  for k, f := range fields {
    val := shape.Attribute(k)
    fmt.Printf("\t%v: %v\n", f.String(), val)
  }

  fmt.Println()
}

func ShapeToJson(shape *shp.ZipReader) (int, []byte) {
  fields := shape.Fields()

  shapes := make([]interface{}, 0)
  nrShapes := 0

  for shape.Next() {
    jo := map[string]interface{}{
    }

    for k, f := range fields {
      jo[f.String()] = k
    }

    shapes = append(shapes, jo)
    nrShapes++
  }

  output, err := json.Marshal(shapes)

  if err != nil {
    panic(err)
  }

  return nrShapes, output
}

func StoreJson(byteJson []byte) {
  e := ioutil.WriteFile("/go/src/update-shelters/skyddsrum.json", byteJson, 0644)

  if e != nil {
    panic(e)
  }
}

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

  if *listFields == true {
    ListFields(shape)
    os.Exit(0)
  }

  if *showExample == true {
    ShowExample(shape)
    os.Exit(0)
  }

  nrShapes, shapes := ShapeToJson(shape)
  StoreJson(shapes)

  fmt.Printf("Successfully wrote %d shapes to JSON", nrShapes)
  fmt.Println()
}
