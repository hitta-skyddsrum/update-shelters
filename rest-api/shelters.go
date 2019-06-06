package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Shelter struct {
  Address           string    `json:"address"`
  City              string    `json:"city"`
  Estate_id         string    `json:"estateId"`
  Filter_type       string    `json:"filterType"`
  Id                string    `json:"id"`
  Municipality      string    `json:"municipality"`
  Position          struct    {
    Lat             string    `json:"lat"`
    Long            string    `json:"long"`
  }                           `json:"position"`
  Shelter_id        string    `json:"shelterId"`
  Slots             int       `json:"slots"`
}

var ShelterFields = "address, city, estate_id, filter_type, municipality, position_lat, position_long, shelter_id, slots"

func getShelterFromRow(row interface{ Scan(...interface{}) (error) }) (Shelter, error) {
  s := Shelter{}
  err := row.Scan(&s.Address, &s.City, &s.Estate_id, &s.Filter_type, &s.Municipality, &s.Position.Lat, &s.Position.Long, &s.Shelter_id, &s.Slots)

	return s, err
}

func getSheltersFromRows(rows *sql.Rows) ([]Shelter, error) {
	shelters := make([]Shelter, 0)
	for rows.Next() {
		s, err := getShelterFromRow(rows)
		if err != nil {
			return nil, err
		}

		shelters = append(shelters, s)
	}

	return shelters, nil
}

func (a *App) GetSheltersByBBox (w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()

  if params.Get("bbox") == "" {
    a.RespondWithError(w, 422, "bbox-parametern är obligatorisk")
    return
	}

  bbox := strings.Split(params.Get("bbox"), ",")
  if len(bbox) != 4 {
    fmt.Printf("%d %v", len(bbox), params["bbox"])
    a.RespondWithError(w, 422, "bbox-parametern har för få värden")
    return
	}

  rows, err := a.DB.Query("SELECT " + ShelterFields + " FROM `shelters` WHERE position_long > ? AND position_lat > ? AND position_long < ? AND position_lat < ?", bbox[0], bbox[1], bbox[2], bbox[3])
	if err != nil {
    fmt.Printf("Selecting shelters failed: %s\n", err)
    a.RespondWithError(w, 500, "Någonting gick fel, försök igen")
    return
	}
	defer rows.Close()

  shelters, err := getSheltersFromRows(rows)
  if err != nil {
    panic(err)
  }

  a.RespondWithJson(w, http.StatusOK, shelters)
}

func (a *App) GetShelterById (w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  if params == nil {
    a.RespondWithError(w, 422, "id krävs som parameter.")
    return
  }

  row := a.DB.QueryRow("SELECT " + ShelterFields + " FROM `shelters` WHERE id = ?", params["id"])

  shelter, err := getShelterFromRow(row)

  if err == sql.ErrNoRows {
    a.RespondWithJson(w, http.StatusNotFound, nil)
    return
  }

  if err != nil {
    panic(err)
  }

  a.RespondWithJson(w, http.StatusOK, shelter)
}
