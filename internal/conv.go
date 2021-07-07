// Package conv provides a group of functions to convert between primitive types, maps, slices and structs.
package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	// StringSliceSep is used by SplitString() as the separator.
	StringSliceSep = "~"
)

var (
	intSize = int(unsafe.Sizeof(0)) * 8 // How many bits is an int.
	typTime = reflect.TypeOf(time.Time{})
)

var primitiveTypes = map[reflect.Kind]bool{
	reflect.Bool:       false,
	reflect.Int8:       false,
	reflect.Int16:      false,
	reflect.Int32:      false,
	reflect.Int64:      false,
	reflect.Int:        false,
	reflect.Uint8:      false,
	reflect.Uint16:     false,
	reflect.Uint32:     false,
	reflect.Uint64:     false,
	reflect.Uint:       false,
	reflect.Float32:    false,
	reflect.Float64:    false,
	reflect.Complex64:  false,
	reflect.Complex128: false,
	reflect.String:     false,
}

// IsPrimitiveKind returns true if the given Kind is bool/int*/uint*/float*/complex*/string .
func IsPrimitiveKind(k reflect.Kind) bool {
	_, ok := primitiveTypes[k]
	return ok
}

// IsPrimitiveType returns true if the given type is bool/int*/uint*/float*/complex*/string .
func IsPrimitiveType(t reflect.Type) bool {
	return t != nil && IsPrimitiveKind(t.Kind())
}

// IsSimpleType returns true if the given type IsPrimitiveType() or is time.Time .
func IsSimpleType(t reflect.Type) bool {
	return t != nil && (IsPrimitiveType(t) || t == typTime)
}

// Conv provides a group of functions to convert between primitive types, time.Time, maps, slices and structs.
// A new instance with default values has default conversion behavior.
//   Conv{}.ConvertType(...)
//
// Don't call functions that does not start with 'Convert' directly, they are for configuration and are called internally
// by other functions with names which start with 'Convert', such as Convert()/ConvertType()/ConvertStringToSlice() .
//
type Conv struct {
	// SplitString is the function used to split the string into elements of the slice, when converting a string to a slice.
	// Set this field if need to customize the procedure.
	// If this field is nil, the function DefaultSplitString() will be used.
	SplitString func(v string) []string

	// NameIndexer is the function used to match names when converting from map to struct or from struct to struct.
	// If the given name is match, the function returns the value from the source map with @ok=true;
	// otherwise returns (nil, false) .
	// If it returns OK, the value from the source map will be converted into the destination struct
	// using Conv.ConvertType() .
	//
	// When converting a map to a struct, each field name of the struct will be indexed using this function.
	// When converting a struct to another, field names and values from the souce struct will be put into a map,
	// then each field name of the destination struct will be indexed with the map.
	//
	// If this function is nil, the Go built-in indexer for maps will be used.
	// The build-in indexer is like:
	//   v, ok := m[name]
	//
	// If a case-insensitive indexer is needed, use the CaseInsensitiveNameIndexer function.
	//
	NameIndexer func(m map[string]interface{}, name string) (v interface{}, ok bool)

	// TimeToString formats the given time.
	// Set this field if need to customize the procedure.
	// If this field is nil, the function DefaultTimeToString() will be used.
	TimeToString func(t time.Time) (string, error)

	// StringToTime parses the given string and returns the time it represends.
	// Set this field if need to customize the procedure.
	// If this field is nil, the function DefaultStringToTime() will be used.
	StringToTime func(v string) (time.Time, error)
}

// DefaultSplitString spilit a string by '~'. This is the default value for Conv.SplitString() when it is nil.
func DefaultSplitString(v string) []string {
	return strings.Split(v, StringSliceSep)
}

// DefaultTimeToString formats time using the time.RFC3339 format.
func DefaultTimeToString(t time.Time) (string, error) {
	return t.Format(time.RFC3339), nil
}

