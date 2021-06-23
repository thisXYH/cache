package conv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const (
	// 若将字符串转换为 slice ，字符串会使用此符号切割为多个元素。
	StringSliceDelimiter = "~"
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

// 一个 int 在当前系统占多少位。
var intSize = int(unsafe.Sizeof(0)) * 8

// IsPrimitiveType 判断给定的类型是否是基础类型，基础类型包括所有的数值类型及字符串。
func IsPrimitiveType(k reflect.Kind) bool {
	_, ok := primitiveTypes[k]
	return ok
}

func canBeNil(k reflect.Kind) bool {
	return k == reflect.Map || k == reflect.Slice || k == reflect.Ptr
}

// ConvertStringToSlice 将字符串转换为 slice ，元素必须是基础类型，其通过 StringToPrimitive() 转换。
// 若字符串包含字符 StringSliceDelimiter (~) ，将根据此分隔符切割为多个元素。
// 若目标类型不是基础类型的 slice ，则 panic(error)。
func ConvertStringToSlice(v string, primitiveSliceType reflect.Type) (interface{}, error) {
	const fnName = "ConvertStringToSlice"

	if primitiveSliceType.Kind() != reflect.Slice {
		panic(fmt.Errorf("the destiniation type must be a slice of a primitive, got %v", primitiveSliceType))
	}

	elemKind := primitiveSliceType.Elem().Kind()
	if !IsPrimitiveType(elemKind) {
		panic(fmt.Errorf("cannot convert from string to %v, the element's type must be primitive", primitiveSliceType))
	}

	parts := strings.Split(v, StringSliceDelimiter)
	dst := reflect.MakeSlice(primitiveSliceType, 0, len(parts))

	for elemIdx, elemIn := range parts {
		elemOut, err := ConvertStringToPrimitive(elemIn, elemKind)
		if err != nil {
			return nil, buildError(fnName, "cannot convert to %v, at index %v : %v", primitiveSliceType, elemIdx, err.Error())
		}

		dst = reflect.Append(dst, reflect.ValueOf(elemOut))
	}

	return dst.Interface(), nil
}

// ConverToBool 尝试将给定值转换为 bool 。
//
// 规则如下：
// 数值型：0 转为 false ，其余转为 true ；
// 字符串：同 strconv.ParseBool() ；
// nil 转为 false 。
// 其他情况返回 false, error 。
func ConverToBool(v interface{}) (bool, error) {
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
		float32, float64, complex64, complex128:
		return !reflect.ValueOf(vv).IsZero(), nil

	default:
		return false, buildError(fnName, "cannot convert %v to bool", reflect.TypeOf(v))
	}
}

// ConvertStringToPrimitive 将字符串转为指定的基础类型。
// 若目标类型不是基础类型，则 panic(error) 。
func ConvertStringToPrimitive(v string, dstKind reflect.Kind) (interface{}, error) {
	switch dstKind {
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

	default:
		// 调用此方法必须给定基础类型。
		panic(fmt.Errorf("%v is not a primitive type", dstKind))
	}
}

// ConvertPrimitiveToString 获取给定的基础类型值的字符串形势。
// 若给定值不止基础类型，或给定值为 nil ，则 panic(error)。
// 特别的， bool 会转为字符串的 0/1 而不是 true/false ，以提高转换为其他数值的兼容性； nil 转为空字符串。
func ConvertPrimitiveToString(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}

	k := reflect.TypeOf(v).Kind()
	if !IsPrimitiveType(k) {
		panic(fmt.Errorf("cannot convert %v to any primitive value", k))
	}

	switch k {
	case reflect.Bool:
		// bool 默认字符串为 true/false ，不便于与数字间的转换，改用 0/1 处理。
		if v.(bool) {
			return "1", nil
		} else {
			return "0", nil
		}

	case reflect.String:
		return v.(string), nil

	default:
		// 其余基础类型直接“偷懒”，利用 fmt 完成。
		return fmt.Sprint(v), nil
	}
}

