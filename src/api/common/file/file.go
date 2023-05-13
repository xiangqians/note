// file
// @author xiangqian
// @date 13:03 2023/05/13
package file

import (
	util_os "note/src/util/os"
	"os"
)

// Clear 清空文件，如果文件存在
func Clear(path string) error {
	if !util_os.IsExist(path) {
		return nil
	}

	file, err := os.OpenFile(path,
		os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
		0666)
	file.Close()
	return err
}
