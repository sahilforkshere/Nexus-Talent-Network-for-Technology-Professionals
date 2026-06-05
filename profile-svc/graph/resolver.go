package graph

import "database/sql"

// Resolver is the root dependency container for all GraphQL resolvers.
// Add any service dependencies here (DB, Kafka, Redis, etc.)
type Resolver struct {
	DB *sql.DB
}
