# Owncast Monitoring Stack

This setup includes Prometheus for metrics collection and Grafana for visualization of Owncast streaming metrics.

## Services

### Owncast
- **URL**: http://localhost:8080
- **Admin**: http://localhost:8080/admin
- **Credentials**: admin / abc123

### Prometheus
- **URL**: http://localhost:9090
- **Metrics Endpoint**: http://localhost:8080/api/admin/prometheus
- **Scrape Interval**: 30 seconds

### Grafana
- **URL**: http://localhost:3000
- **Credentials**: admin / admin
- **Default Dashboard**: Owncast Overview

## Available Metrics

### Raw Metrics
- `owncast_instance_active_viewer_count` - Current number of viewers
- `owncast_instance_active_chat_client_count` - Connected chat clients
- `owncast_instance_total_chat_users` - Total chat users
- `owncast_instance_current_chat_message_count` - Current chat messages
- `owncast_instance_playback_error_count` - Playback errors
- `owncast_instance_cpu_usage` - CPU usage percentage

### Recording Rules (Pre-computed)
- `owncast:viewers:avg_1h` - Average viewers over 1 hour
- `owncast:viewers:max_1h` - Peak viewers in last hour
- `owncast:viewers:max_24h` - Peak viewers in last 24 hours
- `owncast:chat:engagement_ratio` - Chat/viewer ratio
- `owncast:chat:total_engagement` - Total viewers + chat clients
- `owncast:cpu:avg_5m` - 5-minute CPU average
- `owncast:errors:rate_5m` - Playback error rate
- `owncast:stream:is_live` - Boolean if stream has viewers

### Alerts
- **OwncastHighCPU** - Warns when CPU > 80% for 5 minutes
- **OwncastCriticalCPU** - Critical alert when CPU > 95% for 2 minutes
- **OwncastHighPlaybackErrors** - Warns when errors > 10 for 2 minutes
- **OwncastLowChatEngagement** - Info alert when < 5% of viewers chatting (10+ viewers)

## Quick Start

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Stop and remove volumes (reset everything)
docker compose down -v
```

## Accessing Services

1. **Start streaming** to Owncast at `rtmp://localhost:1935/live` with stream key from admin panel
2. **View metrics** in Prometheus at http://localhost:9090
3. **View dashboard** in Grafana at http://localhost:3000
   - Username: `admin`
   - Password: `admin`
   - Dashboard: "Owncast Overview"

## Grafana Dashboard

The "Owncast Overview" dashboard includes:
- Current viewer count with statistics
- Active chat clients
- CPU usage with threshold alerts
- Playback error count
- Viewer trends over time
- CPU usage trends
- Chat activity metrics
- Chat engagement ratio

## Customization

### Prometheus Configuration
Edit `prometheus-config-example.yml` to customize:
- Scrape intervals
- Authentication credentials
- Target labels

### Prometheus Rules
Edit `prometheus-rules.yml` to customize:
- Recording rules
- Alert thresholds
- Alert durations

### Grafana Dashboards
- Access Grafana at http://localhost:3000
- Navigate to dashboard and click "Dashboard settings" (gear icon)
- Modify panels, add queries, change visualizations
- Save changes

## Files Structure

```
.
├── compose.yml                           # Docker Compose configuration
├── prometheus-config-example.yml         # Prometheus configuration
├── prometheus-rules.yml                  # Recording rules and alerts
└── grafana/
    ├── dashboards/
    │   └── owncast-overview.json        # Owncast dashboard
    └── provisioning/
        ├── datasources/
        │   └── prometheus.yml           # Prometheus datasource config
        └── dashboards/
            └── owncast.yml              # Dashboard provisioning config
```

## Troubleshooting

### Check service status
```bash
docker compose ps
```

### View service logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f grafana
docker compose logs -f prometheus
docker compose logs -f owncast
```

### Verify Prometheus is scraping
```bash
curl http://localhost:9090/api/v1/targets
```

### Verify Grafana datasource
```bash
curl -u admin:admin http://localhost:3000/api/datasources
```

### Reset everything
```bash
docker compose down -v
docker compose up -d
```
