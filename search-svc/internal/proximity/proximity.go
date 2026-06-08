package proximity

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var driver neo4j.DriverWithContext

func Init(uri, username, password string) error {
	d, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return err
	}
	driver = d
	return driver.VerifyConnectivity(context.Background())
}

// ConnectedCompanies returns a set of company names where the user has
// 1st or 2nd degree connections who have a WORKS_AT relationship.
// Falls back to empty map if Neo4j unavailable — search still works, just unranked.
func ConnectedCompanies(ctx context.Context, userID string) map[string]int {
	if driver == nil {
		return nil
	}
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx, `
			MATCH (me:Person {user_id: $user_id})-[:CONNECTED_TO]->(c1:Person)
			WHERE c1.company IS NOT NULL
			WITH c1.company AS company, 1 AS degree
			UNION
			MATCH (me:Person {user_id: $user_id})-[:CONNECTED_TO]->(:Person)-[:CONNECTED_TO]->(c2:Person)
			WHERE c2.company IS NOT NULL AND c2.user_id <> $user_id
			WITH c2.company AS company, 2 AS degree
			RETURN company, min(degree) AS degree
		`, map[string]any{"user_id": userID})
		if err != nil {
			return nil, err
		}
		companies := map[string]int{}
		for records.Next(ctx) {
			rec := records.Record()
			company, _ := rec.Get("company")
			degree, _ := rec.Get("degree")
			if c, ok := company.(string); ok {
				if d, ok := degree.(int64); ok {
					companies[c] = int(d)
				}
			}
		}
		return companies, nil
	})
	if err != nil || result == nil {
		return nil
	}
	return result.(map[string]int)
}

// ProximityScore returns a boost score for a company.
// 1st degree connection = 20 points, 2nd degree = 10 points, none = 0.
func ProximityScore(company string, connected map[string]int) int {
	if connected == nil {
		return 0
	}
	degree, ok := connected[company]
	if !ok {
		return 0
	}
	if degree == 1 {
		return 20
	}
	return 10
}
