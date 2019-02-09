// +build !lambda

package main

import (
  "database/sql"
  "fmt"
  _ "github.com/go-sql-driver/mysql"
  "github.com/hitta-skyddsrum/update-shelters/db-flags"
)

func main() {
  dbU, dbP, dbA := dbFlags.GetDbFlags()

  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/", *dbU, *dbP, *dbA))
  if err != nil {
    panic(err)
  }
  fmt.Println("Database connection open.")

  MigrateDb(db)
}
