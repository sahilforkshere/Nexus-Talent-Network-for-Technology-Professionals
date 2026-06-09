package recommend

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var driver neo4j.DriverWithContext

func Init(uri, user, password string) error {
	d, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return err
	}
	if err := d.VerifyConnectivity(context.Background()); err != nil {
		return err
	}
	driver = d
	log.Println("recommend: connected to neo4j")
	return nil
}

// CompaniesForUser returns companies where the user's 1st and 2nd degree connections work.
// Returns map[company]score where 1st degree = 20, 2nd degree = 10.
func CompaniesForUser(ctx context.Context, userID string) map[string]int {
	if driver == nil {
		return nil
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, `
		MATCH (me:Person {user_id: $userID})-[:CONNECTED_TO]->(conn:Person)
		WHERE conn.company IS NOT NULL
		RETURN conn.company AS company, 20 AS score
		UNION
		MATCH (me:Person {user_id: $userID})-[:CONNECTED_TO]->(:Person)-[:CONNECTED_TO]->(conn2:Person)
		WHERE conn2.company IS NOT NULL
		RETURN conn2.company AS company, 10 AS score
	`, map[string]any{"userID": userID})
	if err != nil {
		log.Printf("recommend: neo4j query error: %v", err)
		return nil
	}

	companies := make(map[string]int)
	for result.Next(ctx) {
		record := result.Record()
		company, _ := record.Get("company")
		score, _ := record.Get("score")
		if c, ok := company.(string); ok {
			if s, ok := score.(int64); ok {
				companies[c] = int(s)
			}
		}
	}
	return companies
}
