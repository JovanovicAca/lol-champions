package repository

import (
	"fmt"
)

const (
	host     = "localhost"
	server   = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "lol-champions"
)

var ConnectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
	host, port, user, password, dbname)
