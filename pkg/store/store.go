package store

import (
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

// for test
func (c *Client) openInMemory() error {
	var err error
	opt := badger.DefaultOptions("").WithInMemory(true)
	c.DB, err = badger.Open(opt)
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

func skeyPrefix(key string) []byte {
	return append([]byte(key), []byte(":")...)
}

func skey(key string, value []byte) []byte {
	return append(skeyPrefix(key), value...)
}

func Sadd(txn *badger.Txn, key string, value []byte) error {
	newKey := skey(key, value)

	errr := txn.Set(newKey, nil)
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
	newKey := skey(key, value)

	errr := txn.Delete(newKey)
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
	var values [][]byte

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	it := txn.NewIterator(opts)
	defer it.Close()
	prefix := skeyPrefix(key)
	prefixLen := len(prefix)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		value := k[prefixLen:]
		value2 := make([]byte, len(value))
		copy(value2, value)
		values = append(values, value2)
	}

	return values, nil
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
