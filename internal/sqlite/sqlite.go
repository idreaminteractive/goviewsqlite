package sqlite

import (
	"context"

	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type LocalDB struct {
	dsn        string
	connection *sql.DB
	ctx        context.Context // background context
	cancel     func()          // cancel background context
}

func NewDB(dsn string) *LocalDB {
	db := &LocalDB{
		dsn: dsn,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db

}

func (db *LocalDB) Connection() *sql.DB {
	return db.connection
}

func (db *LocalDB) Open() (err error) {

	conn, err := sql.Open("sqlite", db.dsn)
	if err != nil {
		fmt.Printf("Could not open db")
		return err
	}

	db.connection = conn
	return nil
}

func (d *LocalDB) Ready() bool {
	return d.connection != nil
}

func (db *LocalDB) Close() error {

	// Cancel background context.
	db.cancel()

	// Close database.
	if db.connection != nil {
		return db.connection.Close()
	}
	return nil
}
