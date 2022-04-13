package kvod

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"io/ioutil"
	"os"
)

func serialize[T any](data T) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(data)
	return buffer.Bytes(), err
}

func deserialize[T any](data []byte, ret T) error {
	var buffer bytes.Buffer
	buffer.Write(data)
	dec := gob.NewDecoder(&buffer)
	return dec.Decode(ret)
}

func write(filename string, data []byte) error {
	file, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0660,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func read(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GenerateRandom generates a random of size provided
func GenerateRandom(size int) ([]byte, error) {
	randomData := make([]byte, size)
	_, err := rand.Reader.Read(randomData)
	return randomData, err
}
