import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib.patches import FancyBboxPatch, FancyArrowPatch

fig, ax = plt.subplots(1, 1, figsize=(22, 16))
ax.set_xlim(0, 22)
ax.set_ylim(0, 16)
ax.axis('off')
fig.patch.set_facecolor('#0f1117')
ax.set_facecolor('#0f1117')

# ── colour palette ──────────────────────────────────────────────────────────
C = {
    'client':   '#1e3a5f',
    'svc':      '#1a472a',
    'db':       '#4a1942',
    'kafka':    '#7a3b00',
    'es':       '#1a3a4a',
    'redis':    '#3a1a1a',
    'border':   '#ffffff',
    'text':     '#ffffff',
    'arrow':    '#aaaaaa',
    'kafka_arrow': '#f5a623',
    'direct_arrow': '#50fa7b',
    'label_bg': '#1e1e2e',
}

def box(ax, x, y, w, h, color, label, sublabel='', fontsize=11):
    rect = FancyBboxPatch((x, y), w, h,
                          boxstyle="round,pad=0.1",
                          facecolor=color, edgecolor='#ffffff',
                          linewidth=1.5, zorder=3)
    ax.add_patch(rect)
    cy = y + h/2 + (0.18 if sublabel else 0)
    ax.text(x + w/2, cy, label,
            ha='center', va='center', fontsize=fontsize,
            fontweight='bold', color='white', zorder=4)
    if sublabel:
        ax.text(x + w/2, y + h/2 - 0.28, sublabel,
                ha='center', va='center', fontsize=7.5,
                color='#cccccc', zorder=4)

def arrow(ax, x1, y1, x2, y2, color='#aaaaaa', label='', style='->', lw=1.5):
    ax.annotate('', xy=(x2, y2), xytext=(x1, y1),
                arrowprops=dict(arrowstyle=style, color=color,
                                lw=lw, connectionstyle='arc3,rad=0.0'),
                zorder=5)
    if label:
        mx, my = (x1+x2)/2, (y1+y2)/2
        ax.text(mx+0.08, my+0.12, label, fontsize=7, color=color,
                ha='center', va='center',
                bbox=dict(boxstyle='round,pad=0.15', facecolor='#1e1e2e',
                          edgecolor='none', alpha=0.85), zorder=6)

# ════════════════════════════════════════════════════════════════════════════
# TITLE
# ════════════════════════════════════════════════════════════════════════════
ax.text(11, 15.4, 'Nexus — Talent Network Architecture',
        ha='center', va='center', fontsize=17, fontweight='bold',
        color='white')
ax.text(11, 15.0, 'Services built so far  ·  Days 1 – 9',
        ha='center', va='center', fontsize=10, color='#aaaaaa')

# ════════════════════════════════════════════════════════════════════════════
# ROW 1 — Client / Browser
# ════════════════════════════════════════════════════════════════════════════
box(ax,  4.5, 13.2, 4.0, 0.9, C['client'], 'Browser / GraphQL Playground',
    'localhost:4001  |  4002  |  4003', fontsize=9)

# ════════════════════════════════════════════════════════════════════════════
# ROW 2 — Microservices
# ════════════════════════════════════════════════════════════════════════════
svc_y = 11.2
box(ax,  0.5, svc_y, 3.8, 1.3, C['svc'],
    'profile-svc', ':4001\nregister · login · updateProfile · addSkill')
box(ax,  5.1, svc_y, 3.8, 1.3, C['svc'],
    'network-svc', ':4002\nsendConnectionRequest · acceptConnection\ngetPeopleYouMayKnow')
box(ax,  9.7, svc_y, 3.8, 1.3, C['svc'],
    'jobs-svc', ':4003\npostJob · getJob · listJobs · searchJobs')
box(ax, 14.3, svc_y, 3.8, 1.3, '#2a2a4a',
    'feed-svc (soon)', ':4004  (Day 10 — not built yet)')
box(ax, 18.9, svc_y-0.15, 2.6, 1.1, '#2a2a4a',
    'search-svc (soon)', ':4005\n(Day 11)')

# arrows: browser → services
for sx in [2.4, 7.0, 11.6]:
    arrow(ax, 6.5, 13.2, sx, svc_y+1.3, color=C['direct_arrow'], label='GraphQL\nHTTP POST')

# ════════════════════════════════════════════════════════════════════════════
# ROW 3 — Kafka bus
# ════════════════════════════════════════════════════════════════════════════
kafka_y = 8.8
box(ax,  3.0, kafka_y, 10.0, 1.0, C['kafka'],
    'Apache Kafka  (Event Bus)  —  localhost:9092',
    'Topics:  user_created  ·  job_posted  ·  connection_accepted  (future)  ·  post_created  (future)',
    fontsize=9)

# profile-svc → Kafka (publishes user_created)
arrow(ax, 2.4, svc_y, 4.5, kafka_y+1.0,
      color=C['kafka_arrow'], label='user_created')

# jobs-svc → Kafka (publishes job_posted)
arrow(ax, 11.6, svc_y, 10.5, kafka_y+1.0,
      color=C['kafka_arrow'], label='job_posted')

# Kafka → network-svc (consumes user_created → creates Person node)
arrow(ax, 6.5, kafka_y, 7.0, svc_y,
      color=C['kafka_arrow'], label='consumes\nuser_created')

