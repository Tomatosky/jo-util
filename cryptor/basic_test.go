package cryptor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBase64KnownVector(t *testing.T) {
	const plain = "你好, hello"
	const encoded = "5L2g5aW9LCBoZWxsbw=="
	if got := Base64StdEncode(plain); got != encoded {
		t.Errorf("Base64StdEncode() = %q, want %q", got, encoded)
	}
	if got := Base64StdDecode(encoded); got != plain {
		t.Errorf("Base64StdDecode() = %q, want %q", got, plain)
	}
	if got := Base64StdDecode("%%%"); got != "" {
		t.Errorf("Base64StdDecode(invalid) = %q, want empty string", got)
	}
}

func TestHashKnownVectors(t *testing.T) {
	const input = "hello"
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{name: "MD5 string hex", got: func() string { return Md5String(input) }, want: "5d41402abc4b2a76b9719d911017c592"},
		{name: "MD5 string base64", got: func() string { return Md5StringWithBase64(input) }, want: "XUFAKrxLKna5cZ2REBfFkg=="},
		{name: "MD5 bytes hex", got: func() string { return Md5Byte([]byte(input)) }, want: "5d41402abc4b2a76b9719d911017c592"},
		{name: "MD5 bytes base64", got: func() string { return Md5ByteWithBase64([]byte(input)) }, want: "XUFAKrxLKna5cZ2REBfFkg=="},
		{name: "SHA1 hex", got: func() string { return Sha1(input) }, want: "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{name: "SHA1 base64", got: func() string { return Sha1WithBase64(input) }, want: "qvTGHdzF6KLavt4PO0gs2a6pQ00="},
		{name: "SHA256 hex", got: func() string { return Sha256(input) }, want: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{name: "SHA256 base64", got: func() string { return Sha256WithBase64(input) }, want: "LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="},
		{name: "SHA512 hex", got: func() string { return Sha512(input) }, want: "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"},
		{name: "SHA512 base64", got: func() string { return Sha512WithBase64(input) }, want: "m3HSJL1i83hdltRq0+o9czGb+8KJDKra4t/3JRlnPKcjI8PZm6XBHXx6zG4UuMXaDEZjR1wuXDre9G9zvN7AQw=="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.got(); got != tt.want {
				t.Errorf("hash = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHMACKnownVectors(t *testing.T) {
	const message = "The quick brown fox jumps over the lazy dog"
	const key = "key"
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{name: "HMAC-MD5 hex", got: func() string { return HmacMd5(message, key) }, want: "80070713463e7749b90c2dc24911e275"},
		{name: "HMAC-MD5 base64", got: func() string { return HmacMd5WithBase64(message, key) }, want: "gAcHE0Y+d0m5DC3CSRHidQ=="},
		{name: "HMAC-SHA1 hex", got: func() string { return HmacSha1(message, key) }, want: "de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9"},
		{name: "HMAC-SHA1 base64", got: func() string { return HmacSha1WithBase64(message, key) }, want: "3nybhbi3iqa8ino29wqQcBydtNk="},
		{name: "HMAC-SHA256 hex", got: func() string { return HmacSha256(message, key) }, want: "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8"},
		{name: "HMAC-SHA256 base64", got: func() string { return HmacSha256WithBase64(message, key) }, want: "97yD9DBThCSxMpjmqm+xQ+9NWaFJRhdZl0edvC0aPNg="},
		{name: "HMAC-SHA512 hex", got: func() string { return HmacSha512(message, key) }, want: "b42af09057bac1e2d41708e48a902e09b5ff7f12ab428a4fe86653c73dd248fb82f948a549f7b791a5b41915ee4d1ec3935357e4e2317250d0372afa2ebeeb3a"},
		{name: "HMAC-SHA512 base64", got: func() string { return HmacSha512WithBase64(message, key) }, want: "tCrwkFe6weLUFwjkipAuCbX/fxKrQopP6GZTxz3SSPuC+UilSfe3kaW0GRXuTR7Dk1NX5OIxclDQNyr6Lr7rOg=="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.got(); got != tt.want {
				t.Errorf("HMAC = %q, want %q", got, tt.want)
			}
		})
	}

	if HmacSha256(message, key) == HmacSha256(message, "different-key") {
		t.Error("HMAC-SHA256 did not change when the key changed")
	}
}

func TestFileHashes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hello.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	tests := []struct {
		name string
		hash func(string) (string, error)
		want string
	}{
		{name: "MD5", hash: Md5File, want: "5d41402abc4b2a76b9719d911017c592"},
		{name: "SHA1", hash: Sha1File, want: "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{name: "SHA256", hash: Sha256File, want: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{name: "SHA512", hash: Sha512File, want: "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.hash(path)
			if err != nil {
				t.Fatalf("hash fixture: %v", err)
			}
			if got != tt.want {
				t.Errorf("file hash = %q, want %q", got, tt.want)
			}

			got, err = tt.hash(filepath.Join(t.TempDir(), "missing"))
			if err == nil || got != "" {
				t.Errorf("hash missing file = (%q, %v), want (empty, error)", got, err)
			}

			got, err = tt.hash(t.TempDir())
			if err != nil || got != "" {
				t.Errorf("hash directory = (%q, %v), want (empty, nil)", got, err)
			}
		})
	}
}
