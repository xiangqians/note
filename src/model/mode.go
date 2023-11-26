// @author xiangqian
// @date 10:31 2023/11/26
package model

import (
	"os"
)

// Mode 模式
type Mode string

const (
	ModeDev  Mode = "dev"  // 开发环境
	ModeTest Mode = "test" // 测试环境
	ModeProd Mode = "prod" // 生产环境
)

// 模式
var mode Mode

// 项目目录
var projectDir string

func init() {
	mode = Mode(os.Getenv("GO_ENV_MODE"))
	projectDir, _ = os.Getwd()
}

func GetMode() Mode {
	return mode
}

func GetProjectDir() string {
	return projectDir
}
