---
name: remaining-days
description: "Remaining build plan Days 15-28 for Nexus — what to build next in order"
metadata:
  type: project
---

Remaining work from the original 28-day plan. Pick up from Day 15.

## Day 15 — Jobs Checkpoint
- Post 30 varied test jobs (different titles, companies, locations, descriptions)
- Compare keyword search (`searchJobs`) vs semantic search (`semanticSearchJobs`)
- Confirm semantic results are more relevant for meaning-based queries
- Example: keyword "Go developer" vs semantic "backend engineer interested in distributed systems"

## Day 16 — Neo4j Proximity Boost in Search
- search-svc `search(query)` currently returns ES results flat
- Add Neo4j re-ranking: jobs at companies where searcher has 1st/2nd degree connections get boosted
- Flow: ES results → check Neo4j for connection proximity → re-rank → return
- This is the "three data sources, one ranked result" from the original plan

## Day 17 — Kafka Event Audit
- Verify all events work end-to-end:
  - `user_created` → network-svc creates Person node in Neo4j ✅
  - `job_posted` → feed-svc pushes to all users' feeds ✅
  - `job_posted` → jobs-svc indexes in ES ✅
  - `connection_accepted` → NOT yet implemented (missing Kafka event)
  - `post_created` → NOT yet published to Kafka (feed only uses Redis directly)
- Add missing `connection_accepted` Kafka publish in network-svc
- Add `post_created` Kafka event in feed-svc (for future follower fanout)

## Day 18 — Dockerfiles
- Write multi-stage Dockerfiles for all 5 services
- Stage 1: golang:1.22 builder — compile binary
- Stage 2: scratch or alpine — copy binary only, no source code
- Target: final image under 25MB each
- Add to docker-compose.yml so all services can be started with one command

## Day 19 — GitHub Actions CI
- `.github/workflows/ci.yml`
- On every push to main: build all Go binaries, run `go test ./...`
- Run `rover lint` on all subgraph schemas
- Badge in README showing CI status

## Day 20 — k6 Load Test
- Install k6: `brew install k6`
- Write `k6/load_test.js` with 3 scenarios:
  - 100 VUs doing profile views (getProfile)
  - 50 VUs doing job searches (searchJobs)
  - 50 VUs loading feeds (getFeed)
- Run for 5 minutes, record p95 latency for each scenario
- These numbers fill in the resume bullets: [X]ms, [Y]ms, [Z]ms

## Day 21 — Skill Synonym Normalisation
- Problem: "Golang" and "Go" create duplicate Neo4j skill nodes
- Add `skills_canonical` table in Postgres: synonym → canonical_id mapping
- Every addSkill call normalises through this table first
- Example: "Golang" → "Go", "JS" → "JavaScript", "k8s" → "Kubernetes"

## Day 22 — Cursor-based Feed Pagination
- Current getFeed returns top 50 items, no pagination
- Add cursor to getFeed: client sends last-seen Redis score
- Server returns only posts newer than that score
- Prevents duplicate/missing posts when new items inserted between pages
- This is how Twitter/Instagram implement feeds

## Day 23 — GraphQL Query Complexity Limits
- Protect against deeply nested queries that could scan Neo4j 10 hops deep
- Configure query complexity in Apollo Router or gqlgen
- Assign cost per field, reject queries above total cost limit
- This is how GitHub and Shopify protect their GraphQL APIs

## Day 24 — Graph-aware Recommendation Engine
- "Recommended connections" using Jaccard similarity in Neo4j
- Two users are similar if they share many mutual connections relative to total connections
- Cypher query: count shared CONNECTED_TO neighbours
- Run as background job (goroutine on timer), cache results in Redis
- Add `getRecommendedConnections` query to network-svc

## Day 25 — GraphQL Subscriptions (Real-time Notifications)
- WebSocket-based subscriptions via Apollo Router
- Events: connection accepted, new post from connection, job matching skills
- Add subscription type to relevant subgraph schemas
- Apollo Router v2 supports subscriptions

## Day 26 — Kafka Event Wiring Polish
- Full audit of all Kafka flows after Days 17-25 changes
- Ensure connection_accepted → Neo4j edge + feed update
- Ensure post_created → follower fanout via Redis
- Test each flow manually end-to-end

## Day 27 — Deploy + Tag v1.0.0
- Production docker-compose.yml with all 5 services + infrastructure
- Tag v1.0.0 on GitHub
- Add repo link to resume and LinkedIn

## Day 28 — Final Polish
- Update README with k6 numbers from Day 20
- Final architecture diagram update
- Verify all resume bullets have real numbers filled in
- Push everything, tag release

---

**How to apply:** Start from Day 15 and work in order. Each day builds on the previous.
**Why:** This is the original 28-day plan — completing it gives real performance numbers for resume bullets.
