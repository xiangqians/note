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
)

// 初始化HTML模板
func htmlTemplate(pEngine *gin.Engine) {
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

		// +1
		"No_": func(current int64, size uint8, i any) int64 {
			return util.Add(i, (current-1)*int64(size)) + 1
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

		matches, err := filepath.Glob(templatesDir + "/*")
		if err != nil {
			panic(err)
		}

		coms, err := filepath.Glob(templatesDir + "/com/*")
		if err != nil {
			panic(err)
		}

		getFiles := func(s string) []string {
			files := make([]string, len(coms)+1)
			i := 0
			files[i] = s
			i++
			for _, e := range coms {
				files[i] = e
				i++
			}
			return files
		}

		// Generate our templates map from our layouts/ and includes/ directories
		for _, matche := range matches {
			pFile, ferr := os.Open(matche)
			if ferr != nil {
				continue
			}

			fileInfo, fierr := pFile.Stat()
			if fierr == nil {
				name := filepath.Base(matche)
				// /**/*
				if fileInfo.IsDir() {
					fname := fileInfo.Name()
					subFileInfos, sfierr := pFile.Readdir(-1)
					if sfierr == nil {
						for _, subFileInfo := range subFileInfos {
							subfname := subFileInfo.Name()
							var files []string
							if fname == "com" {
								files = []string{fmt.Sprintf("%s/%s", matche, subfname)}
							} else {
								files = getFiles(fmt.Sprintf("%s/%s", matche, subfname))
							}
							renderer.AddFromFilesFuncs(fmt.Sprintf("%s/%s", fname, subfname), pEngine.FuncMap, files...)
						}
					}
				} else
				// /*
				{
					files := getFiles(matche)
					renderer.AddFromFilesFuncs(name, pEngine.FuncMap, files...)
				}
			}
			pFile.Close()
		}

		return renderer
	}("./templates")
}
