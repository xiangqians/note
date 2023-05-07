// db impl
// @author xiangqian
// @date 12:00 2023/05/07
package db

import "database/sql"

// DbImpl db implement
type DbImpl struct {
	driver string  // driver
	dsn    string  // Data Source Name
	db     *sql.DB // db
	tx     *sql.Tx // tx
	err    error   // error
}

func (db *DbImpl) Open() error {
	db.db, db.err = sql.Open(db.driver, db.dsn)
	return db.err
}

func (db *DbImpl) Begin() error {
	if db.err != nil {
		return db.err
	}

	db.tx, db.err = db.db.Begin()
	return db.err
}

func (db *DbImpl) Add(sql string, args ...any) (int64, error) {
	res, err := db.tx.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DbImpl) exec(sql string, args ...any) (int64, error) {
	if db.err != nil {
		return 0, db.err
	}

	res, err := db.tx.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (db *DbImpl) Del(sql string, args ...any) (int64, error) {
	return db.exec(sql, args...)
}

func (db *DbImpl) Upd(sql string, args ...any) (int64, error) {
	return db.exec(sql, args...)
}

func (db *DbImpl) Qry(sql string, args ...any) (*sql.Rows, error) {
	return db.tx.Query(sql, args...)
}

func (db *DbImpl) Commit() error {
	if db.err == nil && db.tx != nil {
		db.err = db.tx.Commit()
	}
	return db.err
}

func (db *DbImpl) Rollback() error {
	if db.err == nil && db.tx != nil {
		db.err = db.tx.Rollback()
	}
	return db.err
}

func (db *DbImpl) Close() error {
	if db.db != nil {
		db.err = db.db.Close()
	}
	return db.err
}
