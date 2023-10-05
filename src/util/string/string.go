// string
// @author xiangqian
// @date 15:18 2023/10/05
package string

import "fmt"

func String(i any) string {
	if i == nil {
		return ""
	}

	if err, r := i.(error); r {
		return err.Error()
	}

	return fmt.Sprintf("%v", i)
}
