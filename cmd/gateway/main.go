package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mcp-postgres-gateway/pkg/db"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // CRITICAL: Must be explicitly imported here to register the driver
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func formatResult(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func main() {
	err := godotenv.Load(".env")
	// Initialize Postgres Connection
	connStr := os.Getenv("DATABASE_URL")

	exposedTables := strings.Split(os.Getenv("EXPOSED_TABLES"), ",")
	if len(exposedTables) == 0 {
		log.Fatal("EXPOSED_TABLES environment variable is not set")
	}

	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Database initialization failure: %v", err)
	}
	defer database.Close()

	// Establish the MCP Core Server Block
	s := server.NewMCPServer("domainai-gateway", "1.0.0")

	// 1. Tool 1 Implementation: Expose Schema Table Information
	listTablesTool := mcp.NewTool("list_tables",
		mcp.WithDescription("Lists all available database schemas and field structures without exposing raw database records."),
	)

	s.AddTool(listTablesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cols, err := db.InspectExposedSchema(database, exposedTables)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to map system constraints: %s", err.Error())), nil
		}

		return mcp.NewToolResultText(formatResult(cols)), nil
	})

	// Data tools
	queries := db.NewQueries(database)

	metricsTool := mcp.NewTool("get_metrics",
		mcp.WithDescription("Retrieves metrics for specified geographical regions"),

		// Define your slice parameter here
		mcp.WithArray("region",
			mcp.Required(), // <-- This marks the parameter as required in the JSON Schema
			mcp.Description("A list of regions to filter metrics by (e.g. ['US', 'EU'])"),
		),
	)

	s.AddTool(metricsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		regions, err := request.RequireStringSlice("region")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result, err := queries.FindRegionalMetrics(ctx, regions)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(formatResult(result)), nil
	})

	// Start the Server to communicate natively over standard IO channels
	log.Println("MCP Gateway initialized. Establishing communication channel over Stdio...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server crash anomaly: %v\n", err)
		os.Exit(1)
	}
}
