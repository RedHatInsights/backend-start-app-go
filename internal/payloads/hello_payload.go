package payloads

import (
	"consoledot-go-template/internal/models"
	"net/http"

	"github.com/go-chi/render"
)

type HelloPayload struct {
	ID      uint64 `json:"id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type HelloRequest struct {
	HelloPayload
}

type HelloResponse struct {
	HelloPayload
	Recipient string `json:"recipient"`
}

// Bind is called by Chi to adjust the request payload data to your needs.
// all basic binding is done by chi, so unless you need anything special, just return nil here.
func (req HelloRequest) Bind(_ *http.Request) error {
	// ID is read-only field
	// this is to showcase how you'd go about embedding the full model and protect only some of its fields
	// this method has obvious pitfalls when you forget to add this protection.
	req.ID = 0
	return nil
}

// Render is called by Chi to adjust the response payload data.
// Chi does the json marshaling for you, so leave empty unless you need anything special.
func (req HelloResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewHelloResponse(hello *models.Hello) render.Renderer {
	return HelloResponse{
		HelloPayload: HelloPayload{
			Sender:  hello.From,
			Message: hello.Message,
		},
		Recipient: hello.To,
	}
}

func NewHelloListResponse(hellos []*models.Hello) []render.Renderer {
	list := make([]render.Renderer, len(hellos))
	for i, hello := range hellos {
		list[i] = NewHelloResponse(hello)
	}
	return list
}
