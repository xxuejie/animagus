package store

import (
	"bytes"
	"encoding/gob"
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
)

type Client struct {
	DataDir string
	DB      *badger.DB
}

func NewClient(dataDir string) *Client {
	return &Client{
		DataDir: dataDir,
	}
}

func (c *Client) Open() error {
	var err error
	c.DB, err = badger.Open(badger.DefaultOptions(c.DataDir))
	return err
}

func (c *Client) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

func Get(txn *badger.Txn, key []byte) ([]byte, error) {
	var result []byte
	item, err := txn.Get(key)
	if err != nil {
		if err.Error() == "Key not found" {
			return result, nil
		}
		return nil, err
	}

	result, err = item.ValueCopy(nil)

	return result, err
}

func (c *Client) Get(key string) ([]byte, error) {
	var result []byte
	err := c.DB.View(func(txn *badger.Txn) error {
		var errr error
		result, errr = Get(txn, []byte(key))
		return errr
	})

	return result, err
}

func (c *Client) Set(key string, value []byte) error {
	err := c.DB.Update(func(txn *badger.Txn) error {
		errr := txn.Set([]byte(key), value)
		return errr
	})
	return err
}

func (c *Client) Del(key string) error {
	err := c.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	return err
}

func Sadd(txn *badger.Txn, key string, value []byte) error {
	keyByte := []byte(key)
	originValue, errr := Get(txn, keyByte)
	if errr != nil {
		return errr
	}

	var bytesArr [][]byte
	if len(originValue) != 0 {
		bytesArr, errr = Decode(originValue)
		if errr != nil {
			return errr
		}
	}

	if contains(bytesArr, value) {
		return nil
	}

	bytesArr = append(bytesArr, value)
	resultBytes, errr := Encode(bytesArr)
	if errr != nil {
		return errr
	}

	errr = txn.Set(keyByte, resultBytes)
	return errr
}

func (c *Client) Sadd(key string, value []byte) error {
	err := c.DB.Update(func(txn *badger.Txn) error {
		errr := Sadd(txn, key, value)
		return errr
	})
	return err
}

func Srem(txn *badger.Txn, key string, value []byte) error {
	keyByte := []byte(key)
	originValue, errr := Get(txn, keyByte)
	if errr != nil {
		return errr
	}

	bytesArr, errr := Decode(originValue)
	if errr != nil {
		return errr
	}
	if !contains(bytesArr, value) {
		return nil
	}

	bytesArr = sliceRemove(bytesArr, value)
	resultBytes, errr := Encode(bytesArr)
	if errr != nil {
		return errr
	}

	errr = txn.Set(keyByte, resultBytes)
	return errr
}

func (c *Client) Srem(key string, value []byte) error {
	err := c.DB.Update(func(txn *badger.Txn) error {
		errr := Srem(txn, key, value)
		return errr
	})
	return err
}

func Sget(txn *badger.Txn, key string) ([][]byte, error) {
	keyByte := []byte(key)
	var bytesArr [][]byte
	originValue, errr := Get(txn, keyByte)
	if errr != nil {
		return nil, errr
	}

	if len(originValue) != 0 {
		bytesArr, errr = Decode(originValue)
		if errr != nil {
			return nil, errr
		}
	}

	return bytesArr, errr
}

func (c *Client) Sget(key string) ([][]byte, error) {
	var bytesArr [][]byte
	err := c.DB.Update(func(txn *badger.Txn) error {
		var errr error
		bytesArr, errr = Sget(txn, key)
		return errr
	})
	return bytesArr, err
}

func Spop(txn *badger.Txn, key string) ([]byte, error) {
	keyByte := []byte(key)

	values, err := Sget(txn, key)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}
	firstValue := values[0]

	newValues := sliceRemoveIdx(values, 0)

	resultBytes, err := Encode(newValues)
	if err != nil {
		return nil, err
	}

	err = txn.Set(keyByte, resultBytes)
	return firstValue, err
}

func (c *Client) Spop(key string) ([]byte, error) {
	var result []byte = nil
	err := c.DB.Update(func(txn *badger.Txn) error {
		var errr error
		result, errr = Spop(txn, key)
		return errr
	})
	return result, err
}

func Do(txn *badger.Txn, command string, args ...interface{}) error {
	if command == "SADD" {
		key, _ := args[0].(string)
		value, _ := args[1].([]byte)
		return Sadd(txn, key, value)
	}
	if command == "SREM" {
		key, _ := args[0].(string)
		value, _ := args[1].([]byte)
		return Srem(txn, key, value)
	}
	if command == "SET" {
		key, _ := args[0].(string)
		value, _ := args[1].([]byte)
		return txn.Set([]byte(key), value)
	}
	if command == "PUBLISH" {
		key, _ := args[0].(string)
		value, _ := args[1].([]byte)
		Bus.Publish(key, value)
		return nil
	}

	return fmt.Errorf("command %v not found", command)
}

func Encode(value [][]byte) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(value)
	return result.Bytes(), err
}

func Decode(value []byte) ([][]byte, error) {
	from := bytes.NewBuffer(value)
	var result [][]byte
	decoder := gob.NewDecoder(from)
	err := decoder.Decode(&result)
	return result, err
}

func contains(s [][]byte, e []byte) bool {
	for _, a := range s {
		if bytes.Equal(a, e) {
			return true
		}
	}
	return false
}

func sliceRemoveIdx(slice [][]byte, idx int) [][]byte {
	return append(slice[0:idx], slice[idx+1:]...)
}

func sliceRemove(slices [][]byte, toRemove []byte) [][]byte {
	for idx, value := range slices {
		if bytes.Equal(value, toRemove) {
			return sliceRemoveIdx(slices, idx)
		}
	}
	return slices
}
