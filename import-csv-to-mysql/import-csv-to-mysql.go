package main

import (
  "database/sql"
  "fmt"
  "github.com/go-sql-driver/mysql"
  "path/filepath"
  "regexp"
  "strings"
)

func prepareDb(schema string, mgrStmts string, db *sql.DB) {
  _, err := db.Exec("CREATE SCHEMA `" + schema + "` CHARACTER SET utf8 COLLATE utf8_general_ci;")
  if err != nil {
    fmt.Printf("Failed to create schema %s", schema)
    panic(err)
  }

  _, err = db.Exec("USE `" + schema + "`")
  if err != nil {
    panic(err)
  }

  fmt.Println("Running database migration")
  _, err = db.Exec(mgrStmts)
  if err != nil {
    panic(err)
  }
}

func loadCsvToDb(filePath string, db *sql.DB) {
  fmt.Println("Start import CSV to database.")
  mysql.RegisterLocalFile(filePath)
  _, err := db.Exec("LOAD DATA LOCAL INFILE '" + filePath + "' INTO TABLE shelters FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE 1 LINES")
  if err != nil {
    panic(err)
  }
  fmt.Println("End import CSV to database.")
}

func ImportCsvToMysql (db *sql.DB, mgrStmts string, csvPath string) {
  reg, err := regexp.Compile("[^a-zA-Z0-9]+")
  if err != nil {
    panic(err)
  }

  sqlSchema := reg.ReplaceAllString(strings.TrimSuffix(csvPath, filepath.Ext(csvPath)), "")
  prepareDb(sqlSchema, mgrStmts, db)
  loadCsvToDb(csvPath, db)
}
