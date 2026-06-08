---
name: project-nexus-state
description: "Current build state of Nexus — what is built, what works, what is next"
metadata: 
  node_type: memory
  type: project
  originSessionId: 833a7873-3b6c-46a0-aeec-27d18294c065
---

Nexus — Talent Network for Technology Professionals. LinkedIn-like platform built with Go microservices + GraphQL Federation.

**Repo:** github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals
**Working dir:** /Users/sahil/Documents/Nexus
**Go workspace:** go.work links all 5 services

## What is BUILT and WORKING

### Days 1-2 — Infrastructure (Docker Compose)
- PostgreSQL :5432 (nexus/nexus123, db: nexus) — pgvector/pgvector:pg16
- Neo4j :7687 bolt, :7474 browser (neo4j/nexuspassword)
- Kafka :9092 (confluentinc 7.6.1) + Zookeeper
- Redis :6379
- Elasticsearch :9200

### Days 3-7 — profile-svc (:4001)
- register (bcrypt cost 12, Postgres insert, Neo4j Person node, Kafka user_created)
- login (bcrypt verify, JWT access 15min + refresh 7days)
- updateProfile (COALESCE partial updates)
- addSkill (upsert Postgres + Neo4j HAS_SKILL edge)
- getProfile, me queries
- Federation enabled (gqlgen v2, @key on User.user_id)

### Day 8 — network-svc (:4003)
- sendConnectionRequest → PENDING_REQUEST edge in Neo4j
- acceptConnection → deletes PENDING, creates bidirectional CONNECTED_TO
- getPeopleYouMayKnow → 2-hop Cypher traversal, returns friends-of-friends
- Kafka consumer: listens to user_created, creates Person node in Neo4j
- Federation enabled (@key on PersonSuggestion.user_id)

### Day 9 — jobs-svc (:4002)
- postJob → Postgres + Elasticsearch index + Kafka job_posted event + async OpenAI embedding
- getJob, listJobs (Postgres)
- searchJobs → Elasticsearch multi_match (title^3, company^2, description, location)
- semanticSearchJobs → pgvector HNSW index, OpenAI text-embedding-3-small, <-> L2 distance
- job_embeddings table with HNSW index in Postgres
- Federation enabled (@key on Job.job_id)

### Day 10 — feed-svc (:4004)
- createPost → Postgres posts table + Redis sorted set feed:{userID}
- getFeed → Redis ZRevRange (newest first) → Postgres fetch + job entries
- Kafka consumer: listens to job_posted, pushes job:{jobID} to ALL users' feeds
- Federation enabled (@key on Post.post_id)

### Day 11 — GraphQL Federation
- Apollo Router v2.15.0 at :4000
- supergraph.yaml lists all 5 subgraph URLs
- router.yaml has `headers.all.request.propagate: Authorization` (CRITICAL — JWT forwarding)
- Apollo Router endpoint is `/` not `/graphql` — curl must hit `http://localhost:4000/`

**To recompose supergraph:**
```
APOLLO_ELV2_LICENSE=accept ~/.rover/bin/rover supergraph compose --config supergraph.yaml 2>/dev/null > supergraph.graphql
```
**To start router:**
```
APOLLO_ELV2_LICENSE=accept /Users/sahil/Documents/Nexus/router/router --config router.yaml --supergraph supergraph.graphql
```

### Day 12 — pgvector Semantic Search (COMPLETE)
- jobs-svc embeds job text async on postJob via OpenAI text-embedding-3-small
- semanticSearchJobs(query) → embed query → vector similarity search → return jobs
- OPENAI_API_KEY loaded from .env (never hardcode)
- backfill script: scripts/backfill_embeddings/main.go — embeds jobs missing embeddings

### Day 13 — search-svc (:4005) (COMPLETE)
- Unified search across jobs + users via Elasticsearch
- GraphQL union type: SearchResult = JobResult | UserResult
- search(query) returns mixed jobs + users ranked by ES relevance
- profile-svc indexes users in ES on register (async)
- jobs-svc indexes jobs in ES on postJob (sync)

### Day 14 — Polish (COMPLETE)
- README.md with ASCII architecture diagram
- diagrams/architecture.py updated for all 5 services + all data stores
- All entity resolver panics fixed (FindUserByUserID, FindJobByJobID, FindPostByPostID, FindPersonSuggestionByUserID)
- .env file with all secrets (gitignored), .env.example for repo

### Day 15 — Jobs Checkpoint (COMPLETE)
- 30 diverse test jobs seeded via scripts/seed_jobs.sh (uses jq -n for safe JSON, hits `/` not `/graphql`)
- Confirmed: keyword search returns [] for "golang", semantic search returns relevant results
- backfill_embeddings script re-embeds any jobs missing vectors
- Valid enum values: job_type = FULL_TIME|PART_TIME|CONTRACT|INTERNSHIP; experience_level = INTERN|JUNIOR|MID|SENIOR|LEAD

### Day 16 — Neo4j Proximity Boost (COMPLETE)
- search-svc now queries Neo4j for connected companies (1st/2nd degree)
- Jobs at companies where connections work are re-ranked: +20 pts (1st degree), +10 pts (2nd degree)
- search-svc/internal/proximity/proximity.go — ConnectedCompanies() Cypher query
- Gracefully degrades: if Neo4j down or user unauthenticated, returns unranked ES results
- Tested: Sahil → Rahul (Grab) → "Android Engineer @ Grab" ranked first in mobile search

## NEXT TO BUILD
- Day 17: Kafka event audit — add connection_accepted publish in network-svc, post_created in feed-svc
- Day 18: Multi-stage Dockerfiles for all 5 services
- Day 19: GitHub Actions CI
- Day 20: k6 load test (record p95 latency numbers for resume)

## KNOWN ISSUES / FIXES
- Apollo Router endpoint is `/` not `/graphql` — seed scripts must use `http://localhost:4000/`
- JWT tokens expire in 15min — get fresh token before running scripts
- gqlgen regenerate moves helper functions into WARNING comments — manually restore after regenerate
- posts table column is author_id NOT user_id in feed-svc
- supergraph.graphql must be generated with 2>/dev/null — rover logs pollute the file
- Neo4j password: nexuspassword (not nexus1234)
- search-svc Neo4j default password set to nexuspassword in main.go

**Why:** These issues caused bugs during build and will recur.
**How to apply:** Check these whenever touching federation, feed-svc, or running scripts.
