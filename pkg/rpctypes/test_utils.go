package rpctypes

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func loadTestFile(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func assertBytes(t *testing.T, name string, value []byte, expected string) {
	expectedString := expected
	if strings.HasPrefix(expectedString, "0x") {
		expectedString = expectedString[2:]
	}
	expectedBytes, err := hex.DecodeString(expectedString)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(value, expectedBytes) != 0 {
		t.Errorf("Field %s has incorrect value! Expected: %s", name, expected)
	}
}
