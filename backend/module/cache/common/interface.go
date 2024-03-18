package common

import "time"

type Cache interface {
	Check() error
	Set(string, string, time.Duration) error
	Get(string) (string, bool, error)
	Del(string) error
}
