// @author xiangqian
// @date 15:18 2023/10/05
package string

import "fmt"

func String(i any) string {
	if i == nil {
		return ""
	}

	if s, ok := i.(string); ok {
		return s
	}

	if err, ok := i.(error); ok {
		return err.Error()
	}

	return fmt.Sprintf("%v", i)
}
