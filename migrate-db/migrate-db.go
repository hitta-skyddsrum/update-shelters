package main

import (
	"database/sql"
	"fmt"
)

func MigrateDb(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Failed to create transaction.")
		panic(err)
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	schema := "db_metadata"
	fmt.Printf("Creating schema %s\n", schema)
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA `%s` CHARACTER SET utf8 COLLATE utf8_general_ci;", schema))
	if err != nil {
		fmt.Printf("Failed to create schema %s\n", schema)
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf("USE `%s`", schema))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE `datasets` (id INT(10) AUTO_INCREMENT PRIMARY KEY , name VARCHAR(100), published_at VARCHAR(100), is_verified BOOL)")
	if err != nil {
		fmt.Println("Failed to create table `datasets`")
		panic(err)
	}
}
