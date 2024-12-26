package query

import "database/sql"

type Query struct {
	DB *sql.DB
}

func NewQuery(db *sql.DB) Query {
	return Query{
		DB: db,
	}
}
