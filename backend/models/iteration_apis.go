package models

import "time"

type IterationApis struct {
	ID             uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	IterationID    uint   `gorm:"type:bigint;index;not null;comment:迭代id"`
	CollectionID   uint   `gorm:"type:bigint;not null;comment:集合id"`
	CollectionType string `gorm:"type:varchar(255);not null;comment:集合类型:category,doc,http"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewIterationApis(ids ...uint) (*IterationApis, error) {
	if len(ids) > 0 {
		iterationApi := &IterationApis{ID: ids[0]}
		if err := Conn.Take(iterationApi).Error; err != nil {
			return iterationApi, err
		}
		return iterationApi, nil
	}
	return &IterationApis{}, nil
}

func (ia *IterationApis) List(iID ...uint) ([]*IterationApis, error) {
	var iterationApis []*IterationApis

	if len(iID) > 0 {
		return iterationApis, Conn.Where("iteration_id IN ?", iID).Order("created_at desc").Find(&iterationApis).Error
	}

	return iterationApis, Conn.Order("created_at desc").Find(&iterationApis).Error
}

func (ia *IterationApis) GetCollectionIDByIterationID(iID uint) ([]uint, error) {
	var collectionIDs []uint

	return collectionIDs, Conn.Model(&IterationApis{}).Where("iteration_id = ?", iID).Pluck("collection_id", &collectionIDs).Error
}

func IterationApiCount(IterationID uint, cType string) (int64, error) {
	var count int64
	return count, Conn.Model(&IterationApis{}).Where("iteration_id = ?", IterationID).Where("collection_type = ?", cType).Count(&count).Error
}

func (ia *IterationApis) Create() error {
	return Conn.Create(ia).Error
}

func (ia *IterationApis) Delete() error {
	return Conn.Delete(ia).Error
}

func BatchInsertIterationApi(ias []*IterationApis) error {
	return Conn.Create(&ias).Error
}

func BatchDeleteIterationApi(ias []*IterationApis) error {
	return Conn.Delete(&ias).Error
}
