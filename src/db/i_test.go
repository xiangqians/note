// db test
// @author xiangqian
// @date 18:33 2023/03/30
package db

import (
	"fmt"
	"testing"
)

func TestQry(t *testing.T) {
	// db
	db, err := Get("C:\\Users\\xiangqian\\Desktop\\tmp\\note\\data\\1\\database.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// begin
	err = db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	// qry
	rows, err := db.Qry("SELECT DISTINCT(`type`) FROM `note` WHERE `del` = 0")
	if err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	// rows mapper
	types, count, err := RowsMapper[[]string](rows)
	db.Commit()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(count, types)
}
