package main

import (
  "fmt"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "testing"
)

func TestMigrateDbSuccess(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  mock.ExpectBegin()
  mock.ExpectExec("CREATE SCHEMA `db_metadata`").WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec("USE `db_metadata`").WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectExec("CREATE TABLE `datasets`").WillReturnResult(sqlmock.NewResult(1, 1))
  mock.ExpectCommit()

  MigrateDb(db);

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("Tests failed: %s", err)
  }
}

func TestMigrateDbRollback(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  mock.ExpectBegin()
  mock.ExpectExec("CREATE SCHEMA `db_metadata`").WillReturnError(fmt.Errorf(""))
  mock.ExpectRollback()

  defer func() {
    if r := recover(); r != nil {}

    if err := mock.ExpectationsWereMet(); err != nil {
      t.Errorf("Tests failed: %s", err)
    }
  }()
  MigrateDb(db);

  t.Fatalf("Expected MigrateDb to panic")
}
