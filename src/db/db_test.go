// db test
// @author xiangqian
// @date 20:47 2023/06/12
package db

import (
	"log"
	"testing"
)

func TestDb(t *testing.T) {
	db := DB()
	var result int64
	db.Raw("select 1+10").Take(&result)
	log.Println(result)
}
