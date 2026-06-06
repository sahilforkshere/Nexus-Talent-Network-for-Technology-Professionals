package graph

import (
	"context"
	"fmt"

	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/graph/model"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/internal/auth"
	neo4jdb "github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/internal/neo4j"
)

// SendConnectionRequest creates a PENDING relationship from the logged-in user to toUserID.
// Think of it as "Sahil clicks Connect on Rahul's profile".
func (r *mutationResolver) SendConnectionRequest(ctx context.Context, toUserID string) (*model.ConnectionRequest, error) {
	fromUserID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}
	if fromUserID == toUserID {
		return nil, fmt.Errorf("cannot connect to yourself")
	}

	if err := neo4jdb.CreatePendingRequest(ctx, r.Neo4j, fromUserID, toUserID); err != nil {
		return nil, fmt.Errorf("failed to send connection request")
	}

	return &model.ConnectionRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     "PENDING",
	}, nil
}

// AcceptConnection upgrades the PENDING edge to a bidirectional CONNECTED_TO edge.
// Think of it as "Rahul clicks Accept on Sahil's request".
func (r *mutationResolver) AcceptConnection(ctx context.Context, fromUserID string) (bool, error) {
	toUserID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return false, fmt.Errorf("not authenticated")
	}

	if err := neo4jdb.AcceptPendingRequest(ctx, r.Neo4j, fromUserID, toUserID); err != nil {
		return false, fmt.Errorf("failed to accept connection")
	}

	return true, nil
}

// GetPeopleYouMayKnow runs the 2-hop graph traversal:
// "Who are my connections' connections, that I'm not already connected to?"
func (r *queryResolver) GetPeopleYouMayKnow(ctx context.Context) ([]*model.PersonSuggestion, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	rows, err := neo4jdb.PeopleYouMayKnow(ctx, r.Neo4j, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch suggestions")
	}

	var suggestions []*model.PersonSuggestion
	for _, row := range rows {
		uid, _ := row["user_id"].(string)
		name, _ := row["name"].(string)
		location, _ := row["location"].(string)
		mutual, _ := row["mutual_count"].(int64)

		suggestions = append(suggestions, &model.PersonSuggestion{
			UserID:      uid,
			Name:        name,
			Location:    location,
			MutualCount: int(mutual),
		})
	}
	return suggestions, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
