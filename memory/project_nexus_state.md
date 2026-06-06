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
- PostgreSQL :5432 (nexus/nexus, db: nexus) — pgvector/pgvector:pg16
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
- Federation enabled (gqlgen v2, @key on User)

### Day 8 — network-svc (:4002)
- sendConnectionRequest → PENDING_REQUEST edge in Neo4j
- acceptConnection → deletes PENDING, creates bidirectional CONNECTED_TO
- getPeopleYouMayKnow → 2-hop Cypher traversal, returns friends-of-friends
- Kafka consumer: listens to user_created, creates Person node in Neo4j
- Federation enabled (@key on PersonSuggestion)

### Days 9 — jobs-svc (:4003)
- postJob → Postgres + Elasticsearch index + Kafka job_posted event
- getJob, listJobs (Postgres)
- searchJobs → Elasticsearch multi_match (title^3, company^2, description, location)
- Federation enabled (@key on Job)

### Day 10 — feed-svc (:4004)
- createPost → Postgres posts table (column: author_id) + Redis sorted set feed:{userID}
- getFeed → Redis ZRevRange (newest first) → Postgres fetch + job entries
- Kafka consumer: listens to job_posted, pushes job:{jobID} to ALL users' feeds
- Feed items: posts show content, jobs show "[JOB] title at company — location"
- Federation enabled (@key on Post)

### Day 11 Phase 1 — Federation prep
- All 4 services have federation: version: 2 in gqlgen.yml
- All 4 services regenerated with federation.go
- @key directives on: User (user_id), PersonSuggestion (user_id), Job (job_id), Post (post_id)

## NEXT TO BUILD
- Day 11 Phase 2: Apollo Router setup (router.yaml + supergraph) → unified :4000
- Day 11 Phase 3: Test cross-service queries through :4000
- Day 12: pgvector semantic job search with embeddings

## KNOWN ISSUES / FIXES
- posts table column is author_id (NOT user_id) — fixed in feed-svc db layer
- gqlgen regenerate moves helper functions (dbUserToModel, dbJobToModel) into WARNING comments — must manually uncomment after every regenerate
- JWT token in GraphQL Playground headers must be on ONE line, no line breaks

**Why:** These issues caused bugs during build and will recur if gqlgen is re-run.
**How to apply:** After any `go run github.com/99designs/gqlgen generate`, check if helpers got commented out and restore them.
