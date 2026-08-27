package postgres

import (
	"context"
	"cryptocurrency/internal/cases"
	"cryptocurrency/internal/entity"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

func NewRepository(ctx context.Context, dataBaseURL string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, dataBaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "PostgreSQL repository: create connection pool")
	}
	if err := pool.Ping(ctx); err != nil {
		defer pool.Close()
		return nil, errors.Wrap(err, "PostgreSQL repository: ping connection pool")
	}
	return &Repository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (r *Repository) SaveCoinPrices(ctx context.Context, coins []entity.Coin) error {
	if len(coins) == 0 {
		return errors.Wrap(entity.ErrInvalidParams, "PostgresSQL repository SaveCoinPrices: coins is empty")
	}
	query := r.sq.
		Insert("coin_prices").
		Columns("title", "price", "creation_time")

	for _, coin := range coins {
		query = query.Values(
			coin.Title,
			coin.Price,
			coin.CreationTime)
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "PostgresSQL repository SaveCoinPrices: failed to generate sql")
	}
	if _, err := r.pool.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "PostgresSQL repository SaveCoinPrices: failed to execute sql")
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, titles []string, opts ...cases.Option) ([]entity.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entity.ErrInvalidParams, "PostgresSQL repository Get: titles is empty")
	}

	options := cases.NewOptions(opts...)
	var (
		query string
		args  []any
		err   error
	)
	switch options.Mode {
	case cases.ModeLatestPrices:
		query, args, err = r.sq.
			Select("DISTINCT ON (title) title", "price", "creation_time").
			From("coin_prices").
			Where(squirrel.Eq{"title": titles}).
			OrderBy("title", "creation_time DESC").
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetLatestPrices: failed to generate sql")
		}
	case cases.ModeMinPrices:
		since := time.Now().Add(-24 * time.Hour)
		query, args, err = r.sq.
			Select("DISTINCT ON (title) title, price, creation_time").
			From("coin_prices").
			Where(squirrel.Eq{"title": titles}).
			Where(squirrel.GtOrEq{"creation_time": since}).
			OrderBy("title", "price ASC").
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetMinPrices: failed to generate sql")
		}
	case cases.ModeMaxPrices:
		since := time.Now().Add(-24 * time.Hour)
		query, args, err = r.sq.
			Select("DISTINCT ON (title) title, price, creation_time").
			From("coin_prices").
			Where(squirrel.Eq{"title": titles}).
			Where(squirrel.GtOrEq{"creation_time": since}).
			OrderBy("title", "price DESC").
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetMinPrices: failed to generate sql")
		}
	case cases.ModePriceChangePercent:
		since := time.Now().Add(-time.Hour)
		latestQuery, latestArgs, err := r.sq.
			Select("DISTINCT ON (title) title", "price", "creation_time").
			From("coin_prices").
			Where(squirrel.Eq{"title": titles}).
			OrderBy("title", "creation_time DESC").
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetPriceChangePercent: failed to generate sql")
		}
		hourAgoQuery, hourAgoArgs, err := r.sq.
			Select("DISTINCT ON (title) title", "price", "creation_time").
			From("coin_prices").
			Where(squirrel.Eq{"title": titles}).
			Where(squirrel.GtOrEq{"creation_time": since}).
			OrderBy("title", "creation_time ASC").
			ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetPriceChangePercent: failed to generate sql")
		}
		query = `
			WITH latest AS (` + latestQuery + `),
			hour_ago AS (` + hourAgoQuery + `)
			SELECT
				latest.title,
				(latest.price - hour_ago.price)
					/ hour_ago.price * 100,
				latest.creation_time
			FROM latest
			JOIN hour_ago
				ON hour_ago.title = latest.title
			WHERE hour_ago.price <> 0
		`

		args = append(latestArgs, hourAgoArgs...)
	default:
		return nil, errors.Wrap(entity.ErrInvalidParams, "PostgresSQL repository GetLatestPrices: unsupported mode")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "PostgresSQL repository GetPriceChangePercent: failed to execute sql")
	}
	defer rows.Close()

	var coins []entity.Coin
	for rows.Next() {
		var (
			title        string
			price        float64
			creationTime time.Time
		)
		if err := rows.Scan(
			&title,
			&price,
			&creationTime); err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetPriceChangePercent: failed to scan")
		}

		coin, err := entity.NewCoin(title, price, creationTime)
		if err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetPriceChangePercent: failed to create coin")
		}
		coins = append(coins, *coin)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "PostgresSQL repository GetPriceChangePercent: failed to fetch rows")
	}
	return coins, nil
}

func (r *Repository) GetTitles(ctx context.Context) ([]string, error) {
	query, args, err := r.sq.
		Select("title").
		From("coins_tracked").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "PostgresSQL repository GetTitles: failed to generate sql")
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "PostgresSQL repository GetTitles: failed to execute sql")
	}
	defer rows.Close()
	titles := make([]string, 0)
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, errors.Wrap(err, "PostgresSQL repository GetTitles: failed to scan")
		}
		titles = append(titles, title)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "PostgresSQL repository GetTitles: failed to fetch rows")
	}
	return titles, nil
}

func (r *Repository) AddTrackedTitles(ctx context.Context, titles []string) error {
	if len(titles) == 0 {
		return errors.Wrap(entity.ErrInvalidParams, "PostgresSQL repository AddTrackedTitles: titles is empty")
	}
	query := r.sq.
		Insert("coins_tracked").
		Columns("title")

	for _, title := range titles {
		query = query.Values(title)
	}

	sql, args, err := query.
		Suffix("ON CONFLICT (title) DO NOTHING").
		ToSql()

	if err != nil {
		return errors.Wrap(err, "PostgresSQL repository AddTrackedTitles: build insert titles query")
	}
	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "PostgresSQL repository AddTrackedTitles: add tracked titles")
	}
	return nil
}
