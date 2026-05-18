package store

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID  uint   `gorm:"primaryKey"`
	Sub string `gorm:"uniqueIndex"`
}

func (User) TableName() string {
	return "users"
}

type Session struct {
	ID     uint   `gorm:"primaryKey"`
	Token  string `gorm:"index"`
	UserID uint
}

func (Session) TableName() string {
	return "sessions"
}

type gormStore struct {
	db *gorm.DB
}

func New(db *gorm.DB) *gormStore {
	return &gormStore{db: db}
}

type UserStore interface {
	Upsert(sub string) (*User, error)
}

type SessionStore interface {
	Create(userID uint) (string, error)
	Get(token string) (*Session, error)
	Delete(token string) error
}

func (s *gormStore) Upsert(sub string) (*User, error) {
	var user User
	err := s.db.FirstOrCreate(&user, User{Sub: sub}).Error
	return &user, err
}

func (s *gormStore) Create(userID uint) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}
	session := Session{Token: token, UserID: userID}
	err = s.db.Create(&session).Error
	return token, err
}

func (s *gormStore) Get(token string) (*Session, error) {
	var session Session
	err := s.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *gormStore) Delete(token string) error {
	err := s.db.Where("token = ?", token).Delete(&Session{}).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func NewDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=server password=server dbname=server port=5432 sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	return db, nil
}