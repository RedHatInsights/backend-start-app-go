package dao

import (
	"github.com/jackc/pgx/v5"
)

// ErrNoRows is returned when there are no rows in the result
var ErrNoRows = pgx.ErrNoRows
