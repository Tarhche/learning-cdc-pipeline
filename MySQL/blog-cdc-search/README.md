# Blog Application with CDC and Search

A modern blog application built with Go using Domain-Driven Design (DDD) architecture, featuring MySQL database, RabbitMQ message queue, Maxwell CDC, Typesense search engine, and comprehensive unit tests. The application implements Change Data Capture (CDC) to automatically sync database changes to a search index.

## Features

- **CRUD Operations**: Create, read, update, and delete blog posts
- **Rich Content**: Each post includes title, image, excerpt, and body
- **Modern UI**: Clean, responsive web interface with forms for post management
- **RESTful API**: JSON-based API endpoints for programmatic access
- **Full-Text Search**: Advanced search capabilities powered by Typesense
- **Change Data Capture**: Automatic database synchronization using Maxwell and RabbitMQ
- **Services**: Separate services for blog operations and CDC processing
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **No Authentication**: Simple setup without complex user management

## Architecture

This application follows Domain-Driven Design (DDD) principles with a clean layered architecture and microservices pattern:

```
blog-cdc-search/
├── domain/           # Domain entities and business logic
├── application/      # Application services and use cases
├── infrastructure/   # Database, web handlers, and external concerns
├── cmd/
│   ├── blog/        # Main blog application service
│   └── cdc/         # CDC service for data synchronization
└── database/         # Database schemas and configuration
```

### Layer Responsibilities

- **Domain Layer**: Contains the `Post` entity with business rules and validation
- **Application Layer**: Orchestrates business logic through services (PostService, SearchService, CDCService)
- **Infrastructure Layer**: Handles data persistence, web interface, message queues, and search indexing

### Services

1. **Blog Service** (`cmd/blog/`): Main web application handling HTTP requests and blog operations
2. **CDC Service** (`cmd/cdc/`): Processes database change events and syncs to search index
3. **MySQL**: Primary database for blog posts
4. **RabbitMQ**: Message queue for CDC events
5. **Maxwell**: CDC tool that captures MySQL binlog changes
6. **Typesense**: Search engine for full-text search capabilities

## Prerequisites

- Docker and Docker Compose
- Go 1.25 or higher (handled by Docker)
- MySQL 8.0 (handled by Docker)
- RabbitMQ 4.1 (handled by Docker)
- Typesense 29.0 (handled by Docker)

## Quick Start

### 1. Clone and Navigate

```bash
cd blog-cdc-search
```

### 2. Run with Docker Compose

```bash
# Start all services
make up

# Or manually
docker compose up --build -d
```

The application will be available at:

- **Blog Web Interface**: http://localhost:8085
- **MySQL Database**: http://localhost:3306
- **phpMyAdmin**: http://localhost:8081
- **RabbitMQ Management**: http://localhost:15672
- **Typesense Dashboard**: http://localhost:8082
- **Typesense API**: localhost:8108

### 3. Access the Blog

- Visit http://localhost:8085 to see the main blog page
- Click "Create New Post" to add your first blog post
- Use the edit and delete buttons to manage existing posts
- Search functionality is available through the API

### 4. Stop Services

```bash
make down

# Or manually
docker compose down --volumes --remove-orphans
```

## API Endpoints

### Web Interface
- `GET /` - Main blog page
- `GET /post/{id}` - View specific post
- `GET /dashboard` - Admin dashboard
- `GET /dashboard/create` - Create post form
- `GET /dashboard/edit` - Edit post form

### REST API
- `POST /api/posts` - Create a new post
- `GET /api/posts` - Get all posts
- `GET /api/posts?id={id}` - Get a specific post
- `PUT /api/posts?id={id}` - Update a post
- `DELETE /api/posts?id={id}` - Delete a post

### Search API
- `POST /api/search` - Search posts with parameters
- `GET /api/search?query={query}&page={page}&per_page={per_page}` - Search posts

### API Request/Response Format

**Create/Update Post:**
```json
{
  "title": "Post Title",
  "image": "https://example.com/image.jpg",
  "excerpt": "Brief summary",
  "body": "Full post content"
}
```

**Post Response:**
```json
{
  "id": 1,
  "title": "Post Title",
  "image": "https://example.com/image.jpg",
  "excerpt": "Brief summary",
  "body": "Full post content",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Search Request:**
```json
{
  "query": "search term",
  "page": 1,
  "per_page": 10,
  "sort_by": "created_at",
  "filter_by": "category"
}
```

## CDC (Change Data Capture) Architecture

The application uses Maxwell to capture MySQL binlog changes and publish them to RabbitMQ:

1. **MySQL Binlog**: Maxwell reads MySQL binary logs
2. **Event Publishing**: Changes are published to RabbitMQ exchange
3. **CDC Service**: Consumes events and syncs to Typesense search index
4. **Real-time Sync**: Search index stays updated with database changes

### CDC Flow
```
MySQL → Maxwell → RabbitMQ → CDC Service → Typesense
```

## Testing

Run the unit tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run specific test packages:

```bash
go test ./domain/...
go test ./application/...
go test ./infrastructure/...
```

## Docker Configuration

### Services
- **MySQL 8.0**: Database with persistent volume and Maxwell user setup
- **RabbitMQ 4.1**: Message queue with management interface
- **Maxwell**: CDC tool for MySQL binlog processing
- **Typesense 29.0**: Search engine with dashboard
- **Blog App**: Go application on port 8085
- **CDC Service**: Go service for processing CDC events
- **phpMyAdmin**: MySQL management interface
- **Typesense Dashboard**: Search index management

### Adding New Features

1. **Domain Layer**: Add new entities and business rules
2. **Application Layer**: Create services for new use cases
3. **Infrastructure Layer**: Implement data access, message queues, and web handlers
4. **Tests**: Add comprehensive unit tests

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure MySQL container is running: `docker compose ps`
   - Check environment variables match compose.yml
   - Wait for MySQL health check to pass

2. **Port Already in Use**
   - Change port in compose.yml or main.go
   - Kill existing process: `lsof -ti:8085 | xargs kill`

3. **Go Module Issues**
   - Run `go mod tidy` to clean dependencies
   - Delete go.sum and run `go mod download`

4. **CDC Service Not Processing Events**
   - Check RabbitMQ connection: http://localhost:15672
   - Verify Maxwell is running and connected to MySQL
   - Check CDC service logs: `docker compose logs cdc-service`

5. **Search Not Working**
   - Ensure Typesense is running: `docker compose ps typesense`
   - Check Typesense dashboard: http://localhost:8082
   - Verify CDC service is syncing data

### Logs

View application logs:
```bash
docker compose logs blog-app
docker compose logs cdc-service
```

View database logs:
```bash
docker compose logs mysql
```

View message queue logs:
```bash
docker compose logs rabbitmq
```

View CDC logs:
```bash
docker compose logs maxwell
```

## Contributing

1. Follow the existing DDD architecture
2. Add unit tests for new functionality
3. Update documentation as needed
4. Ensure Docker builds successfully
5. Test CDC functionality with database changes

## License

This project is open source and available under the MIT License.
