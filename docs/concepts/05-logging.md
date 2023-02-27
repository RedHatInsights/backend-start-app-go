## Logging

This template introduces [zerolog](https://github.com/rs/zerolog) over platform recommended [logrus](https://github.com/sirupsen/logrus).
Until this is interchangeable Authors feel like zerolog is faster, easier to set up and use and has smaller memory footprint.

The global logger is set up at the start of our service.
For every request we will set up a new context logger and store it into a request context.
We use a logging [middleware](https://drstearns.github.io/tutorials/gomiddleware/) to do this.

This approach makes it easy to add additional fields and these are passed to all log entries.
Context should be kept very small as it gets passed very often through the whole stack.
Zerolog has such a small footprint that this is ok to do.

This approach standardize logging into a pattern of fetching logger from context and logging with help of this logger.
Every function can then be called in any context. We are always sure it logs correct log identifiers.

### Log output - cloudwatch

We need to set up logging output.
The template has two outputs.
Stdout writer used for development and for CI pipelines is very easy to initialize.

```go
stdWriter := zerolog.ConsoleWriter{
    Out:        os.Stdout,
    TimeFormat: time.Kitchen,
}
```

Second for production logging.

In production, we are logging to Amazon Cloudwatch.
From Cloudwatch, the platform team automatically pulls the logs to our Kibana.
The service responsibility is only to deliver logs to Cloudwatch.

To set up Cloudwatch writer, in our `logging.InitializeLogger` function,
we need Cloudwatch credentials, we will get to that in next chapter.


Following snippet initializes Cloudwatch writer assuming the credentials variables.

```go
func newCloudwatchWriter(region string, key string, secret string, session string, logGroup string, logStream string) (*io.Writer, error) {
	cache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(key, secret, session))
	cwClient := cloudwatchlogs.New(cloudwatchlogs.Options{
		Region:      region,
		Credentials: cache,
	})
    cloudWatchWriter, err := cloudwatchwriter2.NewWithClient(cwClient, 500*time.Millisecond, logGroup, logStream)
	return cloudWatchWriter, err
}
```

Now when we will be ready to decide whether we want to use Cloudwatch or Stdout,
we are ready to use the following to initialize our logger.

```go
zerolog.SetGlobalLevel(level)
//nolint:reassign
zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

output := initializeLogOutput() // here we will need to decide on the output used.
logger := zerolog.New(output)

// decorate logger (and thus every log line) with hostname and timestamp
logger = logger.With().Timestamp().Str("hostname", hostname).Logger()
```

### Log middleware

Middleware in Golang is a function that takes `next http.Handler` as parameter and returns http.Handler itself.
The inner http.Handler needs to call `next.ServeHTTP()` to invoke the following handler.
The simplest middleware is a function, that is called on `ServeHTTP()`,
go provides a helper [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) to create such a function

We further wrap that with a higher order function that allows us to pass logger.
In case we want to change logger we want to use, it would be easier to do then if the middleware would use the global logger directly.

The middleware enhances the global (passed in) logger
with the request context fields like remote IP, request path, HTTP method.
Then we log the very first line of every request.

```go
loggerCtx := globalLogger.With().
    Str("remote_ip", r.RemoteAddr).
    Str("url", r.URL.Path).
    Str("method", r.Method)
contextLogger := loggerCtx.Logger()
contextLogger.Debug().Msgf("Started %s request %s", r.Method, r.URL.Path)
```

The major thing we do here is to store our logger into the context that we pass down the middleware stack.
```go
// see logging.WithLogger()
ctx := WithLogger(r.Context(), &contextLogger)
next.ServeHTTP(ww, r.WithContext(ctx))
```

We follow up with deferring a function to log last log line of every request.
We run it deferred, effectively after all following middlewares and thus at the end of the request.

This function has one very special effect, it recovers from panic and logs it.
Every panic that happens further down the middleware stack is thus recovered and logged in our middleware.

```go
t1 := time.Now()
defer func() {
    duration := time.Since(t1)
    afterLogger := contextLogger.With().
        Dur("latency_ms", duration).
        Int("bytes_in", bytesIn).
        Int("bytes_out", ww.BytesWritten()).
        Logger()

    // prevent the application from exiting
    if rec := recover(); rec != nil {
        afterLogger.Error().
            Bool("panic", true).
            Int("status", http.StatusInternalServerError).
            Msgf("Unhandled panic: %s\n%s", rec, debug.Stack())
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }

    log.Info().
        Int("status", ww.Status()).
        Msgf("Completed %s request %s in %s with %d",
            r.Method, r.URL.Path, duration.Round(time.Millisecond).String(), ww.Status())
}()
```

### Add middleware to stack

Now the last thing we need to do is to add our middleware to the API router.

```go
// internal/routes/api_router.go

router.Use(logging.NewMiddleware(log.Logger))
```

And we are all setup.

### Using the logger

Now we have the logger in the context, we will be passing context through the app.
To use the logger, we will do following.

```go
logger := logging.Logger(ctx)

logger.Info().Msg("Message one.")
logger.Debug().Msg("Message two.")
```

Happy logging! :)