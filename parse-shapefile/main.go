package main

import (
  "github.com/jonas-p/go-shp"
  "fmt"
  "flag"
  "os"
  "encoding/json"
  "io/ioutil"
  "bufio"
  "strings"
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

func MapShapeFields(zipShape *shp.ZipReader) map[string]interface{} {
  mappings := map[string]interface{}{
  }
  fields := zipShape.Fields()

  fmt.Print("Map fields from shapefile to desired name:")
  fmt.Println()

  for _, f := range fields {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("%s [%s]:", f.String(), f.String())
    value, _ := reader.ReadString('\n')
    value = strings.TrimSuffix(value, "\n")

    if len(value) > 1 {
      mappings[f.String()] = value
    } else {
      mappings[f.String()] = f.String()
    }

  }

  return mappings
}

type Geometry struct {
  Type        string    `json:"type"`
  Coordinates [2]float64  `json:"coordinates"`
}

func ShapeToJson(zipShape *shp.ZipReader, mappings map[string]interface{}) (int, []byte) {
  fields := zipShape.Fields()

  shapes := make([]interface{}, 0)
  nrShapes := 0

  for zipShape.Next() {
    _, shape := zipShape.Shape()
    jo := map[string]interface{}{
    }

    if shape.BBox().MinX != shape.BBox().MaxX {
      panic("fail")
    }

    if shape.BBox().MinY != shape.BBox().MaxY {
      panic("fail")
    }

    jo["geometry"] = Geometry{
      Type: "Point",
      Coordinates: [...]float64{shape.BBox().MinX, shape.BBox().MinY},
    }

    for k, f := range fields {
      val := zipShape.Attribute(k)
      name := mappings[f.String()].(string)
      jo[name] = val
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

  zipShape, err := shp.OpenZip(shapefile)
  if err != nil {
    panic(err)
  }

  defer zipShape.Close()

  if *listFields == true {
    ListFields(zipShape)
    os.Exit(0)
  }

  if *showExample == true {
    ShowExample(zipShape)
    os.Exit(0)
  }

  mappings := MapShapeFields(zipShape)

  fmt.Print("Will use the following mappings: ", mappings)
  fmt.Println()

  nrShelters, shelters := ShapeToJson(zipShape, mappings)

  StoreJson(shelters)

  fmt.Printf("Successfully wrote %d shapes to JSON", nrShelters)
  fmt.Println()
}
