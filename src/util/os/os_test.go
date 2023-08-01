// os test
// @author xiangqian
// @date 21:55 2023/08/01
package os

import (
	"fmt"
	"testing"
)

func TestOSType(t *testing.T) {
	fmt.Println("IsWindows:", IsWindows())
	fmt.Println("IsLinux:", IsLinux())
}

func TestPath(t *testing.T) {
	path := Path("tmp", "test.txt")
	fmt.Println(path)
}

func TestStat(t *testing.T) {
	paths := []string{
		"C:\\Users\\xiangqian\\Desktop\\tmp",
		"C:\\Users\\xiangqian\\Desktop\\tmp\\apache-maven-3.0.5-bin.tar.gz",
	}

	for _, path := range paths {
		file := Stat(path)
		fmt.Println("path:", path)
		fmt.Println("IsExist:", file.IsExist())
		fmt.Println("IsDir:", file.IsDir())
		fmt.Println("Size:", file.Size())
		fmt.Println()
	}
}