// DefaultStringToTime parses the time using the time.RFC3339Nano format.
func DefaultStringToTime(v string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, v)
}

// CaseInsensitiveNameIndexer indexes a map and compares the keys case-insensitively.
// It compares keys with strings.EqualFold() , and returns on the first key for which EqualFold() is true.
func CaseInsensitiveNameIndexer(m map[string]interface{}, key string) (value interface{}, ok bool) {
	// No build-in method to index a map case-insensitively, we just iterate all keys.
	for k, v := range m {
		if strings.EqualFold(key, k) {
			value = v
			ok = true
			return
		}
	}

	return
}

// ConvertStringToSlice converts a string to a slice.
// The elements of the slice must be simple type, for which IsSimpleType() returns true.
//
// Conv.SplitString() is used to split the string.
//
func (c Conv) ConvertStringToSlice(v string, simpleSliceType reflect.Type) (interface{}, error) {
	const fnName = "ConvertStringToSlice"

	if simpleSliceType.Kind() != reflect.Slice {
		return nil, buildError(fnName, "the destiniation type must be slice, got %v", simpleSliceType)
	}

	elemTyp := simpleSliceType.Elem()
	if !IsSimpleType(elemTyp) {
		return nil, buildError(fnName, "cannot convert from string to %v, the element's type must be a simple type", simpleSliceType)
	}

	var parts []string
	if c.SplitString == nil {
		parts = DefaultSplitString(v)
	} else {
		parts = c.SplitString(v)
	}

	dst := reflect.MakeSlice(simpleSliceType, 0, len(parts))
	elemKind := elemTyp.Kind()

	for i, elemIn := range parts {
		elemOut, err := c.ConvertStringToPrimitive(elemIn, elemKind)
		if err != nil {
			return nil, buildError(fnName, "cannot convert to %v, at index %v: %v", simpleSliceType, i, err)
		}

		dst = reflect.Append(dst, reflect.ValueOf(elemOut))
	}

	return dst.Interface(), nil
}

// ConverToBool converts the value to bool.
//
// Rules:
// nil: as false;
// numbers/time.Time: zero as false, non-zero as true;
// string: same as strconv.ParseBool() ;
// other values are not supported, the function returns false and an error.
func (c Conv) ConverToBool(v interface{}) (bool, error) {
	const fnName = "ConverToBool"

	if v == nil {
		return false, nil
	}

	switch vv := v.(type) {
	case bool:
		return vv, nil

	case string:
		res, err := strconv.ParseBool(vv)
		if err == nil {
			return res, nil
		}
		return res, buildError(fnName, err.Error())

	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128,
		time.Time:
		return !reflect.ValueOf(vv).IsZero(), nil

	default:
		return false, buildError(fnName, "cannot convert %v to bool", reflect.TypeOf(v))
	}
}

// ConvertStringToPrimitive converts a string to a primitive type (which IsPrimitiveType() returns true).
func (c Conv) ConvertStringToPrimitive(v string, dstKind reflect.Kind) (interface{}, error) {
	const fnName = "ConvertStringToPrimitive"

	switch dstKind {
	default:
		return nil, buildError(fnName, "%v is not a primitive type", dstKind)

	case reflect.Bool:
		return strconv.ParseBool(v)

	case reflect.Int:
		res, err := strconv.ParseInt(v, 0, intSize)
		return int(res), err

	case reflect.Int8:
		res, err := strconv.ParseInt(v, 0, 8)
		return int8(res), err

	case reflect.Int16:
		res, err := strconv.ParseInt(v, 0, 16)
		return int16(res), err

	case reflect.Int32:
		res, err := strconv.ParseInt(v, 0, 32)
		return int32(res), err

	case reflect.Int64:
		return strconv.ParseInt(v, 0, 64)

	case reflect.Uint:
		res, err := strconv.ParseUint(v, 0, intSize)
		return uint(res), err

	case reflect.Uint8:
		res, err := strconv.ParseUint(v, 0, 8)
		return uint8(res), err

	case reflect.Uint16:
		res, err := strconv.ParseUint(v, 0, 16)
		return uint16(res), err

	case reflect.Uint32:
		res, err := strconv.ParseUint(v, 0, 32)
		return uint32(res), err

	case reflect.Uint64:
		return strconv.ParseUint(v, 0, 64)

	case reflect.Float32:
		res, err := strconv.ParseFloat(v, 32)
		return float32(res), err

	case reflect.Float64:
		return strconv.ParseFloat(v, 64)

	case reflect.Complex64:
		res, err := strconv.ParseComplex(v, 64)
		return complex64(res), err

	case reflect.Complex128:
		return strconv.ParseComplex(v, 128)

	case reflect.String:
		return v, nil
	}
}

