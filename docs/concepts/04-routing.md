## HTTP Routing

Here we are diving into the implementation of the service itself.

For routing this service template uses [go-chi](https://go-chi.io/#/) library, as it is the most lightweight, but still does all we need.
You can of course choose a different one if you like and the concepts will be very similar.

We will place all relevant routing code under `/internal/routes`.

Here we will not talk about metrics router, that is covered in Metrics Concept.

### Root router

The entry point to the routing is `routes.RootRouter()` this sets up the root path router.

This is `/` path on our server.
This service uses this path for Liveness and Readiness probes of k8s.
We will get to this  later.

For now the important part is that within this router we mount another API router onto a prefixed path.

### API router

Our service will coexist with other services behind a shared public gateway.
The gateway routes traffic by route prefix, but does not strip the matched prefix.
So if our service will match for traffic on `/api/template` prefix, the full path gets forwarded.
This means all our public facing paths need to be exposed with this prefix.

The prefix is determined by `routes.PathPrefix` function.

### Individual routes

All routes have their matching pattern and a [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc).
Handler function is a function that receives a request as a parameter and writes response in a specialized IO.
Our use of handlers is basically Controller from an MVC apps.
We will put all our handlers in a `/internal/services` folder, more on that later.

Chi routes are defined by calling function named by their matching HTTP verb.
For example `router.Get("/hello", Handler)` creates a GET route for pattern `/hello`.

### Grouping routes

We can group routes under one prefix and only define the suffix for them.
To define the route with a prefix we would do following.

```go
router.Route("/prefix", func(subRouter chi.Router) {
	subRouter.Get("/hello", Handler)
})
```

### HTTP server

To put this all together, we will put following code in our `api` entry point and start a http server.

```go
// cmd/api/main.go

rootRouter := routes.RootRouter()

apiServer := http.Server{
    Addr:    fmt.Sprintf(":%d", 8000),
    Handler: rootRouter,
}
```

Later we will make this port configurable.

### Start listening

Following code starts up the server we've set up.
It will write out error message unless the server has been stopped gracefully. 

```go
if err := apiServer.ListenAndServe(); err != nil {
    if !errors.Is(err, http.ErrServerClosed) {
        log.Fatal().Err(err).Msg("Main service listen error")
    }
}
```

### Stop listening

Here we will cover graceful shutdown of our server.
We won't cover all the details of following code.

We set up a channel that will get notified when we receive `SIGINT` or `SIGTERM`.
The signal wait is blocking, so we are waiting for it in separate goroutine.

We have another channel that we wait for in our main and once this gets triggered, we finish.