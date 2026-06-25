package db

import (
	"database/sql"
)

// ColumnMetadata defines a single column in our postgres database
type ColumnMetadata struct {
	TableName  string
	ColumnName string
	DataType   string
}

// InspectExposedSchema reads structural layout data dynamically from the system catalog.
func InspectExposedSchema(db *sql.DB, exposedTables []string) ([]ColumnMetadata, error) {
	query := `
		SELECT table_name, column_name, data_type 
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		AND table_name = ANY($1)
		ORDER BY table_name, column_name;`

	rows, err := db.Query(query, exposedTables)
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metadata, nil
}
