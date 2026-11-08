package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

const (
	BlobMagicHeader    = "1VBLOB\x01"
	BlobNoncePrefixLen = 8
	BlobChunkSizeLen   = 4
	BlobHeaderSize     = 7 + BlobNoncePrefixLen + BlobChunkSizeLen // 19 bytes
	DefaultChunkSize   = 64 * 1024                                 // 64 KB (65536 bytes)
	GCMTagSize         = 16
)

var (
	ErrInvalidBlobHeader = errors.New("invalid encrypted blob header")
	ErrInvalidBlobRange  = errors.New("invalid byte range requested")
	ErrBlobCorrupted     = errors.New("encrypted blob authentication failed or file is corrupted")
)

// GenerateFileKey creates a random 256-bit (32-byte) File Encryption Key encoded as hex
func GenerateFileKey() (string, error) {
	b, err := GenerateRandomBytes(KeyLength)
	if err != nil {
		return "", err
	}
	defer ZeroBytes(b)
	return hex.EncodeToString(b), nil
}

func makeChunkNonce(prefix []byte, chunkIndex uint64) []byte {
	nonce := make([]byte, NonceSize) // 12 bytes
	copy(nonce[0:8], prefix[0:8])
	binary.BigEndian.PutUint32(nonce[8:12], uint32(chunkIndex))
	return nonce
}

func makeChunkAAD(chunkIndex uint64, isFinal bool) []byte {
	aad := make([]byte, 9)
	binary.BigEndian.PutUint64(aad[0:8], chunkIndex)
	if isFinal {
		aad[8] = 0x01
	} else {
		aad[8] = 0x00
	}
	return aad
}

// EncryptStream reads plaintext from r in 64 KB chunks, encrypts each chunk with AES-256-GCM
// using the hex-encoded 256-bit File Encryption Key (fekHex), and writes the 1VBLOB stream to w.
// Returns the total plaintext byte count and hex-encoded SHA-256 checksum of the plaintext.
func EncryptStream(r io.Reader, w io.Writer, fekHex string) (int64, string, error) {
	key, err := hex.DecodeString(fekHex)
	if err != nil || len(key) != KeyLength {
		return 0, "", errors.New("invalid file encryption key")
	}
	defer ZeroBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, "", fmt.Errorf("failed to initialize GCM: %w", err)
	}

	noncePrefix, err := GenerateRandomBytes(BlobNoncePrefixLen)
	if err != nil {
		return 0, "", err
	}

	// Write 19-byte header: [Magic (7B)] [NoncePrefix (8B)] [ChunkSize (4B)]
	header := make([]byte, BlobHeaderSize)
	copy(header[0:7], []byte(BlobMagicHeader))
	copy(header[7:15], noncePrefix)
	binary.BigEndian.PutUint32(header[15:19], uint32(DefaultChunkSize))

	if _, err := w.Write(header); err != nil {
		return 0, "", fmt.Errorf("failed to write blob header: %w", err)
	}

	hasher := sha256.New()
	var totalPlain int64

	currBuf := make([]byte, DefaultChunkSize)
	nextBuf := make([]byte, DefaultChunkSize)
	cipherBuf := make([]byte, 0, DefaultChunkSize+GCMTagSize)

	currN, currErr := io.ReadFull(r, currBuf)
	if currN == 0 && (currErr == io.EOF || currErr == io.ErrUnexpectedEOF) {
		// Empty 0-byte file: write single authenticated empty final chunk
		nonce := makeChunkNonce(noncePrefix, 0)
		aad := makeChunkAAD(0, true)
		sealed := gcm.Seal(cipherBuf[:0], nonce, nil, aad)
		if _, err := w.Write(sealed); err != nil {
			return 0, "", fmt.Errorf("failed to write empty final chunk: %w", err)
		}
		return 0, hex.EncodeToString(hasher.Sum(nil)), nil
	}
	if currErr != nil && currErr != io.EOF && currErr != io.ErrUnexpectedEOF {
		return 0, "", fmt.Errorf("failed reading input stream: %w", currErr)
	}

	var chunkIdx uint64
	for {
		plainSlice := currBuf[:currN]
		hasher.Write(plainSlice)
		totalPlain += int64(currN)

		var isFinal bool
		var nextN int
		var nextErr error

		if currN < DefaultChunkSize || currErr == io.EOF || currErr == io.ErrUnexpectedEOF {
			isFinal = true
		} else {
			// Peek next chunk to know if current chunk is the final chunk
			nextN, nextErr = io.ReadFull(r, nextBuf)
			if nextErr != nil && nextErr != io.EOF && nextErr != io.ErrUnexpectedEOF {
				return totalPlain, "", fmt.Errorf("failed reading input stream: %w", nextErr)
			}
			if nextN == 0 {
				isFinal = true
			}
		}

		nonce := makeChunkNonce(noncePrefix, chunkIdx)
		aad := makeChunkAAD(chunkIdx, isFinal)
		sealed := gcm.Seal(cipherBuf[:0], nonce, plainSlice, aad)

		if _, err := w.Write(sealed); err != nil {
			return totalPlain, "", fmt.Errorf("failed writing encrypted chunk %d: %w", chunkIdx, err)
		}

		if isFinal {
			break
		}

		// Swap buffers and advance to next chunk
		currBuf, nextBuf = nextBuf, currBuf
		currN = nextN
		currErr = nextErr
		chunkIdx++
	}

	return totalPlain, hex.EncodeToString(hasher.Sum(nil)), nil
}

