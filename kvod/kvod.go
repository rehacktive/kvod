package kvod

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"

	"os"
	"path/filepath"
)

const (
	saltFilename = ".salt"
)

// KVod basic struct
type KVod struct {
	path      string
	encrypter crypto
}

// Init KVod struct
func Init(path string, password string) *KVod {
	err := os.MkdirAll(path, 0770)
	if err != nil {
		log.Fatal(err)
	}
	saltFile := filepath.Join(path, saltFilename)

	salt, err := GenerateRandom(saltSize)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(saltFile); os.IsNotExist(err) {
		err = write(saltFile, salt)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		salt, err = read(saltFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	crypto := InitCrypto(password, salt)
	return &KVod{path, *crypto}
}

// Put a struct using a string as key
func (m *KVod) Put(key string, value interface{}) error {
	filename := m.getFilename(key)

	serializedValue, err := serialize(value)
	if err != nil {
		return err
	}
	store := make(map[string][]byte)
	store[key] = serializedValue
	data, err := serialize(store)
	if err != nil {
		return err
	}

	encryptedData, err := m.encrypter.encrypt(data)
	if err != nil {
		return err
	}

	return write(filename, encryptedData)
}

// Get a struct with key
func (m *KVod) Get(key string, value interface{}) error {
	filename := m.getFilename(key)

	data, err := read(filename)
	if err != nil {
		return err
	}
	s, err := m.getMapFromData(data)
	if err != nil {
		return err
	}
	return deserialize(s[key], value)
}

func (m *KVod) getMapFromData(data []byte) (map[string][]byte, error) {
	decrypted, err := m.encrypter.decrypt(data)
	if err != nil {
		return nil, err
	}

	s := make(map[string][]byte)
	err = deserialize(decrypted, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Delete a value by key
func (m *KVod) Delete(key string) error {
	return os.Remove(m.getFilename(key))
}

// GetKeys returns a slice of all keys available
func (m *KVod) GetKeys() ([]string, error) {
	var ret []string
	// read all files, decrypt the content and get the key
	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() != saltFilename {
			data, err := read(filepath.Join(m.path, f.Name()))
			if err != nil {
				return nil, err
			}
			kv, err := m.getMapFromData(data)
			for key := range kv {
				ret = append(ret, key)
			}
		}
	}
	return ret, err
}

func (m *KVod) getFilename(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	filename := hex.EncodeToString(h.Sum(nil))
	return filepath.Join(m.path, filename)
}
