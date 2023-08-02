// os test
// @author xiangqian
// @date 21:55 2023/08/01
package os

import (
	"fmt"
	"testing"
)

func TestOSType(t *testing.T) {
	fmt.Println("IsWindows:\t", IsWindows())
	fmt.Println("IsLinux:\t", IsLinux())
}

func TestPath(t *testing.T) {
	path := Path("tmp", "test.txt")
	fmt.Println(path)
}

func TestStat(t *testing.T) {
	paths := []string{
		"C:\\Users\\xiangqian\\Desktop\\tmp",
		Path("C:\\Users\\xiangqian\\Desktop\\tmp", "apache-maven-3.0.5-bin.tar.gz"),
		Path("C:\\Users\\xiangqian\\Desktop\\tmp", "apache-maven-3.0.5-bin.tar.gz1"),
	}

	for _, path := range paths {
		file := Stat(path)
		fmt.Println("---------", path, "---------")
		fmt.Println("IsExist:", file.IsExist())
		if file.IsExist() {
			fmt.Println("Name:", file.Name())
			fmt.Println("Size:", file.Size())
			fmt.Println("IsDir:", file.IsDir())
		}
		fmt.Println()
	}
}

func TestMkDir(t *testing.T) {
	path := Path("C:\\Users\\xiangqian\\Desktop\\tmp", "note2023")
	err := MkDir(path)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRm(t *testing.T) {
	path := Path("C:\\Users\\xiangqian\\Desktop\\tmp", "note2023")
	err := Rm(path)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCopyFile(t *testing.T) {
	srcPath := Path("C:\\Users\\xiangqian\\Desktop\\tmp", "apache-skywalking-apm-9.0.0.tar.gz")
	dstPath := Path("C:\\Users\\xiangqian\\Desktop\\tmp", "apache-skywalking-apm-9.0.0.tar.gz1")
	written, err := CopyFile(srcPath, dstPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(written)
}
