package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type failingReader struct {
	err error
}

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func requirePanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("function did not panic")
		}
	}()
	f()
}

func decodeHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("DecodeString(%q): %v", value, err)
	}
	return decoded
}

func TestAESModesRoundTripAndCiphertextProperties(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456")
	tests := []struct {
		name       string
		encrypt    func([]byte) ([]byte, error)
		decrypt    func([]byte) ([]byte, error)
		randomized bool
	}{
		{
			name:    "ECB PKCS7",
			encrypt: func(data []byte) ([]byte, error) { return AesEcbEncryptWithErr(data, key, Pkcs7Padding) },
			decrypt: func(data []byte) ([]byte, error) { return AesEcbDecryptWithErr(data, key, Pkcs7Padding) },
		},
		{
			name:    "ECB zero",
			encrypt: func(data []byte) ([]byte, error) { return AesEcbEncryptWithErr(data, key, ZeroPadding) },
			decrypt: func(data []byte) ([]byte, error) { return AesEcbDecryptWithErr(data, key, ZeroPadding) },
		},
		{
			name:       "CBC PKCS7",
			encrypt:    func(data []byte) ([]byte, error) { return AesCbcEncryptWithErr(data, key, Pkcs7Padding) },
			decrypt:    func(data []byte) ([]byte, error) { return AesCbcDecryptWithErr(data, key, Pkcs7Padding) },
			randomized: true,
		},
		{
			name:       "CBC zero",
			encrypt:    func(data []byte) ([]byte, error) { return AesCbcEncryptWithErr(data, key, ZeroPadding) },
			decrypt:    func(data []byte) ([]byte, error) { return AesCbcDecryptWithErr(data, key, ZeroPadding) },
			randomized: true,
		},
		{
			name:       "CFB PKCS7",
			encrypt:    func(data []byte) ([]byte, error) { return AesCfbEncrypt(data, key, Pkcs7Padding) },
			decrypt:    func(data []byte) ([]byte, error) { return AesCfbDecrypt(data, key, Pkcs7Padding) },
			randomized: true,
		},
		{
			name:       "OFB PKCS7",
			encrypt:    func(data []byte) ([]byte, error) { return AesOfbEncrypt(data, key, Pkcs7Padding) },
			decrypt:    func(data []byte) ([]byte, error) { return AesOfbDecrypt(data, key, Pkcs7Padding) },
			randomized: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := tt.encrypt(data)
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			if bytes.Equal(encrypted, data) {
				t.Error("ciphertext equals plaintext")
			}
			second, err := tt.encrypt(data)
			if err != nil {
				t.Fatalf("second encrypt: %v", err)
			}
			if tt.randomized && bytes.Equal(encrypted, second) {
				t.Error("random-IV mode returned identical ciphertext twice")
			}
			if !tt.randomized && !bytes.Equal(encrypted, second) {
				t.Error("deterministic mode returned different ciphertext")
			}

			decrypted, err := tt.decrypt(bytes.Clone(encrypted))
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			if !bytes.Equal(decrypted, data) {
				t.Errorf("round trip = %q, want %q", decrypted, data)
			}
		})
	}
}

