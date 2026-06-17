package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Queries struct {
	db *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
	return &Queries{
		db: db,
	}
}

type RegionalMetricsResult struct {
	Region          string  `json:"region"`
	UnitsSold       int     `json:"unitsSold"`
	TotalRevenue    float64 `json:"totalRevenue"`
	CancelledOrders int     `json:"cancelledOrders"`
}

func (q *Queries) FindRegionalMetrics(ctx context.Context, regions []string) ([]RegionalMetricsResult, error) {
	query := `SELECT compounds.region, 
       				sum(units_sold) AS TOTAL_UNITS_SOLD,
					sum(revenue_egp) AS TOTAL_REVENUE, 
					sum(cancelled_orders) AS TOTAL_CANCELLED
				FROM sales_ledger JOIN compounds ON sales_ledger.compound_id = compounds.id
				WHERE compounds.region = ANY($1)
				GROUP BY compounds.region`

	rows, err := q.db.QueryContext(ctx, query, pq.Array(regions))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []RegionalMetricsResult

	for rows.Next() {
		var col RegionalMetricsResult
		if err := rows.Scan(&col.Region, &col.UnitsSold, &col.TotalRevenue, &col.CancelledOrders); err != nil {
			return nil, err
		}

		result = append(result, col)
	}

	return result, nil
}
