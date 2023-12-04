// @author xiangqian
// @date 21:55 2023/08/01
package os

import (
	"fmt"
	"testing"
)

func TestOSType(t *testing.T) {
	fmt.Println("IsWindows:\t", IsWindows)
	fmt.Println("IsLinux:\t", IsLinux)
}

func TestPath(t *testing.T) {
	path := Path("tmp", "test.txt")
	fmt.Println(path)
}
