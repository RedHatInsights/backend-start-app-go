# OpenAPI requests and responses

_Here we cover how to maintain OpenAPI spec and how to keep the API in accordance to the spec._

Serving OpenAPI spec is a hard requirement for every service on ConsoleDot.
It serves as an API contract, for other services to know what to expect from your service.
This spec should be always backward compatible starting your first release.
If you need to make incompatible changes,
you will need to keep the old spec, and it's corresponding API working to allow for migrations.
Please keep that in mind while designing the spec.
It's convenient to get it somewhat right on the first release.

Note: first release not the first commit, we do not encourage over designing your API upfront,
we know it will change as you implement it, but keep in mind to keep the OpenAPI as tidy as possible. :)

## Generate OpenAPI spec

There are generally three major approaches to the OpenAPI spec to keep it in sync with your API.

1. Generate API service from the spec.
2. Generate OpenAPI spec from the services
3. Hybrid approach

There are advantages to all of these approach, and it's up to you to choose one.
This service template covers hybrid approach as some services might leverage the control it gives to maintainers.

_Note: hybrid approach has more variants, feel free to explore if the choice makes you uncomfortable._

Hybrid approach in our case means we will generate request and response types from Go to OpenAPI.
The rest - paths, various response codes are kept manually.
This gives the highest level of freedom while still taking care of the most of the manual burden
behind keeping the types in sync from Go to OpenAPI.

### Payloads

To maintain our request and response types, we will use payload types that implement
Chi's render interface, see [`go-chi/render`](https://github.com/go-chi/render) for more details.

We will store our payloads in separate folder, these will use models to generate responses.
We can easily use type embedding for the payloads and just copy the model fields, but it is discouraged.
You might accidentally copy and expose fields you don't want. We will be copying the data instead.

We want our users to be able to define a sender of a message and a message,
but we will define recipient statically, and not allow setting it through the create request.
We are also renaming the internal columns to expose them under different names in the API
to showcase what you might leverage the payloads for.

Payload is a go struct that has json tag for bindings from and to JSON.
Following is a basic example where request and response have the same data.
There is more advance example with comments in `hello_payload.go`.

```go
// internal/payloads/hello_payload.go

type HelloPayload struct {
	ID      uint64 `json:"id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
type (
    HelloRequest  HelloPayload
    HelloResponse HelloPayload
)

func (req HelloRequest) Bind(_ *http.Request) error {
    return nil
}

func (req HelloResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
    return nil
}
```

#### Using payloads

We use payloads in the http handlers to bind and render data to and from JSON.

First lets take a look at binding data from json request to Go struct.

```go
// internal/services/hello_service.go

import "github.com/go-chi/render"

func SayHello(w http.ResponseWriter, r *http.Request) {
	payload := payloads.HelloRequest{}

	if err := render.Bind(r, payload); err != nil {
		// Error handling TBD
		return
	}
}
```

And now let see rendering a JSON response.

```go
func ListHellos(w http.ResponseWriter, r *http.Request) {
	helloDao := dao.GetHelloDao(r.Context())
	hellos, err := helloDao.List(r.Context(), 100, 0)
	// error handling TBD

	if renderErr := render.RenderList(w, r, payloads.NewHelloListResponse(hellos)); renderErr != nil {
		// error handling TBD
	}
}
```

### Error handling

In above examples, we have left out the error handling.
Let us dive into it here.
Errors are just another type of payload we want to emit to user.
The main difference here is, that we want to standardize the payload structure across all handlers.

We will have a helper render function to render the error payloads.
As rendering the payload itself can error out,
we want to have a fallback to general 500 error rendering in this helper.

```go
// internal/services/error_renderer.go

func renderError(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	if renderErr := render.Render(w, r, renderer); renderErr != nil {
		writeBasicError(w, r, renderErr) // this is a fallback
	}
}
```

Our helper expects a chi renderer, so we need a payload to pass in.
For this purpose we introduce error payloads in `internal/payloads/error_payload.go`.
The following is a basis for all our error payloads.

```go
// internal/payloads/error_payload.go

// ResponseError is used as a payload for all errors
type ErrorResponse struct {
    // HTTP status code
    HTTPStatusCode int `json:"-"`
    // user facing error message
    Message string `json:"msg"`
    // full root cause
    Error string `json:"error"`
}

func (e ErrorResponse) Render(_ http.ResponseWriter, r *http.Request) error {
    render.Status(r, e.HTTPStatusCode)
    return nil
}

func newErrorResponse(ctx context.Context, status int, userMsg string, err error) ErrorResponse {
    return &ErrorResponse{
        HTTPStatusCode: status,
        Message:        userMsg,
        Error:          err.Error(),
    }
}
```

In previous code snippets, we have seen three types of errors.
Bad request error, which should be status 400.
Database error, which can be not found and should have status 404, or all other DB errors status 500.
General code error and redner error, which are both 500 status.

Let's add payload constructors for these errors.
Notice that all of these have user message.
The go error holds very useful debugging info and can be quite helpful to understand
what went technically wrong.
Usually we want to convey the main message to users and make it more user readable.
That's where the user message comes in.

```go
func NewInvalidRequestError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Invalid request: %s", message)
	return newResponse(ctx, http.StatusBadRequest, message, err)
}

func NewDAOError(ctx context.Context, message string, err error) *ResponseError {
    message = fmt.Sprintf("DAO error: %s", message)
    return newResponse(ctx, http.StatusInternalServerError, message, err)
}

func NewRenderError(ctx context.Context, message string, err error) *ResponseError {
    message = fmt.Sprintf("Rendering error: %s", message)
    return newResponse(ctx, http.StatusInternalServerError, message, err)
}
```

We want to keep these simple, so for the dao error, we will introduce one more render helper.
Here we will wrap the decision whether to render NotFound or Internal error.

```go
// internal/services/error_renderer.go

func renderNotFoundOrDAOError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if errors.Is(err, dao.ErrNoRows) {
		renderError(w, r, payloads.NewNotFoundError(r.Context(), resource, err))
	} else {
		renderError(w, r, payloads.NewDAOError(r.Context(), resource, err))
	}
}
```

Take a look at the code for the full implementation.

### OpenAPI generator

We will add another binary `openapi_spec` for generating our OpenAPI spec.

We will use a [`kin-openapi`](github.com/getkin/kin-openapi/)'s generator to do the heavy lifting for us.

We won't go in details of the generating itself.
There are two notable aspect we need to pay attention to.

We have a file `cmd/openapi_spec/paths.yml` alongside the binary
to manually maintain routes and used request and response schemas.
In the binary `main.go` itself we need to list all the payload schemas in method `addPayloads`.
It is kept as first function for convenience.

The spec will be generated to both JSON and YAML.
It is straight forward to get rid one of the formats if you don't find it useful.

At this point it should be straight forward to add more types or routes to our spec.

## OpenAPI spec endpoint

The app should serve it's current spec for easier integration with some apps.
We will use embedding in the `/api` package and create a http handler,
that will write the generated json spec.

This is all we need:

```go
// api/openapi_handler.go

func ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(embeddedJSONSpec); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"`, err.Error())))
	}
}
```

To return this spec on a path `/api/prefix/v1/openapi.json`
we will create a new route in the api router.

```go
router.Get("/openapi.json", api.ServeOpenAPISpec)
```