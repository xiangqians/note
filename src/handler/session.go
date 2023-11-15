// @author xiangqian
// @date 22:41 2023/11/07
package handler

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

// 会话存储器
var store *sessions.CookieStore

func init1() {
	// keyPairs 用于加密和解密会话 cookie 的密钥对。
	// keyPairs 是一个字节切片（[]byte），可以包含一个或多个密钥。
	// 会话 cookie 是存储在客户端浏览器中的数据，用于标识用户会话并进行身份验证。为了确保安全性，会话 cookie 应该进行加密，以防止被篡改或伪造。
	// keyPairs 参数可以是一个随机生成的字节数组，用于加密和解密会话 cookie。你可以使用不同的密钥对来提高安全性。
	// 根据 github.com/gorilla/securecookie 包的建议，推荐的加密密钥长度为 16、24 或 32 字节，以便与 AES-128、AES-192 和 AES-256 算法相匹配。这意味着你可以选择一个长度为 16、24 或 32 的字节数组作为密钥。
	passwd := "$2a$10$NkWzRTyz1ZNnNfjLmxreaeZ31DCiwCEWJlXJAVDkG8fD9Ble2mg4K"
	keyPairs := []byte(passwd)
	switch len(keyPairs) {
	case 16:
		log.Println("Session Cookie密钥对使用AES-128算法")

	case 24:
		log.Println("Session Cookie密钥对使用AES-192算法")

	case 32:
		log.Println("Session Cookie密钥对使用AES-256算法")

	default:
		panic("session-secret-key（会话密钥）长度不符合")
	}

	// 实例化会话存储器
	store = sessions.NewCookieStore(keyPairs)

	store.Options = &sessions.Options{
		Path:     "/",          // 会话可用的路径
		MaxAge:   60 * 60 * 12, // 设置会话过期时间为12小时，单位：秒
		HttpOnly: true,         // 限制 Cookie 只能通过 HTTP 访问
	}
}

func homeHandler1(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name") // 获取会话

	// 检查会话中的值
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // 未认证，重定向到登录页
		return
	}

	w.Write([]byte("Welcome to the home page!")) // 已认证，显示首页内容
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name") // 获取会话

	// 模拟验证用户的过程
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "admin" && password == "password" {
		session.Values["authenticated"] = true                 // 设置会话值为已认证
		session.Save(r, w)                                     // 保存会话
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther) // 重定向到仪表盘页
		return
	}

	w.Write([]byte("Invalid credentials")) // 验证失败，显示错误信息
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name") // 获取会话

	// 检查会话中的值
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // 未认证，重定向到登录页
		return
	}

	w.Write([]byte("Welcome to the dashboard!")) // 已认证，显示仪表盘内容
}
