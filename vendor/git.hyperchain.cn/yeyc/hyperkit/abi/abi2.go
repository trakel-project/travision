package abi

import (
	"fmt"
	"reflect"
	"strconv"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
)

func (abi ABI) PackJSON(name string, args ...interface{}) ([]byte, error) {
	// Fetch the ABI of the requested method
	var method Method

	if name == "" {
		method = abi.Constructor
	} else {
		m, exist := abi.Methods[name]
		if !exist {
			return nil, fmt.Errorf("method '%s' not found", name)
		}
		method = m
	}
	if len(args) != len(method.Inputs) {
		return nil, fmt.Errorf("`%s` argument count mismatch: %d for %d", name, len(args), len(method.Inputs))
	}
	var vals []interface{}
	for i, a := range args {
		v, err := method.Inputs[i].Type.JsonType2GoType(reflect.ValueOf(a))
		if err != nil {
			return nil, fmt.Errorf("JsonType2GoType: %s", err)
		}
		vals = append(vals, v)
	}
	arguments, err := method.pack(method, vals...)
	if err != nil {
		return nil, err
	}
	// Pack up the method ID too if not a constructor and return
	if name == "" {
		return arguments, nil
	}
	return append(method.Id(), arguments...), nil
}

//func (method Method) ToType(args ...interface{}) ([]interface{}, error){
//	e_args := make([]interface{}, 0, len(args))
//
//	return e_args
//}

func (t Type) JsonType2GoType(val reflect.Value) (interface{}, error) {

	val = indirect(val)

	if (t.IsArray || t.IsSlice) && t.T != BytesTy && t.T != FixedBytesTy {
		err := t.sliceCheck(val)

		if err != nil {
			return nil, err
		}
		var s []interface{}
		for i := 0; i < val.Len(); i++ {
			o, err := t.Elem.JsonType2GoType(val.Index(i))
			if err != nil{
				return nil, err
			}
			s = append(s, o)
		}
		return s, nil
	}

	if t.Kind == reflect.Bool && val.Kind() == reflect.Bool {
		return val.Bool(), nil
	}

	if val.Kind() != reflect.String {
		return nil, fmt.Errorf("abi2: param %v(%v) can't parse to GoType", val, val.Kind())
	}
	return t.jsonType2BasicGoType(val.String())
}

func (t Type) jsonType2BasicGoType(str string) (interface{}, error) {
	if t.T == BytesTy || t.T == FixedBytesTy {
		if len(str) > 1 && str[0:2] == "0x" {
			return common.FromHex(str), nil
		}
		return []byte(str), nil
	}
	switch t.Kind {
	case reflect.Bool:
		i, err := strconv.ParseBool(str)
		if err != nil {
			return nil, err
		}
		return i, nil
	case reflect.Int8:
		i, err := strconv.ParseInt(str, 0, 8)
		if err != nil {
			return nil, err
		}
		return int8(i), nil
	case reflect.Int16:
		i, err := strconv.ParseInt(str, 0, 16)
		if err != nil {
			return nil, err
		}
		return int16(i), nil
	case reflect.Int32:
		i, err := strconv.ParseInt(str, 0, 32)
		if err != nil {
			return nil, err
		}
		return int32(i), nil
	case reflect.Int64:
		i, err := strconv.ParseInt(str, 0, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case reflect.Uint8:
		i, err := strconv.ParseUint(str, 0, 8)
		if err != nil {
			return nil, err
		}
		return uint8(i), nil
	case reflect.Uint16:
		i, err := strconv.ParseUint(str, 0, 16)
		if err != nil {
			return nil, err
		}
		return uint16(i), nil
	case reflect.Uint32:
		i, err := strconv.ParseUint(str, 0, 32)
		if err != nil {
			return nil, err
		}
		return uint32(i), nil
	case reflect.Uint64:
		i, err := strconv.ParseUint(str, 0, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case reflect.Ptr:
		i := new(big.Int)
		if _, ok := i.SetString(str, 0); !ok {
			return nil, fmt.Errorf("abi2: %s can't parse to int256/uint256", str)
		}
		return i, nil
	case reflect.Array:
		return common.HexToAddress(str), nil
	case reflect.String:
		return str, nil
	}
	return nil, fmt.Errorf("abi2: unsupported arg type: %v", t.Kind)
}

func (t Type) sliceCheck(val reflect.Value) error  {
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		return fmt.Errorf("abi2: require slice or array to parse to %v", t.stringKind)
	}
	if t.IsArray && val.Len() != t.SliceSize {
		return fmt.Errorf("abi2: require %v length array/slice to parse to %v", t.SliceSize, t.stringKind)
	}
	return nil
}

// dereferences the value which is packed in interface or reference to ptr
func dereferInterfacePtr(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		return dereferInterfacePtr(v.Elem())
	}
	return v
}