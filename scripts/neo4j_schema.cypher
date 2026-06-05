// Neo4j stores data as NODES (dots) connected by RELATIONSHIPS (lines).
// Unlike SQL tables, there are no columns — just properties on nodes.
//
// Our graph has 3 node types:
//   (:Person)  — mirrors a user from Postgres
//   (:Skill)   — a technology like "Go" or "Kafka"
//   (:Company) — an employer
//
// Relationships:
//   (:Person)-[:CONNECTED_TO]->(:Person)   — social connection
//   (:Person)-[:HAS_SKILL]->(:Skill)       — user knows a skill
//   (:Person)-[:WORKED_AT]->(:Company)     — employment history
//
// CONSTRAINTS enforce uniqueness — like UNIQUE in SQL.
// Without these, registering the same user twice would create 2 Person nodes.

CREATE CONSTRAINT person_user_id IF NOT EXISTS
    FOR (p:Person) REQUIRE p.user_id IS UNIQUE;

CREATE CONSTRAINT skill_name IF NOT EXISTS
    FOR (s:Skill) REQUIRE s.name IS UNIQUE;

CREATE CONSTRAINT company_name IF NOT EXISTS
    FOR (c:Company) REQUIRE c.name IS UNIQUE;

// INDEXES speed up lookups by property (like SQL indexes)
CREATE INDEX person_name IF NOT EXISTS FOR (p:Person) ON (p.name);
CREATE INDEX person_location IF NOT EXISTS FOR (p:Person) ON (p.location);
