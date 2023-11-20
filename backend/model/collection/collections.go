package collection

import (
	"encoding/json"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/global"
	"github.com/apicat/apicat/backend/model/iteration"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/server"
	"github.com/apicat/apicat/backend/model/tag"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/apicat_struct"
	"github.com/apicat/apicat/backend/module/spec"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func init() {
	model.RegMigrate(&Collections{})
}

func NewCollections(ids ...uint) (*Collections, error) {
	if len(ids) > 0 {
		collection := &Collections{ID: ids[0]}
		if err := model.Conn.Take(collection).Error; err != nil {
			return collection, err
		}
		return collection, nil
	}
	return &Collections{}, nil
}

func (c *Collections) List() ([]*Collections, error) {
	collectionsQuery := model.Conn.Where("project_id = ?", c.ProjectId)

	var collections []*Collections
	return collections, collectionsQuery.Order("display_order asc").Find(&collections).Error
}

func (c *Collections) CreateDoc() error {
	var node *Collections
	if err := model.Conn.Where("project_id = ? AND parent_id = ?", c.ProjectId, c.ParentId).Order("display_order desc").First(&node).Error; err == nil {
		c.DisplayOrder = node.DisplayOrder + 1
	}

	return c.Create()
}

func (c *Collections) CreateCategory() error {
	err := model.Conn.Model(&Collections{}).Where("parent_id = ?", c.ParentId).Update("display_order", gorm.Expr("display_order + ?", 1)).Error
	if err != nil {
		return err
	}

	return c.Create()
}

func (c *Collections) Create() error {
	return model.Conn.Create(c).Error
}

func (c *Collections) Update() error {
	return model.Conn.Save(c).Error
}

func BatchUpdateByProjectID(ProjectID uint, c map[string]any) error {
	return model.Conn.Model(&Collections{}).Where("project_id = ?", ProjectID).Updates(c).Error
}

func Deletes(id uint, db *gorm.DB, deletedBy uint) error {
	collection := Collections{}
	if err := model.Conn.Where("id = ?", id).First(&collection).Error; err != nil {
		return err
	}

	collections := []*Collections{}
	if err := model.Conn.Where("parent_id = ?", id).Find(&collections).Error; err != nil {
		return err
	}

	return model.Conn.Transaction(func(tx *gorm.DB) error {
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

		if err := iteration.DeleteIterationApisByCollectionID(collection.ID); err != nil {
			return err
		}

		return nil
	})
}

func (c *Collections) GetSubCollectionsContainsSelf() ([]*Collections, error) {
	var collections []*Collections
	collections, err := c.getSubCollectionsRecursive(&collections)
	if err != nil {
		return nil, err
	}

	collections = append(collections, c)
	return collections, nil
}

func (c *Collections) getSubCollectionsRecursive(collectPtr *[]*Collections) ([]*Collections, error) {
	var subCollections []*Collections

	if err := model.Conn.Where("parent_id = ?", c.ID).Find(&subCollections).Error; err != nil {
		return nil, err
	}

	*collectPtr = append(*collectPtr, subCollections...)

	for _, subColl := range subCollections {
		if subColl.ID != c.ID { // Avoid self-reference
			_, err := subColl.getSubCollectionsRecursive(collectPtr)
			if err != nil {
				return nil, err
			}
		}
	}

	return *collectPtr, nil
}

func (c *Collections) Creator() string {
	u, err := user.NewUsers(c.CreatedBy)
	if err != nil {
		return ""
	}

	return u.Username
}

func (c *Collections) Updater() string {
	u, err := user.NewUsers(c.UpdatedBy)
	if err != nil {
		return ""
	}

	return u.Username
}

func (c *Collections) Deleter() string {
	u, err := user.NewUsers(c.DeletedBy)
	if err != nil {
		return ""
	}

	return u.Username
}

func (c *Collections) TrashList() ([]*Collections, error) {
	var deleteCollections []*Collections
	return deleteCollections, model.Conn.Unscoped().Where("deleted_at is not null AND project_id = ?", c.ProjectId).Find(&deleteCollections).Error
}

func (c *Collections) GetUnscopedCollections() error {
	return model.Conn.Unscoped().Where("id = ? AND project_id = ?", c.ID, c.ProjectId).Take(c).Error
}

func (c *Collections) Restore() error {
	return model.Conn.Unscoped().Model(c).Updates(map[string]interface{}{"project_id": c.ProjectId, "parent_id": c.ParentId, "display_order": 0, "deleted_at": nil, "deleted_by": 0}).Error
}

func CollectionsImport(projectID, parentID uint, collections []*spec.CollectItem, refContentNameToId *model.RefContentVirtualIDToId) []*Collections {
	collectionList := []*Collections{}

	for i, collection := range collections {
		if len(collection.Items) > 0 || collection.Type == "category" {
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
				collectionStr = model.ReplaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionSchemas, "#/definitions/schemas/")
				collectionStr = model.ReplaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionResponses, "#/definitions/responses/")
				collectionStr = model.ReplaceVirtualIDToID(collectionStr, refContentNameToId.DefinitionParameters, "#/definitions/parameters/")
				collectionStr = global.ReplaceGlobalParametersVirtualIDToID(collectionStr, refContentNameToId.GlobalParameters)

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
					tag.TagsImport(projectID, record.ID, collection.Tags)
				}
			}
		}
	}

	return collectionList
}

