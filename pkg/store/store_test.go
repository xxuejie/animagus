package store

import (
	"bytes"
	"testing"

	badger "github.com/dgraph-io/badger/v2"
)

func bytesSliceEqual(a [][]byte, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, aa := range a {
		bb := b[i]
		if !bytes.Equal(aa, bb) {
			return false
		}
	}

	return true
}

func TestGetWithEmtpyKey(t *testing.T) {
	key := "get_key"
	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()
	value, err := client.Get(key)
	if err != nil {
		t.Errorf("get error: %v", err)
	}
	if value != nil {
		t.Errorf("value should be nil, but got %v", value)
	}
}

func TestGetAndSet(t *testing.T) {
	key := "key"
	value := []byte("some_value")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	// set
	err = client.Set(key, value)
	if err != nil {
		t.Errorf("set error: %v", err)
	}

	// get
	getValue, err := client.Get(key)
	if err != nil {
		t.Errorf("get error: %v", err)
	}

	if !bytes.Equal(getValue, value) {
		t.Errorf("expected: %v, got: %v", value, getValue)
	}
}

func TestDel(t *testing.T) {
	key := "key"
	value := []byte("some_value")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	// set
	err = client.Set(key, value)
	if err != nil {
		t.Errorf("set error: %v", err)
	}

	// del
	err = client.Del(key)
	if err != nil {
		t.Errorf("del error: %v", err)
	}

	// get
	getValue, err := client.Get(key)
	if err != nil {
		t.Errorf("get error: %v", err)
	}

	if getValue != nil {
		t.Errorf("expected: nil, got: %v", getValue)
	}
}

func TestSaddAndSget(t *testing.T) {
	key := "key"
	val1 := []byte("A")
	val2 := []byte("B")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	// sadd
	err = client.Sadd(key, val1)
	if err != nil {
		t.Errorf("sadd error: %v", err)
	}
	err = client.Sadd(key, val2)
	if err != nil {
		t.Errorf("sadd error: %v", err)
	}
	err = client.Sadd(key, val1)
	if err != nil {
		t.Errorf("sadd error: %v", err)
	}

	// sget
	values, err := client.Sget(key)
	if err != nil {
		t.Errorf("sget error: %v", err)
	}

	expected := [][]byte{val1, val2}
	if !bytesSliceEqual(values, expected) {
		t.Errorf("expect %v, got %v", expected, values)
	}
}

func TestSrem(t *testing.T) {
	key := "key"
	val1 := []byte("A")
	val2 := []byte("B")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	// sadd
	err = client.Sadd(key, val1)
	if err != nil {
		t.Errorf("sadd error: %v", err)
	}
	err = client.Sadd(key, val2)
	if err != nil {
		t.Errorf("sadd error: %v", err)
	}

	// srem
	err = client.Srem(key, val1)
	if err != nil {
		t.Errorf("srem err: %v", err)
	}

	// sget
	values, err := client.Sget(key)
	if err != nil {
		t.Errorf("sget error: %v", err)
	}

	expected := [][]byte{val2}
	if !bytesSliceEqual(values, expected) {
		t.Errorf("expect %v, got %v", expected, values)
	}
}

func TestDoSet(t *testing.T) {
	key := "key"
	value := []byte("some_value")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	err = client.DB.Update(func(txn *badger.Txn) error {
		return Do(txn, "SET", key, value)
	})
	if err != nil {
		t.Errorf("update error: %v", err)
	}

	getValue, err := client.Get(key)
	if err != nil {
		t.Errorf("get error: %v", err)
	}
	if !bytes.Equal(getValue, value) {
		t.Errorf("expect: %v, get: %v", value, getValue)
	}
}

func TestDoSadd(t *testing.T) {
	key := "key"
	value := []byte("some_value")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	err = client.DB.Update(func(txn *badger.Txn) error {
		return Do(txn, "SADD", key, value)
	})
	if err != nil {
		t.Errorf("update error: %v", err)
	}

	getValue, err := client.Sget(key)
	if err != nil {
		t.Errorf("sget error: %v", err)
	}
	expected := [][]byte{value}
	if !bytesSliceEqual(getValue, expected) {
		t.Errorf("expect: %v, get: %v", expected, getValue)
	}
}

func TestDoRem(t *testing.T) {
	key := "key"
	value := []byte("some_value")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	err = client.DB.Update(func(txn *badger.Txn) error {
		errr := Do(txn, "SADD", key, value)
		if errr != nil {
			return errr
		}
		return Do(txn, "SREM", key, value)
	})
	if err != nil {
		t.Errorf("update error: %v", err)
	}

	getValue, err := client.Sget(key)
	if err != nil {
		t.Errorf("sget error: %v", err)
	}
	expected := [][]byte{}
	if !bytesSliceEqual(getValue, expected) {
		t.Errorf("expect: %v, get: %v", expected, getValue)
	}
}

func TestDoOtherCommand(t *testing.T) {
	key := "key"
	value := []byte("some_value")

	client := NewClient("")
	err := client.openInMemory()
	if err != nil {
		t.Errorf("open error: %v", err)
	}
	defer client.Close()

	err = client.DB.Update(func(txn *badger.Txn) error {
		return Do(txn, "SERR", key, value)
	})
	if err != nil && err.Error() != "command SERR not found" {
		t.Errorf("update error: %v", err)
	}
}
