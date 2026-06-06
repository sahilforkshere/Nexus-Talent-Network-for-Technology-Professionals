package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

var client *es8.Client

const indexName = "jobs"

func Init(url string) error {
	var err error
	client, err = es8.NewClient(es8.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return err
	}
	ensureIndex()
	return nil
}

// ensureIndex creates the jobs index if it doesn't exist
func ensureIndex() {
	mapping := `{
		"mappings": {
			"properties": {
				"job_id":           { "type": "keyword" },
				"title":            { "type": "text" },
				"company":          { "type": "text" },
				"location":         { "type": "text" },
				"description":      { "type": "text" },
				"job_type":         { "type": "keyword" },
				"experience_level": { "type": "keyword" }
			}
		}
	}`

	res, err := client.Indices.Exists([]string{indexName})
	if err != nil || res.StatusCode == 200 {
		return
	}

	res, err = client.Indices.Create(indexName, client.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		log.Printf("es: failed to create index: %v", err)
		return
	}
	defer res.Body.Close()
	log.Println("es: created jobs index")
}

type JobDoc struct {
	JobID           string `json:"job_id"`
	Title           string `json:"title"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Description     string `json:"description"`
	JobType         string `json:"job_type"`
	ExperienceLevel string `json:"experience_level"`
}

// IndexJob adds a job document to Elasticsearch
func IndexJob(ctx context.Context, doc JobDoc) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := client.Index(
		indexName,
		bytes.NewReader(body),
		client.Index.WithDocumentID(doc.JobID),
		client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es index error: %s", res.String())
	}
	log.Printf("es: indexed job %s (%s)", doc.Title, doc.JobID)
	return nil
}

// SearchJobs does a multi-field keyword search across title, company, description, location
func SearchJobs(ctx context.Context, keyword string) ([]string, error) {
	query := fmt.Sprintf(`{
		"query": {
			"multi_match": {
				"query": %q,
				"fields": ["title^3", "company^2", "description", "location"]
			}
		},
		"size": 20
	}`, keyword)

	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithBody(strings.NewReader(query)),
		client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("es search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source JobDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var ids []string
	for _, h := range result.Hits.Hits {
		ids = append(ids, h.Source.JobID)
	}
	return ids, nil
}
