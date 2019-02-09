package main

import (
  "database/sql"
  "fmt"
  "github.com/go-sql-driver/mysql"
  "path/filepath"
  "regexp"
  "strings"
)

func prepareDb(schema string, mgrStmts string, tx *sql.Tx) {
  _, err := tx.Exec("CREATE SCHEMA `" + schema + "` CHARACTER SET utf8 COLLATE utf8_general_ci;")
  if err != nil {
    fmt.Printf("Failed to create schema %s", schema)
    panic(err)
  }

  _, err = tx.Exec("USE `" + schema + "`")
  if err != nil {
    panic(err)
  }

  fmt.Println("Running database migration")
  _, err = tx.Exec(mgrStmts)
  if err != nil {
    panic(err)
  }
}

func loadCsvToDb(filePath string, tx *sql.Tx) {
  fmt.Println("Start import CSV to database.")
  mysql.RegisterLocalFile(filePath)
  _, err := tx.Exec(fmt.Sprintf(
    "LOAD DATA LOCAL INFILE '%s' " +
    "INTO TABLE shelters " +
    "FIELDS TERMINATED BY ',' " +
    "OPTIONALLY ENCLOSED BY '\"' " +
    "LINES TERMINATED BY '\n' " +
    "IGNORE 1 LINES",
  filePath))
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

  sqlSchema := reg.ReplaceAllString(strings.TrimSuffix(filepath.Base(csvPath), filepath.Ext(csvPath)), "")
  tx, err := db.Begin()
  if err != nil {
    fmt.Printf("Failed to create transaction: %s", err)
    panic(err)
  }
  defer func() {
    if r := recover(); r != nil {
      tx.Rollback()
      return
    }

    err = tx.Commit()

    if err != nil {
      fmt.Printf("Commiting DB transaction failed: %s", err)
    }
  }()

  prepareDb(sqlSchema, mgrStmts, tx)
  loadCsvToDb(csvPath, tx)
}
