package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreatePersonNode creates a Person node in Neo4j when a user registers.
// MERGE = create if not exists, do nothing if already there (idempotent).
func CreatePersonNode(ctx context.Context, driver neo4j.DriverWithContext, userID, name, location string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MERGE (p:Person {user_id: $user_id})
			SET p.name = $name, p.location = $location
		`, map[string]any{
			"user_id":  userID,
			"name":     name,
			"location": location,
		})
		return nil, err
	})
	return err
}

// LinkSkillInGraph creates (:Person)-[:HAS_SKILL]->(:Skill) in Neo4j.
// MERGE on both nodes + relationship = fully idempotent, safe to call multiple times.
func LinkSkillInGraph(ctx context.Context, driver neo4j.DriverWithContext, userID, skillName string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MERGE (p:Person {user_id: $user_id})
			MERGE (s:Skill {name: $skill_name})
			MERGE (p)-[:HAS_SKILL]->(s)
		`, map[string]any{
			"user_id":    userID,
			"skill_name": skillName,
		})
		return nil, err
	})
	return err
}