func TestCFBAndOFBNISTVectors(t *testing.T) {
	key := decodeHex(t, "2b7e151628aed2a6abf7158809cf4f3c")
	iv := decodeHex(t, "000102030405060708090a0b0c0d0e0f")
	plaintext := decodeHex(t, "6bc1bee22e409f96e93d7e117393172aae2d8a571e03ac9c9eb76fac45af8e5130c81c46a35ce411e5fbc1191a0a52eff69f2445df4f9b17ad2b417be66c3710")
	tests := []struct {
		name       string
		ciphertext string
		crypt      func(block cipher.Block, dst, src, iv []byte, decrypt bool)
	}{
		{
			name:       "CFB",
			ciphertext: "3b3fd92eb72dad20333449f8e83cfb4ac8a64537a0b3a93fcde3cdad9f1ce58b26751f67a3cbb140b1808cf187a4f4dfc04b05357c5d1c0eeac4c66f9ff7f2e6",
			crypt:      cfbXOR,
		},
		{
			name:       "OFB",
			ciphertext: "3b3fd92eb72dad20333449f8e83cfb4a7789508d16918f03f53c52dac54ed8259740051e9c5fecf64344f7a82260edcc304c6528f659c77866a510d9c1d6ae5e",
			crypt: func(block cipher.Block, dst, src, iv []byte, _ bool) {
				ofbXOR(block, dst, src, iv)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := decodeHex(t, tt.ciphertext)
			for _, size := range []int{len(plaintext), len(plaintext) - 5} {
				t.Run(fmt.Sprintf("length_%d", size), func(t *testing.T) {
					block, err := aes.NewCipher(key)
					if err != nil {
						t.Fatal(err)
					}
					encrypted := make([]byte, size)
					tt.crypt(block, encrypted, plaintext[:size], iv, false)
					if !bytes.Equal(encrypted, expected[:size]) {
						t.Errorf("encrypted = %x, want %x", encrypted, expected[:size])
					}

					decrypted := bytes.Clone(encrypted)
					tt.crypt(block, decrypted, decrypted, iv, true)
					if !bytes.Equal(decrypted, plaintext[:size]) {
						t.Errorf("decrypted = %x, want %x", decrypted, plaintext[:size])
					}
				})
			}
		})
	}
}

func TestAESCTRRoundTrip(t *testing.T) {
	data := []byte("unaligned stream payload")
	key := []byte("1234567890123456")
	encrypted, err := AesCtrCrypt(data, key)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(encrypted, data) || len(encrypted) != len(data) {
		t.Errorf("CTR ciphertext = %v, plaintext = %v", encrypted, data)
	}
	second, err := AesCtrCrypt(data, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encrypted, second) {
		t.Error("CTR output unexpectedly depends on ignored padding argument")
	}
	decrypted, err := AesCtrCrypt(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decrypted, data) {
		t.Errorf("CTR round trip = %q, want %q", decrypted, data)
	}
}

func TestAesCfbEncryptRandomSourceError(t *testing.T) {
	sentinel := errors.New("random source failed")
	originalReader := rand.Reader
	rand.Reader = failingReader{err: sentinel}
	t.Cleanup(func() { rand.Reader = originalReader })

	key := []byte("1234567890123456")
	if _, err := AesCfbEncrypt([]byte("data"), key, Pkcs7Padding); !errors.Is(err, sentinel) {
		t.Errorf("AesCfbEncrypt() error = %v, want %v", err, sentinel)
	}
}

func TestAESErrorPaths(t *testing.T) {
	data := []byte("not block aligned")
	validKey := []byte("1234567890123456")
	invalidKey := []byte("short")
	for name, encrypt := range map[string]func() error{
		"ECB invalid key":          func() error { _, err := AesEcbEncryptWithErr(data, invalidKey, Pkcs7Padding); return err },
		"CBC invalid key":          func() error { _, err := AesCbcEncryptWithErr(data, invalidKey, Pkcs7Padding); return err },
		"CTR invalid key":          func() error { _, err := AesCtrCrypt(data, invalidKey); return err },
		"CFB invalid key":          func() error { _, err := AesCfbEncrypt(data, invalidKey, Pkcs7Padding); return err },
		"OFB invalid key":          func() error { _, err := AesOfbEncrypt(data, invalidKey, Pkcs7Padding); return err },
		"ECB unaligned no padding": func() error { _, err := AesEcbEncryptWithErr(data, validKey, NoPadding); return err },
		"CBC unaligned no padding": func() error { _, err := AesCbcEncryptWithErr(data, validKey, NoPadding); return err },
		"CFB unaligned no padding": func() error { _, err := AesCfbEncrypt(data, validKey, NoPadding); return err },
		"OFB unaligned no padding": func() error { _, err := AesOfbEncrypt(data, validKey, NoPadding); return err },
		"ECB unknown padding":      func() error { _, err := AesEcbEncryptWithErr(data, validKey, PaddingType(99)); return err },
	} {
		t.Run(name, func(t *testing.T) {
			if err := encrypt(); err == nil {
				t.Error("operation returned nil error")
			}
		})
	}

	if got, err := AesCfbDecrypt([]byte("short"), validKey, Pkcs7Padding); err == nil || got != nil {
		t.Errorf("AesCfbDecrypt(short) = (%v, %v), want (nil, error)", got, err)
	}
	if got, err := AesOfbDecrypt([]byte("short"), validKey, Pkcs7Padding); err == nil || got != nil {
		t.Errorf("AesOfbDecrypt(short) = (%v, %v), want (nil, error)", got, err)
	}

	requirePanic(t, func() { AesEcbEncrypt(data, invalidKey, Pkcs7Padding) })
	requirePanic(t, func() { AesCbcEncrypt(data, invalidKey, Pkcs7Padding) })
}

