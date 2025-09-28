package cryptor

import (
	"bytes"
	"runtime/debug"
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

func addPadding(data []byte, blockSize int, paddingType PaddingType) []byte {
	switch paddingType {
	case Pkcs7Padding:
		return pkcs7Padding(data, blockSize)
	case ZeroPadding:
		return zeroPadding(data, blockSize)
	case NoPadding:
		if len(data)%blockSize != 0 {
			debug.PrintStack()
			panic("data length is not aligned to block size")
		}
		return data
	default:
		debug.PrintStack()
		panic("unknown padding type")
	}
}

func removePadding(data []byte, paddingType PaddingType) []byte {
	switch paddingType {
	case Pkcs7Padding:
		return pkcs7UnPadding(data)
	case ZeroPadding:
		return zeroUnPadding(data)
	case NoPadding:
		return data
	default:
		debug.PrintStack()
		panic("unknown padding type")
	}
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
