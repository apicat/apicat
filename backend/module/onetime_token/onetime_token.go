package onetime_token

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/apicat/apicat/backend/module/cache/common"
)

type Helper struct {
	c common.Cache
}

func NewTokenHelper(c common.Cache) *Helper {
	return &Helper{c: c}

}

func (h *Helper) GenerateToken(k string, data any, exp time.Duration) (string, error) {
	hashInBytes := md5.Sum([]byte(k))
	md5String := hex.EncodeToString(hashInBytes[:])

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	if err = h.c.Set(md5String, string(jsonData), exp); err != nil {
		return "", err
	}

	return md5String, nil
}

func (h *Helper) CheckToken(token string, data any) bool {
	jsonData, exist, err := h.c.Get(token)
	if !exist || err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(jsonData), data); err != nil {
		return false
	}

	return true
}

func (h *Helper) DelToken(token string) error {
	return h.c.Del(token)
}
