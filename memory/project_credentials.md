---
name: project-credentials
description: "All user IDs, passwords, ports, and DB credentials for Nexus testing"
metadata: 
  node_type: memory
  type: project
  originSessionId: 833a7873-3b6c-46a0-aeec-27d18294c065
---

## Test Users (all passwords: nexus123)
| Name | Email | user_id |
|------|-------|---------|
| Sahil Pal | sahil@nexus.com | e51c8d53-ff90-4e9f-9c28-b7e9c7623d47 |
| Rahul Sharma | rahul@nexus.com | 2abee1c1-2c9c-4293-a59f-a28f24a073c8 |
| Priya Singh | priya@nexus.com | c2366e2c-4dd3-4fc2-aab2-0c9dbd65b452 |
| Arjun Mehta | arjun@nexus.com | 0a3e4548-e29f-44cc-9cc5-71bdb007b52d |
| Rakhi Dalal | rakhi@nexus.com | ba70578e-a605-416d-8cb4-568605315326 |

Passwords were reset via direct DB UPDATE (bcrypt hash) — all set to nexus123.

## Service Ports
- profile-svc: :4001
- network-svc: :4002
- jobs-svc: :4003
- feed-svc: :4004
- Apollo Router (not built yet): :4000

## Infrastructure Credentials
- PostgreSQL: localhost:5432, user: nexus, password: nexus, db: nexus
- Neo4j: bolt://localhost:7687, user: neo4j, password: nexuspassword, browser: localhost:7474
- Kafka: localhost:9092
- Redis: localhost:6379
- Elasticsearch: localhost:9200

## Neo4j Connections Established
- Sahil ↔ Rahul: CONNECTED_TO (bidirectional, accepted)
- Rahul → Priya: PENDING_REQUEST (not accepted yet)

## Start Commands
```
cd /Users/sahil/Documents/Nexus && docker compose up -d
cd /Users/sahil/Documents/Nexus/profile-svc && go run .
cd /Users/sahil/Documents/Nexus/network-svc && go run .
cd /Users/sahil/Documents/Nexus/jobs-svc && go run .
cd /Users/sahil/Documents/Nexus/feed-svc && go run .
```

## DB Debug Commands
```
docker exec nexus-postgres-1 psql -U nexus -d nexus -c "SELECT user_id, name, email FROM users;"
docker exec nexus-neo4j-1 cypher-shell -u neo4j -p nexuspassword "MATCH (a)-[r]->(b) RETURN a.name, type(r), b.name"
docker exec nexus-redis-1 redis-cli ZREVRANGE "feed:e51c8d53-ff90-4e9f-9c28-b7e9c7623d47" 0 -1
curl -s http://localhost:9200/jobs/_count
```

**Why:** Credentials and user IDs are needed every testing session and are not derivable from code.
**How to apply:** Use these whenever Sahil asks to test something or needs to login.
