// +build !lambda

package main

import (
  "database/sql"
  "flag"
  "fmt"
  "os"
)

func main() {
  dbU := flag.String("db-user", "", "Database user")
  dbP := flag.String("db-password", "", "Database password")
  dbA := flag.String("db-address", "", "Database address, including port number.")
  flag.Parse()

  missReq := false
  flag.VisitAll(func(f *flag.Flag) {
    if f.Value.String() == "" && f.Name != "db-password" {
      fmt.Printf("Flag hasnt been set %s", f.Name)
      fmt.Println()
      missReq = true
    }
  })
  if missReq != false {
    os.Exit(0)
  }

  csvFile := flag.Args()[0]

  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/", *dbU, *dbP, *dbA))
  if err != nil {
    panic(err)
  }
  fmt.Println("Database connection open.")

  ImportCsvToMysql(db, schema, csvFile)
}
