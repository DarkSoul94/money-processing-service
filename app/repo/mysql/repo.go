package mysql

import (
	"database/sql"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/jmoiron/sqlx"
)

type mysqlrepo struct {
	db *sqlx.DB
}

func NewMySQLRepo(db *sql.DB) app.Repository {
	return &mysqlrepo{
		db: sqlx.NewDb(db, "mysql"),
	}
}

func (r *mysqlrepo) Close() error {
	return r.db.Close()
}
