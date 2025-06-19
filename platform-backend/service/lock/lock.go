package lock

import (
	"sync"
	"sync/atomic"
	"time"
)

/*
	单个pod下直接使用内存锁  后续多pod替换为redis
*/

var Locker = &lockManager{}

type lockManager struct {
	locks sync.Map // 存储每个 key 对应的锁
}

type lockWrapper struct {
	mu        sync.Mutex
	isLocked  int32         // 0 表示未锁定，1 表示锁定
	expireAt  int64         // 锁的过期时间（时间戳，毫秒级）
	expireDur time.Duration // 锁的过期持续时间
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// TryLock 尝试获取指定 key 的锁。如果获取成功，返回 true；如果锁过期或被其他 Goroutine 获取，返回 false。
func (lm *lockManager) TryLock(key string, expireDuration time.Duration) bool {
	// 获取或创建一个 lockWrapper
	val, _ := lm.locks.LoadOrStore(key, &lockWrapper{expireDur: expireDuration})
	lock := val.(*lockWrapper)

	// 检查锁是否已经过期，若过期，则重置锁状态
	if lock.expireAt != 0 && currentTimeMillis() > lock.expireAt {
		// 锁已过期，尝试重置
		atomic.StoreInt32(&lock.isLocked, 0)
		lock.expireAt = 0 // 清除过期时间
	}

	// 如果锁未被占用，则设置为已占用并加锁
	if atomic.CompareAndSwapInt32(&lock.isLocked, 0, 1) {
		lock.mu.Lock()
		lock.expireAt = currentTimeMillis() + int64(expireDuration/time.Millisecond) // 设置锁的过期时间
		return true
	}
	return false
}

// Unlock 解锁指定的 key
func (lm *lockManager) Unlock(key string) {
	// 从 Map 中获取锁
	val, exists := lm.locks.Load(key)
	if !exists {
		return // 如果锁不存在，直接返回
	}
	lock := val.(*lockWrapper)

	// 解锁并重置锁状态
	lock.mu.Unlock()
	atomic.StoreInt32(&lock.isLocked, 0)
	lock.expireAt = 0 // 重置过期时间
}
