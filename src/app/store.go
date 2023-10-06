// store
// @author xiangqian
// @date 22:47 2023/05/28
package app

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"fmt"
	gincontrib_sessions "github.com/gin-contrib/sessions"
	"github.com/gorilla/securecookie"
	gorilla_sessions "github.com/gorilla/sessions"
	"github.com/hashicorp/golang-lru/v2"
	"net/http"
	"note/src/session"
	"note/src/typ"
	"strings"
)

// --------------------------------- copy mod\github.com\quasoft\memstore@v0.0.0-20191010062613-2bce066d2b0b\cache.go ---------------------------------

type cache struct {
	data *lru.Cache[string, map[any]any]
}

func (c *cache) value(name string) (map[any]any, bool) {
	//log.Println(c.String())
	return c.data.Get(name)
}

func (c *cache) String() string {
	var arr []string = nil
	keys := c.data.Keys()
	if keys != nil && len(keys) > 0 {
		// len 0, cap ?
		cap := len(keys)
		arr = make([]string, 0, cap)
		arr = append(arr, fmt.Sprintf("cap %d", cap))
		for _, key := range keys {
			if v, ok := c.data.Get(key); ok {
				arr = append(arr, fmt.Sprintf("%s %v", key, v))
			}
		}
	}

	if arr != nil {
		return fmt.Sprintf("%s", strings.Join(arr, "\n\t"))
	}

	return ""
}

func (c *cache) setValue(name string, value map[any]any) {
	// 用户多终端登录限制
	if v, ok := value[session.UserKey]; ok {
		if user, ok := v.(typ.User); ok {
			id := user.Id
			keys := c.data.Keys()
			if keys != nil && len(keys) > 0 {
				//log.Println("before", c.String())
				for _, key := range keys {
					if v, ok = c.data.Get(key); ok {
						if m, ok := v.(map[any]any); ok {
							if v, ok = m[session.UserKey]; ok {
								if user, ok = v.(typ.User); ok {
									if id == user.Id {
										c.data.Remove(key)
									}
								}
							}
						}
					}
				}
				//log.Println("after", c.String())
			}
		}
	}

	c.data.Add(name, value)
	//log.Println(c.String())
}

func (c *cache) delete(name string) {
	c.data.Remove(name)
}

// --------------------------------- copy mod\github.com\quasoft\memstore@v0.0.0-20191010062613-2bce066d2b0b\memstore.go ---------------------------------

// MemStore is an in-memory implementation of gorilla/sessions, suitable
// for use in tests and development environments. Do not use in production.
// Values are cached in a map. The cache is protected and can be used by
// multiple goroutines.
type MemStore struct {
	Codecs  []securecookie.Codec
	Options *gorilla_sessions.Options
	cache   *cache
}

var data *lru.Cache[string, map[any]any]

// NewMemStore returns a new MemStore.
//
// Keys are defined in pairs to allow key rotation, but the common case is
// to set a single authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for
// encryption. The encryption key can be set to nil or omitted in the last
// pair, but the authentication key is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes.
// The encryption key, if set, must be either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256 modes.
//
// Use the convenience function securecookie.GenerateRandomKey() to create
// strong keys.
func NewMemStore(keyPairs ...[]byte) *MemStore {
	data, _ = lru.New[string, map[any]any](16)
	store := MemStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &gorilla_sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		cache: &cache{data: data},
	}
	store.MaxAge(store.Options.MaxAge)
	return &store
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func (m *MemStore) Get(r *http.Request, name string) (*gorilla_sessions.Session, error) {
	return gorilla_sessions.GetRegistry(r).Get(m, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func (m *MemStore) New(r *http.Request, name string) (*gorilla_sessions.Session, error) {
	session := gorilla_sessions.NewSession(m, name)
	options := *m.Options
	session.Options = &options
	session.IsNew = true

	c, err := r.Cookie(name)
	if err != nil {
		// Cookie not found, this is a new session
		return session, nil
	}

	err = securecookie.DecodeMulti(name, c.Value, &session.ID, m.Codecs...)
	if err != nil {
		// Value could not be decrypted, consider this is a new session
		return session, err
	}

	v, ok := m.cache.value(session.ID)
	if !ok {
		// No value found in cache, don't set any values in session object,
		// consider a new session
		return session, nil
	}

	// Values found in session, this is not a new session
	session.Values = m.copy(v)
	session.IsNew = false
	return session, nil
}

// Save adds a single session to the response.
// Set Options.MaxAge to -1 or call MaxAge(-1) before saving the session to delete all values in it.
func (m *MemStore) Save(r *http.Request, w http.ResponseWriter, s *gorilla_sessions.Session) error {
	var cookieValue string
	if s.Options.MaxAge < 0 {
		cookieValue = ""
		m.cache.delete(s.ID)
		for k := range s.Values {
			delete(s.Values, k)
		}
	} else {
		if s.ID == "" {
			s.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		encrypted, err := securecookie.EncodeMulti(s.Name(), s.ID, m.Codecs...)
		if err != nil {
			return err
		}
		cookieValue = encrypted
		m.cache.setValue(s.ID, m.copy(s.Values))
	}
	http.SetCookie(w, gorilla_sessions.NewCookie(s.Name(), cookieValue, s.Options))
	return nil
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (m *MemStore) MaxAge(age int) {
	m.Options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range m.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

func (m *MemStore) copy(v map[any]any) map[any]any {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		panic(fmt.Errorf("could not copy memstore value. Encoding to gob failed: %v", err))
	}

	var val map[any]any
	err = dec.Decode(&val)
	if err != nil {
		panic(fmt.Errorf("could not copy memstore value. Decoding from gob failed: %v", err))
	}
	return val
}

// copy: mod\github.com\gin-contrib\sessions@v0.0.5\memstore\memstore.go

type Store interface {
	gincontrib_sessions.Store
}

// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewStore(keyPairs ...[]byte) Store {
	return &store{NewMemStore(keyPairs...)}
}

type store struct {
	*MemStore
}

func (c *store) Options(options gincontrib_sessions.Options) {
	c.MemStore.Options = options.ToGorillaOptions()
}