// DecryptStreamRange decrypts the requested inclusive plaintext byte range [startByte, endByte]
// from an encrypted 1VBLOB stream r using O(1) chunk seeking and writes the plaintext to w.
// If totalPlainSize == 0, it verifies the empty final chunk and returns nil.
func DecryptStreamRange(r io.ReadSeeker, w io.Writer, fekHex string, totalPlainSize int64, startByte int64, endByte int64) error {
	key, err := hex.DecodeString(fekHex)
	if err != nil || len(key) != KeyLength {
		return errors.New("invalid file encryption key")
	}
	defer ZeroBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to initialize GCM: %w", err)
	}

	// Read and verify 19-byte header
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed seeking blob header: %w", err)
	}
	header := make([]byte, BlobHeaderSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return ErrInvalidBlobHeader
	}
	if subtle.ConstantTimeCompare(header[0:7], []byte(BlobMagicHeader)) != 1 {
		return ErrInvalidBlobHeader
	}
	noncePrefix := header[7:15]
	chunkSize := int64(binary.BigEndian.Uint32(header[15:19]))
	if chunkSize <= 0 || chunkSize > 16*1024*1024 {
		return ErrInvalidBlobHeader
	}

	if totalPlainSize == 0 {
		cipherChunk := make([]byte, GCMTagSize)
		if _, err := io.ReadFull(r, cipherChunk); err != nil {
			return ErrBlobCorrupted
		}
		nonce := makeChunkNonce(noncePrefix, 0)
		aad := makeChunkAAD(0, true)
		if _, err := gcm.Open(nil, nonce, cipherChunk, aad); err != nil {
			return ErrBlobCorrupted
		}
		return nil
	}

	if startByte < 0 || endByte < startByte || endByte >= totalPlainSize {
		return ErrInvalidBlobRange
	}

	totalChunks := (totalPlainSize + chunkSize - 1) / chunkSize
	lastChunkIdx := uint64(totalChunks - 1)

	startChunk := uint64(startByte / chunkSize)
	endChunk := uint64(endByte / chunkSize)

	encryptedChunkSize := chunkSize + GCMTagSize
	seekPos := int64(BlobHeaderSize) + int64(startChunk)*encryptedChunkSize

	if _, err := r.Seek(seekPos, io.SeekStart); err != nil {
		return fmt.Errorf("failed seeking to chunk %d: %w", startChunk, err)
	}

	cipherBuf := make([]byte, encryptedChunkSize)
	plainBuf := make([]byte, 0, chunkSize)

	for chunkIdx := startChunk; chunkIdx <= endChunk; chunkIdx++ {
		isFinal := chunkIdx == lastChunkIdx
		expectedPlainLen := chunkSize
		if isFinal {
			rem := totalPlainSize % chunkSize
			if rem != 0 {
				expectedPlainLen = rem
			}
		}
		expectedCipherLen := expectedPlainLen + GCMTagSize

		if _, err := io.ReadFull(r, cipherBuf[:expectedCipherLen]); err != nil {
			return ErrBlobCorrupted
		}

		nonce := makeChunkNonce(noncePrefix, chunkIdx)
		aad := makeChunkAAD(chunkIdx, isFinal)

		decrypted, err := gcm.Open(plainBuf[:0], nonce, cipherBuf[:expectedCipherLen], aad)
		if err != nil {
			return ErrBlobCorrupted
		}

		// Determine slice bounds within this decrypted chunk
		sliceStart := int64(0)
		if chunkIdx == startChunk {
			sliceStart = startByte % chunkSize
		}
		sliceEnd := expectedPlainLen
		if chunkIdx == endChunk {
			sliceEnd = (endByte % chunkSize) + 1
		}

		if sliceStart < 0 || sliceEnd > int64(len(decrypted)) || sliceStart > sliceEnd {
			return ErrInvalidBlobRange
		}

		if _, err := w.Write(decrypted[sliceStart:sliceEnd]); err != nil {
			return err
		}
	}

	return nil
}

