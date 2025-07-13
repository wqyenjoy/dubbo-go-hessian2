/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hessian

import (
	"io"
	"reflect"

	perrors "github.com/pkg/errors"
)

/////////////////////////////////////////
// map/object
/////////////////////////////////////////

// ::= 'M' type (value value)* 'Z'  # key, value map pairs
// ::= 'H' (value value)* 'Z'       # untyped key, value
func (e *Encoder) encUntypedMap(m map[interface{}]interface{}) error {
	// check ref
	if n, ok := e.checkRefMap(reflect.ValueOf(m)); ok {
		e.buffer = encRef(e.buffer, n)
		return nil
	}

	var err error
	e.buffer = encByte(e.buffer, BC_MAP_UNTYPED)
	for k, v := range m {
		if err = e.Encode(k); err != nil {
			return err
		}
		if err = e.Encode(v); err != nil {
			return err
		}
	}

	e.buffer = encByte(e.buffer, BC_END) // 'Z'

	return nil
}

func getMapKey(key reflect.Value, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Bool:
		return key.Bool(), nil

	case reflect.Int8:
		return int8(key.Int()), nil
	case reflect.Int16:
		return int16(key.Int()), nil
	case reflect.Int32:
		return int32(key.Int()), nil
	case reflect.Int:
		return int(key.Int()), nil
	case reflect.Int64:
		return key.Int(), nil
	case reflect.Uint8:
		return byte(key.Uint()), nil
	case reflect.Uint16:
		return uint16(key.Uint()), nil
	case reflect.Uint32:
		return uint32(key.Uint()), nil
	case reflect.Uint:
		return uint(key.Uint()), nil
	case reflect.Uint64:
		return key.Uint(), nil
	case reflect.Float32:
		return float32(key.Float()), nil
	case reflect.Float64:
		return key.Float(), nil
	case reflect.Uintptr:
		return key.UnsafeAddr(), nil
	case reflect.String:
		return key.String(), nil
	}

	return nil, perrors.Errorf("unsupported map key kind %s", t.Kind().String())
}

func (e *Encoder) encMap(m interface{}) error {
	var (
		err   error
		k     interface{}
		typ   reflect.Type
		value reflect.Value
		keys  []reflect.Value
	)

	value = reflect.ValueOf(m)

	// check ref
	if n, ok := e.checkRefMap(value); ok {
		e.buffer = encRef(e.buffer, n)
		return nil
	}

	// check whether it should encode the map as class.
	if mm, ok := m.(map[string]interface{}); ok {
		if _, ok = mm[ClassKey]; ok {
			return e.EncodeMapClass(mm)
		}
	}

	value = UnpackPtrValue(value)
	// check nil map
	if value.IsNil() || (value.Kind() == reflect.Ptr && !value.Elem().IsValid()) {
		e.buffer = EncNull(e.buffer)
		return nil
	}

	// if pojo, write class name first
	if p, ok := m.(POJO); ok {
		e.buffer = encByte(e.buffer, BC_MAP)
		e.buffer = encString(e.buffer, p.JavaClassName())
	} else {
		e.buffer = encByte(e.buffer, BC_MAP_UNTYPED)
	}

	keys = value.MapKeys()

	if len(keys) > 0 {
		typ = value.Type().Key()
		for i := 0; i < len(keys); i++ {
			k, err = getMapKey(keys[i], typ)
			if err != nil {
				return perrors.Wrapf(err, "getMapKey(idx:%d, key:%+v)", i, keys[i])
			}
			if err = e.Encode(k); err != nil {
				return perrors.Wrapf(err, "failed to encode map key(idx:%d, key:%+v)", i, keys[i])
			}
			entryValueObj := value.MapIndex(keys[i]).Interface()
			if err = e.Encode(entryValueObj); err != nil {
				return perrors.Wrapf(err, "failed to encode map value(idx:%d, key:%+v, value:%+v)", i, k, entryValueObj)
			}
		}
	}

	e.buffer = encByte(e.buffer, BC_END)

	return nil
}

/////////////////////////////////////////
// Map
/////////////////////////////////////////

