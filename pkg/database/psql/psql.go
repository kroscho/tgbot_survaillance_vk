package psql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type client struct {
	db     *sqlx.DB
	logger logrus.FieldLogger
}

// nolint:golint
func New(dsn string, logger logrus.FieldLogger) (*client, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "open postgres connection")
	}
	return &client{
		db:     db,
		logger: logger.WithFields(logrus.Fields{"component": "psql"}),
	}, nil
}

func (c client) GetConnection() DB {
	return c.db
}

func (c client) Close() error {
	return c.db.Close()
}