func TestAesCBCDecryptRejectsInvalidCiphertextLength(t *testing.T) {
	key := []byte("1234567890123456")
	tests := []struct {
		name       string
		ciphertext []byte
		padding    PaddingType
		wantError  string
	}{
		{name: "missing IV", ciphertext: nil, padding: Pkcs7Padding, wantError: "invalid ciphertext size"},
		{name: "short IV", ciphertext: make([]byte, 15), padding: Pkcs7Padding, wantError: "invalid ciphertext size"},
		{name: "unaligned encrypted body", ciphertext: make([]byte, 17), padding: Pkcs7Padding, wantError: "ciphertext is not aligned to block size"},
		{name: "PKCS7 body missing", ciphertext: make([]byte, 16), padding: Pkcs7Padding, wantError: "invalid ciphertext size"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AesCbcDecryptWithErr(tt.ciphertext, key, tt.padding)
			if err == nil || err.Error() != tt.wantError || got != nil {
				t.Errorf("AesCbcDecryptWithErr() = (%v, %v), want (nil, %q)", got, err, tt.wantError)
			}
		})
	}

	emptyCiphertext, err := AesCbcEncryptWithErr(nil, key, NoPadding)
	if err != nil {
		t.Fatalf("AesCbcEncryptWithErr(empty, NoPadding): %v", err)
	}
	decrypted, err := AesCbcDecryptWithErr(bytes.Clone(emptyCiphertext), key, NoPadding)
	if err != nil || len(decrypted) != 0 {
		t.Errorf("empty NoPadding round trip = (%v, %v), want (empty, nil)", decrypted, err)
	}

	requirePanic(t, func() {
		AesCbcDecrypt(make([]byte, 17), key, Pkcs7Padding)
	})
}

func TestAESDecryptRejectsMalformedPKCS7WithoutPanic(t *testing.T) {
	key := []byte("1234567890123456")
	invalidPlaintexts := map[string][]byte{
		"zero padding length":   append(make([]byte, aes.BlockSize-1), 0),
		"padding exceeds block": append(make([]byte, aes.BlockSize-1), aes.BlockSize+1),
		"inconsistent padding":  append(make([]byte, aes.BlockSize-2), 1, 2),
	}
	for name, plaintext := range invalidPlaintexts {
		t.Run(name, func(t *testing.T) {
			block, err := aes.NewCipher(key)
			if err != nil {
				t.Fatal(err)
			}
			iv := make([]byte, aes.BlockSize)
			ciphertext := make([]byte, aes.BlockSize+len(plaintext))
			copy(ciphertext, iv)
			cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

			if got, err := AesCbcDecryptWithErr(ciphertext, key, Pkcs7Padding); err == nil || got != nil {
				t.Errorf("AesCbcDecryptWithErr(malformed padding) = (%v, %v), want (nil, error)", got, err)
			}
		})
	}

	if got, err := AesEcbDecryptWithErr([]byte("short"), key, Pkcs7Padding); err == nil || got != nil {
		t.Errorf("AesEcbDecryptWithErr(unaligned) = (%v, %v), want (nil, error)", got, err)
	}
}