// ConvertSimpleToString converts the given value to a string.
// The value must be a simple type, for which IsSimpleType() returns true.
//
// Conv.StringToTime() is used to format times.
// Specially, booleans are converted to 0/1, not the default foramt true/false.
func (c Conv) ConvertSimpleToString(v interface{}) (string, error) {
	const fnName = "ConvertSimpleToString"

	if v == nil {
		return "", sourceShouldNotBeNilError(fnName)
	}

	t := reflect.TypeOf(v)
	if t == typTime {
		if c.TimeToString == nil {
			return DefaultTimeToString(v.(time.Time))
		}

		return c.TimeToString(v.(time.Time))
	}

	k := t.Kind()
	if !IsPrimitiveKind(k) {
		return "", buildError(fnName, "cannot convert %v to any primitive value", k)
	}

	switch k {
	case reflect.Bool:
		// The default string representation for bools are true/false, which is not compatiable
		// to other number types. To increase compatibility, we use 0/1 instead.
		if v.(bool) {
			return "1", nil
		} else {
			return "0", nil
		}

	case reflect.String:
		return v.(string), nil

	default:
		// Use the default formats for other types.
		return fmt.Sprint(v), nil
	}
}

// ConvertSliceToSlice converts a slice to another slice.
// Each element will be converted using Conv.ConvertType() .
// If the source value is nil, returns nil and an error.
func (c Conv) ConvertSliceToSlice(src interface{}, dstSliceTyp reflect.Type) (interface{}, error) {
	const fnName = "ConvertSliceToSlice"

	if src == nil {
		return nil, sourceShouldNotBeNilError(fnName)
	}

	vsrcSlice := reflect.ValueOf(src)
	if vsrcSlice.Kind() != reflect.Slice {
		return nil, buildError(fnName, "src must be a slice, got %v", vsrcSlice.Kind())
	}

	if dstSliceTyp.Kind() != reflect.Slice {
		return nil, buildError(fnName, "the destination type must be slice, got %v", dstSliceTyp.Kind())
	}

	srcLen := vsrcSlice.Len()
	dstElemTyp := dstSliceTyp.Elem()
	vdstSlice := reflect.MakeSlice(dstSliceTyp, 0, srcLen)

	for i := 0; i < srcLen; i++ {
		vsrcElem := vsrcSlice.Index(i)
		srcElmen := vsrcElem.Interface()
		vdstElem, err := c.ConvertType(srcElmen, dstElemTyp)
		if err != nil {
			return nil, buildError(fnName, "cannot convert to %v, at index %v : %v", dstSliceTyp, i, err.Error())
		}

		vdstSlice = reflect.Append(vdstSlice, reflect.ValueOf(vdstElem))
	}

	return vdstSlice.Interface(), nil
}

