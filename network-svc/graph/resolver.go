package graph

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type Resolver struct {
	Neo4j neo4j.DriverWithContext
}
