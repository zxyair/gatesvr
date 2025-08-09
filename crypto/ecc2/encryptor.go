package ecc2

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"io"
)

// ECCEncryptor implements Encryptor interface
type ECCEncryptor struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	peerPubKey *ecdsa.PublicKey // 对方公钥
}

func NewECCEncryptor(peerPubKey *ecdsa.PublicKey) (*ECCEncryptor, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ECCEncryptor{
		privateKey: priv,
		publicKey:  &priv.PublicKey,
		peerPubKey: peerPubKey,
	}, nil
}

func (e *ECCEncryptor) Name() string {
	return "ecc"
}

// 通过ECDH协商出对称密钥
func (e *ECCEncryptor) deriveSharedKey() ([]byte, error) {
	if e.peerPubKey == nil {
		return nil, errors.New("peer public key is not set")
	}
	x, _ := e.peerPubKey.Curve.ScalarMult(e.peerPubKey.X, e.peerPubKey.Y, e.privateKey.D.Bytes())
	shared := sha256.Sum256(x.Bytes())
	return shared[:], nil
}

// Encrypt 用AES-GCM加密
func (e *ECCEncryptor) Encrypt(data []byte) ([]byte, error) {
	key, err := e.deriveSharedKey()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	// 返回nonce + ciphertext
	return append(nonce, ciphertext...), nil
}

// Decrypt 用AES-GCM解密
func (e *ECCEncryptor) Decrypt(data []byte) ([]byte, error) {
	key, err := e.deriveSharedKey()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// 公钥序列化（可用于传输）
func (e *ECCEncryptor) MarshalPublicKey() ([]byte, error) {
	return asn1.Marshal(*e.publicKey)
}

// 公钥反序列化
func UnmarshalEccPublicKey(data []byte) (*ecdsa.PublicKey, error) {
	var pub ecdsa.PublicKey
	_, err := asn1.Unmarshal(data, &pub)
	if err != nil {
		return nil, err
	}
	return &pub, nil
}
