package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"
)

func testKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("generating test key: %v", err)
	}
	return key
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := testKey(t)
	plaintext := "b42dm850fv8jok7ksuwlh81ozomh19"

	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if !IsEncrypted(encrypted) {
		t.Fatal("expected encrypted value to have enc: prefix")
	}
	if !strings.HasPrefix(encrypted, Prefix+AlgAES256GCM+":") {
		t.Fatalf("unexpected prefix: %s", encrypted)
	}

	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("roundtrip failed: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptProducesUniqueCiphertexts(t *testing.T) {
	key := testKey(t)
	plaintext := "same-value"

	e1, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt #1: %v", err)
	}
	e2, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt #2: %v", err)
	}
	if e1 == e2 {
		t.Fatal("two encryptions of the same plaintext should produce different ciphertexts (random nonce)")
	}
}

func TestEncryptEmptyString(t *testing.T) {
	key := testKey(t)

	encrypted, err := Encrypt(key, "")
	if err != nil {
		t.Fatalf("Encrypt empty: %v", err)
	}

	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Decrypt empty: %v", err)
	}
	if decrypted != "" {
		t.Fatalf("expected empty string, got %q", decrypted)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := testKey(t)
	key2 := testKey(t)

	encrypted, err := Encrypt(key1, "secret")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = Decrypt(key2, encrypted)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptCorruptedPayload(t *testing.T) {
	key := testKey(t)

	corrupted := Prefix + AlgAES256GCM + ":" + base64.StdEncoding.EncodeToString([]byte("not-a-valid-gcm-payload"))
	_, err := Decrypt(key, corrupted)
	if err == nil {
		t.Fatal("expected error for corrupted payload")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	key := testKey(t)

	_, err := Decrypt(key, Prefix+AlgAES256GCM+":!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid Base64")
	}
}

func TestDecryptTruncatedPayload(t *testing.T) {
	key := testKey(t)

	encrypted, err := Encrypt(key, "value")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	// Chop off the last 10 characters to truncate nonce+ciphertext.
	truncated := encrypted[:len(encrypted)-10]
	_, err = Decrypt(key, truncated)
	if err == nil {
		t.Fatal("expected error for truncated payload")
	}
}

func TestDecryptMalformedFormat(t *testing.T) {
	key := testKey(t)

	tests := []struct {
		name    string
		value   string
		wantErr string
	}{
		{"no prefix", "plaintext-value", "does not use the enc: prefix"},
		{"enc only", "enc:", "malformed encrypted value"},
		{"enc algo no payload", "enc:aes-256-gcm", "malformed encrypted value"},
		{"enc algo colon only", "enc:aes-256-gcm:", "malformed encrypted value"},
		{"enc empty algo", "enc::payload", "malformed encrypted value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(key, tt.value)
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestDecryptUnsupportedAlgorithm(t *testing.T) {
	key := testKey(t)

	encoded := Prefix + "chacha20:" + base64.StdEncoding.EncodeToString([]byte("data"))
	_, err := Decrypt(key, encoded)
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
	if !strings.Contains(err.Error(), "unsupported encryption algorithm") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseKey(t *testing.T) {
	validKey := make([]byte, 32)
	validKey[0] = 42
	encoded := base64.StdEncoding.EncodeToString(validKey)

	key, err := ParseKey(encoded)
	if err != nil {
		t.Fatalf("ParseKey valid: %v", err)
	}
	if len(key) != 32 || key[0] != 42 {
		t.Fatalf("unexpected key: %v", key)
	}
}

func TestParseKeyInvalidBase64(t *testing.T) {
	_, err := ParseKey("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid Base64 key")
	}
}

func TestParseKeyWrongLength(t *testing.T) {
	short := base64.StdEncoding.EncodeToString([]byte("short"))
	_, err := ParseKey(short)
	if err == nil {
		t.Fatal("expected error for short key")
	}
	if !strings.Contains(err.Error(), "32 bytes") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsEncrypted(t *testing.T) {
	if !IsEncrypted(Prefix + "aes-256-gcm:abc") {
		t.Fatal("expected true for enc: prefixed value")
	}
	if IsEncrypted("plaintext") {
		t.Fatal("expected false for plaintext")
	}
	if IsEncrypted("enc") {
		t.Fatal("expected false for bare 'enc'")
	}
}

func TestParse(t *testing.T) {
	alg, payload, err := Parse(Prefix + "aes-256-gcm:AbCdEfGh")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if alg != "aes-256-gcm" {
		t.Fatalf("algorithm: got %q", alg)
	}
	if payload != "AbCdEfGh" {
		t.Fatalf("payload: got %q", payload)
	}
}

func TestParseNotEncrypted(t *testing.T) {
	_, _, err := Parse("plaintext")
	if err == nil {
		t.Fatal("expected error for non-enc value")
	}
}
