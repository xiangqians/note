// os
package typ

// OS 操作系统标识
type OS int8

const (
	WindowsOS OS = iota // Windows
	LinuxOS             // Linux
	UnknownOS           // Unknown
)
