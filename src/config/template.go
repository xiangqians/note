// Template
// @author xiangqian
// @date 21:45 2022/12/23
package config

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	html_template "html/template"
	util_os "note/src/util/os"
	util_time "note/src/util/time"
	"os"
	"path/filepath"
	"strings"
)

// 初始化模板（HTML模板）
func initTemplate(engine *gin.Engine) {
	// 自定义模板函数
	customTemplateFunc(engine)

	// 加载html模板
	loadHtmlTemplate(engine, "./res/template")
}

// customTemplateFunc 自定义模板函数
func customTemplateFunc(engine *gin.Engine) {
	// gin内置模板函数
	// go1.19.3/src/text/template/funcs.go:40

	// 自定义模板函数
	engine.SetFuncMap(html_template.FuncMap{
		// 获取i18n文件中key对应的value
		"Localize": i18n.GetMessage,

		// 当前系统时间戳
		"NowUnix": func() int64 {
			return util_time.NowUnix()
		},
	})
}

// loadHtmlTemplate 加载html模板
// templateDir: 模板路径
func loadHtmlTemplate(engine *gin.Engine, templateDir string) {
	// HTML模板
	//engine.LoadHTMLGlob("template/*")
	//engine.LoadHTMLGlob("template/**/*")
	// https://github.com/gin-contrib/multitemplate
	engine.HTMLRender = func(templateDir string) render.HTMLRender {
		// if gin.DebugMode -> NewDynamic()
		renderer := multitemplate.NewRenderer()

		// 获取所有匹配的html模板
		matches, err := filepath.Glob(templateDir + "/*")
		if err != nil {
			panic(err)
		}

		// 获取公共html模板
		commons, err := filepath.Glob(templateDir + "/common/*")
		if err != nil {
			panic(err)
		}

		// Generate our templates map from our layouts/ and includes/ directories
		for _, m := range matches {
			addFromFilesFuncs(renderer, engine.FuncMap, commons, m)
		}

		return renderer
	}(templateDir)
}

func addFromFilesFuncs(renderer multitemplate.Renderer, funcMap html_template.FuncMap, commons []string, name string) {
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
					addFromFilesFuncs(renderer, funcMap, commons, fmt.Sprintf("%s%s%s", name, util_os.FileSeparator(), sfName))
				} else
				// 文件
				{
					var files []string
					if fName == "common" {
						files = []string{fmt.Sprintf("%s%s%s", name, util_os.FileSeparator(), sfName)}
					} else {
						// len 0, cap ?
						files = make([]string, 0, len(commons)+1)
						files = append(files, fmt.Sprintf("%s%s%s", name, util_os.FileSeparator(), sfName))
						files = append(files, commons...)
					}

					renderer.AddFromFilesFuncs(formatTemplateName(fmt.Sprintf("%s/%s", name, sfName)), funcMap, files...)
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
		renderer.AddFromFilesFuncs(formatTemplateName(strings.ReplaceAll(name, "\\", "/")), funcMap, files...)
	}
}

func formatTemplateName(templateName string) string {
	templateName = strings.ReplaceAll(templateName, "\\", "/")

	// name: res/template/user/signin.html -> user/signin.html
	index := strings.Index(templateName, "template")
	templateName = templateName[index+len("template")+1:]

	// 去除后缀名
	return templateName[:len(templateName)-len(".html")]
}
