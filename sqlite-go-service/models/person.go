package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	ID         int
	FIRST_NAME string
	LAST_NAME  string
	EMAIL      string
	IP_ADDRESS string
}
