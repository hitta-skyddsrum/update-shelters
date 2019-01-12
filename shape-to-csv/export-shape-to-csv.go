package main

import (
  "encoding/csv"
  "github.com/hitta-skyddsrum/update-shelters/sweref99-to-latlon"
  "github.com/jonas-p/go-shp"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

func ExportShapeToCSV(filePath string, targetPath string) {
  zipShape := openShape(filePath)
  defer zipShape.Close()

  fields := zipShape.Fields()
  fn := []string{}

  for k := range fields {
    fn = append(fn, string(fields[k].Name[:11]))
  }
  fn = append(fn, "latitude")
  fn = append(fn, "longitude")

  data := [][]string{}
  data = append(data, fn)

  for zipShape.Next() {
    v := []string{}

    for k := range fields {
      v = append(v, zipShape.Attribute(k))
    }

    _, shape := zipShape.Shape()
    coords := coordConv.Sweref99ToLatLon([2]float64{shape.BBox().MinX, shape.BBox().MinY})

    for _, c := range coords {
      v = append(v, strconv.FormatFloat(c, 'f', -1, 64))
    }

    data = append(data, v)
  }

  csvExport(targetPath, data)
}

func openShape(filePath string) *shp.ZipReader {
  zipShape, err := shp.OpenZip(filePath)
  if err != nil {
    panic(err)
  }

  return zipShape
}

func csvExport(fileName string, data [][]string) error {
  file, err := os.Create(strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".csv")
  if err != nil {
    return err
  }

  defer file.Close()

  w := csv.NewWriter(file)

  w.WriteAll(data)

  if err := w.Error(); err != nil {
    return err
  }

  return nil
}

