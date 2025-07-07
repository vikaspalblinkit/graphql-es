# 🚀 Compensation Data API (GraphQL + REST)

This project provides a **GraphQL** and **REST API** for querying compensation data stored in **Elasticsearch**.
Built using **Go**, [gqlgen](https://gqlgen.com/), and the official Elasticsearch Go client.

---

## ✅ Features

* 🔍 GraphQL API with filtering, sorting, and sparse fieldset support
* 📌 REST endpoint to fetch a single record with selected fields
* 📦 Automatically uploads CSV datasets to Elasticsearch on startup
* 📊 Elasticsearch + Kibana setup using **Podman**
* 🧩 Modular Go codebase, clean structure

---

## 🧰 Prerequisites

* Go 1.18+
* [Podman](https://podman.io/) (or Docker if preferred)
* Git

---

## 🛠 Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/vikaspalblinkit/graphql-es
cd graphql-es
```

---

### 2. Start Elasticsearch & Kibana with Podman

```bash
# Create a custom network (run once)
podman network create esnet

# Start Elasticsearch
podman run -d \
  --name elasticsearch \
  --network esnet \
  -p 9200:9200 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  docker.elastic.co/elasticsearch/elasticsearch:8.13.4

# Start Kibana
podman run -d \
  --name kibana \
  --network esnet \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" \
  docker.elastic.co/kibana/kibana:8.13.4
```

> ✅ Elasticsearch: [http://localhost:9200](http://localhost:9200)
> ✅ Kibana: [http://localhost:5601](http://localhost:5601)

---

### 3. Prepare Dataset

* Place all your `.csv` files in the `dataset/` directory.
* Each file should follow a consistent format with fields like `job_title`, `annual_salary`, `location`, etc.

---

### 4. Run the Go API

```bash
gqlgen generate
go mod tidy
go run .
```

> ✅ The API runs on: [http://localhost:8085](http://localhost:8085)

On startup, the app will:

* Create the `compensations` index using `index.json` (if not exists)
* Bulk-upload all `.csv` files from the `dataset/` folder to Elasticsearch

--- 

## DB TradeOff
🔍 Why Elasticsearch Over PostgreSQL 
- Optimized for Search & Analytics
Elasticsearch is built for full-text search and analytical queries, making it ideal for use cases like filtering job titles, locations, and salary ranges.

- High Performance on Large Datasets
Uses inverted indexing and columnar storage to deliver faster query performance over large volumes of data compared to traditional RDBMS.

- Schema Flexibility
Supports semi-structured or evolving data models using dynamic mappings, avoiding the need for rigid schema definitions like in PostgreSQL.

- Advanced Aggregations
Provides out-of-the-box support for complex aggregations (e.g., average, min, max salary per location) that are performant even at scale.

- Real-Time Query Capabilities
Near real-time indexing and search capability allows immediate availability of uploaded or updated data for querying.

## 🖼️ Screenshots

Below are some screenshots of the application in action:

High Level Design 

![Design](https://github.com/user-attachments/assets/87967f41-b6c6-4a22-ba0a-1d2822e3dc55 )
*Elasticsearch and Kibana containers running in Podman*

![Elasticsearch & Kibana Running](https://github.com/user-attachments/assets/9fb19774-a71d-47d9-b836-98f40193304a)
*Elasticsearch and Kibana containers running in Podman*

![Kibana Dashboard](https://github.com/user-attachments/assets/25a77243-59f8-4252-ac68-16e73f41a9e2)


![Kibana Dashboard](https://github.com/user-attachments/assets/67884008-220b-40ef-8657-e480728609d5)
*Kibana dashboard showing loaded compensation data*

![GraphQL Playground](https://github.com/user-attachments/assets/ba668ecf-373b-459d-b1dd-8fb1d9ca31d0)
*GraphQL Playground for querying compensation data*

![REST API Example](https://github.com/user-attachments/assets/9a5f1b9e-9e43-48b3-a716-82a8f490a015)
*Example REST API request and response*

![Elasticsearch Index Mapping](https://github.com/user-attachments/assets/9f946182-7f9c-48c9-bc07-a8756c9e8e05)
*Elasticsearch index mapping for compensation data*

---

## 📡 API Usage

### 🔹 GraphQL

* **Playground**: [http://localhost:8085/](http://localhost:8085/)
* **Endpoint**: `/query`

#### Example GraphQL Query

```graphql
query {
  compensations(
    salaryGte: 120000
    location: "New York"
    sortBy: "annual_salary"
    sortOrder: "desc"
    limit: 5
  ) {
    id
    job_title
    annual_salary
    location
    currency
  }
}
```

You can:

* Filter by fields (e.g., salary, location, title)
* Sort results
* Limit output fields (sparse fieldset)

---

### 🔹 REST Endpoint

#### GET `/compensation_data?id=<doc_id>&fields=field1,field2,...`

##### Example Request:

```bash
curl "http://localhost:8085/compensation_data?id=abc123&fields=job_title,annual_salary"
```

##### Example Response:

```json
{
  "job_title": "Software Engineer",
  "annual_salary": 140000
}
```

> Retrieves a single document by ID with optional field filtering.

---

## 📁 Project Structure

```
graphql-es/
├── main.go                  # App entry point
├── rest_handler.go          # REST endpoint logic
├── index.json               # Elasticsearch index mapping
├── dataset/                 # Place your CSV files here
├── go.mod / go.sum          # Go module dependencies
├── graph/                   # GraphQL schema, resolvers
│   ├── schema.graphqls
│   ├── generated.go
│   ├── model/
│   └── resolver.go
└── internal/
    └── elastic/             # Elasticsearch client logic
```

---

## 🧪 Troubleshooting

* **Elasticsearch crashed with code 137?**
  ➤ Increase Podman memory or reduce `ES_JAVA_OPTS` to `-Xms256m -Xmx256m`.

* **To clean everything and restart:**

```bash
podman stop elasticsearch kibana
podman rm elasticsearch kibana
podman volume prune
```

---

## 👨‍💼 Development Guide

### Update GraphQL Schema

Edit the schema in:

```bash
graph/schema.graphqls
```

Then run:

```bash
go run github.com/99designs/gqlgen generate
```

This regenerates `model` and `resolver` files.

---

## 📜 License

MIT

---

## 🤝 Author

Built with ❤️ by **Vikas Pal**
Feel free to contribute, raise issues, or suggest improvements!
