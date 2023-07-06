package models

import (
	"encoding/json"
	"strings"
	"time"

	"strconv"

	"github.com/apicat/apicat/common/spec"
	"gorm.io/gorm"
)

type Collections struct {
	ID            uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	PublicId      string `gorm:"type:varchar(255);comment:集合公开id"`
	ProjectId     uint   `gorm:"type:bigint;index;not null;comment:项目id"`
	ParentId      uint   `gorm:"type:bigint;not null;comment:父级id"`
	Title         string `gorm:"type:varchar(255);not null;comment:名称"`
	Type          string `gorm:"type:varchar(255);not null;comment:类型:category,doc,http"`
	SharePassword string `gorm:"type:varchar(255);comment:项目分享密码"`
	Content       string `gorm:"type:mediumtext;comment:内容"`
	DisplayOrder  int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedAt     time.Time
	CreatedBy     uint `gorm:"type:bigint;not null;default:0;comment:创建人id"`
	UpdatedAt     time.Time
	UpdatedBy     uint `gorm:"type:bigint;not null;default:0;comment:最后更新人id"`
	DeletedAt     gorm.DeletedAt
	DeletedBy     uint `gorm:"type:bigint;not null;default:0;comment:删除人id"`
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

func Deletes(id uint, db *gorm.DB, deletedBy uint) error {
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
			if err := Deletes(subNode.ID, tx, deletedBy); err != nil {
				return err
			}
		}

		if err := tx.Delete(&collection).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Model(collection).Updates(map[string]interface{}{"deleted_by": deletedBy}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *Collections) Creator() string {
	user, err := NewUsers(c.CreatedBy)
	if err != nil {
		return ""
	}

	return user.Username
}

func (c *Collections) Updater() string {
	user, err := NewUsers(c.UpdatedBy)
	if err != nil {
		return ""
	}

	return user.Username
}

func (c *Collections) Deleter() string {
	user, err := NewUsers(c.DeletedBy)
	if err != nil {
		return ""
	}

	return user.Username
}

func (c *Collections) TrashList() ([]*Collections, error) {
	var deleteCollections []*Collections
	return deleteCollections, Conn.Unscoped().Where("deleted_at is not null AND project_id = ?", c.ProjectId).Find(&deleteCollections).Error
}

func (c *Collections) GetUnscopedCollections() error {
	return Conn.Unscoped().Where("id = ? AND project_id = ?", c.ID, c.ProjectId).Take(c).Error
}

func (c *Collections) Restore() error {
	return Conn.Unscoped().Model(c).Updates(map[string]interface{}{"project_id": c.ProjectId, "parent_id": c.ParentId, "display_order": 0, "deleted_at": nil, "deleted_by": 0}).Error
}

func CollectionsImport(projectID, parentID uint, collections []*spec.CollectItem, refContentNameToId *RefContentVirtualIDToId) []*Collections {
	collectionList := []*Collections{}

	for i, collection := range collections {
		if len(collection.Items) > 0 {
			category := &Collections{
				ProjectId: projectID,
				ParentId:  parentID,
				Title:     collection.Title,
				Type:      "category",
			}
			if err := category.Create(); err == nil {
				collectionList = append(collectionList, category)
				children := CollectionsImport(projectID, category.ID, collection.Items, refContentNameToId)
				collectionList = append(collectionList, children...)
			}
		} else {
			if collectionByte, err := json.Marshal(collection.Content); err == nil {
				collectionStr := string(collectionByte)
				collectionStr = replaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionSchemas, "#/definitions/schemas/")
				collectionStr = replaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionResponses, "#/definitions/responses/")
				collectionStr = replaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionParameters, "#/definitions/parameters/")

				record := &Collections{
					ProjectId:    projectID,
					ParentId:     parentID,
					Title:        collection.Title,
					Type:         "http",
					Content:      collectionStr,
					DisplayOrder: i,
				}
				if err := record.Create(); err == nil {
					collectionList = append(collectionList, record)
					TagsImport(projectID, record.ID, collection.Tags)
				}
			}
		}
	}

	return collectionList
}

func replaceVirtualIDToID(content string, nameIDMap virtualIDToIDMap, prefix string) string {
	for virtualID, id := range nameIDMap {
		oldStr := prefix + strconv.Itoa(int(virtualID))
		newStr := prefix + strconv.Itoa(int(id))

		content = strings.Replace(content, oldStr, newStr, -1)
	}
	return content
}

func CollectionsExport(projectID uint) []*spec.CollectItem {
	collections := []*Collections{}
	collectItems := []*spec.CollectItem{}

	if err := Conn.Where("project_id = ?", projectID).Find(&collections).Error; err == nil {
		parentCollection := &Collections{ID: 0}
		collectItems = collectionsTree(collections, parentCollection, projectID)
	}

	return collectItems
}

func collectionsTree(collections []*Collections, parentCollection *Collections, projectID uint) []*spec.CollectItem {
	collectItems := []*spec.CollectItem{}

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

			collectItem.Items = collectionsTree(collections, collection, projectID)
			collectItems = append(collectItems, collectItem)
		}
	}

	return collectItems
}

// CollectionExport 返回单篇文档导出的 apicat 结构
// project 导出集合所属项目 model
// collection 导出的集合 model
func CollectionExport(project *Projects, collection *Collections) *spec.Spec {
	apicatData := &spec.Spec{}
	apicatData.ApiCat = "apicat"
	apicatData.Info = &spec.Info{
		ID:          project.PublicId,
		Title:       project.Title,
		Description: project.Description,
		Version:     "1.0.0",
	}

	collectItem := &spec.CollectItem{
		ID:       int64(collection.ID),
		ParentID: int64(collection.ParentId),
		Title:    collection.Title,
		Type:     spec.ContentType(collection.Type),
	}
	content := []*spec.NodeProxy{}
	if json.Unmarshal([]byte(collection.Content), &content) == nil {
		collectItem.Content = content
	}

	apicatData.Collections = []*spec.CollectItem{collectItem}
	apicatData.Servers = ServersExport(project.ID)
	apicatData.Globals.Parameters = GlobalParametersExport(project.ID)
	apicatData.Definitions.Schemas = DefinitionSchemasExport(project.ID)
	apicatData.Definitions.Parameters = DefinitionParametersExport(project.ID)
	apicatData.Definitions.Responses = DefinitionResponsesExport(project.ID)

	paths := apicatData.CollectionsMap(true, 2)
	for _, path := range paths {
		for _, v := range path {
			for _, i := range content {
				switch nx := i.Node.(type) {
				case *spec.HTTPNode[spec.HTTPRequestNode]:
					nx.Attrs.Parameters = v.Parameters
					nx.Attrs.Content = v.Content
				case *spec.HTTPNode[spec.HTTPResponsesNode]:
					nx.Attrs.List = v.Responses
				}
			}
		}
	}

	apicatData.Definitions = spec.Definitions{
		Schemas:    spec.Schemas{},
		Parameters: spec.Schemas{},
		Responses:  spec.HTTPResponseDefines{},
	}

	return apicatData
}

func (c *Collections) GetByPublicId() error {
	return Conn.Where("public_id = ?", c.PublicId).First(c).Error
}
