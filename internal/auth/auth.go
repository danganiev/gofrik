package auth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type Middleware struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

type Session struct {
	UserID    int
	Email     string
	ExpiresAt time.Time
}

func NewMiddleware() *Middleware {
	return &Middleware{
		sessions: make(map[string]*Session),
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (am *Middleware) CreateSession(userID int, email string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	am.sessions[token] = &Session{
		UserID:    userID,
		Email:     email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return token, nil
}

func (am *Middleware) GetSession(token string) (*Session, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	session, ok := am.sessions[token]
	if !ok || time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}

func (am *Middleware) DeleteSession(token string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	delete(am.sessions, token)
}

