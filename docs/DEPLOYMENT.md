# Deployment Guide

## Overview

This guide covers deploying the e-commerce platform to production environments, including containerization, cloud deployment, monitoring, and maintenance procedures.

## Production Architecture

### Recommended Infrastructure

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │   Application   │    │    Database     │
│    (Nginx)      │────│     Server      │────│  (PostgreSQL)   │
│                 │    │   (Go App)      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                │
                       ┌─────────────────┐    ┌─────────────────┐
                       │     Cache       │    │   File Storage  │
                       │    (Redis)      │    │  (Cloudinary)   │
                       └─────────────────┘    └─────────────────┘
```

### Environment Requirements

#### Minimum Production Specs

| Component | CPU | Memory | Storage | Network |
|-----------|-----|---------|---------|---------|
| Application Server | 2 vCPU | 4GB RAM | 20GB SSD | 1Gbps |
| Database Server | 2 vCPU | 8GB RAM | 100GB SSD | 1Gbps |
| Cache Server | 1 vCPU | 2GB RAM | 10GB SSD | 1Gbps |
| Load Balancer | 1 vCPU | 1GB RAM | 10GB SSD | 1Gbps |

#### Recommended Production Specs

| Component | CPU | Memory | Storage | Network |
|-----------|-----|---------|---------|---------|
| Application Server | 4 vCPU | 8GB RAM | 50GB SSD | 10Gbps |
| Database Server | 4 vCPU | 16GB RAM | 500GB SSD | 10Gbps |
| Cache Server | 2 vCPU | 4GB RAM | 20GB SSD | 10Gbps |
| Load Balancer | 2 vCPU | 4GB RAM | 20GB SSD | 10Gbps |

## Docker Deployment

### Production Dockerfile

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/web

# Production stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary and static files
COPY --from=builder /app/main /main
COPY --from=builder /app/static /static
COPY --from=builder /app/migrations /migrations

# Set environment variables
ENV GIN_MODE=release
ENV TZ=UTC

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/main", "health"]

# Run the application
ENTRYPOINT ["/main"]
CMD ["api"]
```

### Docker Compose Production

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    image: eshop-api:latest
    container_name: eshop-api
    restart: unless-stopped
    environment:
      - ENV=production
      - PORT=4000
      - DB_URL=${DB_URL}
      - REDIS_URL=${REDIS_URL}
      - CLOUDINARY_URL=${CLOUDINARY_URL}
      - STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
      - SYMMETRIC_KEY=${SYMMETRIC_KEY}
      - ACCESS_TOKEN_DURATION=24h
      - REFRESH_TOKEN_DURATION=720h
    ports:
      - "4000:4000"
    depends_on:
      - postgres
      - redis
    networks:
      - eshop-network
    healthcheck:
      test: ["CMD", "/main", "health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:14-alpine
    container_name: eshop-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=en_US.UTF-8 --lc-ctype=en_US.UTF-8"
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sh:/docker-entrypoint-initdb.d/init-db.sh
    networks:
      - eshop-network
    command: >
      postgres
      -c max_connections=100
      -c shared_buffers=256MB
      -c effective_cache_size=1GB
      -c work_mem=4MB
      -c maintenance_work_mem=64MB
      -c checkpoint_completion_target=0.9
      -c wal_buffers=16MB
      -c default_statistics_target=100
      -c random_page_cost=1.1
      -c effective_io_concurrency=200

  redis:
    image: redis:7-alpine
    container_name: eshop-redis
    restart: unless-stopped
    command: >
      redis-server
      --maxmemory 512mb
      --maxmemory-policy allkeys-lru
      --appendonly yes
      --appendfsync everysec
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - eshop-network

  nginx:
    image: nginx:alpine
    container_name: eshop-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - eshop-network

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  eshop-network:
    driver: bridge
```

### Environment Variables

```bash
# .env.prod
ENV=production
PORT=4000
DOMAIN=yourdomain.com

# Database
DB_URL=postgresql://eshop_user:secure_password@postgres:5432/eshop_prod?sslmode=require
MAX_POOL_SIZE=20
MIGRATION_PATH=file://migrations

# Cache
REDIS_URL=redis:6379

# Authentication
SYMMETRIC_KEY=your-32-character-secret-key-here
ACCESS_TOKEN_DURATION=24h
REFRESH_TOKEN_DURATION=720h

# External Services
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
STRIPE_SECRET_KEY=sk_live_your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=pk_live_your_stripe_publishable_key

# Email
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password

# Database Configuration
POSTGRES_DB=eshop_prod
POSTGRES_USER=eshop_user
POSTGRES_PASSWORD=secure_database_password
```

### Nginx Configuration

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream app_servers {
        server app:4000;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=auth:10m rate=5r/s;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 10240;
    gzip_proxied expired no-cache no-store private must-revalidate;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml
        application/atom+xml
        image/svg+xml;

    server {
        listen 80;
        server_name yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name yourdomain.com;

        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;

        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header Referrer-Policy "no-referrer-when-downgrade" always;
        add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

        # API routes
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            
            proxy_pass http://app_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # Timeouts
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
            
            # Buffer settings
            proxy_buffering on;
            proxy_buffer_size 8k;
            proxy_buffers 16 8k;
        }

        # Auth routes (stricter rate limiting)
        location /api/v1/auth/ {
            limit_req zone=auth burst=10 nodelay;
            
            proxy_pass http://app_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Health check endpoint
        location /health {
            proxy_pass http://app_servers;
            access_log off;
        }

        # Swagger documentation (restrict access)
        location /swagger/ {
            # Restrict to specific IPs or use basic auth
            # allow 192.168.1.0/24;
            # deny all;
            
            proxy_pass http://app_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

## Cloud Deployment

### AWS Deployment

#### Using ECS with Fargate

```yaml
# aws-task-definition.json
{
  "family": "eshop-api",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "1024",
  "memory": "2048",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "eshop-api",
      "image": "your-account.dkr.ecr.region.amazonaws.com/eshop-api:latest",
      "portMappings": [
        {
          "containerPort": 4000,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "ENV",
          "value": "production"
        },
        {
          "name": "PORT",
          "value": "4000"
        }
      ],
      "secrets": [
        {
          "name": "DB_URL",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:eshop/db-url"
        },
        {
          "name": "STRIPE_SECRET_KEY",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:eshop/stripe-key"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/eshop-api",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "healthCheck": {
        "command": ["CMD-SHELL", "/main health"],
        "interval": 30,
        "timeout": 5,
        "retries": 3,
        "startPeriod": 60
      }
    }
  ]
}
```

#### Infrastructure as Code (Terraform)

```hcl
# main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC and Networking
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "eshop-vpc"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "eshop-igw"
  }
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "eshop-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

