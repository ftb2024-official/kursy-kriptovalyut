package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"kursy-kriptovalyut/internal/entities"
	"kursy-kriptovalyut/pkg/logger"
)

var log = logger.NewLogger()

type Postgres struct {
	dbPool *pgxpool.Pool
}

func NewPostgres(connStr string) (*Postgres, error) {
	if connStr == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "empty connection string")
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to create pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to ping DB: %v", err)
	}

	return &Postgres{dbPool: pool}, nil
}

func (p *Postgres) Store(ctx context.Context, coins []entities.Coin) error {
	query := "INSERT INTO coins (title, price) VALUES ($1, $2)" // пересмотреть вторую часть VALUES ($1, $2)

	log.Info("(Store) inserting coins:", zap.Any("coins", coins))
	for _, coin := range coins {
		res, err := p.dbPool.Exec(ctx, query, coin.Title, coin.Price)
		if err != nil {
			log.Error("(Store) failed to insert:", zap.Any("coin", coin))
			return errors.Wrapf(entities.ErrInternal, "failed to execute query: %v", err)
		}

		rows := res.RowsAffected()
		if rows != 1 {
			log.Error("(Store) affected more than one row")
			return errors.Wrapf(entities.ErrInternal, "expected to affect 1 row, affected %v row(s)", rows)
		}
	}

	return nil
}

func (p *Postgres) GetCoinsList(ctx context.Context) ([]string, error) {
	query := "SELECT DISTINCT title FROM coins"

	log.Info("(GetCoinsList) getting list of coin titles")
	rows, err := p.dbPool.Query(ctx, query)
	if err != nil {
		log.Error("(GetCoinsList) failed to get list of coin titles")
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("(GetCoinsList) empty result set")
			return nil, errors.Wrapf(entities.ErrNotFound, "empty result: %v", err.Error())
		}
		return nil, errors.Wrapf(entities.ErrInternal, "failed to execute query: %v", err)
	}
	defer rows.Close()

	var (
		titles []string
		title  string
	)

	for rows.Next() {
		if err := rows.Scan(&title); err != nil {
			log.Error("(GetCoinsList) failed to copy title:", zap.Any("title", title))
			return nil, errors.Wrapf(entities.ErrInternal, "failed to copy title: %v", err)
		}

		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		log.Error("(GetCoinsList) unexpected error:", zap.Any("err", err))
		return nil, errors.Wrapf(entities.ErrInternal, "unexpected error: %v", err)
	}

	return titles, nil
}

func (p *Postgres) GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	query := "SELECT title, price FROM coins WHERE title = $1 AND created_at = CURRENT_DATE ORDER BY created_at DESC LIMIT 1"

	var coin entities.Coin
	coins := make([]entities.Coin, 0, len(titles))

	log.Info("(GetActualCoins) getting coins' last rates", zap.Any("coinTitles", titles))
	for _, title := range titles {
		if err := p.dbPool.QueryRow(ctx, query, title).Scan(&coin.Title, &coin.Price); err != nil {
			log.Info("(GetActualCoins) failed to get coin:", zap.Any("coin", coin))
			if errors.Is(err, pgx.ErrNoRows) {
				log.Info("(GetActualCoins) empty result set")
				continue
			}

			log.Error("(GetActualCoins) failed to copy coin:", zap.Any("coin", coin))
			return nil, errors.Wrapf(entities.ErrInternal, "failed to copy title/price: %v", err)
		}

		coins = append(coins, coin)
	}

	return coins, nil
}

func (p *Postgres) GetAggregateCoins(ctx context.Context, titles []string, aggFuncName string) ([]entities.Coin, error) {
	query := fmt.Sprintf("SELECT title, %v(price) FROM coins WHERE title = $1 AND created_at = CURRENT_DATE GROUP BY title", aggFuncName)

	var coin entities.Coin
	coins := make([]entities.Coin, 0, len(titles))

	log.Info("(GetAggregateCoins) getting coins' aggregated rates", zap.Any("coinTitles", titles))
	for _, title := range titles {
		if err := p.dbPool.QueryRow(ctx, query, title).Scan(&coin.Title, &coin.Price); err != nil {
			log.Info("(GetAggregateCoins) failed to get aggregated coin", zap.Any("aggCoin", coin))
			if errors.Is(err, pgx.ErrNoRows) {
				log.Info("(GetAggregateCoins) empty result set")
				continue
			}
			log.Error("(GetAggregateCoins) failed to copy coin:", zap.Any("coin", coin))
			return nil, errors.Wrapf(entities.ErrInternal, "failed to copy title/price: %v", err)
		}

		coins = append(coins, coin)
	}

	return coins, nil
}