// ConvertMapToStruct converts a map[string]interface{} to a struct.
//
// Each exported field of the struct is indexed from the map by name using Conv.NameIndexer() , if the name exists,
// the corresponding value is converted using Conv.ConvertType() .
//
func (c Conv) ConvertMapToStruct(m map[string]interface{}, typ reflect.Type) (interface{}, error) {
	const fnName = "ConvertMapToStruct"

	if m == nil {
		return nil, sourceShouldNotBeNilError(fnName)
	}

	k := typ.Kind()
	if k != reflect.Struct {
		return nil, buildError(fnName, "the destination type must be stuct, got %v", typ)
	}

	dst := reflect.New(typ).Elem()

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		v, ok := c.indexNameFromMap(m, f.Name)
		if !ok {
			continue
		}

		// Ignore all unexported fields, they can't be set.
		dstF := dst.Field(i)
		if !dstF.CanSet() {
			continue
		}

		fv, err := c.ConvertType(v, f.Type)
		if err != nil {
			return nil, buildError(fnName, "error on converting field '%v': %v", f.Name, err.Error())
		}

		dstF.Set(reflect.ValueOf(fv))
	}

	return dst.Interface(), nil
}

func (c Conv) indexNameFromMap(m map[string]interface{}, k string) (interface{}, bool) {
	if c.NameIndexer == nil {
		v, ok := m[k]
		return v, ok
	}

	return c.NameIndexer(m, k)
}

// ConvertMapToMap converts a map to another map.
// If the source value is nil, the function returns a nil map of the destination type without any error.
//
// All keys and values in the map are converted using Conv.ConvertType() .
//
func (c Conv) ConvertMapToMap(m interface{}, typ reflect.Type) (interface{}, error) {
	const fnName = "ConvertMapToMap"

	src := reflect.ValueOf(m)
	if src.Kind() != reflect.Map {
		return nil, buildError(fnName, "the given value type must be a map, got %v", src.Kind())
	}

	if typ.Kind() != reflect.Map {
		return nil, buildError(fnName, "the destination type must be map, got %v", typ)
	}

	if src.IsNil() {
		return reflect.Zero(typ).Interface(), nil
	}

	dst := reflect.MakeMap(typ)
	dstKeyType := typ.Key()
	dstValueType := typ.Elem()
	iter := src.MapRange()

	for iter.Next() {
		srcKey := iter.Key().Interface()
		dstKey, err := c.ConvertType(srcKey, dstKeyType)
		if err != nil {
			return nil, buildError(fnName, "cannnot covert key '%v' to %v: %v", srcKey, dstKeyType, err.Error())
		}

		srcVal := iter.Value().Interface()
		dstVal, err := c.ConvertType(srcVal, dstValueType)
		if err != nil {
			return nil, buildError(fnName, "cannnot covert value of key '%v' to %v: %v", srcKey, dstValueType, err.Error())
		}

		dst.SetMapIndex(reflect.ValueOf(dstKey), reflect.ValueOf(dstVal))
	}

	return dst.Interface(), nil
}

// ConvertStructToMap is like json.Unmashal(json.Marshal(v), &someMap) . It converts a struct to map[string]interface{} .
//
// Each value of exported field will be processed recursively with an internal function f() , which:
//  - Simple types (which IsSimpleType() returns true) will be cloned into the map directly.
//  - Slices:
//    - A nil/empty slices is converted to an empty slice with cap=0.
//    - A non-empty slice is converted to another slice, each element is process with f() , all elements must be the same type.
//  - Maps:
//    - A nil map are converted to nil of map[string]interface{} .
//    - A non-nil map is converted to map[string]interface{} , keys are processed with Conv.ConvertType() , values with f() .
//  - Structs are converted to map[string]interface{} using Conv.ConvertStructToMap() .
//  - For pointers, the values pointed to are converted with f() .
// Other types not listed are not supported and will result in an error.
//
func (c Conv) ConvertStructToMap(v interface{}) (map[string]interface{}, error) {
	const fnName = "ConvertStructToMap"

	if v == nil {
		return nil, sourceShouldNotBeNilError(fnName)
	}

	srcType := reflect.TypeOf(v)
	if srcType.Kind() != reflect.Struct {
		return nil, buildError(fnName, "the given value must be a struct, got %v", srcType)
	}

	src := reflect.ValueOf(v)
	dst := reflect.MakeMap(reflect.TypeOf(map[string]interface{}(nil)))

	for i := 0; i < src.NumField(); i++ {
		fieldValue := src.Field(i)

		// Ignore unexported fields.
		if !fieldValue.CanInterface() {
			continue
		}

		fieldName := srcType.Field(i).Name
		ff, err := c.convertToMapValue(fieldValue)

		if err != nil {
			return nil, buildError(fnName, "error on converting field %v: %v", fieldName, err.Error())
		}

		dst.SetMapIndex(reflect.ValueOf(fieldName), ff)
	}

	return dst.Interface().(map[string]interface{}), nil
}

