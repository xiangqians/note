// os type
// @author xiangqian
// @date 21:23 2023/03/13
package typ

// OS 操作系统标识
type OS int8

const (
	OSWindows OS = iota // Windows
	OSLinux             // Linux
	OSUnk               // Unknown
)
