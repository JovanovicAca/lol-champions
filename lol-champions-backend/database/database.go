package database

import (
	"context"
	"database/sql"
	"fmt"
)

type Database struct {
	SqlDb *sql.DB
}

var dbContext = context.Background()
var ConnectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
	host, port, user, password, dbname)