// ConvertSliceToSlice 将一个 slice 转到另一个 slice ，对每个元素使用 Convert 方法。
// 若给定值或目标类型不是 slice ，则 panic(error)。
func ConvertSliceToSlice(src interface{}, dstSliceTyp reflect.Type) (interface{}, error) {
	const fnName = "ConvertSliceToSlice"

	vsrcSlice := reflect.ValueOf(src)

	if vsrcSlice.Kind() != reflect.Slice {
		panic(fmt.Errorf("src must be a slice, got %v", vsrcSlice.Kind()))
	}

	if dstSliceTyp.Kind() != reflect.Slice {
		panic(fmt.Errorf("the destination type must be a slice, got %v", dstSliceTyp.Kind()))
	}

	srcLen := vsrcSlice.Len()
	dstElemTyp := dstSliceTyp.Elem()
	vdstSlice := reflect.MakeSlice(dstSliceTyp, 0, srcLen)

	for i := 0; i < srcLen; i++ {
		vsrcElem := vsrcSlice.Index(i)

		vdstElem, err := Convert(vsrcElem.Interface(), dstElemTyp)
		if err != nil {
			return nil, buildError(fnName, "cannot convert to %v, at index %v : %v", dstSliceTyp, i, err.Error())
		}

		vdstSlice = reflect.Append(vdstSlice, reflect.ValueOf(vdstElem))
	}

	return vdstSlice.Interface(), nil
}

// ConvertMapToStruct 通过给定的 map[string]interface{} 创建 struct ，map 的 key 若与 stuct 的字段同名，则会进行赋值。
// 赋值前使用 Convert 方法转换。
// 若目标类型不是 struct ，则 panic(error)。
func ConvertMapToStruct(m map[string]interface{}, typ reflect.Type) (interface{}, error) {
	const fnName = "ConvertMapToStruct"

	k := typ.Kind()
	if k != reflect.Struct {
		panic(fmt.Errorf("the destination type must be stuct, got %v", typ))
	}

	dst := reflect.New(typ)

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		v, ok := m[f.Name]
		if !ok {
			// 字段没出现就不用赋值了。
			continue
		}

		fv, err := Convert(v, f.Type)
		if err != nil {
			return nil, buildError(fnName, "error on converting field '%v': %v", f.Name, err.Error())
		}

		dst.Field(i).Elem().Set(reflect.ValueOf(fv))
	}

	return dst, nil
}

// ConvertMapToMap 通过给定的 map 创建目标 map 。 key 和 value 的转换都使用 Convert 方法。
// 若给定 map 为 nil ，返回 nil 。
// 若目标类型不是 map ，则 panic(error)。
func ConvertMapToMap(m interface{}, typ reflect.Type) (interface{}, error) {
	const fnName = "ConvertMapToMap"

	src := reflect.ValueOf(m)
	if src.Kind() != reflect.Map {
		panic(fmt.Errorf("the given value type must be a map, got %v", src.Kind()))
	}

	if typ.Kind() != reflect.Map {
		panic(fmt.Errorf("the destination type must be map, got %v", typ))
	}

	dst := reflect.MakeMap(typ)
	dstKeyType := typ.Key()
	dstValueType := typ.Elem()
	iter := src.MapRange()

	for iter.Next() {
		srcKey := iter.Key().Interface()
		dstKey, err := Convert(srcKey, dstKeyType)
		if err != nil {
			return nil, buildError(fnName, "cannnot covert key '%v': %v", srcKey, err.Error())
		}

		srcVal := iter.Value().Interface()
		dstVal, err := Convert(srcVal, dstValueType)
		if err != nil {
			return nil, buildError(fnName, "cannnot covert value of key '%v': %v", srcKey, err.Error())
		}

		dst.SetMapIndex(reflect.ValueOf(dstKey), reflect.ValueOf(dstVal))
	}

	return dst, nil
}

