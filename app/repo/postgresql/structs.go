package postgresql

import "github.com/DarkSoul94/money-processing-service/models"

type dbClient struct {
	Id   uint64 `db:"id"`
	Name string `db:"name"`
}

func (r *postgreRepo) toDbClient(client models.Client) dbClient {
	return dbClient{
		Id:   client.Id,
		Name: client.Name,
	}
}