func (c Conv) convertToMapValue(fv reflect.Value) (reflect.Value, error) {
	for fv.Kind() == reflect.Ptr && !fv.IsNil() {
		fv = fv.Elem()
	}

	switch fv.Kind() {
	case reflect.Struct:
		v, err := c.ConvertStructToMap(fv.Interface())
		if err != nil {
			return reflect.Value{}, err
		}

		return reflect.ValueOf(v), nil

	case reflect.Slice:
		var newSlice reflect.Value

		switch {
		case fv.IsNil() || fv.Len() == 0:
			ft := fv.Type()
			sliceType, ok := c.dertermineSliceTypeForMapValue(ft)
			if !ok {
				return reflect.Value{}, fmt.Errorf("cannot convert %v", fv.Type())
			}

			newSlice = reflect.MakeSlice(sliceType, 0, 0)

		default:
			for i := 0; i < fv.Len(); i++ {
				oldVal := fv.Index(i)
				newVal, err := c.convertToMapValue(oldVal)
				if err != nil {
					return reflect.Value{}, fmt.Errorf("index %v: %v", i, err.Error())
				}

				// Lazy initialization. The slice type depends on the type of the first element.
				if i == 0 {
					newSlice = reflect.MakeSlice(reflect.SliceOf(newVal.Type()), 0, fv.Len())
				}

				newSlice = reflect.Append(newSlice, newVal)
			}
		}

		return newSlice, nil

	case reflect.Map:
		newMap := reflect.MakeMap(reflect.TypeOf(map[string]interface{}(nil)))
		iter := fv.MapRange()
		for iter.Next() {
			oldKey := iter.Key()
			oldVal := iter.Value()

			var newKey string
			err := c.Convert(oldKey.Interface(), &newKey)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("key %v: %v", oldKey, err.Error())
			}

			newVal, err := c.convertToMapValue(oldVal)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("value of key %v: %v", newKey, err.Error())
			}

			newMap.SetMapIndex(reflect.ValueOf(newKey), newVal)
		}
		return newMap, nil

	case reflect.Interface:
		// Extract the underlying value.
		v := fv.Interface()
		if v == nil {
			return reflect.ValueOf(nil), nil
		}

		fv = reflect.ValueOf(v)
		return c.convertToMapValue(fv)

	default:
		if !IsSimpleType(fv.Type()) {
			return reflect.Value{}, fmt.Errorf("must be a simple type, got %v", fv.Kind())
		}
		return fv, nil
	}
}

func (c Conv) dertermineSliceTypeForMapValue(srcSliceType reflect.Type) (dstSliceType reflect.Type, ok bool) {
	elemType := srcSliceType.Elem()
	if IsSimpleType(elemType) {
		dstSliceType = srcSliceType
		ok = true
		return
	}

	elemKind := elemType.Kind()
	switch elemKind {
	case reflect.Map, reflect.Struct:
		dstSliceType = reflect.SliceOf(reflect.TypeOf(map[string]interface{}(nil)))
		ok = true
		return

	case reflect.Slice:
		innerSliceType, innerOK := c.dertermineSliceTypeForMapValue(elemType)
		if !innerOK {
			return
		}

		dstSliceType = reflect.SliceOf(innerSliceType)
		ok = true
		return

	default:
		ok = false
		return
	}
}

