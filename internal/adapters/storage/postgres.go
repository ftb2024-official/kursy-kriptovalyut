package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"kursy-kriptovalyut/internal/entities"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connStr string) (*Postgres, error) {
	if connStr == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "empty connection string")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to connect to DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to ping DB: %v", err)
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) Store(ctx context.Context, coins []entities.Coin) error {
	query := "INSERT INTO coins (title, price) VALUES ($1, $2)"
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return errors.Wrapf(entities.ErrInternal, "failed to prepare query: %v", err)
	}
	defer stmt.Close()

	for _, coin := range coins {
		res, err := stmt.ExecContext(ctx, coin.Title, coin.Price)
		if err != nil {
			return errors.Wrapf(entities.ErrInternal, "failed to execute query: %v", err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return errors.Wrapf(entities.ErrInternal, "unexpected error: %v", err)
		}

		if rows != 1 {
			return errors.Wrapf(entities.ErrInternal, "expected to affect 1 row, affected %v row(s)", rows)
		}
	}

	return nil
}

func (p *Postgres) GetCoinsList(ctx context.Context) ([]string, error) {
	query := "SELECT DISTINCT title FROM coins"
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to prepare query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to execute query: %v", err)
	}
	defer rows.Close()

	var (
		titles []string
		title  string
	)

	for rows.Next() {
		if err := rows.Scan(&title); err != nil {
			return nil, errors.Wrapf(entities.ErrInternal, "failed to copy title: %v", err)
		}

		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "unexpected error: %v", err)
	}

	return titles, nil
}

func (p *Postgres) GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	query := "select title, price from coins where title = $1 order by desc limit 1"
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to prepare query: %v", err)
	}
	defer stmt.Close()

	var coin entities.Coin
	coins := make([]entities.Coin, 0, len(titles))

	for _, title := range titles {
		if err := stmt.QueryRowContext(ctx, title).Scan(&coin.Title, &coin.Price); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, errors.Wrapf(entities.ErrInternal, "failed to copy title/price: %v", err)
		}

		coins = append(coins, coin)
	}

	return coins, nil
}

func (p *Postgres) GetAggregateCoins(ctx context.Context, titles []string, aggFuncName string) ([]entities.Coin, error) {
	query := fmt.Sprintf("SELECT title, %v(price) FROM coins WHERE title = $1", aggFuncName)
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to prepare query: %v", err)
	}
	defer stmt.Close()

	var coin entities.Coin
	coins := make([]entities.Coin, 0, len(titles))

	for _, title := range titles {
		if err := stmt.QueryRowContext(ctx, query, title).Scan(&coin.Title, &coin.Price); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, errors.Wrapf(entities.ErrInternal, "failed to copy title/price: %v", err)
		}

		coins = append(coins, coin)
	}

	return coins, nil
}

// rows.Err() суть проверки после успешной итерации