func TestDESModesRoundTripAndCiphertextProperties(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("12345678")
	tests := []struct {
		name       string
		encrypt    func([]byte) []byte
		decrypt    func([]byte) []byte
		randomized bool
	}{
		{name: "ECB", encrypt: func(data []byte) []byte { return DesEcbEncrypt(data, key) }, decrypt: func(data []byte) []byte { return DesEcbDecrypt(data, key) }},
		{name: "CBC", encrypt: func(data []byte) []byte { return DesCbcEncrypt(data, key) }, decrypt: func(data []byte) []byte { return DesCbcDecrypt(data, key) }, randomized: true},
		{name: "CTR", encrypt: func(data []byte) []byte { return DesCtrCrypt(data, key) }, decrypt: func(data []byte) []byte { return DesCtrCrypt(data, key) }},
		{name: "CFB", encrypt: func(data []byte) []byte { return DesCfbEncrypt(data, key) }, decrypt: func(data []byte) []byte { return DesCfbDecrypt(data, key) }, randomized: true},
		{name: "OFB", encrypt: func(data []byte) []byte { return DesOfbEncrypt(data, key) }, decrypt: func(data []byte) []byte { return DesOfbDecrypt(data, key) }, randomized: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted := tt.encrypt(data)
			if bytes.Equal(encrypted, data) {
				t.Error("ciphertext equals plaintext")
			}
			second := tt.encrypt(data)
			if tt.randomized && bytes.Equal(encrypted, second) {
				t.Error("random-IV mode returned identical ciphertext twice")
			}
			if !tt.randomized && !bytes.Equal(encrypted, second) {
				t.Error("deterministic mode returned different ciphertext")
			}
			if decrypted := tt.decrypt(bytes.Clone(encrypted)); !bytes.Equal(decrypted, data) {
				t.Errorf("round trip = %q, want %q", decrypted, data)
			}
		})
	}
}

func TestDESInvalidKeys(t *testing.T) {
	data := []byte("payload")
	badKey := []byte("short")
	for name, f := range map[string]func(){
		"CBC encrypt": func() { DesCbcEncrypt(data, badKey) },
		"CBC decrypt": func() { DesCbcDecrypt(make([]byte, 16), badKey) },
		"CTR":         func() { DesCtrCrypt(data, badKey) },
		"CFB encrypt": func() { DesCfbEncrypt(data, badKey) },
		"CFB decrypt": func() { DesCfbDecrypt(make([]byte, 8), badKey) },
		"OFB encrypt": func() { DesOfbEncrypt(data, badKey) },
		"OFB decrypt": func() { DesOfbDecrypt(make([]byte, 16), badKey) },
	} {
		t.Run(name, func(t *testing.T) { requirePanic(t, f) })
	}
}

func TestGenerateRSAKeyPairAndOAEP(t *testing.T) {
	privateKey, publicKey := GenerateRsaKeyPair(1024)
	if privateKey == nil || publicKey == nil {
		t.Fatal("GenerateRsaKeyPair() returned nil key")
	}
	if err := privateKey.Validate(); err != nil {
		t.Errorf("private key validation: %v", err)
	}
	if privateKey.N.BitLen() != 1024 || publicKey.N.Cmp(privateKey.N) != 0 || publicKey.E != privateKey.E {
		t.Error("generated public/private key pair does not match requested size or modulus")
	}

	data := []byte("Hello, World!")
	label := []byte("test-label")
	encrypted, err := RsaEncryptOAEP(data, label, *publicKey)
	if err != nil {
		t.Fatalf("RsaEncryptOAEP(): %v", err)
	}
	second, err := RsaEncryptOAEP(data, label, *publicKey)
	if err != nil {
		t.Fatalf("second RsaEncryptOAEP(): %v", err)
	}
	if bytes.Equal(encrypted, data) || bytes.Equal(encrypted, second) {
		t.Error("RSA-OAEP ciphertext is plaintext or is not randomized")
	}
	decrypted, err := RsaDecryptOAEP(encrypted, label, *privateKey)
	if err != nil {
		t.Fatalf("RsaDecryptOAEP(): %v", err)
	}
	if !bytes.Equal(decrypted, data) {
		t.Errorf("RSA-OAEP round trip = %q, want %q", decrypted, data)
	}
	if _, err := RsaDecryptOAEP(encrypted, []byte("wrong-label"), *privateKey); err == nil {
		t.Error("RsaDecryptOAEP(wrong label) returned nil error")
	}
	tampered := bytes.Clone(encrypted)
	tampered[len(tampered)-1] ^= 1
	if _, err := RsaDecryptOAEP(tampered, label, *privateKey); err == nil {
		t.Error("RsaDecryptOAEP(tampered ciphertext) returned nil error")
	}
}

