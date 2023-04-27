package models

import (
	"encoding/json"
	"strings"
	"time"

	"strconv"

	"github.com/apicat/apicat/commom/spec"
	"gorm.io/gorm"
)

type Collections struct {
	ID           uint   `gorm:"type:integer primary key autoincrement"`
	ProjectId    uint   `gorm:"index;not null;comment:项目id"`
	ParentId     uint   `gorm:"not null;comment:父级id"`
	Title        string `gorm:"type:varchar(255);not null;comment:名称"`
	Type         string `gorm:"type:varchar(255);not null;comment:类型:category,doc,http"`
	Content      string `gorm:"type:mediumtext;comment:内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt    time.Time
	CreatedBy    uint `gorm:"not null;default:0;comment:创建人id"`
	UpdatedAt    time.Time
	UpdatedBy    uint `gorm:"not null;default:0;comment:最后更新人id"`
	DeletedAt    gorm.DeletedAt
	DeletedBy    uint `gorm:"not null;default:0;comment:删除人id"`
}

func NewCollections(ids ...uint) (*Collections, error) {
	if len(ids) > 0 {
		collection := &Collections{ID: ids[0]}
		if err := Conn.Take(collection).Error; err != nil {
			return collection, err
		}
		return collection, nil
	}
	return &Collections{}, nil
}

func (c *Collections) List() ([]*Collections, error) {
	collectionsQuery := Conn.Where("project_id = ?", c.ProjectId)

	var collections []*Collections
	return collections, collectionsQuery.Order("display_order asc").Order("id desc").Find(&collections).Error
}

func (c *Collections) Create() error {
	return Conn.Create(c).Error
}

func (c *Collections) Update() error {
	return Conn.Save(c).Error
}

func Deletes(id uint, db *gorm.DB) error {
	collection := Collections{}
	if err := Conn.Where("id = ?", id).First(&collection).Error; err != nil {
		return err
	}

	collections := []*Collections{}
	if err := Conn.Where("parent_id = ?", id).Find(&collections).Error; err != nil {
		return err
	}

	return Conn.Transaction(func(tx *gorm.DB) error {
		for _, subNode := range collections {
			if err := Deletes(subNode.ID, tx); err != nil {
				return err
			}
		}

		if err := tx.Delete(&collection).Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *Collections) Creator() string {
	return ""
}

func (c *Collections) Updater() string {
	return ""
}

func (c *Collections) Deleter() string {
	return ""
}

func (c *Collections) TrashList() ([]*Collections, error) {
	var deleteCollections []*Collections
	return deleteCollections, Conn.Unscoped().Where("deleted_at is not null AND project_id = ?", c.ProjectId).Find(&deleteCollections).Error
}

func (c *Collections) GetUnscopedCollections() error {
	return Conn.Unscoped().Where("id = ? AND project_id = ?", c.ID, c.ProjectId).Take(c).Error
}

func (c *Collections) Restore() error {
	return Conn.Unscoped().Model(c).Updates(map[string]interface{}{"project_id": c.ProjectId, "parent_id": c.ParentId, "display_order": 0, "deleted_at": nil}).Error
}

func CollectionsImport(projectID, parentID uint, collections []*spec.CollectItem, definitionSchemas nameToIdMap) {
	for i, collection := range collections {
		if len(collection.Items) > 0 {
			category := &Collections{
				ProjectId: projectID,
				ParentId:  parentID,
				Title:     collection.Title,
				Type:      "category",
			}
			if err := category.Create(); err == nil {
				CollectionsImport(projectID, category.ID, collection.Items, definitionSchemas)
			}
		} else {
			if collectionByte, err := json.Marshal(collection.Content); err == nil {
				collectionStr := string(collectionByte)
				collectionStr = replaceNameToID(collectionStr, definitionSchemas, "#/definitions/schemas/")

				record := &Collections{
					ProjectId:    projectID,
					ParentId:     0,
					Title:        collection.Title,
					Type:         "http",
					Content:      collectionStr,
					DisplayOrder: i,
				}
				if err := record.Create(); err == nil {
					TagsImport(projectID, record.ID, collection.Tags)
				}
			}
		}
	}
}

func replaceNameToID(content string, nameIDMap nameToIdMap, prefix string) string {
	for name, id := range nameIDMap {
		oldStr := prefix + name
		newStr := prefix + strconv.FormatUint(uint64(id), 10)

		content = strings.Replace(content, oldStr, newStr, -1)
	}
	return content
}

// func CollectionsExport(projectID uint) []*spec.CollectItem {
// 	var collections []*Collections
// 	collectItems := make([]*spec.CollectItem, 0)

// 	if err := Conn.Where("project_id = ?", projectID).Find(&collections).Error; err == nil {
// 		parentCollection := &Collections{ID: 0}
// 		collectItems = collectionsTree(collections, parentCollection)
// 	}

// 	return collectItems
// }

func collectionsTree(collections []*Collections, parentCollection *Collections) []*spec.CollectItem {
	collectItems := make([]*spec.CollectItem, 0)
	for _, collection := range collections {
		if collection.ParentId == parentCollection.ID {
			collectItem := &spec.CollectItem{
				ID:       int64(collection.ID),
				ParentID: int64(collection.ParentId),
				Title:    collection.Title,
				Type:     spec.ContentType(collection.Type),
			}

			// 将父级的分类名称也加入Tags中
			if parentCollection.ID > 0 {
				if !collectItem.HasTag(parentCollection.Title) {
					collectItem.Tags = append(collectItem.Tags, parentCollection.Title)
				}
			}

			if tags := TagsExport(collection.ID); len(tags) > 0 {
				collectItem.Tags = append(collectItem.Tags, tags...)
			}

			if collection.Type != "category" {
				content := []*spec.NodeProxy{}
				if json.Unmarshal([]byte(collection.Content), &content) == nil {
					collectItem.Content = content
				}
			}

			collectItem.Items = collectionsTree(collections, collection)
			collectItems = append(collectItems, collectItem)
		}
	}

	return collectItems
}

func CollectionsDdereference(d *Definitions) error {
	ref := "{\"$ref\":\"#/definitions/schemas/" + strconv.FormatUint(uint64(d.ID), 10) + "\"}"

	collections, _ := NewCollections()
	collections.ProjectId = d.ProjectId
	collectionList, err := collections.List()
	if err != nil {
		return err
	}

	for _, collection := range collectionList {
		if strings.Contains(collection.Content, ref) {
			newContent := strings.Replace(collection.Content, ref, d.Schema, -1)
			collection.Content = newContent

			if err := collection.Update(); err != nil {
				return err
			}
		}
	}

	return nil
}
