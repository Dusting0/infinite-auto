package web

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"infinite-calc/defense"
	"infinite-calc/hp"
)

const (
	sessionCookie = "calc_sid"
	sessionTTL    = 24 * time.Hour
)

// session 单个浏览器会话的可变状态
type session struct {
	mu       sync.Mutex
	HP       *hp.HPState
	Defense  *defense.Preset
	lastSeen time.Time
}

func newSession() *session {
	return &session{
		HP:       hp.NewHPState(20),
		Defense:  defense.NewPreset("默认预设"),
		lastSeen: time.Now(),
	}
}

// sessionStore 进程内会话表
type sessionStore struct {
	mu   sync.Mutex
	data map[string]*session
}

func newSessionStore() *sessionStore {
	s := &sessionStore{data: make(map[string]*session)}
	go s.gcLoop()
	return s
}

func (s *sessionStore) gcLoop() {
	t := time.NewTicker(30 * time.Minute)
	defer t.Stop()
	for range t.C {
		s.mu.Lock()
		now := time.Now()
		for id, sess := range s.data {
			sess.mu.Lock()
			idle := now.Sub(sess.lastSeen) > sessionTTL
			sess.mu.Unlock()
			if idle {
				delete(s.data, id)
			}
		}
		s.mu.Unlock()
	}
}

func randomID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// get 根据请求 Cookie 取会话；没有则新建并写入 Cookie
func (s *sessionStore) get(w http.ResponseWriter, r *http.Request) *session {
	var id string
	if c, err := r.Cookie(sessionCookie); err == nil && c.Value != "" {
		id = c.Value
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if id != "" {
		if sess, ok := s.data[id]; ok {
			sess.mu.Lock()
			sess.lastSeen = time.Now()
			sess.mu.Unlock()
			return sess
		}
	}

	id = randomID()
	sess := newSession()
	s.data[id] = sess

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
	return sess
}

// 全局会话仓库
var sessions = newSessionStore()
