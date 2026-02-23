package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/frostyeti/hyprship/apps/api/crypto/hashalgo"
	"golang.org/x/crypto/pbkdf2"
)

type Encryptor interface {
	// Encrypt plaintext with AES-256-GCM. Key is derived from passphrase via PBKDF2.
	Encrypt(passphrase string, plaintext []byte) (string, error)
}

type Decryptor interface {
	// Decrypt a value produced by Encrypt.
	Decrypt(passphrase string, encoded string) ([]byte, error)
}

type EncryptionDriver interface {
	Encryptor
	Decryptor
}

type aesGcmDriver struct {
	Key        []byte
	Iterations uint32
	Algo       uint16
}

func (d *aesGcmDriver) Encrypt(passphrase string, plaintext []byte) (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	keySize := 32
	derivedKey := pbkdf2.Key([]byte(passphrase), salt, int(d.Iterations), keySize, getHashFunc(d.Algo))

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	sealed := gcm.Seal(nil, nonce, plaintext, nil)

	// gcm.Seal appends the tag to the ciphertext
	tagSize := 16
	if len(sealed) < tagSize {
		return "", errors.New("sealed data too small")
	}

	ciphertext := sealed[:len(plaintext)]
	tag := sealed[len(plaintext):]

	// Binary layout:
	// [0:2]              salt size   (u16, 32)
	// [2:4]              key size    (u16, 32)
	// [4:6]              pbkdf2 algo (u16, SHA256=2)
	// [6:8]              nonce size  (u16, 12)
	// [8:10]             tag size    (u16, 16)
	// [10:14]            iterations  (u32, 100000)
	// [14:46]            salt        (32 bytes)
	// [46:58]            nonce       (12 bytes)
	// [58:74]            tag         (16 bytes)
	// [74:]              ciphertext

	bufLen := 74 + len(ciphertext)
	buf := make([]byte, bufLen)

	binary.BigEndian.PutUint16(buf[0:2], uint16(32))
	binary.BigEndian.PutUint16(buf[2:4], uint16(32))
	binary.BigEndian.PutUint16(buf[4:6], d.Algo)
	binary.BigEndian.PutUint16(buf[6:8], uint16(12))
	binary.BigEndian.PutUint16(buf[8:10], uint16(16))
	binary.BigEndian.PutUint32(buf[10:14], d.Iterations)

	copy(buf[14:46], salt)
	copy(buf[46:58], nonce)
	copy(buf[58:74], tag)
	if len(ciphertext) > 0 {
		copy(buf[74:], ciphertext)
	}

	encoded := base64.StdEncoding.EncodeToString(buf)
	return "aesgcm:v1:" + encoded, nil
}

func (d *aesGcmDriver) Decrypt(passphrase string, encoded string) ([]byte, error) {
	if !strings.HasPrefix(encoded, "aesgcm:v1:") {
		return nil, errors.New("invalid encryption prefix")
	}

	b64Data := encoded[len("aesgcm:v1:"):]
	buf, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return nil, err
	}

	if len(buf) < 74 {
		return nil, errors.New("invalid encrypted data length")
	}

	saltSize := binary.BigEndian.Uint16(buf[0:2])
	keySize := binary.BigEndian.Uint16(buf[2:4])
	algo := binary.BigEndian.Uint16(buf[4:6])
	nonceSize := binary.BigEndian.Uint16(buf[6:8])
	tagSize := binary.BigEndian.Uint16(buf[8:10])
	iterations := binary.BigEndian.Uint32(buf[10:14])

	if saltSize != 32 || keySize != 32 || nonceSize != 12 || tagSize != 16 {
		return nil, errors.New("invalid encryption parameters in header")
	}

	salt := buf[14:46]
	nonce := buf[46:58]
	tag := buf[58:74]
	ciphertext := buf[74:]

	derivedKey := pbkdf2.Key([]byte(passphrase), salt, int(iterations), int(keySize), getHashFunc(algo))

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// gcm.Open expects ciphertext + tag
	sealed := make([]byte, 0, len(ciphertext)+len(tag))
	sealed = append(sealed, ciphertext...)
	sealed = append(sealed, tag...)

	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return nil, errors.New("decryption failed: " + err.Error())
	}

	if plaintext == nil {
		return []byte{}, nil
	}

	return plaintext, nil
}

func GetEncryptionDriver(kind string, key []byte) EncryptionDriver {
	switch kind {
	case "", "default", "aesgcm":
		return &aesGcmDriver{
			Key:        key,
			Iterations: uint32(100000),
			Algo:       hashalgo.SHA256,
		}
	default:
		return &aesGcmDriver{
			Key:        key,
			Iterations: uint32(100000),
			Algo:       hashalgo.SHA256,
		}
	}
}
