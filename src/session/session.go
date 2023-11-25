// @author xiangqian
// @date 22:41 2023/11/07
package session

import (
	"github.com/gorilla/sessions"
	"net/http"
	"note/src/model"
)

// 会话存储器
var store *sessions.CookieStore

func init() {
	// keyPairs 用于加密和解密会话 cookie 的密钥对。
	// keyPairs 是一个字节切片（[]byte），可以包含一个或多个密钥。
	// 会话 cookie 是存储在客户端浏览器中的数据，用于标识用户会话并进行身份验证。为了确保安全性，会话 cookie 应该进行加密，以防止被篡改或伪造。
	// keyPairs 参数可以是一个随机生成的字节数组，用于加密和解密会话 cookie。你可以使用不同的密钥对来提高安全性。
	// 根据 github.com/gorilla/securecookie 包的建议，推荐的加密密钥长度为 16、24 或 32 字节，以便与 AES-128、AES-192 和 AES-256 算法相匹配。这意味着你可以选择一个长度为 16、24 或 32 的字节数组作为密钥。
	keyPairs := []byte(model.Ini.Server.SessionSecretKey)
	length := len(keyPairs)
	if length != 16 && length != 24 && length != 32 {
		panic("session-secret-key（会话密钥）长度不符合：推荐加密密钥长度为 16、24 或 32 字节！")
	}

	// 实例化会话存储器
	store = sessions.NewCookieStore(keyPairs)

	// 配置会话存储器
	store.Options = &sessions.Options{
		Path:     "/",          // 会话可用的路径
		MaxAge:   60 * 60 * 12, // 设置会话过期时间为12小时，单位：秒
		HttpOnly: true,         // 限制 Cookie 只能通过 HTTP 访问
	}
}

const (
	userName     = "user"
	languageName = "language"
)

type Session struct {
	session *sessions.Session
	writer  http.ResponseWriter
	request *http.Request
}

func (session *Session) Get(name string) any {
	return session.session.Values[name]
}

func (session *Session) Set(name string, value any) error {
	s := session.session
	s.Values[name] = value
	return s.Save(session.request, session.writer)
}

func (session *Session) GetUser() model.User {
	if user, ok := session.Get(userName).(model.User); ok {
		return user
	}
	return model.User{}
}

func (session *Session) SetUser(user model.User) error {
	return session.Set(userName, user)
}

func (session *Session) GetLanguage() string {
	if language, ok := session.Get(languageName).(string); ok {
		return language
	}
	return ""
}

func (session *Session) SetLanguage(language string) error {
	return session.Set(languageName, language)
}

func GetSession(writer http.ResponseWriter, request *http.Request) *Session {
	session, err := store.Get(request, "note") // cookie名称
	if err != nil {
		panic(err)
	}

	return &Session{
		session: session,
		writer:  writer,
		request: request,
	}
}
