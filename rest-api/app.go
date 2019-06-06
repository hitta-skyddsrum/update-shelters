package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type ErrorResponse struct {
  Message     string `json:"message"`
}

func (a *App) RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) RespondWithError(w http.ResponseWriter, code int, message string) {
  response := ErrorResponse{
    Message: message,
  }

  a.RespondWithJson(w, code, response)
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_ADDRESS"), os.Getenv("MYSQL_DATABASE_NAME")))
	if err != nil {
		panic(err)
	}

	a := App{
		DB: db,
	}

	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/shelters", a.GetSheltersByBBox).Methods("GET")
	a.Router.HandleFunc("/shelters/:id", a.GetShelterById).Methods("GET")

  srv := &http.Server{
    Handler:      a.Router,
    Addr:         "127.0.0.1:8000",
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  }

  log.Fatal(srv.ListenAndServe())
}
