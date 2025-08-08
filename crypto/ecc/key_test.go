package ecc

import (
	"encoding/pem"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"os"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	// 生成 secp256k1 ECDSA 密钥对
	ecdsaKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	// 转为 ECIES 密钥
	eciesKey := ecies.ImportECDSA(ecdsaKey)

	// Marshal 公钥为字节
	pubBytes := crypto.FromECDSAPub(&ecdsaKey.PublicKey)

	// 打印公钥（16进制）
	fmt.Printf("Public Key (hex): %x\n", pubBytes)
	fmt.Printf("eciesKey: %v\n", eciesKey)
}
func TestKey_SaveKeyPair(t *testing.T) {
	// 生成 secp256k1 ECDSA 密钥对
	ecdsaKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	// 保存私钥到 PEM 文件
	privBytes := crypto.FromECDSA(ecdsaKey)
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	})
	err = os.WriteFile("./pem/key.pem", privPem, 0600)
	if err != nil {
		t.Fatal(err)
	}

	// 保存公钥到 PEM 文件
	pubBytes := crypto.FromECDSAPub(&ecdsaKey.PublicKey)
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: pubBytes,
	})
	err = os.WriteFile("./pem/key_pub.pem", pubPem, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
