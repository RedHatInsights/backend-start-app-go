package dao

import (
	"consoledot-go-template/internal/models"
	"context"
)

var GetHelloDao func(ctx context.Context) HelloDao

// HelloDao groups access methods for access to state of hello.
type HelloDao interface {
	List(ctx context.Context, limit, offset int64) ([]*models.Hello, error)
	Record(ctx context.Context, message *models.Hello) error
}
