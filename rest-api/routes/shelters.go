package main

import (
  "fmt"
  "github.com/gorilla/mux"
  "net/http"
  "strings"
)

type Shelter struct {
  address           string
  city              string
  municipality      string
  position_lat      string
  position_long     string
  shelter_id        string
  slots             int
}

func (a *App) GetSheltersByBBox (w http.ResponseWriter, r *http.Request) {
  fmt.Println("Request started")
  params := mux.Vars(r)

  if params == nil {
    panic(fmt.Errorf("bbox-parametern är obligatorisk"))
  }

  bbox := strings.Split(params["bbox"], ",")
  if len(bbox) != 4 {
    panic(fmt.Errorf("bbox-parametern har för få värden"))
  }

  rows, err := a.DB.Query("SELECT address, city, municipality, position_lat, position_long, shelter_id, slots FROM `shelters` WHERE position_long > `?` AND position_lat > `?` AND position_long < `?` AND position_lat < `?`", bbox[0], bbox[1], bbox[2], bbox[3])
  if err != nil {
    fmt.Printf("Selecting shelters failed: %s\n", err)
    panic(err)
  }
  defer rows.Close()

  shelters := make([]Shelter, 0)
  for rows.Next() {
    s := Shelter{}
    if err := rows.Scan(&s.address, &s.city, &s.municipality, &s.position_lat, &s.position_long, &s.shelter_id, &s.slots); err != nil {
      panic(err)
    }

    shelters = append(shelters, s)
  }

  fmt.Println("Finished")

  a.RespondWithJson(w, http.StatusOK, shelters)
}