// ConvertStructToStruct converts a struct to another.
// If the given value is nil, returns nil and an error.
//
// When converting, all fields of the source struct is to be stored in a map[string]interface{} ,
// then each field of the destination struct is indexed from the map by name using Conv.NameIndexer() ,
// if the name exists, the value is converted using Conv.ConvertType() .
//
// This function can be used to deep-clone a struct.
//
func (c Conv) ConvertStructToStruct(v interface{}, typ reflect.Type) (interface{}, error) {
	const fnName = "ConvertStructToStruct"

	if v == nil {
		return nil, sourceShouldNotBeNilError(fnName)
	}

	dstKind := typ.Kind()
	if dstKind != reflect.Struct {
		return nil, buildError(fnName, "the destination type must be struct, got %v", dstKind)
	}

	srcTyp := reflect.TypeOf(v)
	if srcTyp.Kind() != reflect.Struct {
		return nil, buildError(fnName, "the given value must be a struct, got %v", srcTyp)
	}

	m := make(map[string]interface{})
	src := reflect.ValueOf(v)
	for i := 0; i < src.NumField(); i++ {
		f := src.Field(i)
		if !f.CanInterface() {
			continue
		}

		fName := srcTyp.Field(i).Name
		m[fName] = f.Interface()
	}

	dst := reflect.New(typ).Elem()
	for i := 0; i < dst.NumField(); i++ {
		fType := typ.Field(i)
		v, ok := c.indexNameFromMap(m, fType.Name)
		if !ok {
			continue
		}

		f := dst.Field(i)
		if !f.CanSet() {
			continue
		}

		dstValue, err := c.ConvertType(v, fType.Type)
		if err != nil {
			return nil, buildError(fnName, "error on converting field %v: %v", fType.Name, err.Error())
		}

		f.Set(reflect.ValueOf(dstValue))
	}

	return dst.Interface(), nil
}

// ConvertType is the core function of Conv . It converts the given value to the destination type.
//
// Currently these conversions are supported:
//   simple         -> simple                 * use Conv.ConvertSimpleToSimple()
//   string         -> simple                 * use Conv.ConvertStringToPrimitive(), or Conv.StringToTime() for time values.
//   string         -> []simple               * use Conv.ConvertStringToSlice()
//   map[string]any -> struct                 * use Conv.ConvertMapToStruct()
//   map[any]any    -> map[any]any            * use Conv.ConvertMapToMap()
//   []any          -> []any                  * use Conv.ConvertType() recursively
//   struct         -> map[string]interface{} * use Conv.ConvertStructToMap()
//   struct         -> struct                 * use Conv.ConvertStructToStuct()
// 'any' generally means interface{} .
//
// typ can be a type of pointer, the conversion of the underlying type must be supported.
//
// This function can be used to deep-clone a struct, e.g.
//   clone, err := ConvertType(src, reflect.TypeOf(src))
//
func (c Conv) ConvertType(v interface{}, typ reflect.Type) (interface{}, error) {
	const fnName = "ConvertType"

	dstKind := typ.Kind()
	if v == nil {
		if canBeNil(dstKind) {
			return nil, nil
		}

		return nil, buildError(fnName, "cannot convert nil to %v", typ)
	}

	// Try to get the underlying type from a pointer type.
	// It may be a pointer to another pointer...
	ptrStack := make([]reflect.Type, 0)
	for typ.Kind() == reflect.Ptr {
		ptrStack = append(ptrStack, typ)
		typ = typ.Elem()
	}

	dst, err := c.convertToNonPtr(v, typ)
	if err != nil {
		return nil, buildError(fnName, err.Error())
	}

	// Convert to pointer if needed.
	if len(ptrStack) > 0 {
		var prev, current reflect.Value
		for i := len(ptrStack) - 1; i >= 0; i-- {
			if i == len(ptrStack)-1 {
				prev = reflect.ValueOf(dst)
			} else {
				prev = current
			}

			current = reflect.New(ptrStack[i])
			current.Elem().Set(prev)
		}

		dst = current.Interface()
	}

	return dst, nil
}

