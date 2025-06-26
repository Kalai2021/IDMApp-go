# Migration Guide: Java Spring Boot to Go

This guide explains the differences between the original Java Spring Boot IDM application and the Go version, and how to migrate between them.

## Overview

The Go version is a complete rewrite of the Java Spring Boot application, maintaining the same functionality while leveraging Go's strengths.

## Key Differences

### 1. Framework and Architecture

| Aspect | Java Spring Boot | Go Version |
|--------|------------------|------------|
| Framework | Spring Boot 3.2.3 | Gin 1.9.1 |
| Architecture | MVC with Spring DI | Clean Architecture with dependency injection |
| Configuration | YAML files | Environment variables |
| Dependency Injection | Spring IoC container | Manual dependency injection |

### 2. Database Layer

| Aspect | Java Spring Boot | Go Version |
|--------|------------------|------------|
| ORM | Spring Data JPA | GORM |
| Database Driver | PostgreSQL JDBC | PostgreSQL driver |
| Migrations | Flyway/Liquibase | GORM Auto-migration |
| Connection Pool | HikariCP | Built-in connection management |

### 3. Authentication & Authorization

| Aspect | Java Spring Boot | Go Version |
|--------|------------------|------------|
| JWT Library | Auth0 Java JWT | golang-jwt/jwt |
| Middleware | Spring Security | Custom JWT middleware |
| OpenFGA | Java SDK | Go SDK |
| Password Hashing | BCrypt | golang.org/x/crypto/bcrypt |

### 4. Project Structure

#### Java Structure
```
IDMApp/
├── shared-security/
├── idmapp/
├── usermanagement/
├── groupmanagement/
├── rolemanagement/
├── orgmanagement/
└── pom.xml
```

#### Go Structure
```
idmapp-go/
├── config/
├── controllers/
├── database/
├── dto/
├── middleware/
├── models/
├── routes/
├── services/
├── main.go
└── go.mod
```

## Data Model Mapping

### User Entity

#### Java (User.java)
```java
public class User {
    private UUID id;
    private String name;
    private String firstName;
    private String lastName;
    private String email;
    private String password;
    private boolean isActive = true;
    private LocalDateTime createdAt;
    private LocalDateTime modifiedAt;
}
```

#### Go (models/user.go)
```go
type User struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name      string    `json:"name" gorm:"not null"`
    FirstName string    `json:"firstName" gorm:"column:first_name"`
    LastName  string    `json:"lastName" gorm:"column:last_name"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`
    IsActive  bool      `json:"isActive" gorm:"default:true"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

## API Endpoint Mapping

| Java Endpoint | Go Endpoint | Method | Description |
|---------------|-------------|--------|-------------|
| `/api/v1/users` | `/api/v1/users` | GET | Get all users |
| `/api/v1/users/{id}` | `/api/v1/users/{id}` | GET | Get user by ID |
| `/api/v1/users` | `/api/v1/users` | POST | Create user |
| `/api/v1/users/{id}` | `/api/v1/users/{id}` | PUT | Update user |
| `/api/v1/users/{id}` | `/api/v1/users/{id}` | DELETE | Delete user |
| `/api/v1/auth/login` | `/api/v1/auth/login` | POST | User login |

## Configuration Migration

### Java (application.yml)
```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/${DB_NAME}
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}

auth0:
  domain: ${AUTH0_DOMAIN}
  audience: ${AUTH0_AUDIENCE}
  client-id: ${AUTH0_CLIENT_ID}
  client-secret: ${AUTH0_CLIENT_SECRET}
```

### Go (.env)
```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=idmapp
DB_USERNAME=postgres
DB_PASSWORD=password

AUTH0_DOMAIN=your-domain.auth0.com
AUTH0_AUDIENCE=your-audience
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
```

## Database Migration

### Schema Changes
The Go version uses GORM auto-migration, which will create the necessary tables automatically. The schema is compatible with the Java version.

### Data Migration
If you have existing data in the Java application:

1. Export data from Java application:
```sql
-- Export users
COPY (SELECT * FROM users) TO '/tmp/users.csv' CSV HEADER;

-- Export groups
COPY (SELECT * FROM groups) TO '/tmp/groups.csv' CSV HEADER;

-- Export roles
COPY (SELECT * FROM roles) TO '/tmp/roles.csv' CSV HEADER;

-- Export organizations
COPY (SELECT * FROM orgs) TO '/tmp/orgs.csv' CSV HEADER;
```

2. Import data to Go application:
```sql
-- Import users
COPY users FROM '/tmp/users.csv' CSV HEADER;

-- Import groups
COPY groups FROM '/tmp/groups.csv' CSV HEADER;

-- Import roles
COPY roles FROM '/tmp/roles.csv' CSV HEADER;

-- Import organizations
COPY orgs FROM '/tmp/orgs.csv' CSV HEADER;
```

## Deployment Migration

### Java Deployment
```bash
# Build JAR
mvn clean package

# Run with Spring Boot
java -jar target/idmapp-1.0.0.jar
```

### Go Deployment
```bash
# Build binary
go build -o idmapp-go main.go

# Run binary
./idmapp-go
```

## Performance Comparison

### Memory Usage
- **Java**: ~200-500MB (JVM overhead)
- **Go**: ~20-50MB (native binary)

### Startup Time
- **Java**: 10-30 seconds (JVM startup + Spring context)
- **Go**: 1-5 seconds (direct execution)

### Runtime Performance
- **Java**: Good performance with JIT compilation
- **Go**: Excellent performance with native compilation

## Migration Checklist

### Pre-Migration
- [ ] Backup existing Java application data
- [ ] Document current configuration
- [ ] Identify custom business logic
- [ ] Plan deployment strategy

### Migration Steps
- [ ] Set up Go development environment
- [ ] Configure database connection
- [ ] Set up Auth0 integration
- [ ] Configure OpenFGA (if used)
- [ ] Test all API endpoints
- [ ] Validate data integrity
- [ ] Performance testing

### Post-Migration
- [ ] Monitor application performance
- [ ] Verify all functionality works
- [ ] Update documentation
- [ ] Train team on Go codebase
- [ ] Plan Java application retirement

## Benefits of Go Migration

### Performance
- Faster startup time
- Lower memory usage
- Better concurrency handling
- Native binary deployment

### Development
- Simpler dependency management
- Faster compilation
- Better tooling
- Easier deployment

### Operations
- Smaller container images
- Better resource utilization
- Easier monitoring
- Simplified deployment pipeline

## Challenges and Considerations

### Learning Curve
- Team needs to learn Go
- Different programming paradigms
- New tooling and ecosystem

### Ecosystem
- Smaller ecosystem compared to Java
- Fewer enterprise libraries
- Different best practices

### Maintenance
- Need to maintain both versions during transition
- Different debugging approaches
- New monitoring strategies

## Support and Resources

### Go Learning Resources
- [Go Official Documentation](https://golang.org/doc/)
- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)

### Migration Tools
- Use the provided scripts in `scripts/` directory
- Follow the setup guide in README.md
- Refer to API documentation in `docs/API.md`

### Community Support
- [Go Forum](https://forum.golangbridge.org/)
- [Gin GitHub Issues](https://github.com/gin-gonic/gin/issues)
- [GORM GitHub Issues](https://github.com/go-gorm/gorm/issues) 