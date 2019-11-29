package rpctypes

import (
	"encoding/json"
	"fmt"
	"hash"
	"testing"

	blake2b "github.com/minio/blake2b-simd"
)

func newBlake2b() (hash.Hash, error) {
	config := blake2b.Config{
		Size:   32,
		Person: []byte("ckb-default-hash"),
	}
	return blake2b.New(&config)
}

func TestSerializeBlock1(t *testing.T) {
	block1_data := loadTestFile(t, "block1.json")
	var block Block
	err := json.Unmarshal(block1_data, &block)
	if err != nil {
		t.Fatal(err)
	}
	expectedTxHashes := []string{
		"0xd8365de5c11fde09b8ff2142e22696422346b8b77bcde9fbc5a7c2231070e24c",
		"0xf2d0bb9da099956dc131c7f0377f12f90979155c9f5ed856282c1b0a77308442",
		"0x4b8c56426058e3caa789e346b19b5a971de994338bfb0461a6d83df1a2ef270b",
		"0x09ac0a3908e746692df59911e06a7a115ccaded0d84b3f181abe5bce3648d466",
	}
	for i, expectedHash := range expectedTxHashes {
		h, err := newBlake2b()
		if err != nil {
			t.Fatal(err)
		}
		err = block.Transactions[i].RawTransaction.SerializeToCore(h)
		if err != nil {
			t.Fatal(err)
		}
		txHash := h.Sum(nil)
		assertBytes(t, fmt.Sprintf("tx[%d] hash", i), txHash, expectedHash)
	}

	h, err := newBlake2b()
	if err != nil {
		t.Fatal(err)
	}
	err = block.Header.SerializeToCore(h)
	if err != nil {
		t.Fatal(err)
	}
	headerHash := h.Sum(nil)
	assertBytes(t, "header hash", headerHash,
		"0x7da7da17aeb1bec53c2f42364d59534435e741d7ac9cc1dcb694c2c3e37c4e3e")
}

func TestSerializeBlock2(t *testing.T) {
	block2_data := loadTestFile(t, "block2.json")
	var block Block
	err := json.Unmarshal(block2_data, &block)
	if err != nil {
		t.Fatal(err)
	}
	expectedTxHashes := []string{
		"0x5d9e7a4c4d3f90d2249eac4b5efb317b761cd42d878c3d4f65f350a37c2238c8",
	}
	for i, expectedHash := range expectedTxHashes {
		h, err := newBlake2b()
		if err != nil {
			t.Fatal(err)
		}
		err = block.Transactions[i].RawTransaction.SerializeToCore(h)
		if err != nil {
			t.Fatal(err)
		}
		txHash := h.Sum(nil)
		assertBytes(t, fmt.Sprintf("tx[%d] hash", i), txHash, expectedHash)
	}

	h, err := newBlake2b()
	if err != nil {
		t.Fatal(err)
	}
	err = block.Header.SerializeToCore(h)
	if err != nil {
		t.Fatal(err)
	}
	headerHash := h.Sum(nil)
	assertBytes(t, "header hash", headerHash,
		"0xf16ca901832577b55517b02aebd1bcd1d9440104f98d77f8b7406d4dbfa9498a")
}

func TestSerializeBlock3(t *testing.T) {
	block3_data := loadTestFile(t, "block3.json")
	var block Block
	err := json.Unmarshal(block3_data, &block)
	if err != nil {
		t.Fatal(err)
	}

	expectedTxHashes := []string{
		"0x813df66b968769b0fb240689fa1c699b7cdc5eed23cd038136b98432bd7dbaaa",
	}
	for i, expectedHash := range expectedTxHashes {
		h, err := newBlake2b()
		if err != nil {
			t.Fatal(err)
		}
		err = block.Transactions[i].RawTransaction.SerializeToCore(h)
		if err != nil {
			t.Fatal(err)
		}
		txHash := h.Sum(nil)
		assertBytes(t, fmt.Sprintf("tx[%d] hash", i), txHash, expectedHash)
	}

	h, err := newBlake2b()
	if err != nil {
		t.Fatal(err)
	}
	err = block.Header.SerializeToCore(h)
	if err != nil {
		t.Fatal(err)
	}
	headerHash := h.Sum(nil)
	assertBytes(t, "header hash", headerHash,
		"0x61ea89d07fef9470eb9137dcd7a000cf1ef3fdc0630b5a606d59bd5df40aeb50")
}