# RDS Instance
resource "aws_db_instance" "postgres" {
  identifier = "eshop-postgres"
  
  engine         = "postgres"
  engine_version = "14.9"
  instance_class = "db.t3.micro"
  
  allocated_storage     = 20
  max_allocated_storage = 100
  storage_type         = "gp2"
  storage_encrypted    = true
  
  db_name  = "eshop"
  username = "eshop_user"
  password = var.db_password
  
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  skip_final_snapshot = true
  deletion_protection = true

  tags = {
    Name = "eshop-postgres"
  }
}

# ElastiCache Redis
resource "aws_elasticache_subnet_group" "main" {
  name       = "eshop-cache-subnet"
  subnet_ids = aws_subnet.private[*].id
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "eshop-redis"
  engine               = "redis"
  node_type           = "cache.t3.micro"
  num_cache_nodes     = 1
  parameter_group_name = "default.redis7"
  port                = 6379
  subnet_group_name   = aws_elasticache_subnet_group.main.name
  security_group_ids  = [aws_security_group.redis.id]
}

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "eshop-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets           = aws_subnet.public[*].id

  enable_deletion_protection = true

  tags = {
    Name = "eshop-alb"
  }
}
```

### Google Cloud Platform (GCP)

#### Using Cloud Run

```yaml
# cloudbuild.yaml
steps:
  # Build the container image
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/eshop-api:$COMMIT_SHA', '.']

  # Push the container image to Container Registry
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/eshop-api:$COMMIT_SHA']

  # Deploy to Cloud Run
  - name: 'gcr.io/cloud-builders/gcloud'
    args:
    - 'run'
    - 'deploy'
    - 'eshop-api'
    - '--image'
    - 'gcr.io/$PROJECT_ID/eshop-api:$COMMIT_SHA'
    - '--region'
    - 'us-central1'
    - '--platform'
    - 'managed'
    - '--allow-unauthenticated'
    - '--set-env-vars'
    - 'ENV=production'
    - '--set-secrets'
    - 'DB_URL=db_url:latest,STRIPE_SECRET_KEY=stripe_key:latest'
```

### Digital Ocean

#### Using App Platform

```yaml
# .do/app.yaml
name: eshop-api
services:
- name: api
  source_dir: /
  github:
    repo: your-username/go-eshop
    branch: main
  run_command: ./main api
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  
  envs:
  - key: ENV
    value: production
  - key: PORT
    value: "8080"
  - key: DB_URL
    type: SECRET
    value: your_database_connection_string
  - key: REDIS_URL
    type: SECRET
    value: your_redis_connection_string
  
  health_check:
    http_path: /health
  
  http_port: 8080

