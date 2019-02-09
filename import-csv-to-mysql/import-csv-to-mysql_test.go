package main

import (
  "database/sql/driver"
  "fmt"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "regexp"
  "testing"
)

type TimeMatcher  struct{}

func (t TimeMatcher) Match(v driver.Value) bool {
  m, err := regexp.MatchString("((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])", v.(string))
  if err != nil {
    panic(err)
  }

  if m == false {
    fmt.Printf("Asserted %s to match", v)
  }

  return m
}

func TestMigrateImportSuccess(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("An error occured while opening mock DB: %s", err)
  }
  defer db.Close()
  schemaName := "20180909"
  fileName := fmt.Sprintf("/var/data/%s.csv", schemaName)
  mgrStmts := "INSERT cool stuff"

  mock.ExpectBegin()
  mock.ExpectExec(fmt.Sprintf("CREATE SCHEMA `%s`", schemaName)).WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec(fmt.Sprintf("USE `%s`", schemaName)).WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec(mgrStmts).WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec(fmt.Sprintf("LOAD DATA LOCAL INFILE '%s'", fileName)).WillReturnResult(sqlmock.NewResult(1, 1))

  mock.ExpectExec("USE `db_metadata`").WillReturnResult(sqlmock.NewResult(0, 0))
  mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `datasets` (name, published_at, is_verified) VALUES (?, ?, ?)")).WithArgs(schemaName, TimeMatcher{}, false).WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectCommit()

  ImportCsvToMysql(db, mgrStmts, fileName)

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("Tests failed: %s", err)
  }
}

func TestMigrateCommitFailure(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("An error occured while opening mock DB: %s", err)
  }

  mock.ExpectBegin()
  mock.ExpectExec("CREATE SCHEMA").WillReturnError(fmt.Errorf("Bad commit"))
  mock.ExpectRollback()

  ImportCsvToMysql(db, "", "")

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("Tests failed: %s", err)
  }
}
