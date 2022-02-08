package query

import (
	"database/sql"
	"strings"
	"time"

	"github.com/Rican7/retry/jitter"
	"github.com/canonical/go-dqlite/driver"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/lxc/lxd/shared/logger"
)

const maxRetries = 250

// Retry wraps a function that interacts with the database, and retries it in
// case a transient error is hit.
//
// This should by typically used to wrap transactions.
func Retry(f func() error) error {
	// TODO: the retry loop should be configurable.
	var err error
	for i := 0; i < maxRetries; i++ {
		err = f()
		if err != nil {
			// No point in re-trying or logging a no-row error.
			if errors.Cause(err) == sql.ErrNoRows {
				break
			}

			// Process actual errors.
			if IsRetriableError(err) {
				if i == maxRetries {
					logger.Warn("Database error, giving up", "attempt", i, "err", err)
					break
				}
				logger.Debug("Database error, retrying", "attempt", i, "err", err)
				time.Sleep(jitter.Deviation(nil, 0.8)(100 * time.Millisecond))
				continue
			} else {
				logger.Debug("Database error", "err", err)
			}
		}
		break
	}

	return err
}

// IsRetriableError returns true if the given error might be transient and the
// interaction can be safely retried.
func IsRetriableError(err error) bool {
	err = errors.Cause(err)
	if err == nil {
		return false
	}

	if err, ok := err.(driver.Error); ok && err.Code == driver.ErrBusy {
		return true
	}

	if err == sqlite3.ErrLocked || err == sqlite3.ErrBusy {
		return true
	}

	if strings.Contains(err.Error(), "database is locked") {
		return true
	}

	if strings.Contains(err.Error(), "cannot start a transaction within a transaction") {
		return true
	}

	if strings.Contains(err.Error(), "bad connection") {
		return true
	}

	if strings.Contains(err.Error(), "checkpoint in progress") {
		return true
	}

	return false
}
