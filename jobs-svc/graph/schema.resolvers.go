package graph

import (
	"context"
	"fmt"

	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/graph/model"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/internal/auth"
	jobsdb "github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/internal/db"
)

func (r *mutationResolver) PostJob(ctx context.Context, input model.PostJobInput) (*model.Job, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	var salaryMin, salaryMax *int
	if input.SalaryMin != nil {
		v := *input.SalaryMin
		salaryMin = &v
	}
	if input.SalaryMax != nil {
		v := *input.SalaryMax
		salaryMax = &v
	}

	location := ""
	if input.Location != nil {
		location = *input.Location
	}

	j, err := jobsdb.CreateJob(ctx, r.DB, userID, input.Title, input.Company, location,
		input.JobType, input.ExperienceLevel, input.Description, salaryMin, salaryMax)
	if err != nil {
		return nil, fmt.Errorf("failed to create job")
	}

	return dbJobToModel(j), nil
}

func (r *queryResolver) GetJob(ctx context.Context, jobID string) (*model.Job, error) {
	j, err := jobsdb.GetJobByID(ctx, r.DB, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found")
	}
	return dbJobToModel(j), nil
}

func (r *queryResolver) ListJobs(ctx context.Context) ([]*model.Job, error) {
	jobs, err := jobsdb.ListJobs(ctx, r.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs")
	}
	var result []*model.Job
	for _, j := range jobs {
		result = append(result, dbJobToModel(j))
	}
	return result, nil
}

// SearchJobs will be implemented in Phase 3 (Elasticsearch)
func (r *queryResolver) SearchJobs(ctx context.Context, keyword string) ([]*model.Job, error) {
	return nil, fmt.Errorf("search not yet implemented — coming in Phase 3")
}

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

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
