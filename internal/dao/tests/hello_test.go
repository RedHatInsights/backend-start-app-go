//go:build database
// +build database

package tests

import (
	"consoledot-go-template/internal/dao"
	"consoledot-go-template/internal/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupHelloDao(t *testing.T) (dao.HelloDao, context.Context) {
	ctx := context.Background()
	return dao.GetHelloDao(ctx), ctx
}

func newHello() *models.Hello {
	return &models.Hello{
		From:    "test@example.com",
		To:      "another@example.com",
		Message: "Test greeting",
	}
}

func TestHelloRecord(t *testing.T) {
	helloDao, ctx := setupHelloDao(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		hello := newHello()
		err := helloDao.Record(ctx, hello)
		require.NoError(t, err)

		assert.Greater(t, 0, hello.ID)
	})
}
