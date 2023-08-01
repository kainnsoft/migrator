package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kainnsoft/migrator/config"
)

type (
	IMySQL interface {
		DB() *sql.DB
		Close() error
	}
	mySqlDB struct {
		db *sql.DB
	}
)

func New(cfg *config.MySql) (IMySQL, error) {
	var (
		db  *sql.DB
		err error
	)
	if db, err = sql.Open("mysql", dsn(cfg)); err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &mySqlDB{db: db}, nil
}

func dsn(cfg *config.MySql) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
}

func (m *mySqlDB) DB() *sql.DB {
	return m.db
}

func (m *mySqlDB) Close() error {
	if m.db != nil {
		m.db.Close()
	}
	return nil
}
