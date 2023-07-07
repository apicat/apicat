package models

import (
	"time"
)

type ShareTmpTokens struct {
	ID           uint      `gorm:"type:bigint;primaryKey;autoIncrement"`
	ShareToken   string    `gorm:"type:varchar(255);index;not null;comment:md5的分享token"`
	Expiration   time.Time `gorm:"type:datetime;not null;comment:过期时间"`
	ProjectID    uint      `gorm:"type:bigint;index;not null;comment:项目id"`
	CollectionID uint      `gorm:"type:bigint;index;comment:集合id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewShareTmpTokens() *ShareTmpTokens {
	return &ShareTmpTokens{}
}

func (stt *ShareTmpTokens) GetByShareToken() error {
	return Conn.Where("share_token = ?", stt.ShareToken).Take(stt).Error
}

func (stt *ShareTmpTokens) Create() error {
	return Conn.Create(stt).Error
}

func (stt *ShareTmpTokens) Delete() error {
	return Conn.Delete(stt).Error
}

func (stt *ShareTmpTokens) DeleteByProjectID() error {
	return Conn.Where("project_id = ?", stt.ProjectID).Delete(stt).Error
}

func (stt *ShareTmpTokens) DeleteByCollectionID() error {
	return Conn.Where("collection_id = ?", stt.CollectionID).Delete(stt).Error
}
