package search

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
)

var client *elasticsearch.Client

func Init(url string) error {
	var err error
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	return err
}

type JobHit struct {
	JobID           string `json:"job_id"`
	Title           string `json:"title"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	JobType         string `json:"job_type"`
	ExperienceLevel string `json:"experience_level"`
}

type UserHit struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Headline string `json:"headline"`
	Location string `json:"location"`
}

func SearchJobs(ctx context.Context, keyword string) ([]JobHit, error) {
	body := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  keyword,
				"fields": []string{"title^3", "company^2", "description", "location"},
			},
		},
		"size": 10,
	}
	return queryIndex[JobHit](ctx, "jobs", body)
}

func SearchUsers(ctx context.Context, keyword string) ([]UserHit, error) {
	body := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  keyword,
				"fields": []string{"name^3", "headline^2", "location", "skills"},
			},
		},
		"size": 10,
	}
	return queryIndex[UserHit](ctx, "users", body)
}

func queryIndex[T any](ctx context.Context, index string, body map[string]any) ([]T, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source T `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var out []T
	for _, h := range result.Hits.Hits {
		out = append(out, h.Source)
	}
	return out, nil
}
