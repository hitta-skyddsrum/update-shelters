package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestGetSheltersByBBoxSuccess(t *testing.T) {
	x1, x2, y1, y2 := "11", "22", "33", "44"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/shelters?bbox=%s,%s,%s,%s", x1, y1, x2, y2), nil)
	res := httptest.NewRecorder()
	db, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Create db mock failed: %s", err)
	}
	defer db.Close()

	shelter := Shelter{
		Address:      "Bästa svängen 1",
		City:         "Hökarängen",
		Municipality: "Stockholm",
		Position: Position{
			Lat:  "16.505",
			Long: "59.1234",
		},
		Shelter_id: "1",
		Slots:      8,
	}
	rows := sqlmock.NewRows([]string{"address", "city", "estate_id", "filter_type", "municipality", "position_lat", "position_long", "shelter_id", "slots"}).
		AddRow(shelter.Address, shelter.City, shelter.Estate_id, shelter.Filter_type, shelter.Municipality, shelter.Position.Lat, shelter.Position.Long, shelter.Shelter_id, shelter.Slots)

	sm.ExpectQuery("SELECT "+ShelterFields+" FROM `shelters`").WithArgs(x1, y1, x2, y2).WillReturnRows(rows)

	a := App{
		DB: db,
	}
	shelters := make([]Shelter, 0)
	shelters = append(shelters, shelter)
	a.GetSheltersByBBox(res, req)

	if err := sm.ExpectationsWereMet(); err != nil {
		t.Errorf("Tests failed: %s", err)
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected to retrieve status OK. Got: %d", res.Code)
	}

	var jsonBody []Shelter
	json.NewDecoder(res.Result().Body).Decode(&jsonBody)

	if jsonBody[0].Address != shelter.Address {
		t.Errorf("Expected to retrieve shelter address `%v`. Got: %v", shelter.Address, jsonBody[0].Address)
	}
}

func TestGetSheltersByBBoxFailureWithoutParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "/shelters", nil)
	res := httptest.NewRecorder()

	a := App{}

	a.GetSheltersByBBox(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected to retrieve status 422. Got: %d", res.Code)
	}

	var jsonBody ErrorResponse
	json.NewDecoder(res.Result().Body).Decode(&jsonBody)

	if jsonBody.Message != "bbox-parametern är obligatorisk" {
		t.Errorf("Tests failed: Expected an error with accurate message, got: %s", jsonBody.Message)
	}
}

func TestGetSheltersByBBoxFailureWithBadBBoxParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "/shelters", nil)
	res := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"bbox": "1"})

	defer func() {
		if r := recover(); r != nil {
			if m, _ := regexp.MatchString("för få värden", fmt.Sprintf("%s", r)); m == false {
				t.Errorf("Tests failed: Expected an error with accurate message, got: %s", r)
			}
		}
	}()

	a := App{}
	a.GetSheltersByBBox(res, req)

}

func TestGetSingleShelterSuccess(t *testing.T) {
	id := "3223"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/shelters/%s", id), nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})
	res := httptest.NewRecorder()
	db, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Create db mock failed: %s", err)
	}
	defer db.Close()

	shelter := Shelter{
		Address:      "Sätra solfångare ",
		City:         "Skärholmen",
		Municipality: "Stockholm",
		Position: Position{
			Lat:  "16.205",
			Long: "59.4234",
		},
		Shelter_id: id,
		Slots:      91,
	}
	rows := sqlmock.NewRows([]string{"address", "city", "estate_id", "filter_type", "municipality", "position_lat", "position_long", "shelter_id", "slots"}).
		AddRow(shelter.Address, shelter.City, shelter.Estate_id, shelter.Filter_type, shelter.Municipality, shelter.Position.Lat, shelter.Position.Long, shelter.Shelter_id, shelter.Slots)

	sm.ExpectQuery("SELECT " + ShelterFields + " FROM `shelters` WHERE id = ?").WithArgs(id).WillReturnRows(rows)

	a := App{
		DB: db,
	}
	shelters := make([]Shelter, 0)
	shelters = append(shelters, shelter)
	a.GetShelterById(res, req)

	if err := sm.ExpectationsWereMet(); err != nil {
		t.Errorf("Tests failed: %s", err)
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected to retrieve status OK. Got: %d", res.Code)
	}
}

func TestGetSingleShelterNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/shelters/unknown", nil)
	id := "144"
	req = mux.SetURLVars(req, map[string]string{"id": id})
	res := httptest.NewRecorder()
	db, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Create db mock failed: %s", err)
	}
	defer db.Close()

	sm.ExpectQuery("SELECT " + ShelterFields + " FROM `shelters` WHERE id = ?").WithArgs(id).WillReturnError(sql.ErrNoRows)

	a := App{
		DB: db,
	}

	a.GetShelterById(res, req)

	if res.Code != 404 {
		t.Errorf("Expected http status to be 404, got %d", res.Code)
	}
}

func TestGetSingleShelterWithoutId(t *testing.T) {
	req, _ := http.NewRequest("GET", "/shelters/1", nil)
	res := httptest.NewRecorder()

	a := App{}

	a.GetShelterById(res, req)

	var jsonBody ErrorResponse
	json.NewDecoder(res.Result().Body).Decode(&jsonBody)

	if m, _ := regexp.MatchString("id krävs", fmt.Sprintf("%s", jsonBody.Message)); m == false {
		t.Errorf("Tests failed: Expected an error with accurate message, got: %s", jsonBody.Message)
	}
}
