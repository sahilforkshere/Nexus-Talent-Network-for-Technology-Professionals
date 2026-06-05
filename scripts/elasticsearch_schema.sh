#!/bin/bash
# Elasticsearch stores documents (like JSON objects) in an "index"
# (similar to a table in SQL). We define mappings so ES knows
# which fields are searchable text vs exact values vs numbers.

# Jobs index — used for keyword search: title, location, salary, skills
curl -s -X PUT "http://localhost:9200/jobs" -H 'Content-Type: application/json' -d '{
  "mappings": {
    "properties": {
      "job_id":           { "type": "keyword" },
      "title":            { "type": "text", "analyzer": "english" },
      "company":          { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
      "location":         { "type": "keyword" },
      "job_type":         { "type": "keyword" },
      "experience_level": { "type": "keyword" },
      "salary_min":       { "type": "integer" },
      "salary_max":       { "type": "integer" },
      "description":      { "type": "text", "analyzer": "english" },
      "skills":           { "type": "keyword" },
      "created_at":       { "type": "date" }
    }
  }
}' && echo ""

# Profiles index — used for people search by name, skills, location
curl -s -X PUT "http://localhost:9200/profiles" -H 'Content-Type: application/json' -d '{
  "mappings": {
    "properties": {
      "user_id":   { "type": "keyword" },
      "name":      { "type": "text", "analyzer": "english" },
      "headline":  { "type": "text" },
      "location":  { "type": "keyword" },
      "skills":    { "type": "keyword" }
    }
  }
}' && echo ""

echo "Elasticsearch indexes created"
