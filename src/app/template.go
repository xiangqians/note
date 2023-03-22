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
	"note/src/util"
	"os"
	"path/filepath"
	"strings"
)

// 初始化HTML模板
func htmlTemplate(pEngine *gin.Engine) {

	// gin内置模板函数
	// go1.19.3/src/text/template/funcs.go:40

	// 自定义模板函数
	pEngine.SetFuncMap(template.FuncMap{
		// 为了获取 i18n 文件中 key 对应的 value
		"Localize": i18n.GetMessage,

		// 格式化日期时间戳（s）
		"FormatUnix": func(unix int64) string {
			if unix == 0 {
				return "-"
			}

			return util.FormatUnix(unix)
		},

		// 人性化日期时间戳（s）
		"HumanizUnix": func(unix int64) string {
			return util.HumanizUnix(unix)
		},

		// 人性化文件大小
		"HumanizFileSize": func(size int64) string {
			return util.HumanizFileSize(size)
		},

		// No.
		"No_": func(current int64, size uint8, i int) int64 {
			return (current-1)*int64(size) + int64(i) + 1
		},

		// add
		"Add": func(i1 any, i2 any) int64 {
			return util.Int64(i1) + util.Int64(i2)
		},

		"Put": func(h gin.H, key string, value any) string {
			h[key] = value
			return ""
		},
	})

	// HTML模板
	//pEngine.LoadHTMLGlob("templates/*")
	//pEngine.LoadHTMLGlob("templates/**/*")
	// https://github.com/gin-contrib/multitemplate
	pEngine.HTMLRender = func(templatesDir string) render.HTMLRender {
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
			addFromFilesFuncs(renderer, pEngine.FuncMap, commons, m)
		}

		return renderer
	}("./templates")
}

func addFromFilesFuncs(renderer multitemplate.Renderer, funcMap template.FuncMap, commons []string, name string) {
	// 打开文件
	pFile, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer pFile.Close()

	// 获取文件信息
	fInfo, err := pFile.Stat()
	if err != nil {
		panic(err)
	}

	// /**/*
	if fInfo.IsDir() {
		fName := fInfo.Name()
		sfInfos, err := pFile.Readdir(-1) // sub file info
		if err == nil {
			for _, sfInfo := range sfInfos {
				sfName := sfInfo.Name()

				// 目录
				if sfInfo.IsDir() {
					addFromFilesFuncs(renderer, funcMap, commons, fmt.Sprintf("%s%s%s", name, util.FileSeparator, sfName))
				} else
				// 文件
				{
					var files []string
					if fName == "com" {
						files = []string{fmt.Sprintf("%s%s%s", name, util.FileSeparator, sfName)}
					} else {
						// len 0, cap ?
						files = make([]string, 0, len(commons)+1)
						files = append(files, fmt.Sprintf("%s%s%s", name, util.FileSeparator, sfName))
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
