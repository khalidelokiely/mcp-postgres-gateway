package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type ColumnMetadata struct {
	TableName  string
	ColumnName string
	DataType   string
}

// InspectExposedSchema reads structural layout data dynamically from the system catalog.
func InspectExposedSchema(db *sql.DB) ([]ColumnMetadata, error) {
	query := `
		SELECT table_name, column_name, data_type 
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		AND table_name IN ('compounds', 'sales_ledger')
		ORDER BY table_name, column_name;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []ColumnMetadata
	for rows.Next() {
		var col ColumnMetadata
		if err := rows.Scan(&col.TableName, &col.ColumnName, &col.DataType); err != nil {
			return nil, err
		}
		metadata = append(metadata, col)
	}

	return metadata, nil
}
