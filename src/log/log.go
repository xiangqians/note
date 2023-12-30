// @author xiangqian
// @date 22:04 2023/11/19
package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"note/src/model"
	util_os "note/src/util/os"
	util_time "note/src/util/time"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var ini = model.Ini

// 初始化日志写入器
func init() {
	// 时区
	timeZone := ini.Sys.TimeZone
	if timeZone != "" {
		// 加载时区
		loc, err := time.LoadLocation(timeZone)
		// 设置时区
		if err == nil {
			time.Local = loc
		}
	}

	// 文件写入器
	fileWriter := &fileWriter{}
	fileWriter.openFile()

	// 多重写入器：文件写入器（写到文件） & 控制台写入器（写到控制台）
	writer := io.MultiWriter(fileWriter, os.Stdout)

	// 设置gin日志默认写入器
	gin.DefaultWriter = writer

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置日志写入器
	log.SetOutput(writer)

	log.Printf("TimeZone %s\n", time.Local)
	log.Println(ini)

	//for i := 0; i < 5; i++ {
	//	go func() {
	//		for {
	//			log.Println("夫天地者，万物之逆旅也；光阴者，百代之过客也。而浮生若梦，为欢几何？古人秉烛夜游，良有以也。况阳春召我以烟景，大块假我以文章。会桃花之芳园，序天伦之乐事。群季俊秀，皆为惠连；吾人咏歌，独惭康乐。幽赏未已，高谈转清。开琼筵以坐花，飞羽觞而醉月。不有佳咏，何伸雅怀？如诗不成，罚依金谷酒数。")
	//			time.Sleep(100 * time.Millisecond)
	//		}
	//	}()
	//}
}

// 日志文件写入器
type fileWriter struct {
	mutex sync.Mutex // 创建一个互斥锁
	file  *os.File   // 日志文件指针
}

// 打开日志文件（如果不存在则创建）
func (writer *fileWriter) openFile() {
	// 日志文件目录
	dir := ini.Log.Dir
	fileInfo, err := os.Stat(dir)

	// 日志文件目录不存在或者不是文件目录，则创建日志文件目录
	if (err != nil && !os.IsExist(err)) || !fileInfo.IsDir() {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err != nil {
		panic(err)
	}

	// 日志文件路径
	path := util_os.Path(dir, ini.Log.FileName)

	// 打开文件
	writer.file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
}

func (writer *fileWriter) Write(p []byte) (int, error) {
	// defer的作用是把defer关键字之后的函数执行压入一个栈中延迟执行，多个defer的执行顺序是后进先出
	defer writer.mutex.Unlock() // 解锁

	// 加锁
	writer.mutex.Lock()

	// 检查文件大小
	fileInfo, err := writer.file.Stat()
	if err != nil {
		panic(err)
	}
	// 文件超过最大大小，备份日志文件并删除最早的日志文件
	if fileInfo.Size() >= ini.Log.MaxSize {
		// 先关闭当前日志文件句柄
		err = writer.file.Close()
		if err != nil {
			panic(err)
		}

		// 备份当前日志文件
		path := util_os.Path(ini.Log.Dir, ini.Log.FileName)
		err = os.Rename(path, fmt.Sprintf("%s.%s", path, util_time.NowTime().Format("20060102150405")))
		if err != nil {
			panic(err)
		}

		// 历史日志文件
		var historyPaths []string

		// 日志文件目录绝对路径
		err = filepath.Walk(ini.Log.Dir, func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 忽略目录
			if fileInfo.IsDir() {
				return nil
			}

			if strings.HasPrefix(fileInfo.Name(), fmt.Sprintf("%s.", ini.Log.FileName)) {
				historyPaths = append(historyPaths, path)
			}

			return nil
		})

		// 历史日志文件根据日期降序排
		sort.Slice(historyPaths, func(i, j int) bool {
			iHistoryPath := historyPaths[i]
			iTime := iHistoryPath[strings.LastIndex(iHistoryPath, ".")+1:]
			jHistoryPath := historyPaths[j]
			jTime := jHistoryPath[strings.LastIndex(jHistoryPath, ".")+1:]
			return iTime > jTime
		})
		length := len(historyPaths)
		maxHistory := ini.Log.MaxHistory
		if length > maxHistory {
			for i := maxHistory; i < length; i++ {
				os.Remove(historyPaths[i])
			}
		}

		// 重新创建日志文件
		writer.openFile()
	}

	return writer.file.Write(p)
}
