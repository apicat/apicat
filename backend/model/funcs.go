package model

import (
	"errors"

	"gorm.io/gorm"
)

func NotRecord(tx *gorm.DB) error {
	if tx == nil || tx.Error == nil {
		return nil
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return tx.Error
}
