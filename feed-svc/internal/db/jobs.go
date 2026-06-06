package db

import (
	"context"
	"database/sql"
)

type Job struct {
	JobID     string
	PostedBy  string
	Title     string
	Company   string
	Location  string
	CreatedAt string
}

func GetJobByID(ctx context.Context, db *sql.DB, jobID string) (*Job, error) {
	row := db.QueryRowContext(ctx, `
		SELECT job_id, COALESCE(posted_by::text, ''), title, company, COALESCE(location, ''), created_at::text
		FROM jobs WHERE job_id = $1
	`, jobID)

	j := &Job{}
	err := row.Scan(&j.JobID, &j.PostedBy, &j.Title, &j.Company, &j.Location, &j.CreatedAt)
	if err != nil {
		return nil, err
	}
	return j, nil
}
