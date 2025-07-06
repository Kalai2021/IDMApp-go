# IDM Application - Fullstack with Centralized Logging

A complete Identity Management (IDM) system with Go backend, React frontend, and centralized logging using the ELK stack (Elasticsearch, Logstash, Kibana) and Fluentd.

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React App     │    │   Go Backend    │    │   PostgreSQL    │
│   (Port 3000)   │    │   (Port 8080)   │    │   (Port 5432)   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │        Fluentd            │
                    │      (Port 24224)         │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │     Elasticsearch         │
                    │       (Port 9200)         │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │        Kibana             │
                    │       (Port 5601)         │
                    └───────────────────────────┘
```

## Features

- 🔐 **Authentication**: Secure login with token-based authentication
- 👥 **User Management**: Complete CRUD operations for users
- 🔒 **Protected Routes**: Automatic redirection for unauthenticated users
- 🚀 **Docker Ready**: Containerized deployment with Docker
- 📊 **Centralized Logging**: Complete ELK stack integration
- 🔍 **Real-time Monitoring**: Live log aggregation and visualization
- 📈 **Structured Logging**: JSON-formatted logs with metadata

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- PostgreSQL database
- At least 8GB RAM (for ELK stack)

## Quick Start

### Option 1: Fullstack with Logging (Recommended)

1. **Start the complete stack:**
   ```bash
   cd IDMApp-go
   chmod +x setup-fullstack-logging.sh
   ./setup-fullstack-logging.sh
   ```

2. **Access the services:**
   - **Frontend**: http://localhost:3000
   - **Backend API**: http://localhost:8080
   - **Kibana Dashboard**: http://localhost:5601
   - **Elasticsearch**: http://localhost:9200

### Option 2: Manual Docker Compose

```bash
# Start with logging
docker compose -f docker-compose.fullstack-logging.yml up --build

# Or start without logging
docker compose -f docker-compose.fullstack.yml up --build
```

### Option 3: Development Mode

1. **Backend (Go):**
   ```bash
   cd IDMApp-go
   go mod download
   go run main.go
   ```

2. **Frontend (React):**
   ```bash
   cd IDMReactClient
   npm install
   npm start
   ```

## 📊 Logging System

### Overview

The application uses a centralized logging system with:
- **Fluentd**: Log aggregation and routing
- **Elasticsearch**: Log storage and indexing
- **Kibana**: Log visualization and analysis

### Log Format

Both frontend and backend use standardized JSON logging:

```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "info",
  "service": "idmapp-backend|idmreactclient-frontend",
  "message": "User login successful",
  "data": {
    "user_id": "12345",
    "session_id": "abc123",
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
    "request_id": "req-12345",
    "duration_ms": 150,
    "status_code": 200,
    "endpoint": "/api/login",
    "method": "POST"
  }
}
```

### Log Categories

- **Application Logs**: User actions, business logic events
- **Security Logs**: Authentication, authorization, access attempts
- **Performance Logs**: Response times, database queries
- **Error Logs**: Exceptions, failures, debugging information

### Viewing Logs

#### 1. Kibana Dashboard (Recommended)
- **URL**: http://localhost:5601
- **Setup**: Create index patterns for `idm-backend-logs-*` and `idm-frontend-logs-*`
- **Features**: Search, filter, visualize, and analyze logs

#### 2. Docker Container Logs
```bash
# Backend logs
docker logs idm-backend -f

# Frontend logs
docker logs idm-react-client -f

# All services
docker compose -f docker-compose.fullstack-logging.yml logs -f

# Specific service
docker compose -f docker-compose.fullstack-logging.yml logs -f idm-backend
```

#### 3. Elasticsearch Direct
```bash
# Check indices
curl http://localhost:9200/_cat/indices

# Search logs
curl http://localhost:9200/idm-backend-logs-*/_search
```

### Logging Configuration

#### Backend (Go)
Environment variables:
```bash
FLUENT_ENABLED=true
FLUENT_ENDPOINT=http://fluentd:24224
LOG_LEVEL=info
```

#### Frontend (React)
Environment variables:
```bash
REACT_APP_FLUENT_ENDPOINT=http://fluentd:24224
REACT_APP_LOGGING_ENABLED=true
```

### Adding Logging to Your Code

#### Backend (Go)
```go
import "idmapp-go/services"

func (c *Controller) HandleRequest(w http.ResponseWriter, r *http.Request) {
    logger := services.GetFluentLogger()
    
    logger.Info("Request received", map[string]interface{}{
        "endpoint": r.URL.Path,
        "method":   r.Method,
        "user_id":  "12345",
    })
}
```

#### Frontend (React)
```typescript
import { useLogger } from './hooks/useLogger';

function MyComponent() {
  const logger = useLogger();
  
  const handleClick = () => {
    logger.info('Button clicked', {
      component: 'MyComponent',
      action: 'button_click',
      user_id: '12345'
    });
  };
  
  return <button onClick={handleClick}>Click me</button>;
}
```

## Configuration

### Environment Variables

#### Backend (Go)
| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `host.docker.internal` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `iamdb` |
| `DB_USERNAME` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `SERVER_PORT` | Backend port | `8080` |
| `LOG_LEVEL` | Log level | `info` |
| `FLUENT_ENABLED` | Enable Fluentd logging | `true` |
| `FLUENT_ENDPOINT` | Fluentd endpoint | `http://fluentd:24224` |