databases:
- name: eshop-postgres
  engine: PG
  version: "14"
  size: db-s-1vcpu-1gb
  
- name: eshop-redis
  engine: REDIS
  version: "7"
  size: db-s-1vcpu-1gb
```

## Database Migration in Production

### Zero-Downtime Migration Strategy

1. **Backward Compatible Changes**
   ```sql
   -- Step 1: Add new column (nullable)
   ALTER TABLE products ADD COLUMN new_field VARCHAR(255);
   
   -- Step 2: Deploy application code that can handle both old and new schema
   
   -- Step 3: Populate new column
   UPDATE products SET new_field = old_field WHERE new_field IS NULL;
   
   -- Step 4: Make column NOT NULL (after all data is migrated)
   ALTER TABLE products ALTER COLUMN new_field SET NOT NULL;
   
   -- Step 5: Drop old column (in next release)
   ALTER TABLE products DROP COLUMN old_field;
   ```

2. **Breaking Changes (Blue-Green Deployment)**
   ```bash
   # Create new environment
   docker-compose -f docker-compose.blue-green.yml up -d
   
   # Run migrations on new environment
   docker exec app-green /main migrate up
   
   # Switch traffic to new environment
   # Update load balancer configuration
   
   # Verify and cleanup old environment
   ```

### Migration Scripts

```bash
#!/bin/bash
# scripts/migrate-production.sh

set -e

DATABASE_URL=$1
MIGRATION_DIR="./migrations"

if [ -z "$DATABASE_URL" ]; then
    echo "Usage: $0 <database_url>"
    exit 1
fi

echo "Starting production migration..."

# Backup database
echo "Creating database backup..."
pg_dump "$DATABASE_URL" > "backup_$(date +%Y%m%d_%H%M%S).sql"

# Check migration status
echo "Current migration version:"
migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" version

# Run migrations
echo "Running migrations..."
migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" up

# Verify migration
echo "New migration version:"
migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" version

echo "Migration completed successfully!"
```

## Monitoring and Observability

### Application Metrics

```go
// metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    RequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "The total number of processed HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    RequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "The HTTP request latencies in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users_total",
            Help: "The total number of active users",
        },
    )

    DatabaseConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_connections_active",
            Help: "The number of active database connections",
        },
    )
)

// Middleware for metrics collection
func PrometheusMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        status := strconv.Itoa(c.Writer.Status())
        
        RequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        RequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration.Seconds())
    })
}
```

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alert_rules.yml"

scrape_configs:
  - job_name: 'eshop-api'
    static_configs:
      - targets: ['app:4000']
    scrape_interval: 5s
    metrics_path: /metrics

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - 'alertmanager:9093'
```

### Alert Rules

```yaml
# alert_rules.yml
groups:
  - name: eshop_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors per second"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is {{ $value }} seconds"

      - alert: DatabaseConnections
        expr: database_connections_active > 80
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High database connection usage"
          description: "Database connections: {{ $value }}"

      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is down"
          description: "{{ $labels.instance }} service is down"
```

### Logging in Production

```yaml
# docker-compose.logging.yml
version: '3.8'

services:
  app:
    # ... existing configuration
    logging:
      driver: "fluentd"
      options:
        fluentd-address: "fluentd:24224"
        tag: "eshop.api"

  fluentd:
    image: fluentd:v1.14-debian
    ports:
      - "24224:24224"
      - "24224:24224/udp"
    volumes:
      - ./fluentd.conf:/fluentd/etc/fluent.conf
      - fluentd_data:/fluentd/log

  elasticsearch:
    image: elasticsearch:8.5.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data

  kibana:
    image: kibana:8.5.0
    ports:
      - "5601:5601"
    environment:
      ELASTICSEARCH_URL: http://elasticsearch:9200
    depends_on:
      - elasticsearch

volumes:
  fluentd_data:
  elasticsearch_data:
```

## Security in Production

### SSL/TLS Configuration

```bash
#!/bin/bash
# scripts/setup-ssl.sh

# Using Let's Encrypt with Certbot
sudo apt-get update
sudo apt-get install certbot python3-certbot-nginx

# Generate certificates
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Set up auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Security Headers

```go
// middleware/security.go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Next()
    }
}
```

### Rate Limiting

```go
// middleware/rate_limit.go
import (
    "github.com/gin-gonic/gin"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimitMiddleware(store limiter.Store, rate limiter.Rate) gin.HandlerFunc {
    middleware := limiter.NewHTTPMiddleware(limiter.New(store, rate))
    return func(c *gin.Context) {
        middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            c.Next()
        })).ServeHTTP(c.Writer, c.Request)
        
        if c.Writer.Status() == http.StatusTooManyRequests {
            c.Abort()
        }
    }
}

