package cache

type KeyOperation struct {
	cp *Operation

	// 缓存key。
	Key string
}

// Get 获取指定缓存值,
// 如果key存在，value被更新成对应值，
// 反之value值不做改变。
func (keyOp *KeyOperation) Get(value any) error {
	return keyOp.cp.cacheProvider.Get(keyOp.Key, value)
}

// MustGet 是 Get 的 panic 版。
func (keyOp *KeyOperation) MustGet(value any) {
	keyOp.cp.cacheProvider.MustGet(keyOp.Key, value)
}

// TryGet 尝试获取指定缓存，
// 若key存在，value被更新成对应值，返回true，
// 反之value值不做改变，返回false。
func (keyOp *KeyOperation) TryGet(value any) (bool, error) {
	return keyOp.cp.cacheProvider.TryGet(keyOp.Key, value)
}

// Create 仅当缓存键不存在时，创建缓存，
// t 过期时长， 0 表不过期。
// return: true表示创建了缓存；false说明缓存已经存在了。
func (keyOp *KeyOperation) Create(value any) (bool, error) {
	return keyOp.cp.cacheProvider.Create(keyOp.Key, value, keyOp.cp.expireTime.NextExpireTime())
}

// MustCreate 是 Create 的 panic 版。
func (keyOp *KeyOperation) MustCreate(value any) bool {
	return keyOp.cp.cacheProvider.MustCreate(keyOp.Key, value, keyOp.cp.expireTime.NextExpireTime())
}

// Set 设置或者更新缓存，
// t 过期时长， 0 表不过期。
func (keyOp *KeyOperation) Set(value any) error {
	return keyOp.cp.cacheProvider.Set(keyOp.Key, value, keyOp.cp.expireTime.NextExpireTime())
}

// MustSet 是 Set 的 panic 版。
func (keyOp *KeyOperation) MustSet(value any) {
	keyOp.cp.cacheProvider.MustSet(keyOp.Key, value, keyOp.cp.expireTime.NextExpireTime())
}

// Remove 移除指定缓存,
// return: true成功移除，false缓存不存在。
func (keyOp *KeyOperation) Remove() (bool, error) {
	return keyOp.cp.cacheProvider.Remove(keyOp.Key)
}

// MustRemove 是 Remove 的 panic 版。
func (keyOp *KeyOperation) MustRemove() bool {
	return keyOp.cp.cacheProvider.MustRemove(keyOp.Key)
}
