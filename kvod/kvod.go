package kvod

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"sort"
	"sync"

	"os"
	"path/filepath"
)

const (
	dbFilename = "data.db"
)

type dbDetails struct {
	Salt []byte
}

// KVod basic struct
type KVod struct {
	path      string
	encrypter crypto
	details   dbDetails
}

type KVodContainer[T any] struct {
	kvod        *KVod
	path        string
	pathEncoded string
	mu          sync.Mutex
}

// Init KVod struct0
func Init(path string, password string) *KVod {
	err := os.MkdirAll(path, 0770)
	if err != nil {
		log.Fatal(err)
	}
	dbDetailsFile := filepath.Join(path, dbFilename)

	salt, err := GenerateRandom(saltSize)
	if err != nil {
		log.Fatal(err)
	}
	var details dbDetails

	if _, err := os.Stat(dbDetailsFile); os.IsNotExist(err) {
		details = dbDetails{
			Salt: salt,
		}
		serializedDetails, err := serialize(details)
		if err != nil {
			log.Fatal(err)
		}
		err = write(dbDetailsFile, serializedDetails)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		serializedDetails, err := read(dbDetailsFile)
		if err != nil {
			log.Fatal(err)
		}
		err = deserialize(serializedDetails, &details)
		if err != nil {
			log.Fatal(err)
		}
	}

	crypto := InitCrypto(password, details.Salt)
	return &KVod{path, *crypto, details}
}

func CreateContainer[T any](kvod *KVod, containerPath string) *KVodContainer[T] {
	pathEncoded := sha256Hex(containerPath)
	err := os.MkdirAll(filepath.Join(kvod.path, pathEncoded), 0770)
	if err != nil {
		log.Fatal(err)
	}
	return &KVodContainer[T]{
		path:        containerPath,
		pathEncoded: pathEncoded,
		kvod:        kvod,
		mu:          sync.Mutex{},
	}
}

// Put a struct using a string as key
func (m *KVodContainer[T]) Put(key string, value T) error {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.Lock()
	defer m.mu.Unlock()

	return os.Remove(m.getFilename(key))
}

// GetKeys returns a slice of all keys available
func (m *KVodContainer[T]) GetKeys() ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var ret []string
	// read all files, decrypt the content and get the key
	files, err := ioutil.ReadDir(filepath.Join(m.kvod.path, m.pathEncoded))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() != dbFilename {
			data, err := read(filepath.Join(m.kvod.path, m.pathEncoded, f.Name()))
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
	sort.Strings(ret)
	return ret, err
}

// GetAll returns a slice of all data available
func (m *KVodContainer[T]) GetData() ([]T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var ret []T
	// read all files, decrypt the content and get the data
	files, err := ioutil.ReadDir(filepath.Join(m.kvod.path, m.pathEncoded))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() != dbFilename {
			data, err := read(filepath.Join(m.kvod.path, m.pathEncoded, f.Name()))
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
	m.mu.Lock()
	defer m.mu.Unlock()

	ret := make(map[string]T)
	// read all files, decrypt the content and get the data
	files, err := ioutil.ReadDir(filepath.Join(m.kvod.path, m.pathEncoded))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() != dbFilename {
			data, err := read(filepath.Join(m.kvod.path, m.pathEncoded, f.Name()))
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
	filename := sha256Hex(key)
	return filepath.Join(m.kvod.path, m.pathEncoded, filename)
}

func sha256Hex(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}
