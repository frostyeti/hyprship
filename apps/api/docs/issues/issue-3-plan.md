# Issue #3: Crypto Module for API Project

This plan details the implementation of a cryptography module for hashing and encrypting credentials, addressing Issue #3.

## Phase 1: Foundation (hashalgo)
Create the `hashalgo` package to establish a standard set of numeric identifiers for hash algorithms used in binary representations.

1. **Path:** `apps/api/crypto/hashalgo/hashalgo.go`
2. **Implementation:**
   - Define constants: `Invalid (0)`, `SHA256 (2)`, `SHA384 (3)`, `SHA512 (4)`.
   - Implement `ToString(value uint16) string` mapping numeric constants to string representations ("SHA256", etc.).
   - Implement `ToValue(value string) uint16` mapping string representations to numeric constants.

## Phase 2: Password Hashing (pbkdf2)
Implement a password hashing utility using PBKDF2 with SHA-256, a 32-byte salt, and 100,000 iterations.

1. **Path:** `apps/api/crypto/pbkdf2.go`
2. **Implementation:**
   - Define interfaces: `PasswordHash`, `PasswordVerify`, and a composite `PasswordHasher` struct that embeds both.
   - Implement `pbkdf2Hasher` struct satisfying these interfaces.
   - Create the factory function `GetPasswordHasher(kind string) PasswordHasher`.
3. **Hashing Logic (`Hash`):**
   - Generate a 32-byte random salt using `crypto/rand`.
   - Use `golang.org/x/crypto/pbkdf2` to derive a 32-byte key with 100,000 iterations and SHA-256.
   - Pack the binary layout strictly as specified:
     - `[0:2]` Salt size (uint16 = 32)
     - `[2:4]` Hash algo (uint16 = 2)
     - `[4:8]` Iterations (uint32 = 100000)
     - `[8:40]` Salt (32 bytes)
     - `[40:72]` Hash (32 bytes)
     - Total binary length = 72 bytes.
   - Base64 encode the binary buffer and prepend `"pbkdf2:v1:"`.
4. **Verification Logic (`Verify`):**
   - Strip the `"pbkdf2:v1:"` prefix and base64 decode.
   - Extract parameters (salt size, algo, iterations) and the salt/hash from the binary layout.
   - Re-derive the key using the extracted parameters and the provided plaintext password.
   - Compare the derived key with the extracted hash using `crypto/subtle.ConstantTimeCompare`.
5. **Rehash Logic (`NeedsRehash`):**
   - Extract iterations and algo from the hash and return `true` if they differ from the `pbkdf2Hasher`'s current configured defaults.

## Phase 3: Encryption (aesgcm)
Implement a symmetric encryption utility using AES-256-GCM. The key is derived using PBKDF2-SHA256 from a provided passphrase.

1. **Path:** `apps/api/crypto/aesgcm.go` (Note: Issue says `internal/crypto/aesgcm.go` but since we are building an `apps/api/crypto` module and `pbkdf2` is there, we will put it in `apps/api/crypto/aesgcm.go` for consistency unless explicitly directed otherwise. We will use `apps/api/crypto/aesgcm.go`).
2. **Implementation:**
   - Define interfaces: `Encryptor`, `Decryptor`, and a composite `EncryptionDriver`.
   - Implement `aesGcmDriver` struct satisfying these interfaces.
   - Create the factory function `GetEncryptionDriver(kind string, key []byte) EncryptionDriver`.
3. **Encryption Logic (`Encrypt`):**
   - Generate a 32-byte random salt and 12-byte random nonce using `crypto/rand`.
   - Derive a 32-byte AES key from the passphrase and salt using PBKDF2 (SHA-256, 100,000 iterations).
   - Create an AES cipher block and GCM wrapper.
   - Encrypt the plaintext using the GCM wrapper (which appends the 16-byte authentication tag to the ciphertext).
   - Pack the binary layout strictly as specified:
     - `[0:2]` Salt size (uint16 = 32)
     - `[2:4]` Key size (uint16 = 32)
     - `[4:6]` PBKDF2 algo (uint16 = 2)
     - `[6:8]` Nonce size (uint16 = 12)
     - `[8:10]` Tag size (uint16 = 16)
     - `[10:14]` Iterations (uint32 = 100000)
     - `[14:46]` Salt (32 bytes)
     - `[46:58]` Nonce (12 bytes)
     - `[58:]` Ciphertext (including the 16-byte tag)
       - *Correction on binary layout specified in issue:* `crypto/cipher.AEAD.Seal` appends the tag to the ciphertext automatically. We will pack `[58:]` with the sealed output (which is ciphertext + tag) rather than manually separating the tag, while ensuring we map the boundaries correctly during decryption. If strict byte mapping is required for the tag, we will separate it manually before packing.
   - Base64 encode the binary buffer and prepend `"aesgcm:v1:"`.
4. **Decryption Logic (`Decrypt`):**
   - Strip the `"aesgcm:v1:"` prefix and base64 decode.
   - Extract parameters, salt, nonce, and ciphertext (with tag) from the binary layout.
   - Re-derive the AES key from the passphrase and extracted salt using PBKDF2.
   - Create the AES cipher block and GCM wrapper.
   - Decrypt the ciphertext using the GCM wrapper.

## Phase 4: Testing
Write comprehensive unit tests ensuring robustness and correctness.

1. **Path:** `apps/api/crypto/pbkdf2_test.go` and `apps/api/crypto/aesgcm_test.go`
2. **Coverage:**
   - Verify proper prefix generation (`pbkdf2:v1:` and `aesgcm:v1:`).
   - Ensure uniqueness of output for identical inputs (salting/nonce check).
   - Successful verification/decryption with valid credentials.
   - Failed verification/decryption with invalid credentials.
   - Handling of malformed, truncated, or invalid base64 data.
   - Test `NeedsRehash` functionality.
   - Table-driven tests covering empty strings, long strings, and Unicode characters.
   - Explicit binary layout verification tests to ensure byte packing matches the specification exactly.

## Phase 5: Documentation Updates
1. **New Doc:** Create `apps/api/docs/crypto.md` detailing the binary layouts, usage, and parameters for PBKDF2 and AES-GCM.
2. **Update AGENTS.md:** Append a new "Cryptography" section to `apps/api/AGENTS.md` explaining:
   - Use hashing (PBKDF2) for internal application passwords and API keys.
   - Use encryption (AES-GCM) for storing external secrets/passwords that need to be decrypted later.
   - Note that database columns ending in `_digest` should utilize the hasher.