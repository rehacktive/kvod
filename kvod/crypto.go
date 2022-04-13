package kvod

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize    = 32
	saltSize   = 12
	nonceSize  = 12
	iterations = 10000
)

type crypto struct {
	key []byte
}

// InitCrypto inizialize the crypto struct using password and salt to generate the main key
func InitCrypto(password string, salt []byte) *crypto {
	key := generateKey([]byte(password), salt)
	return &crypto{key}
}

func (c *crypto) encrypt(data []byte) ([]byte, error) {
	nonce, err := GenerateRandom(nonceSize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	encrypted := aesgcm.Seal(nonce, nonce, data, nil)

	return encrypted, nil
}

func (c *crypto) decrypt(data []byte) ([]byte, error) {
	nonce := data[:nonceSize]

	cp, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cp)
	if err != nil {
		return nil, err
	}

	out, err := gcm.Open(nil, nonce, data[nonceSize:], nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func generateKey(passphrase []byte, salt []byte) []byte {
	dk := pbkdf2.Key(passphrase, salt, iterations, keySize, sha1.New)
	return dk
}
