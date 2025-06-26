# IDM App Go - Deployment Guide

## Overview

This guide covers different deployment strategies for the IDM App Go application.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker (optional, for containerized deployment)
- Auth0 account (for authentication)
- OpenFGA instance (optional, for authorization)

## Environment Configuration

### Required Environment Variables

Create a `.env` file based on `env.example`:
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=iamdb
DB_USERNAME=postgres
DB_PASSWORD=your-secure-password

# Auth0 Configuration
AUTH0_DOMAIN=your-domain.auth0.com
AUTH0_AUDIENCE=your-audience
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret

# OpenFGA Configuration (Optional)
OPENFGA_API_URL=http://localhost:8080
OPENFGA_STORE_ID=your-store-id
OPENFGA_API_TOKEN=your-api-token

# Server Configuration
SERVER_PORT=8080
LOG_LEVEL=info
```

## Deployment Methods

### 1. Binary Deployment

#### Build the Application
```bash
# Build for current platform
go build -o idmapp-go main.go

# Build for specific platform (Linux)
GOOS=linux GOARCH=amd64 go build -o idmapp-go main.go

# Build with optimizations
go build -ldflags="-s -w" -o idmapp-go main.go
```

#### Run the Application
```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=idmapp
export DB_USERNAME=postgres
export DB_PASSWORD=password
export SERVER_PORT=8080

# Run the application
./idmapp-go
```

### 2. Docker Deployment

#### Build Docker Image
```bash
docker build -t idmapp-go:latest .
```

#### Run Docker Container
```bash
# Run with environment file
docker run -p 8080:8080 --env-file .env idmapp-go:latest

# Run with environment variables
docker run -p 8080:8080 \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_NAME=idmapp \
  -e DB_USERNAME=postgres \
  -e DB_PASSWORD=password \
  idmapp-go:latest
```

### 3. Docker Compose Deployment

#### Start All Services
```bash
docker-compose up -d
```

#### View Logs
```bash
docker-compose logs -f app
```

#### Stop Services
```bash
docker-compose down
```

### 4. Production Deployment

#### Using Systemd (Linux)

Create a systemd service file `/etc/systemd/system/idmapp-go.service`:

```ini
[Unit]
Description=IDM App Go
After=network.target postgresql.service

[Service]
Type=simple
User=idmapp
WorkingDirectory=/opt/idmapp-go
ExecStart=/opt/idmapp-go/idmapp-go
Restart=always
RestartSec=5
EnvironmentFile=/opt/idmapp-go/.env

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable idmapp-go
sudo systemctl start idmapp-go
sudo systemctl status idmapp-go
```

#### Using Nginx as Reverse Proxy

Create nginx configuration `/etc/nginx/sites-available/idmapp-go`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/idmapp-go /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## Database Setup

### PostgreSQL Installation

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### CentOS/RHEL
```bash
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### macOS
```bash
brew install postgresql
brew services start postgresql
```

### Database Configuration

1. Create database and user:
```sql
CREATE DATABASE idmapp;
CREATE USER idmapp_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE idmapp TO idmapp_user;
```

2. Configure PostgreSQL for remote connections (if needed):
   - Edit `/etc/postgresql/*/main/postgresql.conf`
   - Set `listen_addresses = '*'`
   - Edit `/etc/postgresql/*/main/pg_hba.conf`
   - Add `host all all 0.0.0.0/0 md5`

3. Restart PostgreSQL:
```bash
sudo systemctl restart postgresql
```

## Security Considerations

### Environment Variables
- Never commit `.env` files to version control
- Use secure passwords for database
- Rotate secrets regularly
- Use environment-specific configurations

### Network Security
- Use HTTPS in production
- Configure firewall rules
- Use VPN for database access
- Implement rate limiting

### Application Security
- Enable CORS appropriately
- Validate all inputs
- Use prepared statements (handled by GORM)
- Implement proper logging
- Regular security updates

## Monitoring and Logging

### Application Logs
The application uses structured logging with Logrus. Configure log levels appropriately:

```bash
# Development
LOG_LEVEL=debug

# Production
LOG_LEVEL=info
```

### Health Checks
Monitor the health endpoint:
```bash
curl http://localhost:8080/health
```

### Database Monitoring
Monitor PostgreSQL performance:
```sql
-- Check active connections
SELECT count(*) FROM pg_stat_activity;

-- Check slow queries
SELECT query, mean_time FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;
```

## Backup and Recovery

### Database Backup
```bash
# Create backup
pg_dump -h localhost -U postgres idmapp > backup.sql

# Restore backup
psql -h localhost -U postgres idmapp < backup.sql
```

### Application Backup
```bash
# Backup configuration
cp .env .env.backup

# Backup binary
cp idmapp-go idmapp-go.backup
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check database is running
   - Verify connection parameters
   - Check firewall settings

2. **Port Already in Use**
   - Change SERVER_PORT in .env
   - Kill existing process: `lsof -ti:8080 | xargs kill`

3. **Permission Denied**
   - Check file permissions
   - Ensure user has proper access

4. **JWT Token Issues**
   - Verify Auth0 configuration
   - Check token expiration
   - Validate audience and issuer

### Debug Mode
Enable debug logging:
```bash
LOG_LEVEL=debug ./idmapp-go
```

### Performance Tuning
- Monitor database connections
- Use connection pooling
- Optimize queries
- Consider caching strategies 