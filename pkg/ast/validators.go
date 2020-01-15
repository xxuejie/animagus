package ast

import (
	"fmt"
	"math"
)

func IsValidScript(value *Value) error {
	if value.GetT() != Value_SCRIPT {
		return fmt.Errorf("Invalid script type!")
	}
	if len(value.GetChildren()) != 3 {
		return fmt.Errorf("Invalid number of script items!")
	}
	if value.GetChildren()[0].GetT() != Value_BYTES ||
		len(value.GetChildren()[0].GetRaw()) != 32 ||
		value.GetChildren()[1].GetT() != Value_UINT64 ||
		(value.GetChildren()[1].GetU() != 0 && value.GetChildren()[1].GetU() != 1) ||
		value.GetChildren()[2].GetT() != Value_BYTES {
		return fmt.Errorf("Invalid child type!")
	}
	return nil
}

func IsValidOutPoint(value *Value) error {
	if value.GetT() != Value_OUT_POINT {
		return fmt.Errorf("Invalid out point!")
	}
	if len(value.GetChildren()) != 2 {
		return fmt.Errorf("Invalid number of out point items!")
	}
	if value.GetChildren()[0].GetT() != Value_BYTES ||
		len(value.GetChildren()[0].GetRaw()) != 32 ||
		value.GetChildren()[1].GetT() != Value_UINT64 ||
		value.GetChildren()[1].GetU() >= math.MaxUint32 {
		return fmt.Errorf("Invalid child type!")
	}
	return nil
}

func IsValidCell(value *Value) error {
	if value.GetT() != Value_CELL {
		return fmt.Errorf("Invalid cell!")
	}
	if len(value.GetChildren()) != 4 && len(value.GetChildren()) != 5 {
		return fmt.Errorf("Invalid number of out point items!")
	}
	if value.GetChildren()[0].GetT() != Value_UINT64 ||
		(IsValidScript(value.GetChildren()[1]) != nil) ||
		(!(value.GetChildren()[2].GetT() == Value_NIL || IsValidScript(value.GetChildren()[2]) == nil)) ||
		value.GetChildren()[3].GetT() != Value_BYTES ||
		(len(value.GetChildren()) == 5 && (IsValidOutPoint(value.GetChildren()[4]) != nil)) {
		return fmt.Errorf("Invalid child type")
	}
	return nil
}

func IsValidCellInput(value *Value) error {
	if value.GetT() != Value_CELL_INPUT {
		return fmt.Errorf("Invalid cell input!")
	}
	if len(value.GetChildren()) != 2 {
		return fmt.Errorf("Invalid number of cell input items")
	}
	if IsValidOutPoint(value.GetChildren()[0]) != nil ||
		value.GetChildren()[1].GetT() != Value_UINT64 {
		return fmt.Errorf("Invalid child type")
	}
	return nil
}

func IsValidCellDep(value *Value) error {
	if value.GetT() != Value_CELL_DEP {
		return fmt.Errorf("Invalid cell dep!")
	}
	if len(value.GetChildren()) != 2 {
		return fmt.Errorf("Invalid number of cell dep items")
	}
	if IsValidOutPoint(value.GetChildren()[0]) != nil ||
		value.GetChildren()[1].GetT() != Value_UINT64 ||
		(value.GetChildren()[1].GetU() != 0 && value.GetChildren()[1].GetU() != 1) {
		return fmt.Errorf("Invalid child type")
	}
	return nil
}

func IsValidTransaction(value *Value) error {
	if value.GetT() != Value_TRANSACTION {
		return fmt.Errorf("Invalid transaction!")
	}
	if len(value.GetChildren()) != 3 {
		return fmt.Errorf("Invalid number of transaction items")
	}
	if value.GetChildren()[0].GetT() != Value_LIST ||
		value.GetChildren()[1].GetT() != Value_LIST ||
		value.GetChildren()[2].GetT() != Value_LIST {
		return fmt.Errorf("Invalid child type")
	}
	for _, child := range value.GetChildren()[0].GetChildren() {
		if err := IsValidCellInput(child); err != nil {
			return err
		}
	}
	for _, child := range value.GetChildren()[1].GetChildren() {
		if err := IsValidCell(child); err != nil {
			return err
		}
	}
	for _, child := range value.GetChildren()[2].GetChildren() {
		if err := IsValidCellDep(child); err != nil {
			return err
		}
	}
	return nil
}
