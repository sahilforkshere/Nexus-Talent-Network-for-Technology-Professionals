package graph

import (
	"encoding/base64"

	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/graph/model"
	jobsdb "github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/internal/db"
)

func dbJobToModel(j *jobsdb.Job) *model.Job {
	job := &model.Job{
		JobID:           j.JobID,
		Title:           j.Title,
		Company:         j.Company,
		JobType:         j.JobType,
		ExperienceLevel: j.ExperienceLevel,
		Description:     j.Description,
		IsActive:        j.IsActive,
		CreatedAt:       j.CreatedAt,
	}
	if j.PostedBy.Valid {
		job.PostedBy = &j.PostedBy.String
	}
	if j.Location.Valid {
		job.Location = &j.Location.String
	}
	if j.SalaryMin.Valid {
		v := int(j.SalaryMin.Int32)
		job.SalaryMin = &v
	}
	if j.SalaryMax.Valid {
		v := int(j.SalaryMax.Int32)
		job.SalaryMax = &v
	}
	return job
}

// encodeCursor encodes a job_id as a base64 opaque cursor.
func encodeCursor(jobID string) string {
	return base64.StdEncoding.EncodeToString([]byte(jobID))
}

// decodeCursor decodes a cursor back to job_id.
func decodeCursor(cursor string) string {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return cursor
	}
	return string(b)
}
