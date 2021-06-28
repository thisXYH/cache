package cache

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// key = nil的时候，使用的代替表示。
	nilKeyString = "@.nil.@*"
)

// Operation 缓存操作对象。
type Operation struct {
	// 缓存key分三段 <CacheNamespace>:<Prefix>[:unique flag]
	cacheNamespace string

	// KeyBase = <CacheNamespace>:<Prefix>
	keyBase string

	// cacheProvider 缓存提供者
	cacheProvider CacheProvider

	// 过期时间。
	expireTime *ExpireTime

	// [:unique flag] 部分的拼接元素的个数。
	// 不支持的拼接类型：Complex64, Complex128, Array, Chan, Func ,Interface, Map, Slice, Struct, UnsafePointer
	uniqueFlagLen int
}

// NewOperation 创建一个缓存操作对象,
// 缓存key分三段 <CacheNamespace>:<Prefix>[:unique flag]
// expireTime : 过期时长， nil或者CacheExpirationZero 表不过期。
// uniqueFlagLen : 指定用来拼接[:unique flag]部分的元素个数。
func NewOperation(cacheNamespace, keyPrefix string, uniqueFlagLen int, cacheProvider CacheProvider, expireTime *ExpireTime) *Operation {
	if cacheNamespace == "" || keyPrefix == "" {
		panic(fmt.Errorf(`neither 'cacheNamespace' nor 'keyPrefix' can be zero value`))
	}

	if cacheProvider == nil {
		panic(fmt.Errorf(`'cacheProvider' must not be nil`))
	}

	if uniqueFlagLen < 0 {
		panic(fmt.Errorf(`'uniqueFlagLen' must not be letter than 0`))
	}

	cp := &Operation{}
	cp.cacheNamespace = cacheNamespace
	cp.keyBase = cacheNamespace + ":" + keyPrefix
	cp.cacheProvider = cacheProvider

	if expireTime == nil {
		cp.expireTime = CacheExpirationZero
	} else {
		cp.expireTime = expireTime
	}

	cp.uniqueFlagLen = uniqueFlagLen

	return cp
}

// Key 获取指定key的缓存操作对象。
func (c *Operation) Key(keys ...interface{}) *KeyOperation {
	if len(keys) != c.uniqueFlagLen {
		panic(fmt.Errorf("param 'keys' len(%d)  != uniqueFlagLen(%d)", len(keys), c.uniqueFlagLen))
	}

	return &KeyOperation{
		cp:  c,
		Key: c.buildCacheKey(keys...),
	}
}

// buildCacheKey 构建缓存key。
func (c *Operation) buildCacheKey(keys ...interface{}) string {
	if len(keys) == 0 {
		return c.keyBase // key：没有 [:unique flag]
	}
	sb := strings.Builder{}
	sb.WriteString(c.keyBase)

	for _, v := range keys {
		sb.WriteString("_")
		sb.WriteString(c.oneKeyToStr(v))
	}

	return sb.String()
}

/*
	不支持的type
		Complex64
		Complex128
		Array
		Chan
		Func
		Interface
		Map
		Slice
		Struct
		UnsafePointer
*/
func (c *Operation) oneKeyToStr(v interface{}) string {
	v = c.indirect(v)
	if v == nil { //空值替代。
		return nilKeyString
	}
	vs := ""

	switch s := v.(type) {
	case string:
		vs = s
	case bool:
		if s {
			vs = "1"
		} else {
			vs = "0"
		}
	case time.Time:
		vs = UnixTime(s).String() //毫秒级时间戳

	case int:
		vs = strconv.FormatInt(int64(s), 10)
	case int64:
		vs = strconv.FormatInt(s, 10)
	case int32:
		vs = strconv.FormatInt(int64(s), 10)
	case int16:
		vs = strconv.FormatInt(int64(s), 10)
	case int8:
		vs = strconv.FormatInt(int64(s), 10)
	case uint:
		vs = strconv.FormatUint(uint64(s), 10)
	case uint64:
		vs = strconv.FormatUint(s, 10)
	case uint32:
		vs = strconv.FormatUint(uint64(s), 10)
	case uint16:
		vs = strconv.FormatUint(uint64(s), 10)
	case uint8:
		vs = strconv.FormatUint(uint64(s), 10)
	case float64:
		vs = strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		vs = strconv.FormatFloat(float64(s), 'f', -1, 32)

	case fmt.Stringer:
		vs = s.String()

	default:
		panic(fmt.Errorf("can't String %s (implement fmt.Stringer)", reflect.TypeOf(v).Name()))
	}

	return vs
}

// indirect 移除间接引用。
func (*Operation) indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}