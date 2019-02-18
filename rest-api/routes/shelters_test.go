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

  sm.ExpectQuery("SELECT address, city, municipality, position_lat, position_long, shelter_id, slots FROM `shelters`").WithArgs(x1, y1, x2, y2).WillReturnRows(rows)

  a := App{
    DB: db,
  }
  shelters := make([]Shelter, 0)
  shelters = append(shelters, shelter)
  a.On("RespondWithJson", res, http.StatusOK, shelters)
  // a.AssertCalled(t, "RespondWithJson", res, http.StatusOK, mock.AnythingOfType("[]Shelter"))
  a.AssertNumberOfCalls(t, "RespondWithJson", 1)
  a.GetSheltersByBBox(res, req)
  a.AssertExpectations(t)

  return
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

  defer func() {
    r := recover()

    if r == nil {
      t.Errorf("Tests failed: Expected a panic.")
    }

    err := fmt.Sprintf("%s", r)
    if m, _ := regexp.MatchString("bbox-parametern är obligatorisk", err); m == false {
      t.Errorf("Tests failed: Expected a panic with accurate message, got: %s", r)
    }
  }()

  a := App{}
  a.GetSheltersByBBox(res, req)
}

func TestGetSheltersByBBoxFailureWithBadBBoxParams (t *testing.T) {
  req, _ := http.NewRequest("GET", "/shelters", nil)
  res := httptest.NewRecorder()
  req = mux.SetURLVars(req, map[string]string{"bbox": "1"})

  defer func() {
    r := recover()

    if r == nil {
      t.Errorf("Tests failed: Expected a panic.")
    }

    err := fmt.Sprintf("%s", r)
    if m, _ := regexp.MatchString("för få värden", err); m == false {
      t.Errorf("Tests failed: Expected panic with accurate message, got: %s", r)
    }
  }()

  a := App{}
  a.GetSheltersByBBox(res, req)
}
