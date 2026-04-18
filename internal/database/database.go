package database

import (
"context"
"database/sql"
"fmt"

_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	conn *sql.DB
}

func New(connectionURL string) (*DB, error) {
conn, err := sql.Open("pgx", connectionURL)

if err!= nil {
return nil, fmt.Errorf("opening db:%w", err)
}

if err := conn.PingContext(context.Background()); err!= nil {
return nil, fmt.Errorf("connecting to db: %w",err)
}

return &DB{conn:conn}, nil 

}

func (db *DB) CreateSchema(ctx context.Context) error {
_, err := db.conn.ExecContext( ctx, ` 
CREATE TABLE IF NOT EXISTS jobs(
id TEXT PRIMARY KEY,
status TEXT NOT NULL,
pages INTEGER DEFAULT 0,
error TEXT, 
seed TEXT NOT NULL,
depth INTEGER,
concurrency INTEGER,
created_at TIMESTAMPZ DEFAULT NOW()
)
`)
return err 
}
