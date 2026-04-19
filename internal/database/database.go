package database

import (
"context"
"database/sql"
"fmt"
"github.com/ShwetaRoy17/GowebCrawler/internal/models"

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
		created_at TIMESTAMPTZ DEFAULT NOW()
	)
`)
return err 
}

func (db *DB) InsertJob(ctx context.Context, job *models.JobStatus) error {
	_, err := db.conn.ExecContext(ctx, `
	INSERT INTO jobs (id, status, pages, error, seed, depth, concurrency)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,job.ID, job.Status, job.Pages, job.Error, job.Seed, job.Depth, job.Concurrency)
	return err
}


func (db *DB) UpdateJob(ctx context.Context, job *models.JobStatus) error {
	_, err := db.conn.ExecContext(ctx, `
	UPDATE jobs SET status=$1, pages=$2, error=$3 WHERE id=$4`,
 job.Status, job.Pages, job.Error, job.ID)
	return err
}

func (db *DB) GetJob(ctx context.Context, jobID string) (*models.JobStatus, error) {
	row := db.conn.QueryRowContext(ctx, `SELECT id, status, pages, error, seed, depth, concurrency FROM jobs WHERE id=$1`, jobID)
	job := &models.JobStatus{}
	err := row.Scan(&job.ID, &job.Status, &job.Pages, &job.Error, &job.Seed, &job.Depth, &job.Concurrency)
	if err != nil {
		return nil, fmt.Errorf("scanning job row: %w", err)
	}
	return job, nil	
}


func (db *DB) ListJobs(ctx context.Context) ([]*models.JobStatus, error){
	rows, err := db.conn.QueryContext(ctx, `SELECT id, status, pages, error, seed, depth, concurrency FROM jobs ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("querying jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*models.JobStatus
	for rows.Next() {
		job := &models.JobStatus{}
		err := rows.Scan(&job.ID, &job.Status, &job.Pages, &job.Error, &job.Seed, &job.Depth, &job.Concurrency)
		if err != nil {
			return nil, fmt.Errorf("scanning job row: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (db *DB) LengthJobs(ctx context.Context) (int, error) {
	row := db.conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs`)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("counting jobs: %w", err)
	}
	return count, nil
}