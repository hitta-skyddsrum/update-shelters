// +build !lambda

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/hitta-skyddsrum/update-shelters/db-flags"
)

func main() {
	dbU, dbP, dbA := dbFlags.GetDbFlags()

	csvFile := flag.Args()[0]

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/", *dbU, *dbP, *dbA))
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection open.")

	ImportCsvToMysql(db, schema, csvFile)
}
