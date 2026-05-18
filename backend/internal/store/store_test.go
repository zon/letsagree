package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockUserStore_Upsert(t *testing.T) {
	store := StubUsers()
	user, err := store.Upsert("sub1")
	assert.NoError(t, err)
	assert.Equal(t, "sub1", user.Sub)
	assert.Equal(t, uint(1), user.ID)
}

func TestMockUserStore_Upsert_existing(t *testing.T) {
	store := StubUsers()
	store.users["sub1"] = User{ID: 99, Sub: "sub1"}
	user, err := store.Upsert("sub1")
	assert.NoError(t, err)
	assert.Equal(t, uint(99), user.ID)
}

func TestMockUserStore_UpsertedSub(t *testing.T) {
	store := StubUsers()
	store.Upsert("sub1")
	assert.Equal(t, "sub1", store.UpsertedSub(t))
}

func TestMockSessionStore_Create(t *testing.T) {
	store := StubSessions()
	token, err := store.Create(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestMockSessionStore_Get(t *testing.T) {
	store := WithSession(AnySession())
	session, err := store.Get("test-token")
	assert.NoError(t, err)
	assert.Equal(t, "test-token", session.Token)
}

func TestMockSessionStore_Get_notFound(t *testing.T) {
	store := NoSessions()
	_, err := store.Get("nonexistent")
	assert.Error(t, err)
}

func TestMockSessionStore_Delete(t *testing.T) {
	store := WithSession(AnySession())
	err := store.Delete("test-token")
	assert.NoError(t, err)
	assert.False(t, store.Has("test-token"))
}

func TestMockSessionStore_Has(t *testing.T) {
	store := WithSession(AnySession())
	assert.True(t, store.Has("test-token"))
	assert.False(t, store.Has("nonexistent"))
}

func TestGenerateToken(t *testing.T) {
	token1, err := generateToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token1)

	token2, err := generateToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
}