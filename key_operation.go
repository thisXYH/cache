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
	keyOp.p.MustGet(keyOp.Key, value)
}

// TryGet 尝试获取指定缓存。
// 若key存在，value被更新成对应值，返回true，反之value值不做改变，返回false。
func (keyOp *KeyOperation) TryGet(value any) (bool, error) {
	return keyOp.p.TryGet(keyOp.Key, value)
}

// Create 仅当缓存键不存在时，创建缓存。
//  return: true表示创建了缓存；false说明缓存已经存在了。
func (keyOp *KeyOperation) Create(value any) (bool, error) {
	return keyOp.p.Create(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustCreate 是 Create 的 panic 版。
func (keyOp *KeyOperation) MustCreate(value any) bool {
	return keyOp.p.MustCreate(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// Set 设置或者更新缓存。
func (keyOp *KeyOperation) Set(value any) error {
	return keyOp.p.Set(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustSet 是 Set 的 panic 版。
func (keyOp *KeyOperation) MustSet(value any) {
	keyOp.p.MustSet(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// Remove 移除指定缓存,
//  return: true成功移除，false缓存不存在。
func (keyOp *KeyOperation) Remove() (bool, error) {
	return keyOp.p.Remove(keyOp.Key)
}

// MustRemove 是 Remove 的 panic 版。
func (keyOp *KeyOperation) MustRemove() bool {
	return keyOp.p.MustRemove(keyOp.Key)
}

// GenericKeyOperation 是泛型版本的 KeyOperation 。
type GenericKeyOperation[T any] struct {
	p   CacheProvider
	exp *Expiration

	// 缓存key。
	Key string
}

// Get 获取指定缓存值。
// 如果key存在，value被更新成对应值， 反之value值不做改变。
func (keyOp *GenericKeyOperation[T]) Get(value *T) error {
	return keyOp.p.Get(keyOp.Key, value)
}

// MustGet 是 Get 的 panic 版。
func (keyOp *GenericKeyOperation[T]) MustGet(value *T) {
	keyOp.p.MustGet(keyOp.Key, value)
}

// TryGet 尝试获取指定缓存。
// 若key存在，value被更新成对应值，返回true，反之value值不做改变，返回false。
func (keyOp *GenericKeyOperation[T]) TryGet(value *T) (bool, error) {
	return keyOp.p.TryGet(keyOp.Key, value)
}

// Create 仅当缓存键不存在时，创建缓存。
//  return: true表示创建了缓存；false说明缓存已经存在了。
func (keyOp *GenericKeyOperation[T]) Create(value T) (bool, error) {
	return keyOp.p.Create(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustCreate 是 Create 的 panic 版。
func (keyOp *GenericKeyOperation[T]) MustCreate(value T) bool {
	return keyOp.p.MustCreate(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// Set 设置或者更新缓存。
func (keyOp *GenericKeyOperation[T]) Set(value T) error {
	return keyOp.p.Set(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// MustSet 是 Set 的 panic 版。
func (keyOp *GenericKeyOperation[T]) MustSet(value T) {
	keyOp.p.MustSet(keyOp.Key, value, keyOp.exp.NextExpireTime())
}

// Remove 移除指定缓存,
//  return: true成功移除，false缓存不存在。
func (keyOp *GenericKeyOperation[T]) Remove() (bool, error) {
	return keyOp.p.Remove(keyOp.Key)
}

// MustRemove 是 Remove 的 panic 版。
func (keyOp *GenericKeyOperation[T]) MustRemove() bool {
	return keyOp.p.MustRemove(keyOp.Key)
}
