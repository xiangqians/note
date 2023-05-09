// db impl
// @author xiangqian
// @date 12:00 2023/05/07
package db

import database_sql "database/sql"

type result int8

const (
	lastInsertId result = iota
	rowsAffected result = iota
)

// DbImpl Db implement
type DbImpl struct {
	db *database_sql.DB // db
	tx *database_sql.Tx // tx
}

func (db *DbImpl) Begin() (err error) {
	db.tx, err = db.db.Begin()
	return
}

func (db *DbImpl) Add(sql string, args ...any) (int64, error) {
	return db.exec(lastInsertId, sql, args...)
}

func (db *DbImpl) Del(sql string, args ...any) (int64, error) {
	return db.exec(rowsAffected, sql, args...)
}

func (db *DbImpl) Upd(sql string, args ...any) (int64, error) {
	return db.exec(rowsAffected, sql, args...)
}

func (db *DbImpl) exec(result result, sql string, args ...any) (int64, error) {
	res, err := db.tx.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	switch result {
	case lastInsertId:
		return res.LastInsertId()

	case rowsAffected:
		return res.RowsAffected()

	default:
		return 0, nil
	}
}

func (db *DbImpl) Qry(sql string, args ...any) (*database_sql.Rows, error) {
	return db.tx.Query(sql, args...)
}

func (db *DbImpl) Commit() (err error) {
	if db.tx != nil {
		err = db.tx.Commit()
	}
	return
}

func (db *DbImpl) Rollback() (err error) {
	if db.tx != nil {
		err = db.tx.Rollback()
	}
	return
}

func (db *DbImpl) Close() (err error) {
	if db.db != nil {
		err = db.db.Close()
	}
	return
}
