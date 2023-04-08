// test
// @author xiangqian
// @date 21:21 2023/02/16
package test

import (
	"fmt"
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
		//err := os.VerifyName(name)
		//fmt.Println(name, err)
		fmt.Println(name)
	}

}