// Usage
func setupRateLimit() gin.HandlerFunc {
    store := redis.NewStore(redisClient)
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
    }
    return RateLimitMiddleware(store, rate)
}
```

## Backup and Disaster Recovery

### Database Backup Strategy

```bash
#!/bin/bash
# scripts/backup-database.sh

set -e

# Configuration
BACKUP_DIR="/backups"
RETENTION_DAYS=30
DATABASE_URL=$1
BACKUP_NAME="eshop_backup_$(date +%Y%m%d_%H%M%S)"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Create backup
echo "Creating database backup: $BACKUP_NAME"
pg_dump "$DATABASE_URL" | gzip > "$BACKUP_DIR/$BACKUP_NAME.sql.gz"

# Upload to S3 (optional)
if [ -n "$AWS_S3_BUCKET" ]; then
    echo "Uploading backup to S3..."
    aws s3 cp "$BACKUP_DIR/$BACKUP_NAME.sql.gz" "s3://$AWS_S3_BUCKET/backups/"
fi

# Cleanup old backups
echo "Cleaning up backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "eshop_backup_*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: $BACKUP_NAME"
```

### Automated Backup with Cron

```bash
# Add to crontab
0 2 * * * /path/to/scripts/backup-database.sh $DATABASE_URL >> /var/log/backup.log 2>&1
```

### Disaster Recovery Plan

1. **Data Recovery**
   ```bash
   # Restore from backup
   gunzip -c backup.sql.gz | psql $DATABASE_URL
   ```

2. **Infrastructure Recovery**
   ```bash
   # Rebuild from Infrastructure as Code
   terraform apply
   
   # Redeploy application
   docker-compose up -d
   ```

3. **Recovery Testing**
   - Regular recovery drills
   - Automated recovery testing
   - Documentation updates

## Performance Optimization

### Database Optimization

```sql
-- Regular maintenance queries
VACUUM ANALYZE;
REINDEX DATABASE eshop;

-- Monitor slow queries
SELECT query, mean_time, calls, total_time 
FROM pg_stat_statements 
ORDER BY total_time DESC 
LIMIT 10;

-- Index usage analysis
SELECT 
    indexrelname as index_name,
    relname as table_name,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

### Application Optimization

```go
// Connection pool optimization
func setupDatabase() *pgxpool.Pool {
    config, _ := pgxpool.ParseConfig(databaseURL)
    
    // Production settings
    config.MaxConns = 30
    config.MinConns = 10
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = time.Minute * 30
    config.HealthCheckPeriod = time.Minute * 1
    
    return pgxpool.ConnectConfig(context.Background(), config)
}

// Redis optimization
func setupCache() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         redisURL,
        PoolSize:     20,
        MinIdleConns: 5,
        MaxRetries:   3,
        DialTimeout:  time.Second * 5,
        ReadTimeout:  time.Second * 3,
        WriteTimeout: time.Second * 3,
    })
}
```

## Maintenance Procedures

### Regular Maintenance Checklist

#### Daily
- [ ] Check application logs for errors
- [ ] Monitor system resources (CPU, memory, disk)
- [ ] Verify backup completion
- [ ] Review security alerts

#### Weekly
- [ ] Analyze performance metrics
- [ ] Review slow query logs
- [ ] Update security patches
- [ ] Test backup restoration

#### Monthly
- [ ] Database maintenance (VACUUM, REINDEX)
- [ ] Review and rotate logs
- [ ] Security audit
- [ ] Capacity planning review
- [ ] Update dependencies

### Emergency Procedures

#### Service Outage Response

1. **Immediate Response (0-5 minutes)**
   - Check service status
   - Identify affected components
   - Implement quick fixes if available

2. **Investigation (5-30 minutes)**
   - Analyze logs and metrics
   - Identify root cause
   - Communicate with stakeholders

3. **Resolution (30+ minutes)**
   - Implement permanent fix
   - Verify service restoration
   - Post-mortem analysis

#### Scaling Procedures

```bash
# Scale application instances
docker-compose up -d --scale app=3

# Database scaling (read replicas)
# Update connection strings to use read replicas for read operations

# Cache scaling
# Add Redis cluster nodes
```

This deployment guide provides a comprehensive foundation for deploying and maintaining the e-commerce platform in production environments. Regular updates and team training will ensure smooth operations and minimal downtime.