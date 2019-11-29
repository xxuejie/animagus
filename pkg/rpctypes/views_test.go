package rpctypes

import (
	"encoding/json"
	"testing"
)

func TestBlockHash(t *testing.T) {
	block1_data := loadTestFile(t, "block1.json")
	var block BlockView
	err := json.Unmarshal(block1_data, &block)
	if err != nil {
		t.Fatal(err)
	}
	if block.Header.Number != 15081 {
		t.Errorf("Invalid block number: %d, expected: %d",
			block.Header.Number, 15081)
	}
	assertBytes(t, "header hash", block.Header.Hash[:],
		"0x7da7da17aeb1bec53c2f42364d59534435e741d7ac9cc1dcb694c2c3e37c4e3e")
}
