---
name: project-nexus-state
description: "Current build state of Nexus — what is built, what works, what is next"
metadata:
  type: project
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
- register (bcrypt cost 12, Postgres insert, Neo4j Person node, Kafka user_created, ES user index)
- login (bcrypt verify, JWT access 15min + refresh 7days)
- updateProfile (COALESCE partial updates)
- addSkill (upsert Postgres + Neo4j HAS_SKILL edge)
- getProfile, me queries
- Federation enabled (gqlgen v2, @key on User.user_id)

### Day 8 — network-svc (:4002)
- sendConnectionRequest → PENDING_REQUEST edge in Neo4j
- acceptConnection → deletes PENDING, creates bidirectional CONNECTED_TO
- getPeopleYouMayKnow → 2-hop Cypher traversal, returns friends-of-friends
- Kafka consumer: listens to user_created, creates Person node in Neo4j
- Federation enabled (@key on PersonSuggestion.user_id)

### Day 9 — jobs-svc (:4003)
- postJob → Postgres + Elasticsearch index + Kafka job_posted event + async OpenAI embedding
- getJob, listJobs (Postgres)
- searchJobs → Elasticsearch multi_match (title^3, company^2, description, location)
- semanticSearchJobs → OpenAI text-embedding-3-small + pgvector <-> similarity search
- Federation enabled (@key on Job.job_id)

### Day 10 — feed-svc (:4004)
- createPost → Postgres posts table (column: author_id) + Redis sorted set feed:{userID}
- getFeed → Redis ZRevRange (newest first) → Postgres fetch + job entries
- Kafka consumer: listens to job_posted, pushes job:{jobID} to ALL users' feeds
- Feed items: posts show content, jobs show "[JOB] title at company — location"
- Federation enabled (@key on Post.post_id)

### Day 11 — GraphQL Federation (COMPLETE)
- rover CLI at ~/.rover/bin/rover, Apollo Router at ./router/router (v2.15.0)
- supergraph.yaml — lists all 5 subgraph URLs
- router.yaml — listen :4000, introspection, sandbox, CORS, Authorization header propagation
- Unified API at localhost:4000

### Day 12 — Semantic Search (COMPLETE)
- jobs-svc/internal/embedding/embed.go — OpenAI embeddings + pgvector storage + cosine search
- job_embeddings table (vector(1536), HNSW index) created on startup
- semanticSearchJobs(query) resolver in jobs-svc

### Day 13 — search-svc (:4005) (COMPLETE)
- Unified search: search(query) → jobs + people in one response
- GraphQL union type: SearchResult = JobResult | UserResult
- Elasticsearch: jobs index + users index
- profile-svc indexes users to ES on register
- Added to supergraph federation

## STARTING ALL SERVICES

```
Docker:      docker compose up -d
Terminal 1:  cd /Users/sahil/Documents/Nexus && export $(cat .env | grep -v '#' | xargs) && cd profile-svc && go run .
Terminal 2:  cd /Users/sahil/Documents/Nexus/network-svc && go run .
Terminal 3:  cd /Users/sahil/Documents/Nexus && export $(cat .env | grep -v '#' | xargs) && cd jobs-svc && go run .
Terminal 4:  cd /Users/sahil/Documents/Nexus/feed-svc && go run .
Terminal 5:  cd /Users/sahil/Documents/Nexus/search-svc && go run .
Terminal 6:  cd /Users/sahil/Documents/Nexus && APOLLO_ELV2_LICENSE=accept ./router/router --config router.yaml --supergraph supergraph.graphql
```

Recompose (all 5 services running):
```
APOLLO_ELV2_LICENSE=accept ~/.rover/bin/rover supergraph compose --config supergraph.yaml 2>/dev/null > supergraph.graphql
```

## NEXT TO BUILD
- Day 14: Polish, README, resume bullets

## KNOWN ISSUES / FIXES
- posts table column is author_id (NOT user_id)
- gqlgen regenerate moves helper functions into WARNING comment block — manually restore
- JWT token must be ONE line in playground headers
- supergraph.graphql needs 2>/dev/null when generating
- Apollo Router v2: listen under supergraph.listen, sandbox needs introspection: true
- Apollo Router does NOT forward Authorization by default — headers.all.request propagate in router.yaml
- JWT_SECRET must match across all services — use nexus-dev-secret-change-in-production (default)
- OpenAI key gets auto-revoked if shared publicly — always use .env
- profile-svc and jobs-svc need .env exported before starting
