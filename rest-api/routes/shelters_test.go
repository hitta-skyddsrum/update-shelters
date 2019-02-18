package main

import (
  "database/sql"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/stretchr/testify/mock"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "net/http"
  "net/http/httptest"
  "regexp"
  "testing"
)

type App struct {
  mock.Mock
  Router              *mux.Router
  DB                  *sql.DB
}

func (a *App) RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
  a.Called(w, code, payload)
  a.MethodCalled("RespondWithJson", w, code, payload)
}

func TestGetSheltersByBBoxSuccess (t *testing.T) {
  x1, x2, y1, y2 := "11", "22", "33", "44"
  req, _ := http.NewRequest("GET", "/shelters", nil)
  req = mux.SetURLVars(req, map[string]string{"bbox": fmt.Sprintf("%s,%s,%s,%s", x1, y1, x2, y2) })
  res := httptest.NewRecorder()
  db, sm, err := sqlmock.New()
  if err != nil {
    t.Fatalf("Create db mock failed: %s", err)
  }
  defer db.Close()

  shelter := Shelter{
    address: "Bästa svängen 1",
    city: "Hökarängen",
    municipality: "Stockholm",
    position_lat: "16.505",
    position_long: "59.1234",
    shelter_id: "1",
    slots: 8,
  }
  rows := sqlmock.NewRows([]string{"address", "city", "municipality", "position_lat", "position_long", "shelter_id", "slots"}).
  AddRow(shelter.address, shelter.city, shelter.municipality, shelter.position_lat, shelter.position_long, shelter.shelter_id, shelter.slots)

  sm.ExpectQuery("SELECT " + ShelterFields + " FROM `shelters`").WithArgs(x1, y1, x2, y2).WillReturnRows(rows)

  a := App{
    DB: db,
  }
  shelters := make([]Shelter, 0)
  shelters = append(shelters, shelter)
  a.On("RespondWithJson", res, http.StatusOK, shelters)
  a.GetSheltersByBBox(res, req)
  a.AssertExpectations(t)

  if err := sm.ExpectationsWereMet(); err != nil {
    t.Errorf("Tests failed: %s", err)
  }

  if res.Code != http.StatusOK {
    t.Errorf("Expected to retrieve status OK. Got: %d", res.Code)
  }
}

func TestGetSheltersByBBoxFailureWithoutParams (t *testing.T) {
  req, _ := http.NewRequest("GET", "/shelters", nil)
  res := httptest.NewRecorder()


  a := App{}
  err := a.GetSheltersByBBox(res, req)

  if m, _ := regexp.MatchString("bbox-parametern är obligatorisk", fmt.Sprintf("%s", err)); m == false {
    t.Errorf("Tests failed: Expected an error with accurate message, got: %s", err)
  }
}

func TestGetSheltersByBBoxFailureWithBadBBoxParams (t *testing.T) {
  req, _ := http.NewRequest("GET", "/shelters", nil)
  res := httptest.NewRecorder()
  req = mux.SetURLVars(req, map[string]string{"bbox": "1"})

  a := App{}
  err := a.GetSheltersByBBox(res, req)

  if m, _ := regexp.MatchString("för få värden", fmt.Sprintf("%s", err)); m == false {
    t.Errorf("Tests failed: Expected an error with accurate message, got: %s", err)
  }
}

func TestGetSingleShelterSuccess (t *testing.T) {
  id := "3223"
  req, _ := http.NewRequest("GET", fmt.Sprintf("/shelters/%s", id), nil)
  req = mux.SetURLVars(req, map[string]string{"id": id })
  res := httptest.NewRecorder()
  db, sm, err := sqlmock.New()
  if err != nil {
    t.Fatalf("Create db mock failed: %s", err)
  }
  defer db.Close()

  shelter := Shelter{
    address: "Sätra solfångare ",
    city: "Skärholmen",
    municipality: "Stockholm",
    position_lat: "16.205",
    position_long: "59.4234",
    shelter_id: id,
    slots: 91,
  }
  rows := sqlmock.NewRows([]string{"address", "city", "municipality", "position_lat", "position_long", "shelter_id", "slots"}).
  AddRow(shelter.address, shelter.city, shelter.municipality, shelter.position_lat, shelter.position_long, shelter.shelter_id, shelter.slots)

  sm.ExpectQuery("SELECT " + ShelterFields + " FROM `shelters` WHERE id = ?").WithArgs(id).WillReturnRows(rows)

  a := App{
    DB: db,
  }
  shelters := make([]Shelter, 0)
  shelters = append(shelters, shelter)
  a.On("RespondWithJson", res, http.StatusOK, shelter)
  a.GetShelterById(res, req)
  a.AssertExpectations(t)

  if err := sm.ExpectationsWereMet(); err != nil {
    t.Errorf("Tests failed: %s", err)
  }

  if res.Code != http.StatusOK {
    t.Errorf("Expected to retrieve status OK. Got: %d", res.Code)
  }
}

func TestGetSingleShelterNotFound (t *testing.T) {
  req, _ := http.NewRequest("GET", "/shelters/unknown", nil)
  id := "144"
  req = mux.SetURLVars(req, map[string]string{"id": id })
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
  a.On("RespondWithJson", res, http.StatusNotFound, nil)

  err = a.GetShelterById(res, req)
  if err != nil {
    t.Fatalf("Failed: %s", err)
  }

  a.AssertExpectations(t)
}

func TestGetSingleShelterWithoutId (t *testing.T) {
  req, _ := http.NewRequest("GET", "/shelters/1", nil)
  res := httptest.NewRecorder()

  a := App{}
  err := a.GetShelterById(res, req)

  if m, _ := regexp.MatchString("id krävs", fmt.Sprintf("%s", err)); m == false {
    t.Errorf("Tests failed: Expected an error with accurate message, got: %s", err)
  }
}
