package cryptor

import (
	"bytes"
	"crypto/cipher"
	"errors"
)

func generateAesKey(key []byte, size int) []byte {
	genKey := make([]byte, size)
	copy(genKey, key)
	for i := size; i < len(key); {
		for j := 0; j < size && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func generateDesKey(key []byte) []byte {
	genKey := make([]byte, 8)
	copy(genKey, key)
	for i := 8; i < len(key); {
		for j := 0; j < 8 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// cfbXOR preserves the package's legacy CFB byte format. CFB does not authenticate ciphertext.
func cfbXOR(block cipher.Block, dst, src, iv []byte, decrypt bool) {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		panic("cryptor: CFB IV length must equal block size")
	}
	if len(dst) < len(src) {
		panic("cryptor: CFB output smaller than input")
	}

	feedback := bytes.Clone(iv)
	keyStream := make([]byte, blockSize)
	for len(src) > 0 {
		block.Encrypt(keyStream, feedback)
		n := min(len(src), blockSize)
		if decrypt {
			copy(feedback, src[:n])
		}
		for i := 0; i < n; i++ {
			dst[i] = src[i] ^ keyStream[i]
		}
		if !decrypt {
			copy(feedback, dst[:n])
		}
		dst = dst[n:]
		src = src[n:]
	}
}

// ofbXOR preserves the package's legacy OFB byte format. OFB does not authenticate ciphertext.
func ofbXOR(block cipher.Block, dst, src, iv []byte) {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		panic("cryptor: OFB IV length must equal block size")
	}
	if len(dst) < len(src) {
		panic("cryptor: OFB output smaller than input")
	}

	keyStream := bytes.Clone(iv)
	for len(src) > 0 {
		block.Encrypt(keyStream, keyStream)
		n := min(len(src), blockSize)
		for i := 0; i < n; i++ {
			dst[i] = src[i] ^ keyStream[i]
		}
		dst = dst[n:]
		src = src[n:]
	}
}

func addPadding(data []byte, blockSize int, paddingType PaddingType) ([]byte, error) {
	switch paddingType {
	case Pkcs7Padding:
		return pkcs7Padding(data, blockSize), nil
	case ZeroPadding:
		return zeroPadding(data, blockSize), nil
	case NoPadding:
		if len(data)%blockSize != 0 {
			return nil, errors.New("data length is not aligned to block size")
		}
		return data, nil
	default:
		return nil, errors.New("unknown padding type")
	}
}

func removePadding(data []byte, blockSize int, paddingType PaddingType) ([]byte, error) {
	switch paddingType {
	case Pkcs7Padding:
		return pkcs7UnPaddingWithErr(data, blockSize)
	case ZeroPadding:
		return zeroUnPadding(data), nil
	case NoPadding:
		return data, nil
	default:
		return nil, errors.New("unknown padding type")
	}
}

func pkcs7UnPaddingWithErr(src []byte, blockSize int) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("invalid PKCS7 padding")
	}
	padding := int(src[len(src)-1])
	if padding == 0 || padding > blockSize || padding > len(src) {
		return nil, errors.New("invalid PKCS7 padding")
	}
	for _, value := range src[len(src)-padding:] {
		if int(value) != padding {
			return nil, errors.New("invalid PKCS7 padding")
		}
	}
	return src[:len(src)-padding], nil
}

func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}

func zeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{0}, padding)
	return append(data, padText...)
}

func zeroUnPadding(data []byte) []byte {
	length := len(data)
	for i := length - 1; i >= 0; i-- {
		if data[i] != 0 {
			return data[:i+1]
		}
	}
	return data[:0]
}
