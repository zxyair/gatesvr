/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2022/11/1 12:50 上午
 * @Desc: TODO
 */

package ecc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gatesvr/errors"
	"gatesvr/utils/xconv"
	"gatesvr/utils/xpath"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io"
	"os"
	"path"
)

type Key struct {
	prv *ecdsa.PrivateKey
}

// GenerateKey 生成秘钥
func GenerateKey(curve Curve) (*Key, error) {
	prv, err := ecdsa.GenerateKey(curve.New(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Key{prv: prv}, nil
}

// PublicKey 获取公钥
func (k *Key) PublicKey() *ecdsa.PublicKey {
	return &k.prv.PublicKey
}

// PrivateKey 获取私钥
func (k *Key) PrivateKey() *ecdsa.PrivateKey {
	return k.prv
}

// MarshalPublicKey 编码公钥
func (k *Key) MarshalPublicKey() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := k.marshalPublicKey(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 编码私钥
func (k *Key) marshalPublicKey(out io.Writer) error {
	derText, err := x509.MarshalPKIXPublicKey(k.PublicKey())
	if err != nil {
		return err
	}

	return pem.Encode(out, &pem.Block{
		Type:  "ECDSA PUBLIC KEY",
		Bytes: derText,
	})
}

// MarshalPrivateKey 编码私钥
func (k *Key) MarshalPrivateKey() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := k.marshalPrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 编码私钥
func (k *Key) marshalPrivateKey(out io.Writer) error {
	derText, err := x509.MarshalECPrivateKey(k.PrivateKey())
	if err != nil {
		return err
	}

	return pem.Encode(out, &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: derText,
	})
}

// SaveKeyPair 保存秘钥对
func (k *Key) SaveKeyPair(dir string, file string) (err error) {
	if !xpath.IsDir(dir) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	}

	err = k.savePublicKey(dir, file)
	if err != nil {
		return
	}

	return k.savePrivateKey(dir, file)
}

// 保存公钥
func (k *Key) savePrivateKey(dir string, file string) (err error) {
	filepath := path.Join(dir, file)
	defer func() {
		if err != nil {
			_ = os.Remove(filepath)
		}
	}()

	f, err := os.Create(filepath)
	if err != nil {
		return
	}

	return k.marshalPrivateKey(f)
}

// 保存公钥
func (k *Key) savePublicKey(dir string, file string) (err error) {
	base, _, name, ext := xpath.Split(file)
	if ext != "" {
		file = name + ".pub." + ext
	} else {
		file = name + ".pub"
	}

	filepath := path.Join(dir, base, file)
	defer func() {
		if err != nil {
			_ = os.Remove(filepath)
		}
	}()

	f, err := os.Create(filepath)
	if err != nil {
		return
	}

	return k.marshalPublicKey(f)
}

func loadKey(key string) (*pem.Block, error) {
	var (
		err    error
		buffer []byte
	)

	if xpath.IsFile(key) {
		buffer, err = os.ReadFile(key)
		if err != nil {
			return nil, err
		}
	} else {
		buffer = xconv.StringToBytes(key)
	}

	block, _ := pem.Decode(buffer)

	return block, nil
}

func parseECIESPublicKey(publicKey string) (*ecies.PublicKey, error) {
	pub, err := parseECDSAPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return ecies.ImportECDSAPublic(pub), nil
}

func parseECIESPrivateKey(privateKey string) (*ecies.PrivateKey, error) {
	prv, err := parseECDSAPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return ecies.ImportECDSA(prv), nil
}

func parseECDSAPublicKey(publicKey string) (*ecdsa.PublicKey, error) {
	// 我假设 loadKey 返回 *pem.Block 和 error
	// 另外，我将您的变量名 black 改为 block，这更符合习惯
	block, err := loadKey(publicKey)
	if err != nil {
		return nil, err
	}

	if block == nil {
		return nil, errors.New("invalid public key: PEM block not found")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch key := pub.(type) {
	case *ecdsa.PublicKey:
		// 核心修复：在这里手动设置曲线！
		key.Curve = elliptic.P256()
		fmt.Printf("[DEBUG] parseECDSAPublicKey: Curve has been set. key.Curve is: %v\n", key.Curve)

		return key, nil
	default:
		return nil, errors.New("key is not a valid ECDSA public key")
	}
}

func parseECDSAPrivateKey(privateKey string) (*ecdsa.PrivateKey, error) {
	block, err := loadKey(privateKey)
	if err != nil {
		return nil, err
	}

	if block == nil {
		return nil, errors.New("invalid private key: PEM block not found")
	}

	// 1. 先解析出私钥
	privateKeyObj, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	
	// 2. 核心修复：为解析出的私钥手动设置曲线！
	privateKeyObj.Curve = elliptic.P256()

	// 3. 返回修正后的私钥对象
	return privateKeyObj, nil
}
