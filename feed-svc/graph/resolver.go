package graph

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Resolver struct {
	DB    *sql.DB
	Redis *redis.Client
}
