package hasher

import (
	"crypto/rand"
	"crypto/subtle"

	"github.com/esklo/avito-backend-winter-2025/internal/service"
	"golang.org/x/crypto/argon2"
)

var _ service.Hasher = (*Argon2)(nil)

type Argon2 struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func NewArgon2() *Argon2 {
	// can be moved to config
	return &Argon2{
		memory:      32 * 1024,
		iterations:  3,
		parallelism: 8,
		keyLength:   32,
	}
}

func (h *Argon2) Hash(password string) (hash []byte, salt []byte, err error) {
	salt = make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	hash = argon2.IDKey(
		[]byte(password),
		salt,
		h.iterations,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	return hash, salt, nil
}

func (h *Argon2) Verify(password string, hash []byte, salt []byte) bool {
	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		h.iterations,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	return subtle.ConstantTimeCompare(hash, newHash) == 1
}
