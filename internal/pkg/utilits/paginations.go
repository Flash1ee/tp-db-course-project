package utilits

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"tp-db-project/internal/app"
	"tp-db-project/internal/app/models"
)

const (
	queryStat = "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1"
)

func AddPagination(tableName string, pag *models.Pagination, db *pgxpool.Pool) (limit int64, offset int64, err error) {
	var numberRows int64
	if err = db.QueryRow(context.Background(), queryStat, tableName).Scan(&numberRows); err != nil {
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
