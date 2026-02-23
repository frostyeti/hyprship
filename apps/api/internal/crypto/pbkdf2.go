package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash"
	"strings"

	"github.com/frostyeti/hyprship/apps/api/internal/crypto/hashalgo"
	"golang.org/x/crypto/pbkdf2"
)

type PasswordHash interface {
	// Hash a plaintext password.
	Hash(password string) (string, error)
}

type PasswordVerify interface {
	// Verify plaintext against a stored hash. Returns (true, nil) on match.
	Verify(hashed, plaintext string) (bool, error)

	// Returns true when the stored hash was created with different parameters
	// from the current defaults (e.g. after an iterations upgrade).
	NeedsRehash(hashed string) bool
}

type PasswordHasher interface {
	PasswordHash
	PasswordVerify
}

type pbkdf2Hasher struct {
	Iterations uint32
	Algo       uint16
}

func getHashFunc(algo uint16) func() hash.Hash {
	switch algo {
	case hashalgo.SHA384:
		return sha512.New384
	case hashalgo.SHA512:
		return sha512.New
	default:
		return sha256.New
	}
}

func (h *pbkdf2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	keySize := 32
	derivedKey := pbkdf2.Key([]byte(password), salt, int(h.Iterations), keySize, getHashFunc(h.Algo))

	// Binary layout:
	// [0:2]   salt size  (u16, always 32)
	// [2:4]   hash algo  (u16, SHA256=2)
	// [4:8]   iterations (u32)
	// [8:40]  salt       (32 bytes)
	// [40:72] hash       (32 bytes)
	buf := make([]byte, 72)
	binary.BigEndian.PutUint16(buf[0:2], uint16(len(salt)))
	binary.BigEndian.PutUint16(buf[2:4], h.Algo)
	binary.BigEndian.PutUint32(buf[4:8], h.Iterations)
	copy(buf[8:40], salt)
	copy(buf[40:72], derivedKey)

	encoded := base64.StdEncoding.EncodeToString(buf)
	return "pbkdf2:v1:" + encoded, nil
}

func (h *pbkdf2Hasher) extractParts(hashed string) ([]byte, error) {
	if !strings.HasPrefix(hashed, "pbkdf2:v1:") {
		return nil, errors.New("invalid hash prefix")
	}

	encoded := hashed[len("pbkdf2:v1:"):]
	buf, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	if len(buf) != 72 {
		return nil, errors.New("invalid hash length")
	}
	return buf, nil
}

func (h *pbkdf2Hasher) Verify(hashed, plaintext string) (bool, error) {
	buf, err := h.extractParts(hashed)
	if err != nil {
		return false, err
	}

	saltSize := binary.BigEndian.Uint16(buf[0:2])
	if saltSize != 32 {
		return false, errors.New("invalid salt size")
	}

	algo := binary.BigEndian.Uint16(buf[2:4])
	iterations := binary.BigEndian.Uint32(buf[4:8])
	salt := buf[8:40]
	storedHash := buf[40:72]

	keySize := 32
	derivedKey := pbkdf2.Key([]byte(plaintext), salt, int(iterations), keySize, getHashFunc(algo))

	match := subtle.ConstantTimeCompare(storedHash, derivedKey) == 1
	return match, nil
}

func (h *pbkdf2Hasher) NeedsRehash(hashed string) bool {
	buf, err := h.extractParts(hashed)
	if err != nil {
		return false
	}

	algo := binary.BigEndian.Uint16(buf[2:4])
	iterations := binary.BigEndian.Uint32(buf[4:8])

	return algo != h.Algo || iterations != h.Iterations
}

func GetPasswordHasher(kind string) PasswordHasher {
	switch kind {
	case "", "default", "pbkdf2":
		return &pbkdf2Hasher{
			Iterations: uint32(100000),
			Algo:       hashalgo.SHA256,
		}
	default:
		return &pbkdf2Hasher{
			Iterations: uint32(100000),
			Algo:       hashalgo.SHA256,
		}
	}
}
