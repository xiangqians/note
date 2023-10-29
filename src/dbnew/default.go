// @author xiangqian
// @date 17:43 2023/10/29
package dbnew

type DefaultDb struct {
}

func (db *DefaultDb) Begin() error {
	return nil
}

func (db *DefaultDb) Add(sql string, args ...any) (rowsAffected int64, insertId int64, err error) {
	return 0, 0, nil
}

func (db *DefaultDb) Del(sql string, args ...any) (rowsAffected int64, err error) {
	return 0, nil
}

func (db *DefaultDb) Upd(sql string, args ...any) (rowsAffected int64, err error) {
	return 0, nil
}

func (db *DefaultDb) Get(sql string, args ...any) (result Result, err error) {
	return nil, nil
}

func (db *DefaultDb) Page(sql string, current int64, size uint8, args ...any) (result Result, err error) {
	return nil, nil
}

func (db *DefaultDb) Commit() error {
	return nil
}

func (db *DefaultDb) Rollback() error {
	return nil
}

func (db *DefaultDb) Close() error {
	return nil
}
