# MCP Postgres Gateway

This repository implements a Model Context Protocol (MCP) server in Go that connects to a PostgreSQL database to securely expose schema structures and regional metrics to LLMs.

---

## Project Structure & Core Files

* **`seed.sql`**: Defines the database schema for the `compounds` and `sales_ledger` tables, and seeds the initial records for various real estate developments.
* **`postgres.go`**: Contains the `InspectExposedSchema` function, which dynamically reads table column layouts directly from the PostgreSQL system catalog.
* **`queries.go`**: Houses the `FindRegionalMetrics` function used to execute analytical queries tracking units sold, total revenue in EGP, and cancellations across specified regions.
* **`docker-compose.yml`**: Sets up a local PostgreSQL 15 container instance and automatically maps `seed.sql` to handle database initialization on boot.
* **`main.go`**: Serves as the primary application entry point. It parses configurations from the environment, initializes the database connection, maps the core MCP tools, and hosts the Stdio communication channel.
* **`.env.example`**: Provides a baseline template for configuring environment variables like `DATABASE_URL` and `EXPOSED_TABLES`.

---

## Local Setup Instructions

### 1. Configure the Environment
Duplicate the environment template file into a local `.env` file:

```bash
cp .env.example .env
```

Open `.env` and verify that the database credentials and target connection strings align with your environment.

### 2. Launch the Database
Spin up the localized, pre-seeded PostgreSQL instance via Docker Compose:

```bash
docker-compose up -d
```

This command builds the container and runs `seed.sql` to build out your schemas and records seamlessly.

### 3. Run and Verify with the MCP Inspector
To execute the application locally and inspect the exposed protocol features, boot up the official server inspector using the following command:

```bash
npx -y @modelcontextprotocol/inspector go run cmd/gateway/main.go
```

---

## Available Protocol Tools

Once connected, the gateway exposes two core capabilities to your MCP environment:

* **`list_tables`**: Explores the database structural metadata dynamically using `postgres.go` without returning actual raw rows.
* **`get_metrics`**: Consumes an array of regional name filters to run analytical aggregations over sales performance data using `queries.go`.
