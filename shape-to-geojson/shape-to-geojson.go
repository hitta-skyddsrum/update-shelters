package main

import (
  "github.com/hitta-skyddsrum/update-shelters/sweref99-to-latlon"
  "github.com/jonas-p/go-shp"
  "encoding/json"
)

type Geometry struct {
  Type        string    `json:"type"`
  Coordinates [2]float64  `json:"coordinates"`
}

type FeatureCollection struct {
  Type        string        `json:"type"`
  Features    []Feature `json:"features"`
}

type Feature struct {
  Type        string                    `json:"type"`
  Properties  map[string]interface{}    `json:"properties"`
  Geometry    Geometry                  `json:"geometry"`
}

func ShapeToGeoJson(zipShape *shp.ZipReader, mappings map[string]interface{}) (int, []byte) {
  fc := FeatureCollection{
    Type: "FeatureCollection",
    Features: make([]Feature, 0),
  }
  nrShapes := 0

  fields := zipShape.Fields()

  for zipShape.Next() {
    _, shape := zipShape.Shape()
    ft := Feature{
      Type: "Feature",
      Properties: map[string]interface{}{},
    }

    if shape.BBox().MinX != shape.BBox().MaxX {
      panic("fail")
    }

    if shape.BBox().MinY != shape.BBox().MaxY {
      panic("fail")
    }

    ft.Geometry = Geometry{
      Type: "Point",
      Coordinates: coordConv.Sweref99ToLatLon([2]float64{shape.BBox().MinX, shape.BBox().MinY}),
    }

    for k, f := range fields {
      val := zipShape.Attribute(k)
      name := mappings[f.String()].(string)
      ft.Properties[name] = val
    }

    fc.Features = append(fc.Features, ft)
    nrShapes++
  }

  output, err := json.Marshal(fc)

  if err != nil {
    panic(err)
  }

  return nrShapes, output
}

