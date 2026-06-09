package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Job struct {
	JobID           string
	PostedBy        sql.NullString
	Title           string
	Company         string
	Location        sql.NullString
	JobType         string
	ExperienceLevel string
	SalaryMin       sql.NullInt32
	SalaryMax       sql.NullInt32
	Description     string
	IsActive        bool
	CreatedAt       string
}

func CreateJob(ctx context.Context, db *sql.DB, postedBy, title, company, location, jobType, experienceLevel, description string, salaryMin, salaryMax *int) (*Job, error) {
	row := db.QueryRowContext(ctx, `
		INSERT INTO jobs (posted_by, title, company, location, job_type, experience_level, salary_min, salary_max, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING job_id, posted_by, title, company, location, job_type, experience_level, salary_min, salary_max, description, is_active, created_at::text
	`, nullStr(postedBy), title, company, nullStr(location), jobType, experienceLevel, nullInt(salaryMin), nullInt(salaryMax), description)

	return scanJob(row)
}

func GetJobByID(ctx context.Context, db *sql.DB, jobID string) (*Job, error) {
	row := db.QueryRowContext(ctx, `
		SELECT job_id, posted_by, title, company, location, job_type, experience_level, salary_min, salary_max, description, is_active, created_at::text
		FROM jobs WHERE job_id = $1
	`, jobID)
	return scanJob(row)
}

func ListJobs(ctx context.Context, db *sql.DB) ([]*Job, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT job_id, posted_by, title, company, location, job_type, experience_level, salary_min, salary_max, description, is_active, created_at::text
		FROM jobs WHERE is_active = true ORDER BY created_at DESC LIMIT 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		j := &Job{}
		if err := rows.Scan(&j.JobID, &j.PostedBy, &j.Title, &j.Company, &j.Location,
			&j.JobType, &j.ExperienceLevel, &j.SalaryMin, &j.SalaryMax,
			&j.Description, &j.IsActive, &j.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func ListJobsCursor(ctx context.Context, db *sql.DB, limit int, after *string) ([]*Job, error) {
	var rows *sql.Rows
	var err error
	if after != nil && *after != "" {
		// decode cursor — it's just the job_id (base64 decoded by caller)
		rows, err = db.QueryContext(ctx, `
			SELECT job_id, posted_by, title, company, location, job_type, experience_level, salary_min, salary_max, description, is_active, created_at::text
			FROM jobs WHERE is_active = true AND created_at < (SELECT created_at FROM jobs WHERE job_id = $1)
			ORDER BY created_at DESC LIMIT $2
		`, *after, limit)
	} else {
		rows, err = db.QueryContext(ctx, `
			SELECT job_id, posted_by, title, company, location, job_type, experience_level, salary_min, salary_max, description, is_active, created_at::text
			FROM jobs WHERE is_active = true ORDER BY created_at DESC LIMIT $1
		`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []*Job
	for rows.Next() {
		j := &Job{}
		if err := rows.Scan(&j.JobID, &j.PostedBy, &j.Title, &j.Company, &j.Location,
			&j.JobType, &j.ExperienceLevel, &j.SalaryMin, &j.SalaryMax,
			&j.Description, &j.IsActive, &j.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}


func ListJobsByCompanies(ctx context.Context, db *sql.DB, companies []string) ([]*Job, error) {
	if len(companies) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(companies))
	args := make([]any, len(companies))
	for i, c := range companies {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = c
	}
	query := fmt.Sprintf(`SELECT job_id, posted_by, title, company, location, job_type, experience_level, salary_min, salary_max, description, is_active, created_at::text
		FROM jobs WHERE is_active = true AND company IN (%s) ORDER BY created_at DESC LIMIT 20`,
		strings.Join(placeholders, ","))
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []*Job
	for rows.Next() {
		j := &Job{}
		if err := rows.Scan(&j.JobID, &j.PostedBy, &j.Title, &j.Company, &j.Location,
			&j.JobType, &j.ExperienceLevel, &j.SalaryMin, &j.SalaryMax,
			&j.Description, &j.IsActive, &j.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func scanJob(row *sql.Row) (*Job, error) {
	j := &Job{}
	err := row.Scan(&j.JobID, &j.PostedBy, &j.Title, &j.Company, &j.Location,
		&j.JobType, &j.ExperienceLevel, &j.SalaryMin, &j.SalaryMax,
		&j.Description, &j.IsActive, &j.CreatedAt)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullInt(i *int) any {
	if i == nil {
		return nil
	}
	return *i
}