// ::= 'M' type (value value)* 'Z'  # key, value map pairs
// ::= 'H' (value value)* 'Z'       # untyped key, value
func (d *Decoder) decMapByValue(value reflect.Value) error {
	var (
		tag        byte
		err        error
		entryKey   interface{}
		entryValue interface{}
	)

	// tag, _ = d.readBufByte()
	tag, err = d.ReadByte()
	// check error
	if err != nil {
		return perrors.WithStack(err)
	}

	switch tag {
	case BC_NULL:
		// null map tag check
		return nil
	case BC_REF:
		refObj, decErr := d.decRef(int32(tag))
		if decErr != nil {
			return perrors.WithStack(decErr)
		}
		SetValue(value, EnsurePackValue(refObj))
		return nil
	case BC_MAP:
		d.decString(TAG_READ) // read map type , ignored
	case BC_MAP_UNTYPED:
		// do nothing
	default:
		return perrors.Errorf("expect map header, but get %x", tag)
	}

	m := reflect.MakeMap(UnpackPtrType(value.Type()))
	// pack with pointer, so that to ref the same map
	m = PackPtr(m)
	d.appendRefs(m)

	// read key and value
	for {
		entryKey, err = d.DecodeValue()
		if err != nil {
			// EOF means the end flag 'Z' of map is already read
			if perrors.Is(err, io.EOF) {
				break
			} else {
				return perrors.WithStack(err)
			}
		}
		if entryKey == nil {
			break
		}
		entryValue, err = d.DecodeValue()
		// fix: check error
		if err != nil {
			return perrors.WithStack(err)
		}

		// add a layer of conversion to make the map compatible with more types during decoding
		key := EnsurePackValue(entryKey)
		if mKey := m.Elem().Type().Key(); key.Type().ConvertibleTo(mKey) {
			key = key.Convert(mKey)
		}
		val := EnsureRawValue(entryValue)
		if mVal := m.Elem().Type().Elem(); val.Type().ConvertibleTo(mVal) {
			val = val.Convert(mVal)
		}
		m.Elem().SetMapIndex(key, val)
	}

	SetValue(value, m)

	return nil
}

// decode map object
func (d *Decoder) decMap(flag int32) (interface{}, error) {
	var (
		err        error
		tag        byte
		ok         bool
		m          map[interface{}]interface{}
		k          interface{}
		v          interface{}
		instValue  reflect.Value
		fieldName  string
		fieldValue reflect.Value
		typ        reflect.Type
	)

	if flag != TAG_READ {
		tag = byte(flag)
	} else {
		tag, _ = d.ReadByte()
	}

	switch {
	case tag == BC_NULL:
		return nil, nil
	case tag == BC_REF:
		return d.decRef(int32(tag))
	case tag == BC_MAP:
		if typ, err = d.decMapType(); err != nil {
			return nil, err
		}

		if typ.Kind() == reflect.Map {
			// Estimate map size by peeking ahead
			size := 16 // Default size if we can't estimate

			// Try to estimate map size by examining the next few bytes
			// This helps reduce map resizing during decoding
			if d.len() >= 4 {
				peek := d.peek(4)
				if len(peek) >= 2 {
					// Check if the next byte looks like a length indicator
					if peek[0] >= BC_INT_ZERO && peek[0] <= INT_DIRECT_MAX {
						// Simple case: small integer directly encoded in the tag
						size = int(peek[0] - BC_INT_ZERO)
					} else if peek[0] == BC_INT && len(peek) >= 4 {
						// Integer case: read the next 4 bytes as an int32
						// This is just an estimate, so we don't need to be exact
						size = int(peek[1])<<16 | int(peek[2])<<8 | int(peek[3])
						if size < 0 || size > 10000 {
							// Sanity check - if the size looks wrong, use default
							size = 16
						}
					}
				}
			}

			instValue = reflect.MakeMapWithSize(typ, size)
		} else {
			instValue = reflect.New(typ).Elem()
		}

		d.appendRefs(instValue)

		for d.peekByte() != BC_END {
			k, err = d.Decode()
			if err != nil {
				return nil, err
			}
			v, err = d.Decode()
			if err != nil {
				return nil, err
			}

			if typ.Kind() == reflect.Map {
				instValue.SetMapIndex(reflect.ValueOf(k), EnsureRawValue(v))
			} else {
				fieldName, ok = k.(string)
				if !ok {
					return nil, perrors.Errorf("the type of map key must be string, but get %v", k)
				}
				fieldValue = instValue.FieldByName(fieldName)
				if fieldValue.IsValid() {
					fieldValue.Set(EnsureRawValue(v))
				}
			}
		}
		_, err = d.ReadByte()
		if err != nil {
			return nil, perrors.WithStack(err)
		}
		return instValue.Interface(), nil
	case tag == BC_MAP_UNTYPED:
		// Estimate map size to avoid resizing
		size := 16 // Default size

		// Try to peek ahead to estimate size
		if d.len() >= 2 {
			peek := d.peek(2)
			if len(peek) >= 1 {
				// Simple heuristic: if next byte is a small integer,
				// it might be indicating the number of entries
				if peek[0] >= BC_INT_ZERO && peek[0] <= INT_DIRECT_MAX {
					size = int(peek[0] - BC_INT_ZERO)
					// Sanity check
					if size < 1 {
						size = 16
					}
				}
			}
		}

		// Pre-allocate map with the estimated size
		m = make(map[interface{}]interface{}, size)
		d.appendRefs(m)
		for d.peekByte() != BC_END {
			k, err = d.Decode()
			if err != nil {
				return nil, err
			}
			v, err = d.Decode()
			if err != nil {
				return nil, err
			}
			m[k] = v
		}
		_, err = d.ReadByte()
		if err != nil {
			return nil, perrors.WithStack(err)
		}
		return m, nil

	default:
		return nil, perrors.Errorf("illegal map type tag:%+v", tag)
	}
}
