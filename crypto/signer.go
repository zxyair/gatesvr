package crypto

import (
	"github.com/dobyte/due/v2/log"
)

type Signer interface {
	// Name 名称
	Name() string
	// Sign 签名
	Sign(data []byte) ([]byte, error)
	// Verify 验签
	Verify(data []byte, signature []byte) (bool, error)
}

var signers = make(map[string]Signer)

// RegisterSigner 注册签名器
func RegisterSigner(signer Signer) {
	if signer == nil {
		log.Fatal("can't register a invalid signer")
	}

	name := signer.Name()

	if name == "" {
		log.Fatal("can't register a signer without name")
	}

	if _, ok := signers[name]; ok {
		log.Warnf("the old %s signer will be overwritten", name)
	}

	signers[name] = signer
}

// InvokeSigner 调用签名器
func InvokeSigner(name string) Signer {
	signer, ok := signers[name]
	if !ok {
		log.Fatalf("%s signer is not registered", name)
	}

	return signer
}
