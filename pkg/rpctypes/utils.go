package rpctypes

import (
	"hash"

	blake2b "github.com/minio/blake2b-simd"
)

func newBlake2b() (hash.Hash, error) {
	config := blake2b.Config{
		Size:   32,
		Person: []byte("ckb-default-hash"),
	}
	return blake2b.New(&config)
}

func CalculateHash(e CoreSerializer) ([]byte, error) {
	h, err := newBlake2b()
	if err != nil {
		return nil, err
	}
	err = e.SerializeToCore(h)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
