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

type KVodContainer[T any] struct {
	kvod *KVod
	path string
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

func CreateContainer[T any](kvod *KVod, containerPath string) *KVodContainer[T] {
	err := os.MkdirAll(filepath.Join(kvod.path, containerPath), 0770)
	if err != nil {
		log.Fatal(err)
	}
	return &KVodContainer[T]{
		path: containerPath,
		kvod: kvod,
	}
}

// Put a struct using a string as key
func (m *KVodContainer[T]) Put(key string, value T) error {
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

	encryptedData, err := m.kvod.encrypter.encrypt(data)
	if err != nil {
		return err
	}

	return write(filename, encryptedData)
}

// Get a struct with key
func (m *KVodContainer[T]) Get(key string) (*T, error) {
	filename := m.getFilename(key)

	data, err := read(filename)
	if err != nil {
		return nil, err
	}
	s, err := m.kvod.getMapFromData(data)
	if err != nil {
		return nil, err
	}
	var ret T
	err = deserialize(s[key], &ret)
	return &ret, err
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
func (m *KVodContainer[T]) Delete(key string) error {
	return os.Remove(m.getFilename(key))
}

// GetKeys returns a slice of all keys available
func (m *KVodContainer[T]) GetKeys() ([]string, error) {
	var ret []string
	// read all files, decrypt the content and get the key
	files, err := ioutil.ReadDir(filepath.Join(m.kvod.path, m.path))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() != saltFilename {
			data, err := read(filepath.Join(m.kvod.path, m.path, f.Name()))
			if err != nil {
				return nil, err
			}
			kv, err := m.kvod.getMapFromData(data)
			if err != nil {
				return nil, err
			}
			for key := range kv {
				ret = append(ret, key)
			}
		}
	}
	return ret, err
}

// GetAll returns a slice of all data available
func (m *KVodContainer[T]) GetData() ([]T, error) {
	var ret []T
	// read all files, decrypt the content and get the data
	files, err := ioutil.ReadDir(filepath.Join(m.kvod.path, m.path))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() != saltFilename {
			data, err := read(filepath.Join(m.kvod.path, m.path, f.Name()))
			if err != nil {
				return nil, err
			}
			kv, err := m.kvod.getMapFromData(data)
			if err != nil {
				return nil, err
			}
			for key := range kv {
				var val T
				err = deserialize(kv[key], &val)
				if err != nil {
					return nil, err
				}
				ret = append(ret, val)
			}
		}
	}
	return ret, err
}

// GetAll returns a slice of all data available
func (m *KVodContainer[T]) GetAll() (map[string]T, error) {
	ret := make(map[string]T)
	// read all files, decrypt the content and get the data
	files, err := ioutil.ReadDir(filepath.Join(m.kvod.path, m.path))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() != saltFilename {
			data, err := read(filepath.Join(m.kvod.path, m.path, f.Name()))
			if err != nil {
				return nil, err
			}
			kv, err := m.kvod.getMapFromData(data)
			if err != nil {
				return nil, err
			}
			for key := range kv {
				var val T
				err = deserialize(kv[key], &val)
				if err != nil {
					return nil, err
				}
				ret[key] = val
			}
		}
	}
	return ret, err
}

func (m *KVodContainer[T]) getFilename(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	filename := hex.EncodeToString(h.Sum(nil))
	return filepath.Join(m.kvod.path, m.path, filename)
}
