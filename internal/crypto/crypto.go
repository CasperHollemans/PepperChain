package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/CasperHollemans/PepperChain/internal/logging"
	"math/big"
)

type CryptoService interface {
	Hash(data []byte) []byte
	Sign(privateKey *ecdsa.PrivateKey, data []byte) []byte
	GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error)
	VerifySignature(key []byte, hash []byte, bytes []byte) bool
}

type CryptoServiceImpl struct {
	logger logging.Logger
}

func NewCryptoService(logger logging.Logger) *CryptoServiceImpl {
	return &CryptoServiceImpl{
		logger: logger,
	}
}

func (c *CryptoServiceImpl) Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (c *CryptoServiceImpl) Sign(privateKey *ecdsa.PrivateKey, data []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, data)
	if err != nil {
		panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

func (c *CryptoServiceImpl) GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		c.logger.Error("Error generating key pair", err)
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func (c *CryptoServiceImpl) VerifySignature(key []byte, hash []byte, bytes []byte) bool {
	r := new(big.Int).SetBytes(bytes[:len(bytes)/2])
	s := new(big.Int).SetBytes(bytes[len(bytes)/2:])
	x := key[:len(key)/2]
	y := key[len(key)/2:]
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(x),
		Y:     new(big.Int).SetBytes(y),
	}
	return ecdsa.Verify(publicKey, hash, r, s)
}
