// test
// @author xiangqian
// @date 21:21 2023/02/16
package test

import (
	"fmt"
	"note/src/util"
	"testing"
)

func TestVerifyDirName(t *testing.T) {

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

	for _, name := range names {
		err := util.VerifyFileName(name)
		fmt.Println(name, err)
	}

}
