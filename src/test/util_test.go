// test
// @author xiangqian
// @date 21:21 2023/02/16
package test

import (
	"fmt"
	util_validate "note/src/util/validate"
	"testing"
)

func TestValidateName(t *testing.T) {

	// \ / : * ? " < > |

	names := []string{
		"test",
		"test\\",
		"test/",
		"test:",
		"test*",
		"test?",
		"test\"",
		"test<",
		"test>",
		"test|",
		"hello",
		"world",
	}

	for i, name := range names {
		err := util_validate.FileName(name)
		fmt.Println(i, name, err)
	}

}
