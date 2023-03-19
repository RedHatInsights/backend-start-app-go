package services_test

import (
	"bytes"
	"consoledot-go-template/internal/dao"
	"consoledot-go-template/internal/dao/stub"
	"consoledot-go-template/internal/services"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListHellos(t *testing.T) {
	t.Run("handles empty database well", func(t *testing.T) {
		ctx := stub.WithHelloDao(context.Background())

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/template/hellos", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.ListHellos)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Wrong status code")
		assert.Equal(t, "[]\n", rr.Body.String())
	})
}

func TestSayHello(t *testing.T) {
	t.Run("records hello with a static recipient", func(t *testing.T) {
		ctx := stub.WithHelloDao(context.Background())
		hDao := dao.GetHelloDao(ctx)

		values := map[string]interface{}{
			"message": "hello beautiful Open Source world!",
			"sender":  "test@example.com",
		}
		jsonData, err := json.Marshal(values)
		require.NoError(t, err, "unable to marshal input payload to JSON")

		req, err := http.NewRequestWithContext(ctx, "POST", "/api/template/hellos", bytes.NewBuffer(jsonData))
		require.NoError(t, err, "failed to create request")
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.SayHello)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code, "Wrong status code")

		hellos, listErr := hDao.List(ctx, 100, 0)
		require.NoError(t, listErr, "failed to list hellos")

		assert.Equal(t, 1, len(hellos))
		assert.Equal(t, "test@example.com", hellos[0].From)
		assert.Equal(t, services.Recipient, hellos[0].To)
	})
}
