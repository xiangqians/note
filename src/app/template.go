// Template
// @author xiangqian
// @date 21:45 2022/12/23
package app

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	util_num "note/src/util/num"
	util_os "note/src/util/os"
	util_reflect "note/src/util/reflect"
	util_time "note/src/util/time"
	"os"
	"path/filepath"
	"strings"
)

// 初始化HTML模板
func htmlTemplate(engine *gin.Engine) {

	// gin内置模板函数
	// go1.19.3/src/text/template/funcs.go:40

	// 自定义模板函数
	engine.SetFuncMap(template.FuncMap{
		// 为了获取 i18n 文件中 key 对应的 value
		"Localize": i18n.GetMessage,

		// 格式化日期时间戳（s）
		"FormatUnix": func(unix int64) string {
			return util_time.FormatUnix(unix)
		},

		// 人性化日期时间戳（s）
		"HumanizUnix": func(unix int64) string {
			return util_time.HumanizUnix(unix)
		},

		// 人性化文件大小
		"HumanizFileSize": func(size int64) string {
			return util_os.HumanizFileSize(size)
		},

		// No.
		"No_": func(page any, index int) int64 {
			current := util_reflect.CallField[int64](page, "Current")
			size := util_reflect.CallField[int64](page, "Size")
			return (current-1)*size + int64(index) + 1
		},

		// add 两数相加
		"Add": func(i1 any, i2 any) int64 {
			return util_num.Int64(i1) + util_num.Int64(i2)
		},

		// put
		"Put": func(h gin.H, key string, value any) string {
			h[key] = value
			return ""
		},

		// Timestamp
		"Timestamp": func() int64 {
			return util_time.NowUnix()
		},
	})

	// HTML模板
	//engine.LoadHTMLGlob("templates/*")
	//engine.LoadHTMLGlob("templates/**/*")
	// https://github.com/gin-contrib/multitemplate
	engine.HTMLRender = func(templatesDir string) render.HTMLRender {
		// if gin.DebugMode -> NewDynamic()
		renderer := multitemplate.NewRenderer()

		// 获取所有匹配的html模板
		matches, err := filepath.Glob(templatesDir + "/*")
		if err != nil {
			panic(err)
		}

		// 获取公共html模板
		commons, err := filepath.Glob(templatesDir + "/common/*")
		if err != nil {
			panic(err)
		}

		// Generate our templates map from our layouts/ and includes/ directories
		for _, m := range matches {
			addFromFilesFuncs(renderer, engine.FuncMap, commons, m)
		}

		return renderer
	}("./templates")
}

func addFromFilesFuncs(renderer multitemplate.Renderer, funcMap template.FuncMap, commons []string, name string) {
	// 打开文件
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 获取文件信息
	fInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// /**/*
	if fInfo.IsDir() {
		fName := fInfo.Name()
		sfInfos, err := file.Readdir(-1) // sub file info
		if err == nil {
			for _, sfInfo := range sfInfos {
				sfName := sfInfo.Name()

				// 目录
				if sfInfo.IsDir() {
					addFromFilesFuncs(renderer, funcMap, commons, fmt.Sprintf("%s%s%s", name, util_os.FileSeparator, sfName))
				} else
				// 文件
				{
					var files []string
					if fName == "common" {
						files = []string{fmt.Sprintf("%s%s%s", name, util_os.FileSeparator, sfName)}
					} else {
						// len 0, cap ?
						files = make([]string, 0, len(commons)+1)
						files = append(files, fmt.Sprintf("%s%s%s", name, util_os.FileSeparator, sfName))
						files = append(files, commons...)
					}

					renderer.AddFromFilesFuncs(strings.ReplaceAll(fmt.Sprintf("%s/%s", name, sfName), "\\", "/")[len("templates/"):], funcMap, files...)
				}
			}
		}
	} else
	// /*
	{
		// len 0, cap ?
		files := make([]string, 0, len(commons)+1)
		files = append(files, name)
		files = append(files, commons...)
		//renderer.AddFromFilesFuncs(filepath.Base(name), funcMap, files...)
		renderer.AddFromFilesFuncs(strings.ReplaceAll(name, "\\", "/")[len("templates/"):], funcMap, files...)
	}
}
