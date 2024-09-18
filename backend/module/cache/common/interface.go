package common

import "time"

type Cache interface {
	Check() error
	Set(k string, data string, du time.Duration) error
	Get(k string) (string, bool, error)
	LPush(k string, values ...interface{}) error
	RPop(k string) (string, bool, error)
	LLen(k string) (int64, error)
	Expire(k string, du time.Duration) error
	Del(string) error
}
