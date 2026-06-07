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
	if err != nil {
		return err
	}
	return ensureIndex()
}

func ensureIndex() error {
	mapping := `{
		"mappings": {
			"properties": {
				"user_id":   { "type": "keyword" },
				"name":      { "type": "text" },
				"headline":  { "type": "text" },
				"location":  { "type": "text" },
				"skills":    { "type": "text" }
			}
		}
	}`
	res, err := client.Indices.Create("users",
		client.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

type UserDoc struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Headline string `json:"headline"`
	Location string `json:"location"`
	Skills   string `json:"skills"`
}

func IndexUser(ctx context.Context, doc UserDoc) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	res, err := client.Index("users",
		bytes.NewReader(data),
		client.Index.WithDocumentID(doc.UserID),
		client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
