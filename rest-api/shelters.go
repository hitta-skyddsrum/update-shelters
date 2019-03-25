package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Shelter struct {
	address       string
	city          string
	municipality  string
	position_lat  string
	position_long string
	shelter_id    string
	slots         int
}

var ShelterFields = "address, city, municipality, position_lat, position_long, shelter_id, slots"

func getShelterFromRow(row interface{ Scan(...interface{}) error }) (Shelter, error) {
	s := Shelter{}
	err := row.Scan(&s.address, &s.city, &s.municipality, &s.position_lat, &s.position_long, &s.shelter_id, &s.slots)

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

func (a *App) GetSheltersByBBox(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params == nil {
		return fmt.Errorf("bbox-parametern är obligatorisk")
	}

	bbox := strings.Split(params["bbox"], ",")
	if len(bbox) != 4 {
		return fmt.Errorf("bbox-parametern har för få värden")
	}

	rows, err := a.DB.Query("SELECT "+ShelterFields+" FROM `shelters` WHERE position_long > `?` AND position_lat > `?` AND position_long < `?` AND position_lat < `?`", bbox[0], bbox[1], bbox[2], bbox[3])
	if err != nil {
		return fmt.Errorf("Selecting shelters failed: %s\n", err)
	}
	defer rows.Close()

	shelters, err := getSheltersFromRows(rows)
	if err != nil {
		return err
	}

	a.RespondWithJson(w, http.StatusOK, shelters)

	return nil
}

func (a *App) GetShelterById(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	if params == nil {
		return fmt.Errorf("id krävs som parameter.")
	}

	row := a.DB.QueryRow("SELECT "+ShelterFields+" FROM `shelters` WHERE id = ?", params["id"])

	shelter, err := getShelterFromRow(row)
	if err == sql.ErrNoRows {
		a.RespondWithJson(w, http.StatusNotFound, nil)
		return nil
	}
	if err != nil {
		return err
	}

	a.RespondWithJson(w, http.StatusOK, shelter)

	return nil
}
