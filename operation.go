package cache

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cmstar/go-conv"
)

// Operation 缓存操作对象。
type Operation struct {
	// 缓存key分三段 <CacheNamespace>:<Prefix>[:unique flag]。
	cacheNamespace string

	// KeyBase = <CacheNamespace>:<Prefix> .
	keyBase string

	// cacheProvider 缓存提供者。
	cacheProvider CacheProvider

	// 过期时间。
	expireTime *Expiration

	// [:unique flag] 部分的拼接元素的个数。
	// 受支持的 [:unique flag] 类型: bool, int*, uint*, float*, string, time.time, UnixTime 。
	uniqueFlagLen int
}

// NewOperation 创建一个缓存操作对象。
// 缓存key分三段 <CacheNamespace>:<Prefix>[:unique flag]。
// expireTime: 过期时长， nil 或者 CacheExpirationZero 表不过期。
// uniqueFlagLen: 指定用来拼接 [:unique flag] 部分的元素个数。
// 受支持的 [:unique flag] 类型: bool, int*, uint*, float*, string, time.time, UnixTime 。
func NewOperation(cacheNamespace, keyPrefix string, uniqueFlagLen int, cacheProvider CacheProvider, expireTime *Expiration) *Operation {
	if cacheNamespace == "" || keyPrefix == "" {
		panic(fmt.Errorf(`neither 'cacheNamespace' nor 'keyPrefix' can be zero value`))
	}

	if cacheProvider == nil {
		panic(fmt.Errorf(`'cacheProvider' must not be nil`))
	}

	if uniqueFlagLen < 0 {
		panic(fmt.Errorf(`'uniqueFlagLen' must not be less than 0`))
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

// Key 获取指定key的缓存 key 操作对象。
//  受支持的 key 类型: bool, int*, uint*, float*, string, time.time, UnixTime 。
func (c *Operation) Key(keys ...interface{}) *KeyOperation {
	if len(keys) != c.uniqueFlagLen {
		panic(fmt.Errorf("param 'keys' len(%d)  != uniqueFlagLen(%d)", len(keys), c.uniqueFlagLen))
	}

	return &KeyOperation{
		p:   c.cacheProvider,
		exp: c.expireTime,
		Key: c.buildCacheKey(keys...),
	}
}

// buildCacheKey 构建缓存key。
func (c *Operation) buildCacheKey(keys ...interface{}) string {
	if len(keys) == 0 {
		return c.keyBase // key：没有 [:unique flag]。
	}
	sb := strings.Builder{}
	sb.WriteString(c.keyBase)

	for _, v := range keys {
		sb.WriteString("_")
		sb.WriteString(oneKeyToStr(v))
	}

	return sb.String()
}

// Operation0 表示 key 只由0个元素组成的缓存操作对象。
type Operation0[TRes any] struct {
	op Operation
}

// NewOperation0 类似 NewOperation ，但创建一个 key 只由0个元素组成的缓存操作对象。
func NewOperation0[TRes any](
	cacheNamespace, keyPrefix string,
	cacheProvider CacheProvider,
	expireTime *Expiration,
) *Operation0[TRes] {
	return &Operation0[TRes]{
		*NewOperation(cacheNamespace, keyPrefix, 0, cacheProvider, expireTime),
	}
}

// Key 获取指定key的缓存操作对象。
func (c *Operation0[TRes]) Key() *GenericKeyOperation[TRes] {
	return &GenericKeyOperation[TRes]{
		p:   c.op.cacheProvider,
		exp: c.op.expireTime,
		Key: c.op.buildCacheKey(),
	}
}

// Operation1 表示 key 只由0个元素组成的缓存操作对象。
type Operation1[TKey, TRes any] struct {
	op Operation
}

// NewOperation1 类似 NewOperation ，但创建一个 key 只由1个元素组成的缓存操作对象。
func NewOperation1[TKey, TRes any](
	cacheNamespace, keyPrefix string,
	cacheProvider CacheProvider,
	expireTime *Expiration,
) *Operation1[TKey, TRes] {
	return &Operation1[TKey, TRes]{
		*NewOperation(cacheNamespace, keyPrefix, 1, cacheProvider, expireTime),
	}
}

// Key 获取指定key的缓存操作对象。
func (c *Operation1[TKey, TRes]) Key(v TKey) *GenericKeyOperation[TRes] {
	return &GenericKeyOperation[TRes]{
		p:   c.op.cacheProvider,
		exp: c.op.expireTime,
		Key: c.op.buildCacheKey(v),
	}
}

// Operation2 表示 key 只由2个元素组成的缓存操作对象。
type Operation2[TKey1, TKey2, TRes any] struct {
	op Operation
}

// NewOperation2 类似 NewOperation ，但创建一个 key 只由2个元素组成的缓存操作对象。
func NewOperation2[TKey1, TKey2, TRes any](
	cacheNamespace, keyPrefix string,
	cacheProvider CacheProvider,
	expireTime *Expiration,
) *Operation2[TKey1, TKey2, TRes] {
	return &Operation2[TKey1, TKey2, TRes]{
		*NewOperation(cacheNamespace, keyPrefix, 2, cacheProvider, expireTime),
	}
}

// Key 获取指定key的缓存操作对象。
func (c *Operation2[TKey1, TKey2, TRes]) Key(v1 TKey1, v2 TKey2) *GenericKeyOperation[TRes] {
	return &GenericKeyOperation[TRes]{
		p:   c.op.cacheProvider,
		exp: c.op.expireTime,
		Key: c.op.buildCacheKey(v1, v2),
	}
}

// Operation3 表示 key 只由3个元素组成的缓存操作对象。
type Operation3[TKey1, TKey2, TKey3, TRes any] struct {
	op Operation
}

// NewOperation3 类似 NewOperation ，但创建一个 key 只由3个元素组成的缓存操作对象。
func NewOperation3[TKey1, TKey2, TKey3, TRes any](
	cacheNamespace, keyPrefix string,
	cacheProvider CacheProvider,
	expireTime *Expiration,
) *Operation3[TKey1, TKey2, TKey3, TRes] {
	return &Operation3[TKey1, TKey2, TKey3, TRes]{
		*NewOperation(cacheNamespace, keyPrefix, 3, cacheProvider, expireTime),
	}
}

// Key 获取指定key的缓存操作对象。
func (c *Operation3[TKey1, TKey2, TKey3, TRes]) Key(v1 TKey1, v2 TKey2, v3 TKey3) *GenericKeyOperation[TRes] {
	return &GenericKeyOperation[TRes]{
		p:   c.op.cacheProvider,
		exp: c.op.expireTime,
		Key: c.op.buildCacheKey(v1, v2, v3),
	}
}

// Operation4 表示 key 只由4个元素组成的缓存操作对象。
type Operation4[TKey1, TKey2, TKey3, TKey4, TRes any] struct {
	op Operation
}

// NewOperation4 类似 NewOperation ，但创建一个 key 只由4个元素组成的缓存操作对象。
func NewOperation4[TKey1, TKey2, TKey3, TKey4, TRes any](
	cacheNamespace, keyPrefix string,
	cacheProvider CacheProvider,
	expireTime *Expiration,
) *Operation4[TKey1, TKey2, TKey3, TKey4, TRes] {
	return &Operation4[TKey1, TKey2, TKey3, TKey4, TRes]{
		*NewOperation(cacheNamespace, keyPrefix, 4, cacheProvider, expireTime),
	}
}

// Key 获取指定key的缓存操作对象。
func (c *Operation4[TKey1, TKey2, TKey3, TKey4, TRes]) Key(v1 TKey1, v2 TKey2, v3 TKey3, v4 TKey4) *GenericKeyOperation[TRes] {
	return &GenericKeyOperation[TRes]{
		p:   c.op.cacheProvider,
		exp: c.op.expireTime,
		Key: c.op.buildCacheKey(v1, v2, v3, v4),
	}
}

/*
	不支持的type
		Array
		Chan
		Func
		Interface
		Map
		Slice
		Struct
		UnsafePointer
*/
func oneKeyToStr(v interface{}) string {
	v = indirect(v)
	if v == nil {
		panic(fmt.Errorf("key flag must not be nil pointer"))
	}
	vs := ""

	// 基础类型。
	if conv.IsPrimitiveType(reflect.TypeOf(v)) {
		conv.Convert(v, &vs)
		return vs
	}

	switch s := v.(type) {
	case time.Time:
		vs = UnixTime(v.(time.Time)).String() //毫秒级时间戳。
	case fmt.Stringer:
		vs = s.String()

	default:
		panic(fmt.Errorf("can't String %s (implement fmt.Stringer)", reflect.TypeOf(v).Name()))
	}

	return vs
}

// indirect 移除间接引用。
func indirect(a interface{}) interface{} {
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
