package main

import (
  "github.com/jonas-p/go-shp"
  "fmt"
  "flag"
  "os"
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

func ShowExample(zipShape *shp.ZipReader) {
  fields := zipShape.Fields()
  zipShape.Next()
  _, shape := zipShape.Shape()

  for k, f := range fields {
    val := zipShape.Attribute(k)
    fmt.Printf("\t%v: %v\n", f.String(), val)
  }

  coordinates := Sweref99ToLatLon([2]float64{shape.BBox().MinX, shape.BBox().MinY})
  fmt.Printf("\tLatitude: %v\n", coordinates[0])
  fmt.Printf("\tLongitude: %v\n", coordinates[1])

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

func StoreJson(byteJson []byte) {
  e := ioutil.WriteFile("skyddsrum.geojson", byteJson, 0644)

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

  nrShelters, shelters := ShapeToGeoJson(zipShape, mappings)

  StoreJson(shelters)

  fmt.Printf("Successfully wrote %d shapes to JSON", nrShelters)
  fmt.Println()
}
