package graph

import (
	"database/sql"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Resolver is the root dependency container for all GraphQL resolvers.
// Add any service dependencies here (DB, Kafka, Redis, etc.)
type Resolver struct {
	DB     *sql.DB
	Neo4j  neo4j.DriverWithContext
}
