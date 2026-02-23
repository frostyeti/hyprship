package crypto

import (
	"encoding/base64"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/frostyeti/hyprship/apps/api/crypto/hashalgo"
	"github.com/stretchr/testify/assert"
)

func TestAESGCM_EncryptDecrypt(t *testing.T) {
	driver := GetEncryptionDriver("aesgcm", []byte("unused"))

	cases := []struct {
		passphrase string
		plaintext  []byte
	}{
		{"mypassword123", []byte("hello world")},
		{"secret", []byte("")},
		{"longpassphrase" + strings.Repeat("a", 100), []byte("complex 🔐 unicode")},
		{"12345", make([]byte, 1024)},
	}

	for _, tc := range cases {
		t.Run("Encrypt_"+tc.passphrase, func(t *testing.T) {
			encrypted, err := driver.Encrypt(tc.passphrase, tc.plaintext)
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(encrypted, "aesgcm:v1:"))

			decrypted, err := driver.Decrypt(tc.passphrase, encrypted)
			assert.NoError(t, err)
			assert.Equal(t, tc.plaintext, decrypted)

			_, errWrongPw := driver.Decrypt(tc.passphrase+"wrong", encrypted)
			assert.Error(t, errWrongPw, "wrong passphrase should fail decryption")
		})
	}
}

func TestAESGCM_UniqueOutput(t *testing.T) {
	driver := GetEncryptionDriver("aesgcm", []byte("unused"))

	enc1, _ := driver.Encrypt("pass", []byte("same text"))
	enc2, _ := driver.Encrypt("pass", []byte("same text"))

	assert.NotEqual(t, enc1, enc2, "encryption with same pass/text should produce different outputs due to random salt/nonce")
}

func TestAESGCM_Errors(t *testing.T) {
	driver := GetEncryptionDriver("aesgcm", []byte("unused"))

	_, err := driver.Decrypt("pass", "badprefix:v1:something")
	assert.Error(t, err, "bad prefix should error")

	_, err = driver.Decrypt("pass", "aesgcm:v1:!!!invalidbase64")
	assert.Error(t, err, "invalid base64 should error")

	shortBuf := make([]byte, 70)
	shortEnc := "aesgcm:v1:" + base64.StdEncoding.EncodeToString(shortBuf)
	_, err = driver.Decrypt("pass", shortEnc)
	assert.Error(t, err, "short buffer should error")
}

func TestAESGCM_BinaryLayout(t *testing.T) {
	driver := GetEncryptionDriver("aesgcm", []byte("unused"))

	plaintext := []byte("secretdata")
	encrypted, _ := driver.Encrypt("passphrase", plaintext)

	encoded := encrypted[10:] // remove prefix
	buf, err := base64.StdEncoding.DecodeString(encoded)
	assert.NoError(t, err)

	// Base layout is 74 bytes + plaintext
	assert.Equal(t, 74+len(plaintext), len(buf))

	assert.Equal(t, uint16(32), binary.BigEndian.Uint16(buf[0:2]), "salt size")
	assert.Equal(t, uint16(32), binary.BigEndian.Uint16(buf[2:4]), "key size")
	assert.Equal(t, hashalgo.SHA256, binary.BigEndian.Uint16(buf[4:6]), "algo")
	assert.Equal(t, uint16(12), binary.BigEndian.Uint16(buf[6:8]), "nonce size")
	assert.Equal(t, uint16(16), binary.BigEndian.Uint16(buf[8:10]), "tag size")
	assert.Equal(t, uint32(100000), binary.BigEndian.Uint32(buf[10:14]), "iterations")

	salt := buf[14:46]
	assert.Equal(t, 32, len(salt))

	nonce := buf[46:58]
	assert.Equal(t, 12, len(nonce))

	tag := buf[58:74]
	assert.Equal(t, 16, len(tag))

	ciphertext := buf[74:]
	assert.Equal(t, len(plaintext), len(ciphertext))
}
