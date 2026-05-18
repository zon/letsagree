package store

import (
	"sync"

	"gorm.io/gorm"
)

type MockUserStore struct {
	mu       sync.Mutex
	users    map[string]User
	lastUser *User
	upserted string
}

func StubUsers() *MockUserStore {
	return &MockUserStore{users: make(map[string]User)}
}

func (m *MockUserStore) Upsert(sub string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.upserted = sub
	if user, ok := m.users[sub]; ok {
		return &user, nil
	}
	user := User{ID: uint(len(m.users) + 1), Sub: sub}
	m.users[sub] = user
	m.lastUser = &user
	return &user, nil
}

func (m *MockUserStore) UpsertedSub(_ interface{}) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.upserted
}

func (m *MockUserStore) Users() map[string]User {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.users
}

func AnySession() *Session {
	return &Session{
		ID:     1,
		Token:  "test-token",
		UserID: 1,
	}
}

func WithSession(session *Session) *MockSessionStore {
	return &MockSessionStore{sessions: map[string]*Session{session.Token: session}}
}

type MockSessionStore struct {
	mu      sync.Mutex
	sessions map[string]*Session
}

func NoSessions() *MockSessionStore {
	return &MockSessionStore{sessions: map[string]*Session{}}
}

func StubSessions() *MockSessionStore {
	return &MockSessionStore{sessions: make(map[string]*Session)}
}

func (m *MockSessionStore) Create(userID uint) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	token := "stub-token"
	m.sessions[token] = &Session{ID: uint(len(m.sessions)+1), Token: token, UserID: userID}
	return token, nil
}

func (m *MockSessionStore) Get(token string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, ok := m.sessions[token]; ok {
		return session, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockSessionStore) Delete(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, token)
	return nil
}

func (m *MockSessionStore) Has(token string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.sessions[token]
	return ok
}

func (m *MockSessionStore) Seed(token string, userID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[token] = &Session{ID: 1, Token: token, UserID: userID}
}