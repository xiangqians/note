// golang-lru test
// https://pkg.go.dev/github.com/hashicorp/golang-lru/v2
// https://github.com/hashicorp/golang-lru
// 多个协程同时读写同一个 map 的情况下，会直接 panic，通过读写锁（sync.RWMutex）来解决。
//
// @author xiangqian
// @date 20:45 2023/07/17
package server

import (
	"fmt"
	"github.com/hashicorp/golang-lru/v2"
	"testing"
)

func Test1(t *testing.T) {
	l, _ := lru.New[string, string](4)
	l.Add("k1", "v1")
	l.Add("k2", "v2")
	l.Add("k3", "v3")
	l.Add("k4", "v4")
	fmt.Println(l.Values())
	l.Add("k5", "v5")
	fmt.Println(l.Values())
	l.Get("v3")
	fmt.Println(l.Values())
	l.Add("k2", "v22")
	fmt.Println(l.Values())
}

func Test2(t *testing.T) {
	l, _ := lru.New[string, string](4)
	l.Add("k1", "v1")
	l.Add("k2", "v2")
	l.Add("k3", "v3")
	l.Add("k4", "v4")
	fmt.Println(l.Values())
	l.Get("k1")
	l.Add("k5", "v5")
	l.Add("k6", "v6")
	fmt.Println(l.Values())
}
