package services

import (
	"consoledot-go-template/internal/dao"
	"consoledot-go-template/internal/logging"
	"consoledot-go-template/internal/payloads"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// writeBasicError returns an error code without utilizing the Chi rendering stack. It can
// be used for fatal errors which happens during rendering pipeline (e.g. JSON errors).
func writeBasicError(w http.ResponseWriter, r *http.Request, err error) {
	if logger := logging.Logger(r.Context()); logger != nil {
		logger.Error().Msgf("unable to render error %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	wrappedMessage := ""
	if errors.Unwrap(err) != nil {
		wrappedMessage = errors.Unwrap(err).Error()
	}
	// write generic error body by string
	_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s", "error": "%s"}`, err.Error(), wrappedMessage)))
}

func renderError(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	if errRender := render.Render(w, r, renderer); errRender != nil {
		writeBasicError(w, r, errRender)
	}
}

func renderNotFoundOrDAOError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if errors.Is(err, dao.ErrNoRows) {
		renderError(w, r, payloads.NewNotFoundError(r.Context(), resource, err))
	} else {
		renderError(w, r, payloads.NewDAOError(r.Context(), resource, err))
	}
}
