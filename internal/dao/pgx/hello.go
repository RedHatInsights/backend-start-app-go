package pgx

import (
	"consoledot-go-template/internal/dao"
	"consoledot-go-template/internal/db"
	"consoledot-go-template/internal/models"
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func init() {
	dao.GetHelloDao = getHelloDao
}

type helloDaoPgx struct{}

func getHelloDao(ctx context.Context) dao.HelloDao {
	return &helloDaoPgx{}
}

func (x *helloDaoPgx) List(ctx context.Context, limit, offset int64) ([]*models.Hello, error) {
	query := `SELECT * FROM hellos ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query hellos error: %w", err)
	}

	var result []*models.Hello
	if err = pgxscan.ScanAll(&result, rows); err != nil {
		return nil, fmt.Errorf("scanning hello rows error: %w", err)
	}
	return result, nil
}

func (x *helloDaoPgx) Record(ctx context.Context, hello *models.Hello) error {
	query := `
		INSERT INTO hellos (from, to, message)
		VALUES ($1, $2, $3) RETURNING id`

	err := db.Pool.QueryRow(ctx, query, hello.From, hello.To, hello.Message).Scan(&hello.ID)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	return nil
}
