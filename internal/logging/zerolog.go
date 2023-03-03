package logging

import (
	"consoledot-go-template/internal/config"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/lzap/cloudwatchwriter2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var hostname string

func init() {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown-hostname"
	}
	hostname = h
}

// InitializeLogger initializes the global logger with an output set by config.
// It panics on any setup error as it is hard to debug any further problems when logger is not set up.
// It returns a close function to close all the IO writers as second return parameter.
func InitializeLogger() (zerolog.Logger, func()) {
	level, err := zerolog.ParseLevel(config.Logging.Level)
	if err != nil {
		panic(fmt.Errorf("cannot parse log level '%s': %w", config.Logging.Level, err))
	}
	zerolog.SetGlobalLevel(level)
	//nolint:reassign
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	output, closeWriter := initializeLogOutput()
	logger := zerolog.New(output)

	// decorate logger with hostname and timestamp
	logger = logger.With().Timestamp().Str("hostname", hostname).Logger()

	return logger, closeWriter
}

func initializeLogOutput() (io.Writer, func()) {
	stdWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Kitchen,
	}
	// Create stdout output
	if config.Cloudwatch.Enabled {
		// see internal/config/helpers.go you can choose better identifier for you stream
		// this field is customizable, but when you start to deploy multiple containers,
		// binary is a good distinguishing factor which part of your service sent the log.
		stream := config.BinaryName()
		cwClient := newCloudwatchClient(config.Cloudwatch.Region, config.Cloudwatch.Key, config.Cloudwatch.Secret, config.Cloudwatch.Session)
		cloudWatchWriter, err := cloudwatchwriter2.NewWithClient(cwClient, 500*time.Millisecond, config.Cloudwatch.Group, stream)
		if err != nil {
			panic(fmt.Errorf("cannot initialize cloudwatch: %w", err))
		}
		return cloudWatchWriter, cloudWatchWriter.Close
	}
	return stdWriter, func() {}
}

func newCloudwatchClient(region string, key string, secret string, session string) *cloudwatchlogs.Client {
	cache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(key, secret, session))
	cwClient := cloudwatchlogs.New(cloudwatchlogs.Options{
		Region:      region,
		Credentials: cache,
	})
	return cwClient
}
