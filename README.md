# Nexus — Talent Network for Technology Professionals

A LinkedIn-inspired platform built with Go microservices, GraphQL Federation, and event-driven architecture. Designed to demonstrate production-grade backend engineering patterns.

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
   │            │ │  :4003  │ │        │ │         │ │          │
   └─────┬──────┘ └────┬────┘ └───┬────┘ └────┬────┘ └────┬─────┘
         │             │          │            │           │
         ▼             ▼          ▼            ▼           ▼
      Postgres       Neo4j    Postgres      Postgres  Elasticsearch
      Neo4j          Neo4j    Elasticsearch  Redis
      Kafka ──────────────────► Kafka ◄──────┘
                                  │
                              feed-svc
                           consumes events
```

### Data Stores
| Store | Used By | Purpose |
|-------|---------|---------|
| PostgreSQL | profile, jobs, feed | Primary data storage |
| Neo4j | profile, network | Social graph (connections, skills) |
| Redis | feed | Feed sorted sets (newest-first) |
| Elasticsearch | jobs, search | Full-text search across jobs + users |
| pgvector | jobs | Vector embeddings for semantic search |
| Kafka | profile → network, jobs → feed | Async event streaming |

---

## Services

### profile-svc (`:4001`)
User registration, authentication, and profile management.
- JWT auth (HS256) — 15min access token, 7-day refresh token
- bcrypt password hashing (cost 12)
- Dual writes: PostgreSQL (relational) + Neo4j (graph node)
- Publishes `user_created` events to Kafka
- Indexes users to Elasticsearch on register

### network-svc (`:4003`)
Professional connections via graph traversal.
- Send/accept connection requests → Neo4j `PENDING_REQUEST` / `CONNECTED_TO` edges
- "People You May Know" — 2-hop Cypher query (friends of friends)
- Consumes `user_created` Kafka events to mirror Person nodes
- Publishes `connection_accepted` events to Kafka

### jobs-svc (`:4002`)
Job posting with dual-mode search.
- Keyword search via Elasticsearch (`multi_match` across title, company, description)
- **Semantic search** via OpenAI `text-embedding-3-small` + pgvector (`<->` L2 distance, HNSW index)
- Publishes `job_posted` events to Kafka

### feed-svc (`:4004`)
Real-time activity feed.
- Posts stored in PostgreSQL, feed order maintained in Redis sorted sets (score = timestamp)
- Consumes `job_posted` Kafka events — distributes job to all users' feeds automatically
- Feed returns mixed content: user posts + job listings
- Publishes `post_created` events to Kafka

### search-svc (`:4005`)
Unified search across the platform.
- Single `search(query)` returns `JobResult | UserResult` (GraphQL union type)
- Queries Elasticsearch `jobs` and `users` indices in parallel
- **Neo4j proximity boost** — jobs at companies where connections work ranked +20pts (1st degree) / +10pts (2nd degree)

---

## Tech Stack

| Technology | Version | Role |
|-----------|---------|------|
| Go | 1.22 | All microservices |
| GraphQL (gqlgen) | v0.17 | API layer, code generation |
| Apollo Federation | v2 | Schema composition + gateway |
| Apollo Router | v2.15 | GraphQL gateway |
| PostgreSQL + pgvector | pg16 | Relational DB + vector search |
| Neo4j | 5.x | Graph database |
| Apache Kafka | 7.6 | Event streaming |
| Redis | latest | Feed cache |
| Elasticsearch | 8.x | Full-text search |
| OpenAI API | — | Text embeddings (semantic search) |
| Docker | — | Multi-stage builds, final image under 11MB |
| GitHub Actions | — | CI — build + vet on every push |
| k6 | — | Load testing |

---

## Getting Started

### Prerequisites
- Go 1.22+
- Docker + Docker Compose
- OpenAI API key

### 1. Clone and configure
```bash
git clone https://github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals
cd Nexus
cp .env.example .env
# Edit .env and add your OPENAI_API_KEY
```

### 2. Start all services
```bash
# Start infrastructure + all 5 Go services via Docker
docker compose up -d

# Start Apollo Router (unified GraphQL gateway)
APOLLO_ELV2_LICENSE=accept ./router/router --config router.yaml --supergraph supergraph.graphql
```

### 4. Access the API
- **Unified GraphQL API:** http://localhost:4000 (Apollo Sandbox)
- Individual services: :4001 through :4005

---

## Key API Operations

### Register & Login
```graphql
mutation {
  register(input: { email: "you@example.com", password: "pass", name: "Your Name", location: "City" }) {
    access_token
    user { user_id }
  }
}
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

### Semantic Job Search (AI-powered)
```graphql
query {
  semanticSearchJobs(query: "machine learning infrastructure") {
    title company location
  }
}
```

### Unified Search
```graphql
query {
  search(query: "engineer bangalore") {
    __typename
    ... on JobResult { title company }
    ... on UserResult { name headline }
  }
}
```

### Activity Feed
```graphql
query {
  getFeed { post_id user_id content created_at }
}
```

---

## Kafka Events

| Event | Published by | Consumed by | Payload |
|---|---|---|---|
| `user_created` | profile-svc | network-svc | user_id, name, location |
| `job_posted` | jobs-svc | feed-svc | job_id, title, company, location |
| `connection_accepted` | network-svc | — | from_user_id, to_user_id |
| `post_created` | feed-svc | — | post_id, user_id, content |

---

## Load Test Results (k6)

Tested against all 5 services running in Docker on a local machine.

| Scenario | Virtual Users | p95 Latency | Throughput | Pass Rate |
|---|---|---|---|---|
| List Jobs (jobs-svc) | 100 | **19.2ms** | 197 req/s | 100% |
| Unified Search (search-svc) | 50 | **48.5ms** | 96 req/s | 100% |
| Get Feed (feed-svc) | 50 | **21.3ms** | 98 req/s | 100% |

Run the tests yourself:
```bash
# Install k6
brew install k6

# Get a token and run
TOKEN=$(curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation{login(input:{email:\"you@example.com\",password:\"yourpass\"}){access_token}}"}' \
  | jq -r '.data.login.access_token')

k6 run -e TOKEN=$TOKEN k6/profile.js
k6 run -e TOKEN=$TOKEN k6/jobs.js
k6 run -e TOKEN=$TOKEN k6/feed.js
```

---

## CI

GitHub Actions runs on every push to `main` — builds and vets all 5 services.

[![CI](https://github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/actions/workflows/ci.yml/badge.svg)](https://github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/actions/workflows/ci.yml)

---

## Project Structure

```
Nexus/
├── profile-svc/       # Auth, user profiles, skills
├── network-svc/       # Connections, graph traversal
├── jobs-svc/          # Job listings, search, embeddings
├── feed-svc/          # Activity feed
├── search-svc/        # Unified search
├── router/            # Apollo Router binary
├── k6/                # Load test scripts
├── .github/workflows/ # GitHub Actions CI
├── supergraph.yaml    # Federation config
├── supergraph.graphql # Composed schema (auto-generated)
├── router.yaml        # Router config
├── docker-compose.yml # Infrastructure + all 5 services
└── go.work            # Go workspace
```

---

## Author

**Sahil Pal** — Final year, IIITM Gwalior  
[GitHub](https://github.com/sahilpal) · [LinkedIn](https://linkedin.com/in/sahilpal)
