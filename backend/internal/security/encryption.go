package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

// EncryptionManager handles encryption/decryption for sensitive data
type EncryptionManager struct {
	masterKey         []byte
	cipher            cipher.Block
	encryptionAlgo    string
	keyDerivationAlgo string
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager(masterKeyPath string) (*EncryptionManager, error) {
	log.Printf("[ENCRYPTION] Initializing encryption manager")

	// Read master key from file
	keyBytes, err := os.ReadFile(masterKeyPath)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to read master key: %v", err)
		return nil, fmt.Errorf("failed to read master key: %w", err)
	}

	// Validate key length (should be 32 bytes for AES-256)
	if len(keyBytes) != 32 {
		log.Printf("[ENCRYPTION] ERROR: Invalid key length: %d (expected 32)", len(keyBytes))
		return nil, fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(keyBytes))
	}

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to create cipher: %v", err)
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	em := &EncryptionManager{
		masterKey:         keyBytes,
		cipher:            block,
		encryptionAlgo:    "AES-256-GCM",
		keyDerivationAlgo: "SHA-256",
	}

	log.Printf("[ENCRYPTION] Encryption manager initialized with %s", em.encryptionAlgo)
	return em, nil
}

// Encrypt encrypts data using AES-256-GCM
func (em *EncryptionManager) Encrypt(plaintext string) (string, error) {
	log.Printf("[ENCRYPTION] Encrypting %d bytes", len(plaintext))

	// Create GCM cipher
	gcm, err := cipher.NewGCM(em.cipher)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to create GCM: %v", err)
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to generate nonce: %v", err)
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	log.Printf("[ENCRYPTION] Encryption successful: %d bytes → %d bytes", len(plaintext), len(encoded))

	return encoded, nil
}

// Decrypt decrypts data using AES-256-GCM
func (em *EncryptionManager) Decrypt(ciphertext string) (string, error) {
	log.Printf("[ENCRYPTION] Decrypting")

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to decode base64: %v", err)
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create GCM cipher
	gcm, err := cipher.NewGCM(em.cipher)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to create GCM: %v", err)
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		log.Printf("[ENCRYPTION] ERROR: Ciphertext too short")
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext_ := decoded[:nonceSize], decoded[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext_, nil)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Decryption failed: %v", err)
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	log.Printf("[ENCRYPTION] Decryption successful")
	return string(plaintext), nil
}

// DeriveKey derives an encryption key from a password
func (em *EncryptionManager) DeriveKey(password string, salt []byte) []byte {
	log.Printf("[ENCRYPTION] Deriving key from password")

	// Simple key derivation using SHA-256
	// In production, use PBKDF2, bcrypt, or Argon2
	h := sha256.New()
	h.Write(append([]byte(password), salt...))
	h.Write(em.masterKey)

	derived := h.Sum(nil)
	log.Printf("[ENCRYPTION] Key derived: %d bytes", len(derived))

	return derived
}

// EncryptSensitiveField encrypts a sensitive database field
func (em *EncryptionManager) EncryptSensitiveField(fieldType string, value string) (string, error) {
	log.Printf("[ENCRYPTION] Encrypting sensitive field: %s", fieldType)

	switch fieldType {
	case "email":
		return em.Encrypt(value)
	case "phone":
		return em.Encrypt(value)
	case "ssn":
		return em.Encrypt(value)
	case "credit_card":
		return em.EncryptCreditCard(value)
	case "password_hash":
		// Don't encrypt password hashes - they're already hashed
		return value, nil
	default:
		log.Printf("[ENCRYPTION] WARNING: Unknown field type: %s", fieldType)
		return em.Encrypt(value)
	}
}

// DecryptSensitiveField decrypts a sensitive database field
func (em *EncryptionManager) DecryptSensitiveField(fieldType string, encryptedValue string) (string, error) {
	log.Printf("[ENCRYPTION] Decrypting sensitive field: %s", fieldType)

	switch fieldType {
	case "email":
		return em.Decrypt(encryptedValue)
	case "phone":
		return em.Decrypt(encryptedValue)
	case "ssn":
		return em.Decrypt(encryptedValue)
	case "credit_card":
		return em.DecryptCreditCard(encryptedValue)
	case "password_hash":
		// Don't decrypt password hashes
		return encryptedValue, nil
	default:
		log.Printf("[ENCRYPTION] WARNING: Unknown field type: %s", fieldType)
		return em.Decrypt(encryptedValue)
	}
}

