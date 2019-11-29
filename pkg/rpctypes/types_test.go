package rpctypes

import (
	"encoding/json"
	"testing"
)

func TestDeserializeBlock1(t *testing.T) {
	block1_data := loadTestFile(t, "block1.json")
	var block Block
	err := json.Unmarshal(block1_data, &block)
	if err != nil {
		t.Fatal(err)
	}
	if block.Header.Number != 15081 {
		t.Errorf("Invalid block number: %d, expected: %d",
			block.Header.Number, 15081)
	}
}

func TestDeserializeBlock2(t *testing.T) {
	block2_data := loadTestFile(t, "block2.json")
	var block Block
	err := json.Unmarshal(block2_data, &block)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeserializeBlock3(t *testing.T) {
	block3_data := loadTestFile(t, "block3.json")
	var block Block
	err := json.Unmarshal(block3_data, &block)
	if err != nil {
		t.Fatal(err)
	}
}
