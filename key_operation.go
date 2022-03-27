package cache

type KeyOperation struct {
	p   CacheProvider
	exp *Expiration

	// 缓存key。
	Key string
}

// Get 获取指定缓存值。
// 如果key存在，value被更新成对应值， 反之value值不做改变。
func (keyOp *KeyOperation) Get(value any) error {
	return keyOp.p.Get(keyOp.Key, value)
}

// MustGet 是 Get 的 panic 版。
func (keyOp *KeyOperation) MustGet(value any) {
	err := keyOp.p.Get(keyOp.Key, value)
	if err != nil {
		panic(err)
	}
}

// TryGet 尝试获取指定缓存。
// 若key存在，value被更新成对应值，返回true，反之value值不做改变，返回false。
func (keyOp *KeyOperation) TryGet(value any) (bool, error) {
	return keyOp.p.TryGet(keyOp.Key, value)
}

// MustTryGet 是 TryGet 的 panic 版。
func (keyOp *KeyOperation) MustTryGet(value any) bool {
	result, err := keyOp.p.TryGet(keyOp.Key, value)
	if err != nil {
		panic(err)
	}
	return result
}

// Create 仅当缓存键不存在时，创建缓存。
//  return: true表示创建了缓存；false说明缓存已经存在了。
func (keyOp *KeyOperation) Create(value any) (bool, error) {
	return keyOp.p.Create(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustCreate 是 Create 的 panic 版。
func (keyOp *KeyOperation) MustCreate(value any) bool {
	result, err := keyOp.p.Create(keyOp.Key, value, keyOp.exp.NextExpireTime())
	if err != nil {
		panic(err)
	}
	return result
}

// Set 设置或者更新缓存。
func (keyOp *KeyOperation) Set(value any) error {
	return keyOp.p.Set(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustSet 是 Set 的 panic 版。
func (keyOp *KeyOperation) MustSet(value any) {
	err := keyOp.p.Set(keyOp.Key, value, keyOp.exp.NextExpireTime())
	if err != nil {
		panic(err)
	}
}

// Remove 移除指定缓存,
//  return: true成功移除，false缓存不存在。
func (keyOp *KeyOperation) Remove() (bool, error) {
	return keyOp.p.Remove(keyOp.Key)
}

// MustRemove 是 Remove 的 panic 版。
func (keyOp *KeyOperation) MustRemove() bool {
	result, err := keyOp.p.Remove(keyOp.Key)
	if err != nil {
		panic(err)
	}
	return result
}

// KeyOperationT 是泛型版本的 KeyOperation 。
type KeyOperationT[T any] struct {
	p   CacheProvider
	exp *Expiration

	// 缓存key。
	Key string
}

// Get 获取指定缓存值。
func (keyOp *KeyOperationT[T]) Get() (T, error) {
	var v T
	err := keyOp.p.Get(keyOp.Key, &v)
	return v, err
}

// MustGet 是 Get 的 panic 版。
func (keyOp *KeyOperationT[T]) MustGet() T {
	v, err := keyOp.Get()
	if err != nil {
		panic(err)
	}
	return v
}

// TryGet 尝试获取指定缓存。
// 若key存在，value被更新成对应值，返回true，反之value值不做改变，返回false。
func (keyOp *KeyOperationT[T]) TryGet() (T, bool, error) {
	var v T
	result, err := keyOp.p.TryGet(keyOp.Key, &v)
	return v, result, err

}

// MustTryGet 是 TryGet 的 panic 版。
func (keyOp *KeyOperationT[T]) MustTryGet() (T, bool) {
	var v T
	result, err := keyOp.p.TryGet(keyOp.Key, &v)
	if err != nil {
		panic(err)
	}
	return v, result
}

// Create 仅当缓存键不存在时，创建缓存。
//  return: true表示创建了缓存；false说明缓存已经存在了。
func (keyOp *KeyOperationT[T]) Create(value T) (bool, error) {
	return keyOp.p.Create(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustCreate 是 Create 的 panic 版。
func (keyOp *KeyOperationT[T]) MustCreate(value T) bool {
	result, err := keyOp.Create(value)
	if err != nil {
		panic(err)
	}
	return result
}

// Set 设置或者更新缓存。
func (keyOp *KeyOperationT[T]) Set(value T) error {
	return keyOp.p.Set(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustSet 是 Set 的 panic 版。
func (keyOp *KeyOperationT[T]) MustSet(value T) {
	err := keyOp.Set(value)
	if err != nil {
		panic(err)
	}
}

// Remove 移除指定缓存,
//  return: true成功移除，false缓存不存在。
func (keyOp *KeyOperationT[T]) Remove() (bool, error) {
	return keyOp.p.Remove(keyOp.Key)
}

// MustRemove 是 Remove 的 panic 版。
func (keyOp *KeyOperationT[T]) MustRemove() bool {
	result, err := keyOp.Remove()
	if err != nil {
		panic(err)
	}
	return result
}
