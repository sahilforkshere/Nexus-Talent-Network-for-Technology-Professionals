# Nexus — Talent Network for Technology Professionals

A LinkedIn-inspired platform built with Go microservices, GraphQL Federation, and event-driven architecture. Designed to demonstrate production-grade backend engineering patterns.

[![CI](https://github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/actions/workflows/ci.yml/badge.svg)](https://github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/actions/workflows/ci.yml)

---

## Architecture

```
                        ┌─────────────────────┐
                        │   Apollo Router      │
                        │   localhost:4000     │
                        │  (GraphQL Gateway)   │
                        └──────────┬──────────┘
                                   │
          ┌────────────┬───────────┼───────────┬────────────┐
          ▼            ▼           ▼           ▼            ▼
   ┌────────────┐ ┌─────────┐ ┌────────┐ ┌─────────┐ ┌──────────┐
   │profile-svc │ │network  │ │jobs-svc│ │feed-svc │ │search-svc│
   │  :4001     │ │  -svc   │ │ :4002  │ │  :4004  │ │  :4005   │
   └─────┬──────┘ └────┬────┘ └───┬────┘ └────┬────┘ └────┬─────┘
         │             │          │            │           │
         ▼             ▼          ▼            ▼           ▼
      Postgres       Neo4j    Postgres      Postgres  Elasticsearch
      Neo4j (node)   (graph)  ES + pgvector  Redis     Neo4j (boost)
         │                        │            │
         └────────────────────────┴────────────┘
                              Kafka
```

### Data Stores
| Store | Used By | Purpose |
|-------|---------|---------|
| PostgreSQL | profile, jobs, feed | Primary relational storage |
| Neo4j | profile, network, search | Social graph — connections, skills, proximity boost |
| Redis | feed | Feed sorted sets (newest-first, O(1) reads) |
| Elasticsearch | jobs, search | Full-text search across jobs + users |
| pgvector | jobs | Vector embeddings for semantic search (HNSW index) |
| Kafka | all services | Async event streaming (4 topics) |

---

## Services

### profile-svc (`:4001`)
User registration, authentication, and profile management.
- JWT auth (HS256) — 15min access token, 7-day refresh token
- bcrypt password hashing (cost 12)
- Dual writes: PostgreSQL + Neo4j Person node
- Publishes `user_created` to Kafka → network-svc creates graph node
- Indexes users in Elasticsearch on register

### network-svc (`:4003`)
Professional connections via Neo4j graph traversal.
- Send/accept connection requests → `PENDING_REQUEST` / `CONNECTED_TO` edges in Neo4j
- "People You May Know" — 2-hop Cypher traversal (friends of friends)
- Consumes `user_created` Kafka events to mirror Person nodes
- Publishes `connection_accepted` to Kafka

### jobs-svc (`:4002`)
Job posting with multi-modal search and graph-based recommendations.
- Keyword search via Elasticsearch (`multi_match` across title, company, description)
- Semantic search via OpenAI `text-embedding-3-small` + pgvector (HNSW index, `<->` L2 distance)
- **`recommendJobs`** — Neo4j 2-hop traversal finds companies where connections work, returns those jobs ranked by degree (1st = score 20, 2nd = score 10)
- Cursor-based pagination on `listJobs` (opaque base64 cursors, max 50/page)
- Publishes `job_posted` to Kafka

### feed-svc (`:4004`)
Real-time activity feed with cursor pagination.
- Posts stored in PostgreSQL, feed order in Redis sorted sets (score = timestamp)
- Consumes `job_posted` Kafka events — distributes job to all users' feeds
- Cursor-based pagination on `getFeed` (max 50 items/page)
- Publishes `post_created` to Kafka

### search-svc (`:4005`)
Unified search with LLM query expansion and proximity ranking.
- `search(query)` returns `JobResult | UserResult` (GraphQL union type)
- **LLM query expansion** via GPT-4o-mini — `"js"` → `"javascript, js, node.js, react, typescript"` before hitting ES
- **Neo4j proximity boost** — jobs at 1st degree connections' companies get +20pts, 2nd degree +10pts
- 3s timeout on LLM call, gracefully falls back to original query

---

## Tech Stack

| Technology | Version | Role |
|-----------|---------|------|
| Go | 1.22+ | All 5 microservices |
| GraphQL (gqlgen) | v0.17.90 | API layer, code generation |
| Apollo Federation | v2 | Schema composition |
| Apollo Router | v2.15.0 | GraphQL gateway (single endpoint) |
| PostgreSQL + pgvector | pg16 | Relational DB + vector search |
| Neo4j | 5.x | Social graph database |
| Apache Kafka | 7.6.1 | Event streaming (4 topics) |
| Redis | 7-alpine | Feed cache (sorted sets) |
| Elasticsearch | 8.13.4 | Full-text search |
| OpenAI API | — | Embeddings (text-embedding-3-small) + query expansion (GPT-4o-mini) |
| Docker | — | Multi-stage builds — all images under 11MB |
| GitHub Actions | — | CI — build + vet all 5 services on every push |
| k6 | — | Load testing |

---

## Getting Started

### Prerequisites
- Docker + Docker Compose
- OpenAI API key
- Apollo Router binary (included in `router/`)

### 1. Clone and configure
```bash
git clone https://github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals
cd Nexus
cp .env.example .env
# Edit .env — set JWT_SECRET and OPENAI_API_KEY
```

### 2. Start all services
```bash
# Start infrastructure + all 5 Go services
docker compose up -d

# Start Apollo Router (unified GraphQL gateway)
APOLLO_ELV2_LICENSE=accept ./router/router --config router.yaml --supergraph supergraph.graphql
```

### 3. Access the API
- **Unified GraphQL API:** http://localhost:4000 (Apollo Sandbox)
- Individual services: :4001 (profile) :4002 (jobs) :4003 (network) :4004 (feed) :4005 (search)

---

## Key API Operations

### Register & Login
```graphql
mutation {
  register(input: {
    email: "you@example.com", password: "pass",
    name: "Your Name", location: "City"
  }) {
    access_token
    user { user_id }
  }
}

mutation {
  login(input: { email: "you@example.com", password: "pass" }) {
    access_token
  }
}
```

### Cursor Pagination
```graphql
# First page
query {
  listJobs(first: 10) {
    edges { cursor node { title company } }
    pageInfo { hasNextPage endCursor }
  }
}

# Next page (paste endCursor from above)
query {
  listJobs(first: 10, after: "endCursorFromPreviousPage") {
    edges { cursor node { title company } }
    pageInfo { hasNextPage endCursor }
  }
}
```

### Graph-Based Job Recommendations
```graphql
query {
  recommendJobs {
    title company job_type experience_level
  }
}
# Returns jobs at companies where your connections work
# 1st degree connections' companies ranked first
```

### Semantic Job Search + LLM Query Expansion
```graphql
query {
  search(query: "js developer") {
    ... on JobResult { title company }
    ... on UserResult { name headline }
  }
}
# GPT-4o-mini expands "js" → "javascript, js, node.js, react, typescript"
# Results re-ranked by Neo4j connection proximity
```

### Post a Job
```graphql
mutation {
  postJob(input: {
    title: "Backend Engineer", company: "Acme", location: "Remote",
    job_type: "FULL_TIME", experience_level: "MID",
    description: "Build Go microservices at scale"
  }) { job_id }
}
```

---

## Kafka Events

| Event | Published by | Consumed by | Effect |
|---|---|---|---|
| `user_created` | profile-svc | network-svc | Creates Person node in Neo4j |
| `job_posted` | jobs-svc | feed-svc | Pushes job to all users' Redis feeds |
| `connection_accepted` | network-svc | — | Published for downstream consumers |
| `post_created` | feed-svc | — | Published for downstream consumers |

---

## Load Test Results (k6)

Tested against all 5 services running in Docker on a local machine (200 total VUs).

| Scenario | Virtual Users | Duration | p95 Latency | Throughput | Pass Rate |
|---|---|---|---|---|---|
| List Jobs | 100 VUs | 30s | **19.2ms** | 197 req/s | 100% |
| Unified Search | 50 VUs | 30s | **48.5ms** | 96 req/s | 100% |
| Get Feed | 50 VUs | 30s | **21.3ms** | 98 req/s | 100% |

```bash
# Run load tests
brew install k6

TOKEN=$(curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation{login(input:{email:\"you@example.com\",password:\"pass\"}){access_token}}"}' \
  | jq -r '.data.login.access_token')

k6 run -e TOKEN=$TOKEN k6/profile.js
k6 run -e TOKEN=$TOKEN k6/jobs.js
k6 run -e TOKEN=$TOKEN k6/feed.js
```

---

## Docker Image Sizes

Multi-stage builds: `golang:1.24-alpine` builder → `alpine:3.19` final.

| Service | Image Size |
|---|---|
| profile-svc | 10.9 MB |
| jobs-svc | 10.3 MB |
| network-svc | 8.1 MB |
| feed-svc | 9.3 MB |
| search-svc | 10.0 MB |

---

## Project Structure

```
Nexus/
├── profile-svc/          # Auth, user profiles, skills
├── network-svc/          # Connections, graph traversal
├── jobs-svc/             # Job listings, search, recommendations
├── feed-svc/             # Activity feed
├── search-svc/           # Unified search + LLM expansion
├── router/               # Apollo Router binary
├── k6/                   # Load test scripts
├── scripts/              # Seed jobs, backfill embeddings
├── .github/workflows/    # GitHub Actions CI
├── supergraph.yaml       # Federation subgraph config
├── supergraph.graphql    # Composed schema (auto-generated)
├── router.yaml           # Router config (JWT header forwarding)
├── docker-compose.yml    # Infrastructure + all 5 services
└── go.work               # Go workspace
```

---

## Resume Bullets

- Built **Nexus**, a professional network with GraphQL Federation across 5 Go microservices using Apollo Router; each subgraph owns its domain independently and deploys without touching the gateway.
- Social graph stored in **Neo4j**; 2nd-degree connection traversal using Cypher pattern matching — powers "People You May Know" and graph-based job recommendations.
- **Semantic job matching** via OpenAI embeddings + pgvector HNSW index combined with Elasticsearch structured filters and Neo4j proximity re-ranking; p95 latency 48.5ms at 50 VUs.
- **LLM query expansion** via GPT-4o-mini expands search terms before hitting Elasticsearch — `"js"` → `"javascript, js, node.js, react, typescript"` with 3s timeout fallback.
- Redis fanout-on-write delivers posts to all users' feeds via Kafka consumer; feed loads at p95 21.3ms.
- Load-tested at 200 concurrent users; all endpoints within p95 latency targets at 0% error rate.
- Multi-stage Docker builds produce sub-11MB images; GitHub Actions CI builds and vets all 5 services on every push.
- Cursor-based pagination (opaque base64 cursors) and query complexity limits (max 50 items/page, `FixedComplexityLimit(200)`) prevent abuse.

---

## Author

**Sahil Pal** — Final year, IIITM Gwalior  
[GitHub](https://github.com/sahilpal) · [LinkedIn](https://linkedin.com/in/sahilpal)
