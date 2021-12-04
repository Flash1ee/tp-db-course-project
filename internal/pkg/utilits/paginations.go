package utilits

import (
	"database/sql"
	"tp-db-project/internal/app"
	"tp-db-project/internal/app/models"
)

const (
	queryStat = "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1"
)

func AddPagination(tableName string, pag *models.Pagination, db *sql.DB) (limit int64, offset int64, err error) {
	var numberRows int64
	if err = db.QueryRow(queryStat, tableName).Scan(&numberRows); err != nil {
		return app.InvalidInt, app.InvalidInt, err
	}

	numberRows -= pag.Limit
	if pag.Offset < numberRows {
		numberRows = pag.Offset
	}
	if numberRows < 0 {
		numberRows = 0
	}
	return pag.Limit, numberRows, nil
}