#### Frontend (React)
| Variable | Description | Default |
|----------|-------------|---------|
| `REACT_APP_API_URL` | Backend API URL | `http://localhost:8080` |
| `REACT_APP_FLUENT_ENDPOINT` | Fluentd endpoint | `http://fluentd:24224` |
| `REACT_APP_LOGGING_ENABLED` | Enable logging | `true` |

### Database Configuration

The application connects to a PostgreSQL database:
- **Host**: Your local PostgreSQL instance
- **Database**: `iamdb`
- **Username**: `postgres`
- **Password**: `postgres`

## Project Structure

```
IDMApp-go/
├── docker-compose.fullstack-logging.yml  # Complete stack with logging
├── docker-compose.fullstack.yml          # Stack without logging
├── services/
│   └── fluent_logger.go                  # Fluentd logger service
├── middleware/
│   └── logging.go                        # HTTP logging middleware
├── fluentd/
│   ├── Dockerfile                        # Fluentd with plugins
│   └── conf/
│       └── fluent.conf                   # Fluentd configuration
└── setup-fullstack-logging.sh            # Automated setup script

IDMReactClient/
├── Dockerfile                            # React app container
├── nginx.conf                            # Nginx configuration
├── src/
│   ├── services/
│   │   └── logger.ts                     # Frontend logger service
│   └── hooks/
│       └── useLogger.ts                  # React logging hook
└── package.json
```

## Authentication

The app uses token-based authentication:

1. **Login**: Users authenticate with email/password
2. **Token Storage**: JWT tokens are stored in localStorage
3. **Auto-refresh**: Tokens are automatically included in API requests
4. **Logout**: Tokens are cleared on logout
5. **Audit Logging**: All authentication events are logged

## User Management

### Features
- **View Users**: Display all users in a responsive table
- **Create User**: Add new users with validation
- **Edit User**: Update existing user information
- **Delete User**: Remove users with confirmation
- **Status Management**: Toggle user active/inactive status
- **Activity Logging**: All user operations are logged

### User Fields
- Email (unique, required)
- First Name (required)
- Last Name (required)
- Display Name (required)
- Status (active/inactive)

## Development

### Available Scripts

#### Backend (Go)
```bash
go run main.go                    # Run in development
go build                         # Build binary
go test ./...                    # Run tests
```

#### Frontend (React)
```bash
npm start                        # Start development server
npm run build                    # Build for production
npm test                        # Run tests
```

### Code Style

The project uses:
- TypeScript for frontend type safety
- Go modules for backend dependency management
- ESLint for code linting
- Prettier for code formatting
- Functional components with hooks

## Docker Deployment

### Development
```bash
# With logging
docker compose -f docker-compose.fullstack-logging.yml up --build

# Without logging
docker compose -f docker-compose.fullstack.yml up --build
```

### Production
```bash
# Build and run with logging
docker compose -f docker-compose.fullstack-logging.yml up -d

# Or individual services
docker build -t idm-backend .
docker build -t idm-frontend ../IDMReactClient
```

### Docker Compose Files

- **`docker-compose.fullstack-logging.yml`**: Complete stack with ELK logging
- **`docker-compose.fullstack.yml`**: Stack without logging infrastructure

## API Integration

The app integrates with the IDM Go backend API:

### Endpoints Used
- `POST /api/v1/auth/login` - User authentication
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get specific user
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Error Handling
- Automatic token refresh on 401 errors
- User-friendly error messages
- Loading states for better UX
- Comprehensive error logging

## Monitoring and Observability

### Health Checks
- **Backend**: http://localhost:8080/health
- **Frontend**: http://localhost:3000
- **Elasticsearch**: http://localhost:9200/_cluster/health
- **Kibana**: http://localhost:5601/api/status

### Metrics
- Request/response times
- Error rates
- User activity patterns
- Database performance
- System resource usage

### Alerts
- High error rates
- Slow response times
- Authentication failures
- Database connection issues

## Security

- HTTPS enforcement in production
- Secure headers via Nginx
- Token-based authentication
- Input validation and sanitization
- XSS protection
- Comprehensive security logging
- Audit trails for all operations

## Troubleshooting

### Common Issues

1. **API Connection Failed**
   - Check if the backend is running
   - Verify `REACT_APP_API_URL` in environment
   - Check network connectivity
   - Review backend logs in Kibana

2. **Authentication Issues**
   - Clear browser localStorage
   - Check token expiration
   - Verify backend authentication
   - Check authentication logs

3. **Docker Issues**
   - Ensure Docker is running
   - Check port conflicts
   - Verify Docker Compose version
   - Check container logs

4. **Logging Issues**
   - Verify Fluentd is running: `docker logs idm-fluentd`
   - Check Elasticsearch health: `curl http://localhost:9200/_cluster/health`
   - Verify Kibana is accessible: http://localhost:5601
   - Check index patterns in Kibana

5. **Database Issues**
   - Verify PostgreSQL is running
   - Check database credentials
   - Ensure database `iamdb` exists
   - Review database logs

### Performance Optimization

1. **Reduce Memory Usage**
   - Edit `docker-compose.yml` and reduce memory limits
   - Use `docker-compose.override.yml` for development

2. **Disable Unused Services**
   - Comment out services you don't need in `docker-compose.yml`

3. **Use Lightweight Images**
   - Alpine-based images are used where possible

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add appropriate logging
5. Test thoroughly
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 