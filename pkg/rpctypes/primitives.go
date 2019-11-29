package rpctypes

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Bytes []byte

func (data *Bytes) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !strings.HasPrefix(s, "0x") {
		return fmt.Errorf("Bytes should be a hex string started with 0x!")
	}
	b, err := hex.DecodeString(s[2:])
	if err != nil {
		return err
	}
	*data = Bytes(b)
	return nil
}

func (data Bytes) MarshalJSON() ([]byte, error) {
	result := make([]byte, hex.EncodedLen(len(data))+2)
	copy(result[0:2], []byte("0x"))
	hex.Encode(result[2:], data)
	return json.Marshal(string(result))
}

type VarBytes []byte

func (data *VarBytes) UnmarshalJSON(b []byte) error {
	var d Bytes
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	*data = VarBytes(d)
	return nil
}

func (data VarBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(Bytes(data))
}

type Uint128 struct {
	V *big.Int
}

func (u *Uint128) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i := new(big.Int)
	_, success := i.SetString(s, 0)
	if !success {
		return fmt.Errorf("Setting uint128 failure!")
	}
	u.V = i
	return nil
}

func (u Uint128) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%s", u.V.Text(16)))
}

type Uint64 uint64

func (u *Uint64) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*u = Uint64(v)
	return nil
}

func (u Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%x", u))
}

type Uint32 uint32

func (u *Uint32) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return err
	}
	*u = Uint32(v)
	return nil
}

func (u Uint32) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%x", u))
}

type Hash [32]byte

func (h *Hash) UnmarshalJSON(b []byte) error {
	var data Bytes
	err := data.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	if len(data) != 32 {
		return fmt.Errorf("Hash should be exactly 32 bytes!")
	}
	copy(h[:], data)
	return nil
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return Bytes(h[:]).MarshalJSON()
}

type ProposalShortId [10]byte

func (h *ProposalShortId) UnmarshalJSON(b []byte) error {
	var data Bytes
	err := data.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	if len(data) != 10 {
		return fmt.Errorf("Proposal short ID should be exactly 10 bytes!")
	}
	copy(h[:], data)
	return nil
}

func (h ProposalShortId) MarshalJSON() ([]byte, error) {
	return Bytes(h[:]).MarshalJSON()
}

type ScriptHashType byte

const (
	Data ScriptHashType = 0
	Type ScriptHashType = 1
)

func (t *ScriptHashType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "data":
		*t = Data
	case "type":
		*t = Type
	default:
		return fmt.Errorf("Invalid script hash type!")
	}
	return nil
}

func (t ScriptHashType) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	case Data:
		s = "data"
	case Type:
		s = "type"
	default:
		return nil, fmt.Errorf("Invalid script hash type!")
	}
	return json.Marshal(s)
}

type DepType byte

const (
	Code     DepType = 0
	DepGroup DepType = 1
)

func (t *DepType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "code":
		*t = Code
	case "dep_group":
		*t = DepGroup
	default:
		return fmt.Errorf("Invalid script hash type!")
	}
	return nil
}

func (t DepType) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	case Code:
		s = "code"
	case DepGroup:
		s = "dep_group"
	default:
		return nil, fmt.Errorf("Invalid script hash type!")
	}
	return json.Marshal(s)
}
