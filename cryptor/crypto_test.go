package cryptor

import (
	"bytes"
	"os"
	"testing"
)

func TestAesEcbEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456") // 16 bytes key

	// Test with PKCS7 padding
	encrypted := AesEcbEncrypt(data, key, Pkcs7Padding)
	decrypted := AesEcbDecrypt(encrypted, key, Pkcs7Padding)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES ECB PKCS7 padding failed: expected %s, got %s", data, decrypted)
	}

	// Test with Zero padding
	encrypted = AesEcbEncrypt(data, key, ZeroPadding)
	decrypted = AesEcbDecrypt(encrypted, key, ZeroPadding)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES ECB Zero padding failed: expected %s, got %s", data, decrypted)
	}
}

func TestAesEcbEncryptDecryptWithErr(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456") // 16 bytes key

	// Test with PKCS7 padding
	encrypted, err := AesEcbEncryptWithErr(data, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES ECB EncryptWithErr PKCS7 padding failed: %v", err)
	}
	decrypted, err := AesEcbDecryptWithErr(encrypted, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES ECB DecryptWithErr PKCS7 padding failed: %v", err)
	}
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES ECB WithErr PKCS7 padding failed: expected %s, got %s", data, decrypted)
	}
}

func TestAesCbcEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456") // 16 bytes key

	// Test with PKCS7 padding
	encrypted := AesCbcEncrypt(data, key, Pkcs7Padding)
	decrypted := AesCbcDecrypt(encrypted, key, Pkcs7Padding)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES CBC PKCS7 padding failed: expected %s, got %s", data, decrypted)
	}

	// Test with Zero padding
	encrypted = AesCbcEncrypt(data, key, ZeroPadding)
	decrypted = AesCbcDecrypt(encrypted, key, ZeroPadding)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES CBC Zero padding failed: expected %s, got %s", data, decrypted)
	}
}

func TestAesCbcEncryptDecryptWithErr(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456") // 16 bytes key

	// Test with PKCS7 padding
	encrypted, err := AesCbcEncryptWithErr(data, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES CBC EncryptWithErr PKCS7 padding failed: %v", err)
	}
	decrypted, err := AesCbcDecryptWithErr(encrypted, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES CBC DecryptWithErr PKCS7 padding failed: %v", err)
	}
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES CBC WithErr PKCS7 padding failed: expected %s, got %s", data, decrypted)
	}
}

func TestAesCtrCrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456") // 16 bytes key

	// Test CTR mode encryption and decryption (symmetric operation)
	encrypted, err := AesCtrCrypt(data, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES CTR Encrypt failed: %v", err)
	}
	decrypted, err := AesCtrCrypt(encrypted, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES CTR Decrypt failed: %v", err)
	}
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES CTR failed: expected %s, got %s", data, decrypted)
	}
}

func TestAesCfbEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456") // 16 bytes key

	// Test with PKCS7 padding
	encrypted, err := AesCfbEncrypt(data, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES CFB Encrypt PKCS7 padding failed: %v", err)
	}
	decrypted, err := AesCfbDecrypt(encrypted, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES CFB Decrypt PKCS7 padding failed: %v", err)
	}
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES CFB PKCS7 padding failed: expected %s, got %s", data, decrypted)
	}
}

func TestAesOfbEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("1234567890123456") // 16 bytes key

	// Test with PKCS7 padding
	encrypted, err := AesOfbEncrypt(data, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES OFB Encrypt PKCS7 padding failed: %v", err)
	}
	decrypted, err := AesOfbDecrypt(encrypted, key, Pkcs7Padding)
	if err != nil {
		t.Errorf("AES OFB Decrypt PKCS7 padding failed: %v", err)
	}
	if !bytes.Equal(data, decrypted) {
		t.Errorf("AES OFB PKCS7 padding failed: expected %s, got %s", data, decrypted)
	}
}

func TestDesEcbEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("12345678") // 8 bytes key

	encrypted := DesEcbEncrypt(data, key)
	decrypted := DesEcbDecrypt(encrypted, key)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("DES ECB failed: expected %s, got %s", data, decrypted)
	}
}

func TestDesCbcEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("12345678") // 8 bytes key

	encrypted := DesCbcEncrypt(data, key)
	decrypted := DesCbcDecrypt(encrypted, key)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("DES CBC failed: expected %s, got %s", data, decrypted)
	}
}

func TestDesCtrCrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("12345678") // 8 bytes key

	encrypted := DesCtrCrypt(data, key)
	decrypted := DesCtrCrypt(encrypted, key)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("DES CTR failed: expected %s, got %s", data, decrypted)
	}
}

func TestDesCfbEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("12345678") // 8 bytes key

	encrypted := DesCfbEncrypt(data, key)
	decrypted := DesCfbDecrypt(encrypted, key)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("DES CFB failed: expected %s, got %s", data, decrypted)
	}
}

func TestDesOfbEncryptDecrypt(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("12345678") // 8 bytes key

	encrypted := DesOfbEncrypt(data, key)
	decrypted := DesOfbDecrypt(encrypted, key)
	if !bytes.Equal(data, decrypted) {
		t.Errorf("DES OFB failed: expected %s, got %s", data, decrypted)
	}
}

func TestGenerateRsaKeyPair(t *testing.T) {
	privateKey, publicKey := GenerateRsaKeyPair(2048)
	if privateKey == nil || publicKey == nil {
		t.Error("GenerateRsaKeyPair failed: returned nil key")
	}
}

func TestRsaEncryptOAEP(t *testing.T) {
	privateKey, publicKey := GenerateRsaKeyPair(2048)
	data := []byte("Hello, World!")
	label := []byte("test")

	encrypted, err := RsaEncryptOAEP(data, label, *publicKey)
	if err != nil {
		t.Errorf("RsaEncryptOAEP failed: %v", err)
	}

	decrypted, err := RsaDecryptOAEP(encrypted, label, *privateKey)
	if err != nil {
		t.Errorf("RsaDecryptOAEP failed: %v", err)
	}

	if !bytes.Equal(data, decrypted) {
		t.Errorf("RSA OAEP failed: expected %s, got %s", data, decrypted)
	}
}

func TestGenerateRsaKey(t *testing.T) {
	// This test will create temporary files
	priKeyFile := "./test_private.pem"
	pubKeyFile := "./test_public.pem"

	err := GenerateRsaKey(2048, priKeyFile, pubKeyFile)
	if err != nil {
		t.Errorf("GenerateRsaKey failed: %v", err)
	}

	// Clean up
	os.Remove(priKeyFile)
	os.Remove(pubKeyFile)
}

func TestRsaEncryptDecrypt(t *testing.T) {
	// First generate RSA keys
	priKeyFile := "./test_private.pem"
	pubKeyFile := "./test_public.pem"

	err := GenerateRsaKey(2048, priKeyFile, pubKeyFile)
	if err != nil {
		t.Errorf("GenerateRsaKey failed: %v", err)
		return
	}

	defer func() {
		os.Remove(priKeyFile)
		os.Remove(pubKeyFile)
	}()

	data := []byte("Hello, World!")

	// Test RSA encryption and decryption
	encrypted := RsaEncrypt(data, pubKeyFile)
	decrypted := RsaDecrypt(encrypted, priKeyFile)

	if !bytes.Equal(data, decrypted) {
		t.Errorf("RSA Encrypt/Decrypt failed: expected %s, got %s", data, decrypted)
	}
}
