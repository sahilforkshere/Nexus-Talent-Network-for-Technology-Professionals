import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib.patches import FancyBboxPatch

fig, ax = plt.subplots(1, 1, figsize=(26, 18))
ax.set_xlim(0, 26)
ax.set_ylim(0, 18)
ax.axis('off')
fig.patch.set_facecolor('#0f1117')
ax.set_facecolor('#0f1117')

C = {
    'client':       '#1e3a5f',
    'gateway':      '#1a3a5f',
    'svc':          '#1a472a',
    'db':           '#4a1942',
    'kafka':        '#7a3b00',
    'es':           '#1a3a4a',
    'redis':        '#3a1a1a',
    'neo4j':        '#2a3a2a',
    'pgvector':     '#3a3a1a',
    'openai':       '#2a1a4a',
    'border':       '#ffffff',
    'text':         '#ffffff',
    'arrow':        '#aaaaaa',
    'kafka_arrow':  '#f5a623',
    'green_arrow':  '#50fa7b',
    'blue_arrow':   '#8be9fd',
    'purple_arrow': '#bd93f9',
}

def box(ax, x, y, w, h, color, label, sublabel='', fontsize=10):
    rect = FancyBboxPatch((x, y), w, h,
                          boxstyle="round,pad=0.1",
                          facecolor=color, edgecolor='#ffffff',
                          linewidth=1.5, zorder=3)
    ax.add_patch(rect)
    cy = y + h/2 + (0.2 if sublabel else 0)
    ax.text(x + w/2, cy, label,
            ha='center', va='center', fontsize=fontsize,
            fontweight='bold', color='white', zorder=4)
    if sublabel:
        ax.text(x + w/2, y + h/2 - 0.3, sublabel,
                ha='center', va='center', fontsize=7.5,
                color='#cccccc', zorder=4)

def arrow(ax, x1, y1, x2, y2, color='#aaaaaa', label='', lw=1.5):
    ax.annotate('', xy=(x2, y2), xytext=(x1, y1),
                arrowprops=dict(arrowstyle='->', color=color,
                                lw=lw, connectionstyle='arc3,rad=0.0'),
                zorder=5)
    if label:
        mx, my = (x1+x2)/2, (y1+y2)/2
        ax.text(mx+0.1, my+0.12, label, fontsize=7, color=color,
                ha='center', va='center',
                bbox=dict(boxstyle='round,pad=0.15', facecolor='#1e1e2e',
                          edgecolor='none', alpha=0.85), zorder=6)

# ── TITLE ────────────────────────────────────────────────────────────────────
ax.text(13, 17.4, 'Nexus — Talent Network for Technology Professionals',
        ha='center', va='center', fontsize=18, fontweight='bold', color='white')
ax.text(13, 16.95, 'Go Microservices  ·  GraphQL Federation  ·  Event-Driven Architecture  ·  Days 1–13',
        ha='center', va='center', fontsize=10, color='#aaaaaa')

# ── ROW 1: Client ─────────────────────────────────────────────────────────
box(ax, 9.5, 15.8, 7.0, 0.8, C['client'],
    'Browser / Apollo Sandbox', 'localhost:4000', fontsize=10)

# ── ROW 2: Gateway ───────────────────────────────────────────────────────
box(ax, 9.5, 14.2, 7.0, 1.0, C['gateway'],
    'Apollo Router  (GraphQL Gateway)',
    ':4000  ·  Federation v2  ·  Header propagation  ·  CORS', fontsize=10)
arrow(ax, 13.0, 15.8, 13.0, 15.2, color=C['green_arrow'], label='HTTP')

# ── ROW 3: Services ──────────────────────────────────────────────────────
svc_y = 11.8
svc_w = 4.2
svc_h = 1.8

services = [
    (0.4,  'profile-svc', ':4001\nregister · login\nupdateProfile · addSkill'),
    (5.0,  'network-svc', ':4002\nsendConnectionRequest\nacceptConnection\ngetPeopleYouMayKnow'),
    (9.6,  'jobs-svc',    ':4003\npostJob · listJobs\nsearchJobs (ES)\nsemanticSearchJobs (AI)'),
    (14.2, 'feed-svc',    ':4004\ncreatePost · getFeed\nRedis sorted set\nKafka consumer'),
    (18.8, 'search-svc',  ':4005\nsearch(query)\nJobResult | UserResult\nES multi-index'),
]

svc_centers = []
for x, name, sub in services:
    box(ax, x, svc_y, svc_w, svc_h, C['svc'], name, sub, fontsize=8)
    cx = x + svc_w / 2
    svc_centers.append(cx)
    arrow(ax, 13.0, 14.2, cx, svc_y + svc_h, color=C['blue_arrow'])

# ── ROW 4: Kafka ──────────────────────────────────────────────────────────
kafka_y = 9.8
box(ax, 2.5, kafka_y, 13.0, 0.9, C['kafka'],
    'Apache Kafka  (Event Bus)  —  localhost:9092',
    'Topics:  user_created  ·  job_posted', fontsize=9)

# profile-svc → Kafka
arrow(ax, svc_centers[0], svc_y, 4.5, kafka_y+0.9, color=C['kafka_arrow'], label='user_created')
# jobs-svc → Kafka
arrow(ax, svc_centers[2], svc_y, 11.5, kafka_y+0.9, color=C['kafka_arrow'], label='job_posted')
# Kafka → network-svc
arrow(ax, 6.5, kafka_y, svc_centers[1], svc_y, color=C['kafka_arrow'], label='user_created')
# Kafka → feed-svc
arrow(ax, 12.5, kafka_y, svc_centers[3], svc_y, color=C['kafka_arrow'], label='job_posted')