// string      -> primitive/[]primitive
// primitive   -> primitive
// map[any]any -> map[any]any any->any 必须满足对应类型的转换条件。
// []any       -> []any 对于每个元素 any->Convert(any)。
// struct      -> map[string]any/struct
func Convert(v interface{}, typ reflect.Type) (interface{}, error) {
	const fnConvert = "Convert"

	dstKind := typ.Kind()
	if v == nil && canBeNil(dstKind) {
		return nil, nil
	}

	isPtr := dstKind == reflect.Ptr
	if isPtr {
		typ = typ.Elem()
	}

	dst, err := doConvert(v, typ)
	if err != nil {
		return nil, buildError(fnConvert, err.Error())
	}

	if isPtr {
		dst = &dst
	}

	return dst, nil
}

func doConvert(v interface{}, typ reflect.Type) (interface{}, error) {
	dstKind := typ.Kind()
	if IsPrimitiveType(dstKind) {
		return convertPrimitiveToPrimitive(v, dstKind)
	}

	srcTyp := reflect.TypeOf(v)
	srcKind := srcTyp.Kind()

	// map[key]value
	if srcKind == reflect.Map {
		// map[string]any { "_": value } -> Convert(value)
		if underlyingValue := tryFlattenEmptyKeyMap(v); underlyingValue != nil {
			return Convert(underlyingValue, typ)
		}

		switch dstKind {
		// map -> map
		case reflect.Map:
			return ConvertMapToMap(v, typ)

		// map[string]any -> struct
		case reflect.Struct:
			mm, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("when converting to struct, the given map must be map[string]interface{}, got %v", srcTyp)
			}
			return ConvertMapToStruct(mm, typ)

		default:
			return nil, fmt.Errorf("cannot convert %v to %v", srcTyp, typ)
		}
	}

	if dstKind == reflect.Slice {
		switch srcKind {
		// string -> []primitive
		case reflect.String:
			return ConvertStringToSlice(v.(string), typ)

		// []any -> []any
		case reflect.Slice:
			return ConvertSliceToSlice(v, typ)

		default:
			return nil, fmt.Errorf("cannot convert %v to %v", srcTyp, typ)
		}
	}

	return nil, fmt.Errorf("cannot convert %v to %v", srcTyp, typ)
}

// tryFlattenEmptyKeyMap 判断待转换的值是否是 map[string]interface{} ，且是否只有一个 key 且名称为“_”。
// 若是，返回该字段的值；否则返回 nil 。
// 这种 map 是用来包装其他值的，是一个特殊的约定。
func tryFlattenEmptyKeyMap(v interface{}) interface{} {
	m, ok := v.(map[string]interface{})
	if !ok || len(m) != 1 {
		return nil
	}

	for k, v := range m {
		if k == "_" {
			return v
		}
	}

	return nil
}

func convertPrimitiveToPrimitive(v interface{}, dstKind reflect.Kind) (interface{}, error) {
	// 对于 bool 和其他数值类型需单独对待，因为数值转到 string 再转回 bool 会产生不同的结果。
	// 如 33 -> "33" -> error ，而实际上应该是 33 -> true 。
	if dstKind == reflect.Bool {
		return ConverToBool(v)
	}

	// 其余基础类型，统一转字符串再转回来。利用字符串做胶水类型，比较取巧，可以不用实现 M*N 种转换。
	sv, err := ConvertPrimitiveToString(v)
	if err != nil {
		return nil, err
	}

	dst, err := ConvertStringToPrimitive(sv, dstKind)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func buildError(fn, msgFormat string, a ...interface{}) error {
	return errors.New(buildErrorMessage(fn, msgFormat, a...))
}

func buildErrorMessage(fn, msgFormat string, a ...interface{}) string {
	return "conv." + fn + ": " + fmt.Sprintf(msgFormat, a...)
}
