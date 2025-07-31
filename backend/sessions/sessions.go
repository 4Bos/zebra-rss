package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

var sessions map[string]*Session = make(map[string]*Session)
var mu sync.Mutex

type Session struct {
	Id     string
	data   map[string]interface{}
	mu     sync.Mutex
	closed bool
}

func (s *Session) SetString(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *Session) GetString(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key].(string)

	return value, ok
}

func (s *Session) SetInt64(key string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *Session) GetInt64(key string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key].(int64)

	return value, ok
}

func (s *Session) Close() {
	mu.Lock()
	defer mu.Unlock()

	s.closed = true

	delete(sessions, s.Id)
}

func GetSession(sessionId string) *Session {
	if session, ok := sessions[sessionId]; ok {
		return session
	}

	return nil
}

func StartSession() (*Session, error) {
	mu.Lock()
	defer mu.Unlock()

	sessionId, err := generateSessionID()

	if err != nil {
		return nil, err
	}

	session := createSession(sessionId)
	sessions[sessionId] = session

	return session, nil
}

func createSession(sessionId string) *Session {
	return &Session{
		Id:     sessionId,
		data:   make(map[string]interface{}),
		closed: false,
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
