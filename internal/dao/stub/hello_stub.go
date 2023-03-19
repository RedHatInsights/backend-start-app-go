package stub

import (
	"consoledot-go-template/internal/dao"
	"consoledot-go-template/internal/models"
	"context"
)

func init() {
	dao.GetHelloDao = getHelloDao
}

type helloDaoStub struct {
	store []*models.Hello
}

func getHelloDao(ctx context.Context) dao.HelloDao {
	return getHelloDaoStub(ctx)
}

func (x *helloDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Hello, error) {
	return x.store, nil
}

func (x *helloDaoStub) Record(ctx context.Context, hello *models.Hello) error {
	hello.ID = int64(len(x.store))
	x.store = append(x.store, hello)
	return nil
}
