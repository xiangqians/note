// db pool
// @author xiangqian
// @date 13:31 2023/05/07
package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Get 获取db
// dsn: data source name
func Get(dsn string) (Db, error) {
	dn := "sqlite3" // driver name
	db, err := sql.Open(dn, dsn)
	if err != nil {
		return nil, err
	}

	return &DbImpl{db: db}, nil
}