# ── ROW 5: Databases ─────────────────────────────────────────────────────
db_y = 5.8
db_h = 3.4

box(ax, 0.3, db_y, 4.0, db_h, C['db'],
    'PostgreSQL', ':5432  ·  db: nexus\n\nTables:\nusers · skills · user_skills\njobs · posts\nrefresh_tokens', fontsize=8)

box(ax, 4.7, db_y, 3.8, db_h, C['neo4j'],
    'Neo4j', ':7687  Bolt\n\nNodes:\nPerson · Skill\n\nEdges:\nCONNECTED_TO\nPENDING_REQUEST\nHAS_SKILL', fontsize=8)

box(ax, 8.9, db_y, 3.8, db_h, C['es'],
    'Elasticsearch', ':9200\n\nIndexes:\njobs\n(title · company · desc)\n\nusers\n(name · headline)', fontsize=8)

box(ax, 13.1, db_y, 3.8, db_h, C['redis'],
    'Redis', ':6379\n\nSorted sets:\nfeed:{user_id}\n\nScore = timestamp\nNewest-first\nCap 50 items', fontsize=8)

box(ax, 17.3, db_y, 4.0, db_h, C['pgvector'],
    'pgvector', '(PostgreSQL ext)\n\njob_embeddings\nvector(1536)\nHNSW index\n\nL2 distance <->', fontsize=8)

box(ax, 21.7, db_y, 3.9, db_h, C['openai'],
    'OpenAI API', 'text-embedding\n-3-small\n\n1536 dimensions\n\nCalled on:\n· postJob (async)\n· semanticSearch', fontsize=8)

# service → DB arrows
arrow(ax, svc_centers[0], svc_y, 2.3, db_y+db_h, color=C['arrow'], label='SQL')
arrow(ax, svc_centers[0], svc_y, 6.6, db_y+db_h, color=C['green_arrow'], label='Cypher')
arrow(ax, svc_centers[0], svc_y, 10.8, db_y+db_h, color=C['blue_arrow'], label='index users')
arrow(ax, svc_centers[1], svc_y, 6.6, db_y+db_h, color=C['green_arrow'], label='graph ops')
arrow(ax, svc_centers[2], svc_y, 2.3, db_y+db_h, color=C['arrow'])
arrow(ax, svc_centers[2], svc_y, 10.8, db_y+db_h, color=C['blue_arrow'], label='index jobs')
arrow(ax, svc_centers[2], svc_y, 19.3, db_y+db_h, color=C['purple_arrow'], label='store vec')
arrow(ax, svc_centers[2], svc_y, 23.65, db_y+db_h, color=C['purple_arrow'], label='embed')
arrow(ax, svc_centers[3], svc_y, 2.3, db_y+db_h, color=C['arrow'])
arrow(ax, svc_centers[3], svc_y, 15.0, db_y+db_h, color='#ff6b6b', label='ZADD/ZRANGE')
arrow(ax, svc_centers[4], svc_y, 10.8, db_y+db_h, color=C['blue_arrow'], label='multi-search')

# ── Auth note ─────────────────────────────────────────────────────────────
ax.text(13, 5.2,
        'JWT (HS256)  ·  access 15 min  ·  refresh 7 days  ·  bcrypt cost 12  ·  propagated by Apollo Router to all subgraphs',
        ha='center', va='center', fontsize=8.5, color='#f1fa8c',
        bbox=dict(boxstyle='round,pad=0.3', facecolor='#2a2a00',
                  edgecolor='#f1fa8c', linewidth=1, alpha=0.9))

# ── Legend ────────────────────────────────────────────────────────────────
legend_items = [
    (C['gateway'],      'Apollo Router (gateway)'),
    (C['svc'],          'Microservice (Go + gqlgen)'),
    (C['db'],           'PostgreSQL — relational'),
    (C['neo4j'],        'Neo4j — social graph'),
    (C['es'],           'Elasticsearch — full-text'),
    (C['redis'],        'Redis — feed cache'),
    (C['pgvector'],     'pgvector — vector search'),
    (C['openai'],       'OpenAI — embeddings'),
    (C['kafka'],        'Kafka — event streaming'),
]
lx, ly = 0.2, 4.6
ax.text(lx, ly+0.5, 'Legend', fontsize=9, fontweight='bold', color='white')
for i, (col, lbl) in enumerate(legend_items):
    rect = FancyBboxPatch((lx, ly - i*0.45), 0.28, 0.28,
                          boxstyle='round,pad=0.04',
                          facecolor=col, edgecolor='white', linewidth=0.8, zorder=6)
    ax.add_patch(rect)
    ax.text(lx+0.42, ly - i*0.45 + 0.14, lbl,
            va='center', fontsize=8, color='white', zorder=6)

arrow_legends = [
    (C['green_arrow'],  'GraphQL / direct'),
    (C['kafka_arrow'],  'Kafka event (async)'),
    (C['blue_arrow'],   'Elasticsearch'),
    (C['purple_arrow'], 'OpenAI / pgvector'),
    ('#ff6b6b',         'Redis'),
]
base = ly - len(legend_items)*0.45 - 0.3
for i, (col, lbl) in enumerate(arrow_legends):
    y = base - i*0.42
    ax.annotate('', xy=(lx+0.28, y+0.14), xytext=(lx, y+0.14),
                arrowprops=dict(arrowstyle='->', color=col, lw=1.5))
    ax.text(lx+0.42, y+0.14, lbl, va='center', fontsize=8, color='white')

plt.tight_layout()
plt.savefig('/Users/sahil/Documents/Nexus/diagrams/architecture.png',
            dpi=180, bbox_inches='tight', facecolor='#0f1117')
print("saved → diagrams/architecture.png")