func CollectionsExport(projectID uint) []*spec.CollectItem {
	collections := []*Collections{}
	collectItems := []*spec.CollectItem{}

	if err := model.Conn.Where("project_id = ?", projectID).Find(&collections).Error; err == nil {
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

			if tags := tag.TagsExport(collection.ID); len(tags) > 0 {
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

func (c *Collections) GetByPublicId() error {
	return model.Conn.Where("public_id = ?", c.PublicId).First(c).Error
}

func DefinitionsSchemaUnRefByCollections(d *definition.DefinitionSchemas, isUnRef int) error {
	ref := "\"$ref\":\"#/definitions/schemas/" + strconv.FormatUint(uint64(d.ID), 10) + "\""

	collections, _ := NewCollections()
	collections.ProjectId = d.ProjectId
	collectionList, err := collections.List()
	if err != nil {
		return err
	}

	sourceJson := map[string]interface{}{}
	if err := json.Unmarshal([]byte(d.Schema), &sourceJson); err != nil {
		return err
	}
	typeEmptyStructure := apicat_struct.TypeEmptyStructure()

	for _, c := range collectionList {
		if strings.Contains(c.Content, ref) {
			newStr := typeEmptyStructure[sourceJson["type"].(string)]
			if isUnRef == 1 {
				newStr = d.Schema
			}

			newContent := strings.Replace(c.Content, ref, newStr[1:len(newStr)-1], -1)
			c.Content = newContent

			if err := c.Update(); err != nil {
				return err
			}
		}
	}

	return nil
}

func DefinitionsResponseUnRef(dr *definition.DefinitionResponses) error {
	ref := "\"$ref\":\"#/definitions/responses/" + strconv.Itoa(int(dr.ID)) + "\""

	collections, _ := NewCollections()
	collections.ProjectId = dr.ProjectID
	collectionList, err := collections.List()
	if err != nil {
		return err
	}

	header := []interface{}{}
	if err := json.Unmarshal([]byte(dr.Header), &header); err != nil {
		return err
	}

	content := map[string]interface{}{}
	if err := json.Unmarshal([]byte(dr.Content), &content); err != nil {
		return err
	}

	data := map[string]interface{}{
		"name":        dr.Name,
		"description": dr.Description,
		"header":      header,
		"content":     content,
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	newStr := string(dataJson)[1 : len(string(dataJson))-1]

	for _, c := range collectionList {
		if strings.Contains(c.Content, ref) {
			newContent := strings.Replace(c.Content, ref, newStr, -1)
			c.Content = newContent

			if err := c.Update(); err != nil {
				return err
			}
		}
	}

	return nil
}

func DefinitionsResponseDelRef(dr *definition.DefinitionResponses) error {
	re1 := regexp.MustCompile(`,{"code":\d+,"\$ref":"#/definitions/responses/` + strconv.Itoa(int(dr.ID)) + `"}`)
	re2 := regexp.MustCompile(`{"code":\d+,"\$ref":"#/definitions/responses/` + strconv.Itoa(int(dr.ID)) + `"}`)

	collections, _ := NewCollections()
	collections.ProjectId = dr.ProjectID
	collectionList, err := collections.List()
	if err != nil {
		return err
	}

	emptyResponse := apicat_struct.TypeEmptyStructure()["response"]

	for _, c := range collectionList {
		matchRe1 := re1.FindString(c.Content)
		if matchRe1 != "" {
			newContent := strings.Replace(c.Content, matchRe1, "", -1)
			c.Content = newContent
		} else {
			matchRe2 := re2.FindString(c.Content)
			if matchRe2 != "" {
				newContent := strings.Replace(c.Content, matchRe2, emptyResponse, -1)
				c.Content = newContent
			}
		}

		if err := c.Update(); err != nil {
			return err
		}
	}
	return nil
}

// CollectionExport 返回单篇文档导出的 apicat 结构
// project 导出集合所属项目 model
// collection 导出的集合 model
func CollectionExport(project *project.Projects, collection *Collections) *spec.Spec {
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
	apicatData.Servers = server.ServersExport(project.ID)
	apicatData.Globals.Parameters = global.GlobalParametersExport(project.ID)
	apicatData.Definitions.Schemas = definition.DefinitionSchemasExport(project.ID)
	apicatData.Definitions.Parameters = definition.DefinitionParametersExport(project.ID)
	apicatData.Definitions.Responses = definition.DefinitionResponsesExport(project.ID)

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

func PlanningIterationApi(cIDs []uint, i *iteration.Iterations) error {

	iterationApi, _ := iteration.NewIterationApis()
	iterationApis, err := iterationApi.List(i.ID)
	if err != nil {
		return err
	}

	if len(cIDs) == 0 && len(iterationApis) != 0 {
		if err := iteration.BatchDeleteIterationApi(iterationApis); err != nil {
			return err
		}
	}

	c, _ := NewCollections()
	c.ProjectId = i.ProjectID
	collections, err := c.List()
	if err != nil {
		return err
	}
	collectionDict := map[uint]*Collections{}
	for _, v := range collections {
		collectionDict[v.ID] = v
	}

	wantPop := []*iteration.IterationApis{}
	wantPush := []*iteration.IterationApis{}

	// 找出iterationApis中存在但cIDs中不存在的元素
	for _, iterationApi := range iterationApis {
		found := false
		for _, cid := range cIDs {
			if iterationApi.CollectionID == cid {
				found = true
				break
			}
		}
		if !found {
			wantPop = append(wantPop, iterationApi)
		}
	}

	// 找出cIDs中存在但iterationApis中不存在的元素
	for _, cid := range cIDs {
		if _, ok := collectionDict[cid]; ok {
			found := false
			for _, iterationApi := range iterationApis {
				if cid == iterationApi.CollectionID {
					found = true
					break
				}
			}
			if !found {
				wantPush = append(wantPush, &iteration.IterationApis{
					IterationID:    i.ID,
					CollectionID:   collectionDict[cid].ID,
					CollectionType: collectionDict[cid].Type,
				})
			}
		}
	}

	if len(wantPop) > 0 {
		if err := iteration.BatchDeleteIterationApi(wantPop); err != nil {
			return err
		}
	}
	if len(wantPush) > 0 {
		if err := iteration.BatchInsertIterationApi(wantPush); err != nil {
			return err
		}
	}
	return nil
}
