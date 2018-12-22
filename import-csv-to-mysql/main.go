package main

import (
  "database/sql"
  "flag"
  "github.com/go-sql-driver/mysql"
  "io/ioutil"
  "path/filepath"
  "os"
  "strings"
)

func prepareDb(schema string, sqlFile string, db *sql.DB) {
  _, err := db.Exec("CREATE SCHEMA `" + schema + "`")
  if err != nil {
    panic(err)
  }

  _, err = db.Exec("USE `" + schema + "`")
  if err != nil {
    panic(err)
  }

  b, err := ioutil.ReadFile(sqlFile)
  if err != nil {
    panic(err)
  }

  s := string(b)

  _, err = db.Exec(s)
  if err != nil {
    panic(err)
  }

}

func loadCsvToDb(filePath string, db *sql.DB) {
  mysql.RegisterLocalFile(filePath)
  _, err := db.Exec("LOAD DATA LOCAL INFILE '" + filePath + "' INTO TABLE shelters FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE 1 LINES")
  if err != nil {
    panic(err)
  }
}

func main() {
  sqlMigration := flag.String("sql-migration", "", "SQL file containing migration for shelters table")
  flag.Parse()

  if *sqlMigration == "" {
    os.Exit(0)
  }

  csvFile := flag.Args()[0]

  db, err := sql.Open("mysql", "root:shelters@(127.0.0.1:33060)/shelters")
  if err != nil {
    panic(err)
  }

  sqlSchema := strings.TrimSuffix(csvFile, filepath.Ext(csvFile))
  prepareDb(sqlSchema, *sqlMigration, db)
  loadCsvToDb(csvFile, db)
}
