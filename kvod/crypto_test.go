package kvod

import (
	"reflect"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	passphrase, _ := GenerateRandom(10)
	salt, _ := GenerateRandom(saltSize)
	crypto := InitCrypto(string(passphrase), salt)

	data, _ := GenerateRandom(100)
	encryptedData, _ := crypto.encrypt(data)
	decryptedData, _ := crypto.decrypt(encryptedData)
	if !reflect.DeepEqual(data, decryptedData) {
		t.Error("Encrypt/Decrypt failed")
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	passphrase := []byte("test")
	salt, _ := GenerateRandom(saltSize)
	for n := 0; n < b.N; n++ {
		generateKey(passphrase, salt)
	}
}

func BenchmarkEncrypt(b *testing.B) {
	salt, _ := GenerateRandom(saltSize)
	crypto := InitCrypto("test", salt)

	data, _ := GenerateRandom(100)
	for n := 0; n < b.N; n++ {
		crypto.encrypt(data)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	salt, _ := GenerateRandom(saltSize)
	crypto := InitCrypto("test", salt)

	data, _ := GenerateRandom(100)
	encryptedData, _ := crypto.encrypt(data)
	for n := 0; n < b.N; n++ {
		crypto.decrypt(encryptedData)
	}
}
