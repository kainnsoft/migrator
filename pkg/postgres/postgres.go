package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"

	"github.com/kainnsoft/migrator/config"
)

type (
	IPgDB interface {
		DB() *pgxpool.Pool
		Close() error
		MigrationUp() error
	}
	postgres struct {
		pool *pgxpool.Pool
		cfg  *config.PG
	}
)

func New(
	pgCfg *config.PG,
) (IPgDB, error) {
	var p = &postgres{
		cfg: pgCfg,
	}

	if err := p.get(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *postgres) DB() *pgxpool.Pool {
	return p.pool
}

func (p *postgres) get() (err error) {
	// если нет текущего коннекта, то создаем
	if p.pool == nil {
		if p.pool, err = p.connect(dsn(
			fmt.Sprintf("%s:%s", p.cfg.Host, p.cfg.Port),
			p.cfg.Username,
			p.cfg.Password,
			p.cfg.DBName,
		)); err != nil {
			return err
		}
	}
	if err = p.pool.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}

// dsn generate conn string to database
func dsn(endpoint, username, password, database string) string {
	var dsn = &url.URL{
		Scheme: "postgresql",
		Host:   endpoint,
		Path:   database,
	}

	var q = dsn.Query()
	q.Add("sslmode", "disable")
	// q.Add("binary_parameters", "yes") - not required for pgx
	dsn.RawQuery = q.Encode()

	if username == "" {
		return dsn.String()
	}

	if password == "" {
		dsn.User = url.User(username)
		return dsn.String()
	}

	dsn.User = url.UserPassword(username, password)
	return dsn.String()
}

func (p *postgres) connect(connString string) (db *pgxpool.Pool, err error) {
	if db, err = pgxpool.New(context.Background(), connString); err != nil {
		return nil, errors.Wrap(err, "error create conn")
	}

	if p.cfg.MaxOpenConn > 0 {
		db.Config().MaxConns = p.cfg.MaxOpenConn
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "error ping conn")
	}

	return db, nil
}

func (p *postgres) Close() error {
	if p.pool != nil {
		p.pool.Close()
	}
	return nil
}

func (p *postgres) MigrationUp() error {
	var (
		conn   *sql.DB
		err    error
		strurl = dsn(
			fmt.Sprintf("%s:%s", p.cfg.Host, p.cfg.Port),
			p.cfg.Username,
			p.cfg.Password,
			p.cfg.DBName,
		)
	)
	conn, err = sql.Open("postgres", strurl)
	if err != nil {
		return fmt.Errorf("can't sql.Open migrarion: %v", err)
	}
	err = goose.Up(conn, "migrations")
	if err != nil {
		return fmt.Errorf("can't create migrarion: %v", err)
	}

	return nil
}
