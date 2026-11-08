package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	MagicHeader    = "1VAULT\x01"
	SaltSize       = 16
	NonceSize      = 12
	KeyLength      = 32 // 256 bits for AES-256
	ArgonTime      = 3
	ArgonMemory    = 64 * 1024 // 64 MB
	ArgonThreads   = 4
)

var (
	ErrInvalidHeader = errors.New("invalid vault file header or unsupported version")
	ErrDecryptFailed = errors.New("decryption failed: incorrect password or corrupted vault")
)

// GenerateRandomBytes creates cryptographically secure random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %w", err)
	}
	return b, nil
}

// GenerateSecretKey creates a 1Password-style 128-bit Secret Key (e.g., A3-XXXXXX-XXXXXX-XXXXXX-XXXXXX)
func GenerateSecretKey() (string, error) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}
	// Format into readable grouped chunks
	hexStr := fmt.Sprintf("%X", b)
	return fmt.Sprintf("A3-%s-%s-%s-%s", hexStr[0:6], hexStr[6:12], hexStr[12:18], hexStr[18:24]), nil
}

// DeriveKey computes a 256-bit symmetric encryption key using Argon2id
func DeriveKey(password string, secretKey string, salt []byte) []byte {
	secretKey = strings.TrimSpace(secretKey)
	var composite []byte
	if secretKey != "" {
		composite = []byte(password + "|" + secretKey)
	} else {
		composite = []byte(password)
	}
	return argon2.IDKey(composite, salt, ArgonTime, ArgonMemory, ArgonThreads, KeyLength)
}

// Encrypt encrypts plaintext using AES-256-GCM and prepends the magic header, salt, and nonce.
func Encrypt(plaintext []byte, password string, secretKey string) ([]byte, error) {
	salt, err := GenerateRandomBytes(SaltSize)
	if err != nil {
		return nil, err
	}

	key := DeriveKey(password, secretKey, salt)
	defer ZeroBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce, err := GenerateRandomBytes(NonceSize)
	if err != nil {
		return nil, err
	}

	// Encrypt & authenticate (AEAD)
	ciphertext := gcm.Seal(nil, nonce, plaintext, []byte(MagicHeader))

	// Package: [MagicHeader (7B)] + [Salt (16B)] + [Nonce (12B)] + [Ciphertext + Tag]
	header := []byte(MagicHeader)
	out := make([]byte, 0, len(header)+SaltSize+NonceSize+len(ciphertext))
	out = append(out, header...)
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return out, nil
}

// Decrypt verifies the header, derives the key with Argon2id, and decrypts the AES-256-GCM payload.
func Decrypt(data []byte, password string, secretKey string) ([]byte, error) {
	headerLen := len(MagicHeader)
	minLen := headerLen + SaltSize + NonceSize + 16 // 16 bytes minimum for GCM tag
	if len(data) < minLen {
		return nil, ErrInvalidHeader
	}

	if subtle.ConstantTimeCompare(data[:headerLen], []byte(MagicHeader)) != 1 {
		return nil, ErrInvalidHeader
	}

	salt := data[headerLen : headerLen+SaltSize]
	nonce := data[headerLen+SaltSize : headerLen+SaltSize+NonceSize]
	ciphertext := data[headerLen+SaltSize+NonceSize:]

	key := DeriveKey(password, secretKey, salt)
	defer ZeroBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, []byte(MagicHeader))
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

// ZeroBytes overwrites memory buffer with zeroes
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