// Convert is like Conv.ConvertType() , but receives a pointer instead of a type.
// It stores the result in the value pointed to by dst.
// If dst is not a pointer, the function panics an error.
func (c Conv) Convert(src interface{}, dst interface{}) error {
	const fnName = "Convert"

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		panic(buildError(fnName, "the destination value must be a pointer"))
	}

	if src == nil {
		return nil
	}

	for dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
	}

	value, err := c.convertToNonPtr(src, dstValue.Type())
	if err != nil {
		return buildError(fnName, err.Error())
	}

	dstValue.Set(reflect.ValueOf(value))
	return nil
}

func (c Conv) convertToNonPtr(v interface{}, dstTyp reflect.Type) (interface{}, error) {
	srcTyp := reflect.TypeOf(v)
	srcKind := srcTyp.Kind()
	dstKind := dstTyp.Kind()
	if IsPrimitiveKind(srcKind) && IsPrimitiveKind(dstKind) {
		return c.convertPrimitiveToPrimitive(v, dstKind)
	}

	if srcTyp == typTime {
		tm := v.(time.Time)

		switch {
		case dstTyp == typTime:
			return tm, nil

		case dstKind == reflect.String:
			return c.TimeToString(tm)

		case IsPrimitiveKind(dstKind):
			timestamp := tm.Unix()
			res, err := c.convertPrimitiveToPrimitive(timestamp, dstKind)
			if err != nil {
				return nil, fmt.Errorf("failed on converting timestamp %v to %v: %v", timestamp, dstKind, err.Error())
			}
			return res, nil
		}
	} else if srcKind == reflect.Map {
		// map[string]any { "_": value } -> Convert(value)
		if underlyingValue := c.tryFlattenEmptyKeyMap(v); underlyingValue != nil {
			return c.ConvertType(underlyingValue, dstTyp)
		}

		switch dstKind {
		// map -> map
		case reflect.Map:
			return c.ConvertMapToMap(v, dstTyp)

		// map[string]any -> struct
		case reflect.Struct:
			mm, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("when converting to struct, the given map must be map[string]interface{}, got %v", srcTyp)
			}
			return c.ConvertMapToStruct(mm, dstTyp)
		}
	} else if srcKind == reflect.Struct {
		switch dstKind {
		case reflect.Map:
			return c.ConvertStructToMap(v)

		case reflect.Struct:
			return c.ConvertStructToStruct(v, dstTyp)
		}
	} else if dstKind == reflect.Slice {
		switch srcKind {
		// string -> []primitive
		case reflect.String:
			return c.ConvertStringToSlice(v.(string), dstTyp)

		// []any -> []any
		case reflect.Slice:
			return c.ConvertSliceToSlice(v, dstTyp)
		}
	}

	return nil, fmt.Errorf("cannot convert %v to %v", srcTyp, dstTyp)
}

// tryFlattenEmptyKeyMap check the value.
// When:
//   - the map is map[string]interface{}
//   - the map has only one key
//   - the key is an empty string
// the function returns the value of the key; otherwise it returns nil.
//
// Such map is a special contract, it's used when converting a map to a simple type.
//
func (c Conv) tryFlattenEmptyKeyMap(v interface{}) interface{} {
	m, ok := v.(map[string]interface{})
	if !ok || len(m) != 1 {
		return nil
	}

	for k, v := range m {
		if k == "" {
			return v
		}
	}

	return nil
}

