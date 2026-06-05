package db

import (
	"context"
	"database/sql"
)

type Skill struct {
	SkillID  string
	Name     string
	Endorsed int
}

// UpsertSkill inserts a skill if it doesn't exist, returns the skill_id either way.
// ON CONFLICT DO NOTHING + RETURNING handles the "already exists" case cleanly.
func UpsertSkill(ctx context.Context, db *sql.DB, name string) (string, error) {
	var skillID string

	// Try to insert — if skill name already exists, do nothing and fetch existing ID
	err := db.QueryRowContext(ctx, `
		INSERT INTO skills (name) VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING skill_id
	`, name).Scan(&skillID)
	return skillID, err
}

// LinkSkillToUser creates the user_skills row — links a user to a skill
// ON CONFLICT DO NOTHING means adding same skill twice won't error
func LinkSkillToUser(ctx context.Context, db *sql.DB, userID, skillID string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO user_skills (user_id, skill_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, skill_id) DO NOTHING
	`, userID, skillID)
	return err
}

// GetUserSkills fetches all skills for a user with endorsement counts
func GetUserSkills(ctx context.Context, db *sql.DB, userID string) ([]*Skill, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT s.skill_id, s.name, us.endorsed
		FROM skills s
		JOIN user_skills us ON s.skill_id = us.skill_id
		WHERE us.user_id = $1
		ORDER BY us.endorsed DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*Skill
	for rows.Next() {
		s := &Skill{}
		if err := rows.Scan(&s.SkillID, &s.Name, &s.Endorsed); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}