// EncryptCreditCard encrypts credit card number with masking
func (em *EncryptionManager) EncryptCreditCard(cardNumber string) (string, error) {
	log.Printf("[ENCRYPTION] Encrypting credit card")

	// Mask all but last 4 digits
	if len(cardNumber) < 4 {
		return "", fmt.Errorf("invalid card number length")
	}

	masked := "****-****-****-" + cardNumber[len(cardNumber)-4:]

	// Encrypt the full number for storage
	encrypted, err := em.Encrypt(cardNumber)
	if err != nil {
		return "", err
	}

	// Return format: masked:encrypted
	return masked + ":" + encrypted, nil
}

// DecryptCreditCard decrypts credit card number
func (em *EncryptionManager) DecryptCreditCard(maskedEncrypted string) (string, error) {
	log.Printf("[ENCRYPTION] Decrypting credit card")

	// Extract encrypted part (after the colon)
	parts := make([]string, 2)
	for i, p := range maskedEncrypted {
		if p == ':' {
			parts[0] = maskedEncrypted[:i]
			parts[1] = maskedEncrypted[i+1:]
			break
		}
	}

	if len(parts) != 2 || parts[1] == "" {
		return "", fmt.Errorf("invalid encrypted card format")
	}

	// Decrypt
	return em.Decrypt(parts[1])
}

// GenerateMasterKey generates a new 256-bit master key
func GenerateMasterKey() ([]byte, error) {
	log.Printf("[ENCRYPTION] Generating new master key")

	key := make([]byte, 32) // 256 bits
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to generate key: %v", err)
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(key)
	log.Printf("[ENCRYPTION] Master key generated: %s (first 16 chars)", encoded[:16])

	return key, nil
}

// SaveMasterKeyToFile saves master key to a file with restricted permissions
func SaveMasterKeyToFile(key []byte, filePath string) error {
	log.Printf("[ENCRYPTION] Saving master key to file: %s", filePath)

	encoded := base64.StdEncoding.EncodeToString(key)

	// Write key to file
	err := os.WriteFile(filePath, []byte(encoded), 0600) // Read/write for owner only
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to write key file: %v", err)
		return fmt.Errorf("failed to write key file: %w", err)
	}

	log.Printf("[ENCRYPTION] Master key saved successfully")
	return nil
}

// ============================================================================
// TLS/TRANSPORT ENCRYPTION
// ============================================================================

// TLSConfig holds TLS configuration
type TLSConfig struct {
	CertFile     string
	KeyFile      string
	CAFile       string
	MinVersion   string
	CipherSuites []string
}

// LoadTLSConfig loads TLS configuration from files
func LoadTLSConfig(config TLSConfig) (*tls.Config, error) {
	log.Printf("[TLS] Loading TLS configuration")

	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		log.Printf("[TLS] ERROR: Failed to load certificate: %v", err)
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
	}

	log.Printf("[TLS] TLS configuration loaded successfully")
	return tlsConfig, nil
}

// ============================================================================
// ENCRYPTION KEY ROTATION
// ============================================================================

// RotateKeys rotates encryption keys
func (em *EncryptionManager) RotateKeys(newMasterKeyPath string) error {
	log.Printf("[ENCRYPTION] Starting key rotation")

	// Read new master key
	newKeyBytes, err := os.ReadFile(newMasterKeyPath)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to read new master key: %v", err)
		return fmt.Errorf("failed to read new master key: %w", err)
	}

	if len(newKeyBytes) != 32 {
		return fmt.Errorf("invalid new key length: expected 32 bytes, got %d", len(newKeyBytes))
	}

	// Create new cipher
	newCipher, err := aes.NewCipher(newKeyBytes)
	if err != nil {
		log.Printf("[ENCRYPTION] ERROR: Failed to create new cipher: %v", err)
		return fmt.Errorf("failed to create new cipher: %w", err)
	}

	// Update manager with new key
	em.masterKey = newKeyBytes
	em.cipher = newCipher

	log.Printf("[ENCRYPTION] Key rotation completed successfully")
	return nil
}
