package main

import (
  "fmt"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "testing"
)

func TestMigrateImportSuccess(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("An error occured while opening mock DB: %s", err)
  }
  defer db.Close()
  schemaName := "20180909"
  fileName := fmt.Sprintf("%s.csv", schemaName)
  mgrStmts := "INSERT cool stuff"

  mock.ExpectExec(fmt.Sprintf("CREATE SCHEMA `%s`", schemaName)).WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec(fmt.Sprintf("USE `%s`", schemaName)).WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec(mgrStmts).WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec(fmt.Sprintf("LOAD DATA LOCAL INFILE '%s'", fileName)).WillReturnResult(sqlmock.NewResult(1, 1))

  ImportCsvToMysql(db, mgrStmts, fmt.Sprintf("%s.csv", schemaName))

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("Tests failed: %s", err)
  }
}
