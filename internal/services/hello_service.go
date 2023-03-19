package services

import (
	"consoledot-go-template/internal/dao"
	"consoledot-go-template/internal/models"
	"consoledot-go-template/internal/payloads"
	"net/http"

	"github.com/go-chi/render"
)

// static Recipient
const Recipient = "Ondrej Ezr<oezr@redhat.com"

func ListHellos(w http.ResponseWriter, r *http.Request) {
	helloDao := dao.GetHelloDao(r.Context())
	hellos, err := helloDao.List(r.Context(), 100, 0)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "list hellos", err))
		return
	}

	if renderErr := render.RenderList(w, r, payloads.NewHelloListResponse(hellos)); renderErr != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render hello list", renderErr))
	}
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	payload := payloads.HelloRequest{}
	if err := render.Bind(r, &payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "say hello", err))
		return
	}

	hello := models.Hello{To: Recipient, From: payload.Sender, Message: payload.Message}

	helloDao := dao.GetHelloDao(r.Context())
	if err := helloDao.Record(r.Context(), &hello); err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "record hello", err))
		return
	}

	render.Status(r, http.StatusCreated)
	if rndrErr := render.Render(w, r, payloads.NewHelloResponse(&hello)); rndrErr != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render hello", rndrErr))
	}
}
