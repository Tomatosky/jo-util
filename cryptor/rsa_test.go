package cryptor

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"strings"
	"testing"
)

func rsaPEMFixture(t *testing.T) (string, string, *rsa.PrivateKey) {
	t.Helper()
	privateKey, publicKey := GenerateRsaKeyPair(1024)
	if privateKey == nil || publicKey == nil {
		t.Fatal("GenerateRsaKeyPair() returned nil")
	}
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	publicDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicDER,
	})
	return string(publicPEM), string(privatePEM), privateKey
}

func configuredRSASecurity(t *testing.T) *RSASecurity {
	t.Helper()
	publicPEM, privatePEM, _ := rsaPEMFixture(t)
	security := &RSASecurity{}
	if err := security.SetPublicKey("\n  " + publicPEM + "  \n"); err != nil {
		t.Fatalf("SetPublicKey(): %v", err)
	}
	if err := security.SetPrivateKey(privatePEM); err != nil {
		t.Fatalf("SetPrivateKey(): %v", err)
	}
	return security
}

func TestRSASetAndGetKeys(t *testing.T) {
	publicPEM, privatePEM, privateKey := rsaPEMFixture(t)
	security := &RSASecurity{}
	if err := security.SetPublicKey(publicPEM); err != nil {
		t.Fatalf("SetPublicKey(): %v", err)
	}
	if err := security.SetPrivateKey(privatePEM); err != nil {
		t.Fatalf("SetPrivateKey(): %v", err)
	}

	publicKey, err := security.GetPublickey()
	if err != nil {
		t.Fatalf("GetPublickey(): %v", err)
	}
	loadedPrivateKey, err := security.GetPrivatekey()
	if err != nil {
		t.Fatalf("GetPrivatekey(): %v", err)
	}
	if publicKey.N.Cmp(privateKey.N) != 0 || loadedPrivateKey.N.Cmp(privateKey.N) != 0 {
		t.Error("loaded keys do not match input key pair")
	}

	block, _ := pem.Decode([]byte(publicPEM))
	if block == nil {
		t.Fatal("decode public key fixture")
	}
	rawBase64 := base64.StdEncoding.EncodeToString(block.Bytes)
	wrappedRawBase64 := rawBase64[:len(rawBase64)/2] + "\n  " + rawBase64[len(rawBase64)/2:]
	rawSecurity := &RSASecurity{}
	if err := rawSecurity.SetPublicKey(wrappedRawBase64); err != nil {
		t.Fatalf("SetPublicKey(raw Base64 DER): %v", err)
	}
	rawPublicKey, err := rawSecurity.GetPublickey()
	if err != nil {
		t.Fatalf("GetPublickey() after raw Base64 DER: %v", err)
	}
	if rawPublicKey.N.Cmp(privateKey.N) != 0 || rawPublicKey.E != privateKey.E {
		t.Error("raw Base64 DER public key does not match input key")
	}

	pkcs1PublicDER := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	for name, encoded := range map[string]string{
		"PKCS1 PEM":            string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pkcs1PublicDER})),
		"raw PKCS1 Base64 DER": base64.StdEncoding.EncodeToString(pkcs1PublicDER),
	} {
		t.Run(name, func(t *testing.T) {
			security := &RSASecurity{}
			if err := security.SetPublicKey(encoded); err != nil {
				t.Fatalf("SetPublicKey(): %v", err)
			}
			if security.pubkey.N.Cmp(privateKey.N) != 0 || security.pubkey.E != privateKey.E {
				t.Error("loaded PKCS1 public key does not match input key")
			}
		})
	}

	pkcs8PrivateDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal PKCS8 private key: %v", err)
	}
	for name, encoded := range map[string]string{
		"PKCS8 PEM":            string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8PrivateDER})),
		"raw PKCS8 Base64 DER": base64.StdEncoding.EncodeToString(pkcs8PrivateDER),
		"raw PKCS1 Base64 DER": base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey)),
	} {
		t.Run(name, func(t *testing.T) {
			security := &RSASecurity{}
			if err := security.SetPrivateKey(encoded); err != nil {
				t.Fatalf("SetPrivateKey(): %v", err)
			}
			if security.prikey.N.Cmp(privateKey.N) != 0 {
				t.Error("loaded private key does not match input key")
			}
		})
	}

	ecPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate EC key: %v", err)
	}
	ecPublicDER, err := x509.MarshalPKIXPublicKey(&ecPrivateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal EC public key: %v", err)
	}
	ecPrivateDER, err := x509.MarshalPKCS8PrivateKey(ecPrivateKey)
	if err != nil {
		t.Fatalf("marshal EC private key: %v", err)
	}
	if err := (&RSASecurity{}).SetPublicKey(string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: ecPublicDER}))); err == nil {
		t.Error("SetPublicKey(EC key) returned nil error")
	}
	if err := (&RSASecurity{}).SetPrivateKey(string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: ecPrivateDER}))); err == nil {
		t.Error("SetPrivateKey(EC key) returned nil error")
	}

	if err := (&RSASecurity{}).SetPublicKey(" "); err == nil || err.Error() != "public key is empty" {
		t.Errorf("SetPublicKey(empty) error = %v", err)
	}
	if err := (&RSASecurity{}).SetPublicKey("not a public key"); err == nil {
		t.Error("SetPublicKey(malformed) returned nil error")
	}
	if err := (&RSASecurity{}).SetPrivateKey("not a private key"); err == nil {
		t.Error("SetPrivateKey(malformed) returned nil error")
	}
	if err := (&RSASecurity{}).SetPrivateKey(" "); err == nil || err.Error() != "private key is empty" {
		t.Errorf("SetPrivateKey(empty) error = %v", err)
	}
}

