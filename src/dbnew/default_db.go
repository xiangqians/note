// @author xiangqian
// @date 17:43 2023/10/29
package dbnew

type Default struct {
}

func (db *Default) Begin() error {
	return nil
}

func (db *Default) Add(sql string, args ...any) (rowsAffected int64, insertId int64, err error) {
	return 0, 0, nil
}

func (db *Default) Del(sql string, args ...any) (rowsAffected int64, err error) {
	return 0, nil
}

func (db *Default) Upd(sql string, args ...any) (rowsAffected int64, err error) {
	return 0, nil
}

func (db *Default) Get(sql string, args ...any) (result Result, err error) {
	return nil, nil
}

func (db *Default) Page(sql string, current int64, size uint8, args ...any) (result Result, err error) {
	return nil, nil
}

func (db *Default) Commit() error {
	return nil
}

func (db *Default) Rollback() error {
	return nil
}

func (db *Default) Close() error {
	return nil
}
