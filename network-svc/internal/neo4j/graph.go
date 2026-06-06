package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreatePersonNode creates a Person node when a user registers.
// MERGE = create if not exists, safe to call multiple times (idempotent).
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

// CreateConnectionEdge creates a bidirectional CONNECTED_TO edge between two users.
// Called when a connection request is accepted.
// One transaction writes both directions atomically — either both succeed or neither does.
func CreateConnectionEdge(ctx context.Context, driver neo4j.DriverWithContext, userIDA, userIDB string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (a:Person {user_id: $user_id_a})
			MATCH (b:Person {user_id: $user_id_b})
			MERGE (a)-[:CONNECTED_TO]->(b)
			MERGE (b)-[:CONNECTED_TO]->(a)
		`, map[string]any{
			"user_id_a": userIDA,
			"user_id_b": userIDB,
		})
		return nil, err
	})
	return err
}

// PeopleYouMayKnow returns users 2 hops away who are not already connected.
// This is the core graph query — impossible to do efficiently in SQL at scale.
func PeopleYouMayKnow(ctx context.Context, driver neo4j.DriverWithContext, userID string) ([]map[string]any, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx, `
			MATCH (me:Person {user_id: $user_id})-[:CONNECTED_TO]->(friend)-[:CONNECTED_TO]->(suggestion)
			WHERE suggestion.user_id <> $user_id
			AND NOT (me)-[:CONNECTED_TO]->(suggestion)
			WITH suggestion, count(friend) AS mutual_count
			RETURN suggestion.user_id AS user_id,
			       suggestion.name    AS name,
			       suggestion.location AS location,
			       mutual_count
			ORDER BY mutual_count DESC
			LIMIT 10
		`, map[string]any{"user_id": userID})
		if err != nil {
			return nil, err
		}

		var suggestions []map[string]any
		for records.Next(ctx) {
			record := records.Record()
			suggestions = append(suggestions, record.AsMap())
		}
		return suggestions, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]map[string]any), nil
}

// MutualConnectionCount returns how many mutual connections two users share.
// Shown as "12 mutual connections" on profile cards.
func MutualConnectionCount(ctx context.Context, driver neo4j.DriverWithContext, userIDA, userIDB string) (int64, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		record, err := tx.Run(ctx, `
			MATCH (a:Person {user_id: $user_id_a})-[:CONNECTED_TO]->(mutual)<-[:CONNECTED_TO]-(b:Person {user_id: $user_id_b})
			RETURN count(mutual) AS mutual_count
		`, map[string]any{
			"user_id_a": userIDA,
			"user_id_b": userIDB,
		})
		if err != nil {
			return int64(0), err
		}
		if record.Next(ctx) {
			count, _ := record.Record().Get("mutual_count")
			return count.(int64), nil
		}
		return int64(0), nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}
