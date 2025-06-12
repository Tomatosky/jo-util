package cryptor

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"io"
	"strings"
)

type RSASecurity struct {
	pubStr string          //公钥字符串
	priStr string          //私钥字符串
	pubkey *rsa.PublicKey  //公钥
	prikey *rsa.PrivateKey //私钥
}

// SetPublicKey 设置公钥
func (r *RSASecurity) SetPublicKey(pubStr string) (err error) {
	pubStr = strings.TrimSpace(pubStr)
	if pubStr == "" {
		return errors.New("public key is empty")
	}

	// 检查是否已经是 PEM 格式
	if !strings.Contains(pubStr, "BEGIN PUBLIC KEY") {
		// 处理非 PEM 格式的公钥
		pubStr = strings.ReplaceAll(pubStr, "\r\n", "")
		pubStr = strings.ReplaceAll(pubStr, "\n", "")
		pubStr = strings.ReplaceAll(pubStr, "\r", "")
		pubStr = strings.ReplaceAll(pubStr, " ", "")
		var builder strings.Builder
		builder.WriteString("-----BEGIN PUBLIC KEY-----\n")
		pubStr = builder.String()
	}

	r.pubStr = pubStr
	r.pubkey, err = r.GetPublickey()
	return err
}

// SetPrivateKey 设置私钥
func (r *RSASecurity) SetPrivateKey(priStr string) (err error) {
	r.priStr = priStr
	r.prikey, err = r.GetPrivatekey()
	return err
}

// GetPrivatekey *rsa.PublicKey
func (r *RSASecurity) GetPrivatekey() (*rsa.PrivateKey, error) {
	return getPriKey([]byte(r.priStr))
}

// GetPublickey *rsa.PrivateKey
func (r *RSASecurity) GetPublickey() (*rsa.PublicKey, error) {
	return getPubKey([]byte(r.pubStr))
}

// PubKeyENCTYPT 公钥加密
func (r *RSASecurity) PubKeyENCTYPT(input []byte) ([]byte, error) {
	if r.pubkey == nil {
		return []byte(""), errors.New(`please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(r.pubkey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// PubKeyDECRYPT 公钥解密
func (r *RSASecurity) PubKeyDECRYPT(input []byte) ([]byte, error) {
	if r.pubkey == nil {
		return []byte(""), errors.New(`please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(r.pubkey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// PriKeyENCTYPT 私钥加密
func (r *RSASecurity) PriKeyENCTYPT(input []byte) ([]byte, error) {
	if r.prikey == nil {
		return []byte(""), errors.New(`please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(r.prikey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// PriKeyDECRYPT 私钥解密
func (r *RSASecurity) PriKeyDECRYPT(input []byte) ([]byte, error) {
	if r.prikey == nil {
		return []byte(""), errors.New(`please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(r.prikey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}

	return io.ReadAll(output)
}
