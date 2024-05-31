package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type DB struct {
	db     *sql.DB
	ctx    context.Context
	cancel func()
	Dsn    string
	Now    func() time.Time
}

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

func NewDB(dsn string) *DB {
	ctx, cancel := context.WithCancel(context.Background())
	db := &DB{
		Now:    time.Now,
		Dsn:    dsn,
		ctx:    ctx,
		cancel: cancel,
	}
	return db
}

func (database *DB) Open() error {
	// Check Dsn is defined
	if database.Dsn == "" {
		return fmt.Errorf("dsn required")
	}

	conn, err := sql.Open("postgres", database.Dsn)
	if err != nil {
		return err
	}
	database.db = conn
	return nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}
