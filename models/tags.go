package models

import "time"

type Tags struct {
	ID           uint   `gorm:"type:integer primary key autoincrement"`
	ProjectId    uint   `gorm:"index;not null;comment:项目id"`
	Name         string `gorm:"type:varchar(255);not null;comment:名称"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewTags() *Tags {
	return &Tags{}
}

func (t *Tags) Create() error {
	return Conn.Create(t).Error
}

func GetByName(projectID uint, name string) (*Tags, error) {
	t := NewTags()
	err := Conn.Where("project_id = ? and name = ?", projectID, name).Take(t).Error
	return t, err
}

func TagsImport(projectID uint, collectionID uint, tags []string) {
	if len(tags) > 0 {
		for _, tag := range tags {
			t, err := GetByName(projectID, tag)
			if err != nil {
				t.ProjectId = projectID
				t.Name = tag
				err = t.Create()
			}

			if err == nil {
				ttc := NewTagToCollections()
				ttc.TagId = t.ID
				ttc.CollectionId = collectionID
				ttc.Create()
			}
		}
	}
}

func TagsExport(collectionID uint) []string {
	var tagNames []string

	tagIds := CollectionToTagIds(collectionID)
	if len(tagIds) > 0 {
		var tags []Tags
		if err := Conn.Where("id IN ?", tagIds).Find(&tags).Error; err == nil {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}
	}

	return tagNames
}
