package dbutil

import "gorm.io/gorm"

// QuoteTableName quotes a table name based on database dialect to prevent SQL injection
// and handle reserved keywords properly.
//
// MySQL uses backticks: `table_name`
// PostgreSQL uses double quotes: "table_name"
// Other databases use unquoted names
func QuoteTableName(db *gorm.DB, tableName string) string {
	if db == nil || tableName == "" {
		return tableName
	}
	switch db.Name() {
	case "mysql":
		return "`" + tableName + "`"
	case "postgres":
		return `"` + tableName + `"`
	default:
		return tableName
	}
}
