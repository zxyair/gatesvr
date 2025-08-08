package ecc_test

import (
	"encoding/pem"
	"gatesvr/core/hash"
	"gatesvr/crypto/ecc"
	"gatesvr/errors"
	"gatesvr/utils/xrand"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"

	"testing"
)

const (
	publicKey  = "./pem/key_pub.pem"
	privateKey = "./pem/key.pem"
)

var (
	encryptor *ecc.Encryptor
	signer    *ecc.Signer
)

func init() {
	encryptor = ecc.NewEncryptor(
		ecc.WithEncryptorPublicKey(publicKey),
		ecc.WithEncryptorPrivateKey(privateKey),
	)

	signer = ecc.NewSigner(
		ecc.WithSignerHash(hash.SHA256),
		ecc.WithSignerPublicKey(publicKey),
		ecc.WithSignerPrivateKey(privateKey),
	)

}

func Test_Encrypt_Decrypt(t *testing.T) {
	str := xrand.Letters(200000)
	bytes := []byte(str)

	plaintext, err := encryptor.Encrypt(bytes)
	if err != nil {
		t.Fatal(err)
	}

	data, err := encryptor.Decrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(str == string(data))
}

func Benchmark_Encrypt(b *testing.B) {
	text := []byte(xrand.Letters(20000))

	for i := 0; i < b.N; i++ {
		_, err := encryptor.Encrypt(text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Decrypt(b *testing.B) {
	text := []byte(xrand.Letters(20000))
	plaintext, _ := encryptor.Encrypt(text)

	for i := 0; i < b.N; i++ {
		_, err := encryptor.Decrypt(plaintext)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_Sign_Verify(t *testing.T) {
	str := xrand.Letters(300000)
	bytes := []byte(str)

	signature, err := signer.Sign(bytes)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := signer.Verify(bytes, signature)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ok)
}

func Benchmark_Sign(b *testing.B) {
	bytes := []byte(xrand.Letters(20000))

	for i := 0; i < b.N; i++ {
		_, err := signer.Sign(bytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Verify(b *testing.B) {
	bytes := []byte(xrand.Letters(20000))
	signature, _ := signer.Sign(bytes)

	for i := 0; i < b.N; i++ {
		_, err := signer.Verify(bytes, signature)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func LoadPrivateKeyFromPEM(path string) (*ecies.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	// 这里用 go-ethereum 的 FromECDSA/ToECDSA
	priv, err := crypto.ToECDSA(block.Bytes)
	if err != nil {
		return nil, err
	}
	return ecies.ImportECDSA(priv), nil
}

func LoadPublicKeyFromPEM(path string) (*ecies.PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	pub, err := crypto.UnmarshalPubkey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return ecies.ImportECDSAPublic(pub), nil
}
