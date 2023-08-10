package models

import "time"

type ProjectFollows struct {
	ID        uint `gorm:"type:bigint;primaryKey;autoIncrement"`
	UserID    uint `gorm:"type:bigint;index;not null;comment:用户id"`
	ProjectID uint `gorm:"type:bigint;not null;comment:项目id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProjectFollows(ids ...uint) (*ProjectFollows, error) {
	if len(ids) > 0 {
		projectFollow := &ProjectFollows{ID: ids[0]}
		if err := Conn.Take(projectFollow).Error; err != nil {
			return projectFollow, err
		}
		return projectFollow, nil
	}
	return &ProjectFollows{}, nil
}

func (pf *ProjectFollows) List(uID uint) ([]*ProjectFollows, error) {
	var projectFollows []*ProjectFollows

	return projectFollows, Conn.Where("user_id = ?", uID).Order("created_at desc").Find(&projectFollows).Error
}

func (pf *ProjectFollows) GetByUserIDAndProjectID() error {
	return Conn.Where("user_id = ?", pf.UserID).Where("project_id = ?", pf.ProjectID).Take(pf).Error
}

func (pf *ProjectFollows) Create() error {
	return Conn.Create(pf).Error
}

func (pf *ProjectFollows) Delete() error {
	return Conn.Delete(pf).Error
}
