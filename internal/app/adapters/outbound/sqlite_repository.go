package outbound

import (
	"codea-it/internal/app/domain/models"
	"codea-it/internal/app/ports"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() ports.Repository {
	db, err := sql.Open("sqlite3", "./data/mydb.sqlite")
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		return nil
	}

	// Crear las tablas
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS records (
            pairs_names TEXT PRIMARY KEY NOT NULL,
            expiration_date DATETIME,
            status TEXT
        );
        CREATE TABLE IF NOT EXISTS data (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            pairs_name TEXT,
            pair TEXT,
            ltp TEXT,
            FOREIGN KEY (pairs_name) REFERENCES records(pairs_names)
        );
    `)
	if err != nil {
		log.Fatalf("Error creating tables: %q", err)
		return nil
	}

	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) SaveRecord(pairsNames string, expirationDate time.Time, status string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO records (pairs_names, expiration_date, status) VALUES (?, ?, ?)`, pairsNames, expirationDate, status)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLiteRepository) GetRecord(pairsNames string) (models.Record, error) {
	rows := r.db.QueryRow(`SELECT * FROM records WHERE pairs_names = ?`, pairsNames)

	var result models.Record
	if err := rows.Scan(&result.PairsNames, &result.ExpirationDate, &result.Status); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return models.Record{}, nil
		}
		return models.Record{}, err
	}

	return result, nil
}

func (r *SQLiteRepository) UpdateRecord(pairsNames string, expirationDate time.Time, status string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE records SET expiration_date = ?, status = ? WHERE pairs_names = ?`, expirationDate, status, pairsNames)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLiteRepository) SavePairs(pairsNames string, pair string, ltp string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO data (pairs_name, pair, ltp) VALUES (?, ?, ?)`, pairsNames, pair, ltp)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLiteRepository) GetPairs(pairsNames string) ([]models.Pair, error) {
	rows, err := r.db.Query(`SELECT pair, ltp FROM data WHERE pairs_name = ?`, pairsNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Pair
	for rows.Next() {
		var pair string
		var ltp string
		if err := rows.Scan(&pair, &ltp); err != nil {
			return nil, err
		}

		result = append(result, models.Pair{
			Pair: pair,
			Ltp:  ltp,
		})
	}

	return result, nil
}

func (r *SQLiteRepository) UpdatePairs(pairsNames string, pair string, ltp string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE data SET ltp = ? WHERE pairs_name = ? AND pair = ?`, ltp, pairsNames, pair)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
