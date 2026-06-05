package graph

import (
	"context"
	"fmt"

	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/profile-svc/graph/model"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/profile-svc/internal/auth"
	userdb "github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/profile-svc/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user account and returns JWT tokens
func (r *mutationResolver) Register(ctx context.Context, input model.RegisterInput) (*model.AuthPayload, error) {
	// Hash the password — never store plain text passwords
	// bcrypt cost 12 = slow enough to resist brute force, fast enough for UX
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	location := ""
	if input.Location != nil {
		location = *input.Location
	}

	// Insert user into Postgres
	u, err := userdb.CreateUser(ctx, r.DB, input.Email, string(hash), input.Name, location)
	if err != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Generate JWT tokens
	accessToken, err := auth.GenerateAccessToken(u.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	refreshToken, err := auth.GenerateRefreshToken(u.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	return &model.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dbUserToModel(u),
	}, nil
}

// Login verifies credentials and returns JWT tokens
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	u, err := userdb.GetUserByEmail(ctx, r.DB, input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare the provided password against the stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	accessToken, err := auth.GenerateAccessToken(u.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	refreshToken, err := auth.GenerateRefreshToken(u.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	return &model.AuthPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dbUserToModel(u),
	}, nil
}

// UpdateProfile updates the logged-in user's profile fields
func (r *mutationResolver) UpdateProfile(ctx context.Context, input model.UpdateProfileInput) (*model.User, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	// Use empty string for unset optional fields — DB uses COALESCE to keep existing value
	name, headline, bio, location, avatarURL := "", "", "", "", ""
	if input.Name != nil {
		name = *input.Name
	}
	if input.Headline != nil {
		headline = *input.Headline
	}
	if input.Bio != nil {
		bio = *input.Bio
	}
	if input.Location != nil {
		location = *input.Location
	}
	if input.AvatarURL != nil {
		avatarURL = *input.AvatarURL
	}

	u, err := userdb.UpdateUser(ctx, r.DB, userID, name, headline, bio, location, avatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile")
	}
	return dbUserToModel(u), nil
}

// AddSkill — Day 6
func (r *mutationResolver) AddSkill(ctx context.Context, name string) (*model.Skill, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// GetProfile returns a user's public profile by ID
func (r *queryResolver) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	u, err := userdb.GetUserByID(ctx, r.DB, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return dbUserToModel(u), nil
}

// Me returns the currently logged-in user's profile (reads user_id from JWT in context)
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}
	u, err := userdb.GetUserByID(ctx, r.DB, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return dbUserToModel(u), nil
}

// dbUserToModel converts the internal DB struct to the GraphQL model
func dbUserToModel(u *userdb.User) *model.User {
	user := &model.User{
		UserID:    u.UserID,
		Email:     u.Email,
		Name:      u.Name,
		Skills:    []*model.Skill{},
		CreatedAt: u.CreatedAt,
	}
	if u.Headline.String != "" {
		user.Headline = &u.Headline.String
	}
	if u.Bio.String != "" {
		user.Bio = &u.Bio.String
	}
	if u.Location.String != "" {
		user.Location = &u.Location.String
	}
	if u.AvatarURL.String != "" {
		user.AvatarURL = &u.AvatarURL.String
	}
	return user
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
