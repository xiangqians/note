// Package session
//
// 什么是 session？
// session是用来在服务端存储相关数据的，以便在同一个用户的多次请求之间保存用户的状态，比如登录的状态。
// 因为 HTTP 协议是无状态的，要想让客户端（一般浏览器代指一个客户端或用户）的前、后请求关联在一起，就需要给客户端一个唯一的标识来告诉服务端请求是来自于同一个用户，这个标识就是所谓的 sessionid。
// 该 sessionid 由服务端生成，并存储客户端（cookie、url）中。 当客户端再次发起请求的时候，就会携带该标识，服务端根据该标识就能查找到存在服务端上的相关数据。
//
// @author xiangqian
// @date 22:41 2023/11/07
package session

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"note/src/model"
	"strings"
	"sync"
)

// 会话存储器
var store sessions.Store

func init() {
	// 设置会话过期时间，单位：秒
	maxAge := int(model.Ini.Server.SessionMaxAge.Seconds())

	// github.com/gorilla/securecookie 提供了一种安全的cookie，通过在服务端给cookie加密，让其内容不可读，也不可伪造。当然，敏感信息还是强烈建议不要放在cookie中。
	// 会话cookie是存储在客户端浏览器中的数据，用于标识用户会话并进行身份验证。为了确保安全性，会话cookie应该进行加密，以防止被篡改或伪造。
	// 根据 github.com/gorilla/securecookie 包的建议，推荐的加密密钥长度为 16、24 或 32 字节，以便与 AES-128、AES-192 和 AES-256 算法相匹配。这意味着你可以选择一个长度为 16、24 或 32 的字节数组作为密钥。
	hashKey := securecookie.GenerateRandomKey(32)
	blockKey := securecookie.GenerateRandomKey(32)
	codecs := securecookie.CodecsFromPairs(hashKey, blockKey)

	// MaxAge sets the maximum age for the store and the underlying cookie implementation.
	// Individual sessions can be deleted by setting Options.MaxAge = -1 for that session.
	// Set the maxAge for each securecookie instance.
	for _, codec := range codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(maxAge)
		}
	}

	// 实例化会话存储器
	memStore := MemStore{
		// gorilla/securecookie 提供了一种安全的cookie，通过在服务端给cookie加密，让其内容不可读，也不可伪造。当然，敏感信息还是强烈建议不要放在cookie中。
		Codecs: codecs,
		// 配置会话存储器
		Options: &sessions.Options{
			Path:     "/",    // 会话可用的路径
			MaxAge:   maxAge, // 设置会话过期时间，单位：秒
			HttpOnly: true,   // 限制 Cookie 只能通过 HTTP 访问
		},
		// 会话数据
		data: make(map[string]map[any]any),
	}

	store = &memStore
}

const (
	systemName   = "system"
	languageName = "language"
	msgName      = "msg"
)

type Session struct {
	request *http.Request
	writer  http.ResponseWriter
	session *sessions.Session
}

func (session *Session) Get(name string) any {
	return session.session.Values[name]
}

func (session *Session) Set(name string, value any) error {
	s := session.session
	s.Values[name] = value
	return s.Save(session.request, session.writer)
}

func (session *Session) Del(name string) error {
	s := session.session
	delete(s.Values, name)
	return s.Save(session.request, session.writer)
}

func (session *Session) Clear() error {
	s := session.session
	s.Options.MaxAge = -1
	delete(s.Values, systemName)
	return s.Save(session.request, session.writer)
}

func (session *Session) GetSystem() model.System {
	if system, ok := session.Get(systemName).(model.System); ok {
		return system
	}
	return model.System{}
}

func (session *Session) SetSystem(system model.System) error {
	err := session.Set(systemName, system)

	memStore := store.(*MemStore)
	memStore.mutex.Lock()
	defer memStore.mutex.Unlock()

	// 用户多终端登录限制
	data := memStore.data
	for id, m := range data {
		if id == session.session.ID {
			continue
		}
		if m != nil {
			for name, value := range m {
				if name == systemName {
					if _, ok := value.(model.System); ok {
						data[id] = nil
						delete(data, id)
					}
					break
				}
			}
		}
	}
	return err
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

func (session *Session) GetMsg() string {
	if msg, ok := session.Get(msgName).(string); ok {
		session.Del(msgName)
		return msg
	}
	return ""
}

func (session *Session) SetMsg(msg string) error {
	return session.Set(msgName, msg)
}

func GetSession(request *http.Request, writer http.ResponseWriter) *Session {
	session, err := store.Get(request, "note") // cookie名称 -- sessionid
	if err != nil {
		// 密钥对变更导致异常
		if strings.Contains(err.Error(), "the value is not valid") {
			session.Save(request, writer)
		} else {
			panic(err)
		}
	}

	return &Session{
		request: request,
		writer:  writer,
		session: session,
	}
}

// copy: mod\github.com\quasoft\memstore@v0.0.0-20191010062613-2bce066d2b0b\memstore.go

// MemStore is an in-memory implementation of gorilla/sessions, suitable for use in tests and development environments. Do not use in production.
// Values are cached in a map. The cache is protected and can be used by multiple goroutines.
type MemStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options
	data    map[string]map[any]any
	mutex   sync.RWMutex
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func (store *MemStore) Get(request *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(request).Get(store, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func (store *MemStore) New(request *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(store, name)
	options := *store.Options
	session.Options = &options
	session.IsNew = true

	cookie, err := request.Cookie(name)
	if err != nil {
		// Cookie not found, this is a new session
		return session, nil
	}

	err = securecookie.DecodeMulti(name, cookie.Value, &session.ID, store.Codecs...)
	if err != nil {
		// Value could not be decrypted, consider this is a new session
		return session, err
	}

	value, ok := store.get(session.ID)
	if !ok {
		// No value found in cache, don't set any values in session object,
		// consider a new session
		return session, nil
	}

	// Values found in session, this is not a new session
	session.Values = store.copy(value)
	session.IsNew = false
	return session, nil
}

// Save adds a single session to the response.
// Set Options.MaxAge to -1 or call MaxAge(-1) before saving the session to delete all values in it.
func (store *MemStore) Save(request *http.Request, writer http.ResponseWriter, session *sessions.Session) error {
	var cookieValue string
	if session.Options.MaxAge < 0 {
		cookieValue = ""
		store.del(session.ID)
		for k := range session.Values {
			delete(session.Values, k)
		}
	} else {
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		encrypted, err := securecookie.EncodeMulti(session.Name(), session.ID, store.Codecs...)
		if err != nil {
			return err
		}
		cookieValue = encrypted
		store.set(session.ID, store.copy(session.Values))
	}
	http.SetCookie(writer, sessions.NewCookie(session.Name(), cookieValue, session.Options))
	return nil
}

func (store *MemStore) copy(v map[any]any) map[any]any {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		panic(fmt.Errorf("could not copy memstore value. Encoding to gob failed: %v", err))
	}

	var value map[any]any
	err = dec.Decode(&value)
	if err != nil {
		panic(fmt.Errorf("could not copy memstore value. Decoding from gob failed: %v", err))
	}
	return value
}

func (store *MemStore) get(name string) (value map[any]any, ok bool) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	value, ok = store.data[name]
	return
}

func (store *MemStore) set(name string, value map[any]any) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.data[name] = value
}

func (store *MemStore) del(name string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if _, ok := store.data[name]; ok {
		delete(store.data, name)
	}
}
