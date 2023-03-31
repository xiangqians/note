// os type
// @author xiangqian
// @date 21:23 2023/03/13
package os

// OS 操作系统标识
type OS int8

const (
	Windows OS = iota // Windows
	Linux             // Linux
	Unk               // Unknown
)
