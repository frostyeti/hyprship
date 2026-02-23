package crypto

import (
	"encoding/base64"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/frostyeti/hyprship/apps/api/crypto/hashalgo"
	"github.com/stretchr/testify/assert"
)

func TestPBKDF2_HashAndVerify(t *testing.T) {
	hasher := GetPasswordHasher("pbkdf2")

	passwords := []string{
		"password123",
		"",
		"p@ssw0rd!#$%",
		"世界",
		strings.Repeat("a", 1000),
	}

	for _, pw := range passwords {
		t.Run("Verify_"+pw, func(t *testing.T) {
			hashed, err := hasher.Hash(pw)
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(hashed, "pbkdf2:v1:"))

			match, err := hasher.Verify(hashed, pw)
			assert.NoError(t, err)
			assert.True(t, match, "password should match")

			matchWrong, err := hasher.Verify(hashed, pw+"wrong")
			assert.NoError(t, err)
			assert.False(t, matchWrong, "wrong password should not match")
		})
	}
}

func TestPBKDF2_UniqueSalts(t *testing.T) {
	hasher := GetPasswordHasher("pbkdf2")
	hash1, _ := hasher.Hash("samepassword")
	hash2, _ := hasher.Hash("samepassword")
	assert.NotEqual(t, hash1, hash2, "hashes should be unique due to random salt")
}

func TestPBKDF2_NeedsRehash(t *testing.T) {
	hasher := GetPasswordHasher("pbkdf2")
	hashed, _ := hasher.Hash("password")

	assert.False(t, hasher.NeedsRehash(hashed), "fresh hash should not need rehash")

	// Create a modified hash string with different iterations to simulate an old version
	buf, _ := base64.StdEncoding.DecodeString(hashed[10:])
	binary.BigEndian.PutUint32(buf[4:8], 50000) // Lower iterations
	oldHash := "pbkdf2:v1:" + base64.StdEncoding.EncodeToString(buf)

	assert.True(t, hasher.NeedsRehash(oldHash), "modified iterations should require rehash")
}

func TestPBKDF2_Errors(t *testing.T) {
	hasher := GetPasswordHasher("pbkdf2")

	_, err := hasher.Verify("badprefix:v1:something", "pass")
	assert.Error(t, err, "bad prefix should error")

	_, err = hasher.Verify("pbkdf2:v1:!!!invalidbase64", "pass")
	assert.Error(t, err, "invalid base64 should error")

	shortBuf := make([]byte, 70)
	shortHash := "pbkdf2:v1:" + base64.StdEncoding.EncodeToString(shortBuf)
	_, err = hasher.Verify(shortHash, "pass")
	assert.Error(t, err, "wrong length should error")
}

func TestPBKDF2_BinaryLayout(t *testing.T) {
	hasher := GetPasswordHasher("pbkdf2")
	hashed, _ := hasher.Hash("secret")

	encoded := hashed[10:] // remove prefix
	buf, err := base64.StdEncoding.DecodeString(encoded)
	assert.NoError(t, err)

	assert.Equal(t, 72, len(buf))

	saltSize := binary.BigEndian.Uint16(buf[0:2])
	assert.Equal(t, uint16(32), saltSize)

	algo := binary.BigEndian.Uint16(buf[2:4])
	assert.Equal(t, hashalgo.SHA256, algo)

	iterations := binary.BigEndian.Uint32(buf[4:8])
	assert.Equal(t, uint32(100000), iterations)
}
