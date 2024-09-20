package local

import (
	"sync"
	"time"
)

type localCache struct {
	dataMap *sync.Map
}

type cacheData struct {
	data        string
	expiredTime time.Time
}

var dataMap *sync.Map

func NewLocal() (*localCache, error) {
	return &localCache{
		dataMap: dataMap,
	}, nil
}

func (l *localCache) Check() error {
	dataMap = &sync.Map{}
	return nil
}

func (l *localCache) Set(k string, data string, du time.Duration) error {
	item := &cacheData{
		data:        data,
		expiredTime: time.Now().Add(du),
	}
	l.dataMap.Store(k, item)
	return nil
}

func (l *localCache) Get(k string) (string, bool, error) {
	v, ok := l.dataMap.Load(k)
	if !ok {
		return "", false, nil
	}
	x := v.(*cacheData)
	if x.expiredTime.After(time.Now()) {
		return x.data, true, nil
	}
	l.dataMap.Delete(k)
	return "", false, nil
}

func (l *localCache) Del(k string) error {
	l.dataMap.Delete(k)
	return nil
}

func (l *localCache) LPush(k string, values ...interface{}) error {
	return nil
}

func (l *localCache) RPop(k string) (string, bool, error) {
	return "", false, nil
}

func (l *localCache) LLen(k string) (int64, error) {
	return 0, nil
}

func (l *localCache) Expire(k string, du time.Duration) error {
	return nil
}

// func NewLocalCache[T any]() cache.Cache[T] {
// 	c := &localCache{}
// 	go func() {
// 		t := time.NewTicker(time.Second * 30)
// 		for range t.C {
// 			n := time.Now()
// 			clearn := 0
// 			c.dataMap.Range(func(key, value any) bool {
// 				x := value.(*cacheData)
// 				if x.expiredTime.Before(n) {
// 					clearn++
// 					c.dataMap.Delete(key)
// 				}
// 				// 一次最多清理1000个 避免大面积GC
// 				return clearn < 1000
// 			})
// 		}
// 	}()
// 	return c
// }
