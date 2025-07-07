# Compensation Data GraphQL & REST API

This project provides a GraphQL and REST API for querying compensation data stored in Elasticsearch. Built with Go, gqlgen, and the official Elasticsearch Go client.

## Features
- GraphQL API with filtering, sorting, and sparse fieldset support
- REST endpoint for fetching a single record with sparse fieldset
- Elasticsearch as the backend data store
- Dockerized setup for easy local development

## Prerequisites
- Go 1.18+
- Docker & Docker Compose
- (Optional) Elasticsearch and Kibana (if not using Docker Compose)

## Setup

### 1. Clone the repository
```sh
git clone <your-repo-url>
cd graphql-es
```

### 2. Start Elasticsearch (and Kibana) with Docker Compose
```sh
docker-compose up -d
```
- Elasticsearch will be available at http://localhost:9200
- Kibana (optional) at http://localhost:5601

### 3. Build and Run the Go API
```sh
go mod tidy
go run main.go
```
- The API will start at http://localhost:8087

## API Usage

### GraphQL
- Playground: http://localhost:8087/
- Endpoint: http://localhost:8087/query

Example query:
```graphql
query {
  compensations(id: "<doc_id>", salaryGte: 120000, location: "New York") {
    id
    job_title
    annual_salary
  }
}
```
- You can request only the fields you need (sparse fieldset).

### REST
- Endpoint: `GET /compensation_data?id=<doc_id>&fields=field1,field2`

Example:
```
curl "http://localhost:8087/compensation_data?id=abc123&fields=timestamp,job_title,salary"
```
- Returns only the requested fields for the given record.

## Data Loading
- Place your CSV files in the `dataset/` directory.
- The API will auto-upload data to Elasticsearch on startup.
- Index mapping is defined in `index.json`.

## Troubleshooting
- If Elasticsearch runs out of memory, increase Docker resources or adjust JVM settings in `docker-compose.yaml`.
- To reset data: `docker-compose down -v && docker-compose up -d`

## Development
- Update GraphQL schema in `graph/schema.graphqls` and run `go run github.com/99designs/gqlgen generate` to regenerate models/resolvers.
- Main code locations:
  - GraphQL: `graph/`
  - Elasticsearch client: `internal/elastic/`
  - REST handler: `rest_handlers.go`

## License
MIT



  podman run -d \
  --name elasticsearch \
  --network esnet \
  -p 9200:9200 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  docker.elastic.co/elasticsearch/elasticsearch:8.13.4


  podman run -d \
  --name kibana \
  --network esnet \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" \
  docker.elastic.co/kibana/kibana:8.13.4