# ════════════════════════════════════════════════════════════════════════════
# ROW 4 — Databases
# ════════════════════════════════════════════════════════════════════════════
db_y = 5.8

box(ax,  0.3, db_y, 4.2, 2.4, C['db'],
    'PostgreSQL',
    ':5432  |  db: nexus\n'
    'Tables:\n'
    'users · skills · user_skills\n'
    'jobs · job_skills · applications\n'
    'posts · refresh_tokens',
    fontsize=8)

box(ax,  5.1, db_y, 3.8, 2.4, '#2a3a2a',
    'Neo4j',
    ':7687  |  Bolt\n'
    'Nodes: Person · Skill · Company\n'
    'Edges:\n'
    'CONNECTED_TO  (bidirectional)\n'
    'PENDING_REQUEST\n'
    'HAS_SKILL',
    fontsize=8)

box(ax,  9.5, db_y, 3.8, 2.4, C['es'],
    'Elasticsearch',
    ':9200\n'
    'Indexes:\n'
    'jobs  →  full-text search\n'
    '         (title · company · desc)\n'
    'profiles  (future)',
    fontsize=8)

box(ax, 13.9, db_y, 3.4, 2.4, C['redis'],
    'Redis',
    ':6379\n'
    'Sorted sets:\n'
    'feed:{user_id}  (future)\n'
    'Caches home feed\n'
    'ordered by timestamp',
    fontsize=8)

box(ax, 17.9, db_y, 3.7, 2.4, '#3a3a1a',
    'pgvector',
    '(PostgreSQL extension)\n'
    'Table: job_embeddings\n'
    'vector(1536)  HNSW index\n'
    'Semantic job search\n'
    '(future — Day 12)',
    fontsize=8)

# service → DB arrows
# profile-svc → Postgres
arrow(ax, 2.4, svc_y, 2.4, db_y+2.4, color='#aaaaaa', label='SQL\nread/write')
# profile-svc → Neo4j
arrow(ax, 2.4, svc_y, 7.0, db_y+2.4, color='#69ff94', label='MERGE\nPerson node')
# network-svc → Neo4j
arrow(ax, 7.0, svc_y, 7.0, db_y+2.4, color='#69ff94', label='Cypher\ngraph ops')
# jobs-svc → Postgres
arrow(ax, 11.6, svc_y, 2.4, db_y+2.4, color='#aaaaaa')
# jobs-svc → Elasticsearch
arrow(ax, 11.6, svc_y, 11.4, db_y+2.4, color='#8be9fd', label='index\n& search')

# ════════════════════════════════════════════════════════════════════════════
# ROW 5 — Auth layer note
# ════════════════════════════════════════════════════════════════════════════
ax.text(11, 5.1,
        'JWT Auth (HS256)  ·  access token 15 min  ·  refresh token 7 days  ·  bcrypt cost 12  ·  same secret across all services',
        ha='center', va='center', fontsize=8.5, color='#f1fa8c',
        bbox=dict(boxstyle='round,pad=0.3', facecolor='#2a2a00',
                  edgecolor='#f1fa8c', linewidth=1, alpha=0.9))

# ════════════════════════════════════════════════════════════════════════════
# LEGEND
# ════════════════════════════════════════════════════════════════════════════
legend_items = [
    (C['svc'],          'Microservice (Go + gqlgen)'),
    (C['db'],           'PostgreSQL — relational facts'),
    ('#2a3a2a',         'Neo4j — social graph'),
    (C['es'],           'Elasticsearch — keyword search'),
    (C['redis'],        'Redis — feed cache'),
    (C['kafka'],        'Kafka — event bus'),
    ('#2a2a4a',         'Not built yet'),
]
lx, ly = 0.3, 4.5
for i, (col, lbl) in enumerate(legend_items):
    rect = FancyBboxPatch((lx, ly - i*0.45), 0.3, 0.3,
                          boxstyle='round,pad=0.05',
                          facecolor=col, edgecolor='white', linewidth=0.8, zorder=6)
    ax.add_patch(rect)
    ax.text(lx+0.45, ly - i*0.45 + 0.15, lbl,
            va='center', fontsize=8, color='white', zorder=6)

ax.text(lx, ly+0.55, 'Legend', fontsize=9, fontweight='bold',
        color='white')

# Arrow legend
ax.annotate('', xy=(lx+0.3, ly-3.65), xytext=(lx, ly-3.65),
            arrowprops=dict(arrowstyle='->', color=C['direct_arrow'], lw=1.5))
ax.text(lx+0.45, ly-3.65, 'GraphQL / direct call', va='center',
        fontsize=8, color='white')

ax.annotate('', xy=(lx+0.3, ly-4.1), xytext=(lx, ly-4.1),
            arrowprops=dict(arrowstyle='->', color=C['kafka_arrow'], lw=1.5))
ax.text(lx+0.45, ly-4.1, 'Kafka event (async)', va='center',
        fontsize=8, color='white')

ax.annotate('', xy=(lx+0.3, ly-4.55), xytext=(lx, ly-4.55),
            arrowprops=dict(arrowstyle='->', color='#69ff94', lw=1.5))
ax.text(lx+0.45, ly-4.55, 'Neo4j Cypher', va='center',
        fontsize=8, color='white')

plt.tight_layout()
plt.savefig('/Users/sahil/Documents/Nexus/diagrams/architecture.png',
            dpi=180, bbox_inches='tight', facecolor='#0f1117')
print("saved → diagrams/architecture.png")