func TestRSASecurityMissingKeyErrors(t *testing.T) {
	security := &RSASecurity{}
	tests := []struct {
		name string
		call func() ([]byte, error)
		want string
	}{
		{name: "public encrypt", call: func() ([]byte, error) { return security.PubKeyENCTYPT([]byte("x")) }, want: "please set the public key in advance"},
		{name: "public decrypt", call: func() ([]byte, error) { return security.PubKeyDECRYPT([]byte("x")) }, want: "please set the public key in advance"},
		{name: "private encrypt", call: func() ([]byte, error) { return security.PriKeyENCTYPT([]byte("x")) }, want: "please set the private key in advance"},
		{name: "private decrypt", call: func() ([]byte, error) { return security.PriKeyDECRYPT([]byte("x")) }, want: "please set the private key in advance"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.call()
			if err == nil || err.Error() != tt.want || len(got) != 0 {
				t.Errorf("call = (%v, %v), want (empty, %q)", got, err, tt.want)
			}
		})
	}
}

func TestRSASecurityPublicEncryptPrivateDecrypt(t *testing.T) {
	security := configuredRSASecurity(t)
	data := []byte(strings.Repeat("chunked payload ", 40))
	encrypted, err := security.PubKeyENCTYPT(data)
	if err != nil {
		t.Fatalf("PubKeyENCTYPT(): %v", err)
	}
	second, err := security.PubKeyENCTYPT(data)
	if err != nil {
		t.Fatalf("second PubKeyENCTYPT(): %v", err)
	}
	if bytes.Equal(encrypted, data) || bytes.Equal(encrypted, second) {
		t.Error("public-key ciphertext is plaintext or is not randomized")
	}
	decrypted, err := security.PriKeyDECRYPT(encrypted)
	if err != nil {
		t.Fatalf("PriKeyDECRYPT(): %v", err)
	}
	if !bytes.Equal(decrypted, data) {
		t.Errorf("round trip length/content mismatch: got %d bytes, want %d", len(decrypted), len(data))
	}

	if _, err := security.PriKeyDECRYPT(make([]byte, security.prikey.Size())); err == nil {
		t.Error("PriKeyDECRYPT(deterministic invalid block) returned nil error")
	}
}

func TestRSASecurityPrivateEncryptPublicDecrypt(t *testing.T) {
	security := configuredRSASecurity(t)
	data := []byte(strings.Repeat("chunked payload ", 40))
	encrypted, err := security.PriKeyENCTYPT(data)
	if err != nil {
		t.Fatalf("PriKeyENCTYPT(): %v", err)
	}
	if bytes.Equal(encrypted, data) {
		t.Error("private-key ciphertext equals plaintext")
	}
	decrypted, err := security.PubKeyDECRYPT(encrypted)
	if err != nil {
		t.Fatalf("PubKeyDECRYPT(): %v", err)
	}
	if !bytes.Equal(decrypted, data) {
		t.Errorf("round trip length/content mismatch: got %d bytes, want %d", len(decrypted), len(data))
	}

}

func TestPubKeyDecryptRejectsMalformedEncodedMessages(t *testing.T) {
	_, _, privateKey := rsaPEMFixture(t)
	publicKey := &privateKey.PublicKey
	keySize := publicKey.Size()

	if _, err := pubKeyDecrypt(publicKey, leftPad(publicKey.N.Bytes(), keySize)); err != ErrDataToLarge {
		t.Errorf("pubKeyDecrypt(N) error = %v, want %v", err, ErrDataToLarge)
	}

	privateTransform := func(encoded []byte) []byte {
		representative := new(big.Int).SetBytes(encoded)
		representative.Exp(representative, privateKey.D, privateKey.N)
		return leftPad(representative.Bytes(), keySize)
	}
	malformed := map[string][]byte{
		"wrong block type": func() []byte {
			encoded := make([]byte, keySize)
			encoded[1] = 0
			return encoded
		}(),
		"missing separator": func() []byte {
			encoded := bytes.Repeat([]byte{0xff}, keySize)
			encoded[0], encoded[1] = 0, 1
			return encoded
		}(),
		"short padding": func() []byte {
			encoded := make([]byte, keySize)
			encoded[1] = 1
			for i := 2; i < 6; i++ {
				encoded[i] = 0xff
			}
			return encoded
		}(),
		"invalid padding byte": func() []byte {
			encoded := make([]byte, keySize)
			encoded[1] = 1
			for i := 2; i < 12; i++ {
				encoded[i] = 0xff
			}
			encoded[5] = 0xfe
			return encoded
		}(),
	}
	for name, encoded := range malformed {
		t.Run(name, func(t *testing.T) {
			if got, err := pubKeyDecrypt(publicKey, privateTransform(encoded)); err == nil || got != nil {
				t.Errorf("pubKeyDecrypt(malformed) = (%v, %v), want (nil, error)", got, err)
			}
		})
	}
}
