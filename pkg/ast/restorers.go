package ast

import (
	"github.com/xxuejie/animagus/pkg/rpctypes"
)

func RestoreScript(value *Value, validate bool) (result rpctypes.Script, err error) {
	if validate {
		err = IsValidScript(value)
		if err != nil {
			return
		}
	}
	copy(result.CodeHash[:], value.GetChildren()[0].GetRaw())
	if value.GetChildren()[1].GetU() == 1 {
		result.HashType = rpctypes.Type
	} else {
		result.HashType = rpctypes.Data
	}
	result.Args = make([]byte, len(value.GetChildren()[2].GetRaw()))
	copy(result.Args, value.GetChildren()[2].GetRaw())
	return
}

func RestoreOutPoint(value *Value, validate bool) (result rpctypes.OutPoint, err error) {
	if validate {
		err = IsValidOutPoint(value)
		if err != nil {
			return
		}
	}
	copy(result.TxHash[:], value.GetChildren()[0].GetRaw())
	result.Index = rpctypes.Uint32(value.GetChildren()[1].GetU())
	return
}

func RestoreCell(value *Value, validate bool) (cell rpctypes.CellOutput, cellData rpctypes.Bytes, outPoint *rpctypes.OutPoint, err error) {
	if validate {
		err = IsValidCell(value)
		if err != nil {
			return
		}
	}
	cell.Capacity = rpctypes.Uint64(value.GetChildren()[0].GetU())
	cell.Lock, err = RestoreScript(value.GetChildren()[1], false)
	if err != nil {
		return
	}
	if value.GetChildren()[2].GetT() != Value_NIL {
		var t rpctypes.Script
		t, err = RestoreScript(value.GetChildren()[2], false)
		if err != nil {
			return
		}
		cell.Type = &t
	}
	cellData = make([]byte, len(value.GetChildren()[3].GetRaw()))
	copy(cellData, value.GetChildren()[3].GetRaw())
	if len(value.GetChildren()) == 5 {
		var o rpctypes.OutPoint
		o, err = RestoreOutPoint(value.GetChildren()[4], false)
		if err != nil {
			return
		}
		outPoint = &o
	}
	return
}

func RestoreCellInput(value *Value, validate bool) (cellInput rpctypes.CellInput, err error) {
	if validate {
		err = IsValidCellInput(value)
		if err != nil {
			return
		}
	}
	cellInput.PreviousOutput, err = RestoreOutPoint(value.GetChildren()[0], false)
	if err != nil {
		return
	}
	cellInput.Since = rpctypes.Uint64(value.GetChildren()[1].GetU())
	return
}

func RestoreCellDep(value *Value, validate bool) (cellDep rpctypes.CellDep, err error) {
	if validate {
		err = IsValidCellDep(value)
		if err != nil {
			return
		}
	}
	cellDep.OutPoint, err = RestoreOutPoint(value.GetChildren()[0], false)
	if err != nil {
		return
	}
	if value.GetChildren()[1].GetU() == 1 {
		cellDep.DepType = rpctypes.DepGroup
	} else {
		cellDep.DepType = rpctypes.Code
	}
	return
}

func RestoreTransaction(value *Value, validate bool) (tx rpctypes.Transaction, err error) {
	if validate {
		err = IsValidTransaction(value)
		if err != nil {
			return
		}
	}
	tx.HeaderDeps = []rpctypes.Hash{}
	for _, input := range value.GetChildren()[0].GetChildren() {
		var restoredInput rpctypes.CellInput
		restoredInput, err = RestoreCellInput(input, false)
		if err != nil {
			return
		}
		tx.Inputs = append(tx.Inputs, restoredInput)
		tx.Witnesses = append(tx.Witnesses, []byte{})
	}
	for _, output := range value.GetChildren()[1].GetChildren() {
		var cell rpctypes.CellOutput
		var cellData rpctypes.Bytes
		cell, cellData, _, err = RestoreCell(output, false)
		if err != nil {
			return
		}
		tx.Outputs = append(tx.Outputs, cell)
		tx.OutputsData = append(tx.OutputsData, cellData)
	}
	for _, dep := range value.GetChildren()[2].GetChildren() {
		var restoredDep rpctypes.CellDep
		restoredDep, err = RestoreCellDep(dep, false)
		if err != nil {
			return
		}
		tx.CellDeps = append(tx.CellDeps, restoredDep)
	}
	return
}
