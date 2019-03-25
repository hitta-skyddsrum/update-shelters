package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_ADDRESS")))
	if err != nil {
		panic(err)
	}

	a := App{
		DB: db,
	}

	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/shelters", a.GetSheltersByBBox).Methods("GET")
	a.Router.HandleFunc("/shelters/:id", a.GetShelterById).Methods("GET")
}
