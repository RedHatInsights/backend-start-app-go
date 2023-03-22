package seeds

import (
	"consoledot-go-template/internal/config"
	"consoledot-go-template/internal/db"
	"context"
	"embed"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

//go:embed *.sql
var EmbeddedSeeds embed.FS

var ErrSeedProduction = errors.New("seed in production")

// Seed executes embedded SQL scripts from internal/db/seeds
func Seed(ctx context.Context, seedScript string) error {
	logger := log.Logger.With().Bool("seed", true).Logger()
	logger.Debug().Msgf("Started execution of seed script %s", seedScript)

	// Prevent from accidental execution of drop_all seed in production
	if seedScript == "drop_all" && config.InClowder() {
		return fmt.Errorf("%w: an attempt to run drop_all seed script in clowder environment", ErrSeedProduction)
	}
	file, err := EmbeddedSeeds.Open(fmt.Sprintf("%s.sql", seedScript))
	if err != nil {
		return fmt.Errorf("unable to open seed script %s: %w", seedScript, err)
	}
	defer file.Close()
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read seed script %s: %w", seedScript, err)
	}
	_, err = db.Pool.Exec(ctx, string(buffer))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			logger.Fatal().Err(pgErr).Msgf("a PG error: %s", pgErr.Detail)
		} else {
			return fmt.Errorf("unable to execute script %s: %w", seedScript, err)
		}
	}

	logger.Info().Msgf("Executed seed script %s", seedScript)
	return nil
}
