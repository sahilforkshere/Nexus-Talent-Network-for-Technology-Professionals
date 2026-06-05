package db

import (
	"context"
	"database/sql"
)

type User struct {
	UserID       string
	Email        string
	PasswordHash string
	Name         string
	Headline     sql.NullString
	Bio          sql.NullString
	Location     sql.NullString
	AvatarURL    sql.NullString
	CreatedAt    string
}

// CreateUser inserts a new user row and returns the generated user_id
func CreateUser(ctx context.Context, db *sql.DB, email, passwordHash, name, location string) (*User, error) {
	u := &User{}
	err := db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, name, location)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, email, password_hash, name,
		          COALESCE(headline,''), COALESCE(bio,''), COALESCE(location,''),
		          COALESCE(avatar_url,''), created_at::text
	`, email, passwordHash, name, location).Scan(
		&u.UserID, &u.Email, &u.PasswordHash, &u.Name,
		&u.Headline.String, &u.Bio.String, &u.Location.String,
		&u.AvatarURL.String, &u.CreatedAt,
	)
	return u, err
}

// GetUserByEmail fetches a user by email — used during login
func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (*User, error) {
	u := &User{}
	err := db.QueryRowContext(ctx, `
		SELECT user_id, email, password_hash, name,
		       COALESCE(headline,''), COALESCE(bio,''), COALESCE(location,''),
		       COALESCE(avatar_url,''), created_at::text
		FROM users WHERE email = $1
	`, email).Scan(
		&u.UserID, &u.Email, &u.PasswordHash, &u.Name,
		&u.Headline.String, &u.Bio.String, &u.Location.String,
		&u.AvatarURL.String, &u.CreatedAt,
	)
	return u, err
}

// UpdateUser updates editable profile fields, returns the updated user
func UpdateUser(ctx context.Context, db *sql.DB, userID, name, headline, bio, location, avatarURL string) (*User, error) {
	u := &User{}
	err := db.QueryRowContext(ctx, `
		UPDATE users SET
			name       = COALESCE(NULLIF($2,''), name),
			headline   = COALESCE(NULLIF($3,''), headline),
			bio        = COALESCE(NULLIF($4,''), bio),
			location   = COALESCE(NULLIF($5,''), location),
			avatar_url = COALESCE(NULLIF($6,''), avatar_url),
			updated_at = NOW()
		WHERE user_id = $1
		RETURNING user_id, email, password_hash, name,
		          COALESCE(headline,''), COALESCE(bio,''), COALESCE(location,''),
		          COALESCE(avatar_url,''), created_at::text
	`, userID, name, headline, bio, location, avatarURL).Scan(
		&u.UserID, &u.Email, &u.PasswordHash, &u.Name,
		&u.Headline.String, &u.Bio.String, &u.Location.String,
		&u.AvatarURL.String, &u.CreatedAt,
	)
	return u, err
}

// GetUserByID fetches a user by their UUID — used in getProfile and me queries
func GetUserByID(ctx context.Context, db *sql.DB, userID string) (*User, error) {
	u := &User{}
	err := db.QueryRowContext(ctx, `
		SELECT user_id, email, password_hash, name,
		       COALESCE(headline,''), COALESCE(bio,''), COALESCE(location,''),
		       COALESCE(avatar_url,''), created_at::text
		FROM users WHERE user_id = $1
	`, userID).Scan(
		&u.UserID, &u.Email, &u.PasswordHash, &u.Name,
		&u.Headline.String, &u.Bio.String, &u.Location.String,
		&u.AvatarURL.String, &u.CreatedAt,
	)
	return u, err
}