func TestRSAKeyFilesRoundTrip(t *testing.T) {
	root := t.TempDir()
	privatePath := filepath.Join(root, "private.pem")
	publicPath := filepath.Join(root, "public.pem")
	if runtime.GOOS != "windows" {
		if err := os.WriteFile(privatePath, []byte("old key"), 0o666); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(privatePath, 0o666); err != nil {
			t.Fatal(err)
		}
	}
	if err := GenerateRsaKey(1024, privatePath, publicPath); err != nil {
		t.Fatalf("GenerateRsaKey(): %v", err)
	}
	for _, path := range []string{privatePath, publicPath} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat(%q): %v", path, err)
		}
		if info.Size() == 0 {
			t.Errorf("generated key file %q is empty", path)
		}
	}
	if runtime.GOOS != "windows" {
		if info, err := os.Stat(privatePath); err != nil {
			t.Fatal(err)
		} else if info.Mode().Perm()&0o077 != 0 {
			t.Errorf("private key permissions = %04o, want no group/other permissions", info.Mode().Perm())
		}
	}

	data := []byte("Hello, World!")
	encrypted := RsaEncrypt(data, publicPath)
	second := RsaEncrypt(data, publicPath)
	if bytes.Equal(encrypted, data) || bytes.Equal(encrypted, second) {
		t.Error("RSA PKCS#1 v1.5 ciphertext is plaintext or is not randomized")
	}
	if decrypted := RsaDecrypt(encrypted, privatePath); !bytes.Equal(decrypted, data) {
		t.Errorf("RSA file round trip = %q, want %q", decrypted, data)
	}

	publicPEM, err := os.ReadFile(publicPath)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := getPubKey(publicPEM)
	if err != nil {
		t.Fatal(err)
	}
	privatePEM, err := os.ReadFile(privatePath)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := getPriKey(privatePEM)
	if err != nil {
		t.Fatal(err)
	}

	standardPlaintext, err := privateKey.Decrypt(rand.Reader, encrypted, nil)
	if err != nil || !bytes.Equal(standardPlaintext, data) {
		t.Errorf("standard PKCS#1 v1.5 decrypt = (%q, %v), want (%q, nil)", standardPlaintext, err, data)
	}
	//noinspection GoDeprecation // Verify interoperability with the legacy standard-library format.
	standardCiphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		t.Fatal(err)
	}
	if decrypted := RsaDecrypt(standardCiphertext, privatePath); !bytes.Equal(decrypted, data) {
		t.Errorf("RsaDecrypt(standard PKCS#1 v1.5 ciphertext) = %q, want %q", decrypted, data)
	}

	maxMessageSize := publicKey.Size() - 11
	maxMessage := bytes.Repeat([]byte{1}, maxMessageSize)
	if decrypted := RsaDecrypt(RsaEncrypt(maxMessage, publicPath), privatePath); !bytes.Equal(decrypted, maxMessage) {
		t.Errorf("maximum-length PKCS#1 v1.5 round trip returned %d bytes, want %d", len(decrypted), len(maxMessage))
	}
	requirePanic(t, func() {
		RsaEncrypt(bytes.Repeat([]byte{1}, maxMessageSize+1), publicPath)
	})
}

func TestGenerateRSAKeyReturnsFileErrors(t *testing.T) {
	t.Run("private key path", func(t *testing.T) {
		root := t.TempDir()
		err := GenerateRsaKey(
			1024,
			filepath.Join(root, "missing", "private.pem"),
			filepath.Join(root, "public.pem"),
		)
		if err == nil {
			t.Fatal("GenerateRsaKey() returned nil for an invalid private key path")
		}
	})

	t.Run("public key path", func(t *testing.T) {
		root := t.TempDir()
		err := GenerateRsaKey(
			1024,
			filepath.Join(root, "private.pem"),
			filepath.Join(root, "missing", "public.pem"),
		)
		if err == nil {
			t.Fatal("GenerateRsaKey() returned nil for an invalid public key path")
		}
	})
}
