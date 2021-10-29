package model

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type User struct {
	ID        string      `json:"id,omitempty" firestore:"id,omitempty"`
	CreatedAt interface{} `json:"createdAt,omitempty" firestore:"createdAt,omitempty"`
	Name      string      `json:"name,omitempty" firestore:"name,omitempty"`
	Email     string      `json:"email,omitempty" firestore:"email,omitempty"`
	Password  string      `json:"password,omitempty" firestore:"password,omitempty"`
	Posts     []Post      `json:"posts,omitempty" firestore:"posts,omitempty"`
}

type passwordConfig struct {
	time      uint32
	memory    uint32
	threads   uint8
	keyLength uint32
}

func (user *User) ComparePassword(password string) (bool, error) {
	parts := strings.Split(user.Password, "$")
	if len(parts) < 4 {
		return false, errors.New("invalid string")
	}

	c := &passwordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	c.keyLength = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLength)
	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}

func (user *User) HashAndSetPassword() error {
	c := &passwordConfig{
		time:      1,
		memory:    64 * 1024,
		threads:   4,
		keyLength: 32,
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(user.Password), salt, c.time, c.memory, c.threads, c.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	user.Password = fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return nil
}
