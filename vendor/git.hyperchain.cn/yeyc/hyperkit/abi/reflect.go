package abi

import (
	"fmt"
	"reflect"
)

// indirect recursively dereferences the value until it either gets the value
// or finds a big.Int
func indirect(v reflect.Value) reflect.Value {
	// add
	if v.Kind() == reflect.Interface {
		return indirect(v.Elem())
	}
	// add
	if v.Kind() == reflect.Ptr && v.Elem().Type() != big_t {
		return indirect(v.Elem())
	}
	return v
}

// reflectIntKind returns the reflect using the given size and
// unsignedness.
func reflectIntKind(unsigned bool, size int) reflect.Kind {
	switch size {
	case 8:
		if unsigned {
			return reflect.Uint8
		}
		return reflect.Int8
	case 16:
		if unsigned {
			return reflect.Uint16
		}
		return reflect.Int16
	case 32:
		if unsigned {
			return reflect.Uint32
		}
		return reflect.Int32
	case 64:
		if unsigned {
			return reflect.Uint64
		}
		return reflect.Int64
	}
	return reflect.Ptr
}

// mustArrayToBytesSlice creates a new byte slice with the exact same size as value
// and copies the bytes in value to the new slice.
func mustArrayToByteSlice(value reflect.Value) reflect.Value {
	slice := reflect.MakeSlice(reflect.TypeOf([]byte{}), value.Len(), value.Len())
	reflect.Copy(slice, value)
	return slice
}

// set attempts to assign src to dst by either setting, copying or otherwise.
//
// set is a bit more lenient when it comes to assignment and doesn't force an as
// strict ruleset as bare `reflect` does.
func set(dst, src reflect.Value, output Argument) error {
	dstType := dst.Type()
	srcType := src.Type()

	switch {
	case dstType.AssignableTo(src.Type()):
		dst.Set(src)
	case dstType.Kind() == reflect.Array && srcType.Kind() == reflect.Slice:
		if dst.Len() < output.Type.SliceSize {
			return fmt.Errorf("abi: cannot unmarshal src (len=%d) in to dst (len=%d)", output.Type.SliceSize, dst.Len())
		}
		reflect.Copy(dst, src)
	case dstType.Kind() == reflect.Interface:
		dst.Set(src)
	case dstType.Kind() == reflect.Ptr:
		return set(dst.Elem(), src, output)
	default:
		return fmt.Errorf("abi: cannot unmarshal %v in to %v", src.Type(), dst.Type())
	}
	return nil
}