func (c Conv) convertPrimitiveToPrimitive(v interface{}, dstKind reflect.Kind) (interface{}, error) {
	// For most primitive types, we use string as a middleware. A value is converted to
	// it's string representaion, then is converted to the destination value.
	// The fmt and strconv package help us to deal with strings. So we can avoid to implement
	// such M*N converstions between different types.
	//
	// bool must be treated separately. A string representation of a number, such as '33', can't be converted
	// to a boolean, strconv.ParseBool() returns an error.

	if dstKind == reflect.Bool {
		return c.ConverToBool(v)
	}

	sv, err := c.ConvertSimpleToString(v)
	if err != nil {
		return nil, err
	}

	dst, err := c.ConvertStringToPrimitive(sv, dstKind)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func canBeNil(k reflect.Kind) bool {
	return k == reflect.Map || k == reflect.Slice || k == reflect.Ptr
}

func buildError(fn, msgFormat string, a ...interface{}) error {
	return errors.New(buildErrorMessage(fn, msgFormat, a...))
}

func buildErrorMessage(fn, msgFormat string, a ...interface{}) string {
	return "conv." + fn + ": " + fmt.Sprintf(msgFormat, a...)
}

func sourceShouldNotBeNilError(fn string) error {
	return buildError(fn, "the source value should not be nil")
}

// ConvertStringToSlice is equivalent to Conv{}.ConvertStringToSlice() .
func ConvertStringToSlice(v string, primitiveSliceType reflect.Type) (interface{}, error) {
	return Conv{}.ConvertStringToSlice(v, primitiveSliceType)
}

// ConverToBool is equivalent to Conv{}.ConverToBool() .
func ConverToBool(v interface{}) (bool, error) {
	return Conv{}.ConverToBool(v)
}

// ConvertStringToPrimitive is equivalent to Conv{}.ConvertStringToPrimitive() .
func ConvertStringToPrimitive(v string, dstKind reflect.Kind) (interface{}, error) {
	return Conv{}.ConvertStringToPrimitive(v, dstKind)
}

// ConvertSimpleToString is equivalent to Conv{}.ConvertSimpleToString() .
func ConvertSimpleToString(v interface{}) (string, error) {
	return Conv{}.ConvertSimpleToString(v)
}

// ConvertSliceToSlice is equivalent to Conv{}.ConvertSliceToSlice() .
func ConvertSliceToSlice(src interface{}, dstSliceTyp reflect.Type) (interface{}, error) {
	return Conv{}.ConvertSliceToSlice(src, dstSliceTyp)
}

// ConvertMapToStruct is equivalent to Conv{}.ConvertMapToStruct() .
func ConvertMapToStruct(m map[string]interface{}, typ reflect.Type) (interface{}, error) {
	return Conv{}.ConvertMapToStruct(m, typ)
}

// ConvertMapToMap is equivalent to Conv{}.ConvertMapToMap() .
func ConvertMapToMap(m interface{}, typ reflect.Type) (interface{}, error) {
	return Conv{}.ConvertMapToMap(m, typ)
}

// ConvertStructToMap is equivalent to Conv{}.ConvertStructToMap() .
func ConvertStructToMap(v interface{}) (map[string]interface{}, error) {
	return Conv{}.ConvertStructToMap(v)
}

// ConvertStructToStruct is equivalent to Conv{}.ConvertStructToStruct() .
func ConvertStructToStruct(v interface{}, typ reflect.Type) (interface{}, error) {
	return Conv{}.ConvertStructToStruct(v, typ)
}

// ConvertType is equivalent to Conv{}.ConvertType() .
func ConvertType(v interface{}, typ reflect.Type) (interface{}, error) {
	return Conv{}.ConvertType(v, typ)
}

// Convert is equivalent to Conv{}.Convert() .
func Convert(src interface{}, dst interface{}) error {
	return Conv{}.Convert(src, dst)
}
