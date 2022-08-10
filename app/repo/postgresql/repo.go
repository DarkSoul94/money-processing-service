package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/DarkSoul94/money-processing-service/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type postgreRepo struct {
	db *sqlx.DB
}

func NewPostgreRepo(db *sql.DB) app.Repository {
	return &postgreRepo{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (r *postgreRepo) CreateClient(ctx context.Context, mClient models.Client) (uint64, error) {
	var (
		client dbClient
		query  string
		id     int64
		err    error
	)

	client = r.toDbClient(mClient)

	query = `INSERT INTO clients (name) VALUES ($1) RETURNING id`

	err = r.db.GetContext(ctx, &id, query, client.Name)
	if err != nil {
		logger.LogError(
			"Create client",
			"app/repo/postgres/repo",
			fmt.Sprintf("name: %s", client.Name),
			err,
		)
		return 0, errors.New("failed create client")
	}

	fmt.Println(id)
	return uint64(id), nil
}

func (r *postgreRepo) Close() error {
	return r.db.Close()
}
