package collection

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/team"
	"github.com/apicat/apicat/backend/model/user"

	"github.com/apicat/apicat/backend/module/spec"

	"gorm.io/gorm"
)

type IdToNameMap map[uint]string
type VirtualIDToIDMap map[int64]uint

type RefContentVirtualIDToId struct {
	DefinitionSchemas    VirtualIDToIDMap
	DefinitionResponses  VirtualIDToIDMap
	DefinitionParameters VirtualIDToIDMap
	GlobalParameters     VirtualIDToIDMap
}

func ReplaceVirtualIDToID(content string, nameIDMap VirtualIDToIDMap, prefix string) string {
	for virtualID, id := range nameIDMap {
		oldStr := prefix + strconv.Itoa(int(virtualID)) + "\""
		newStr := prefix + strconv.Itoa(int(id)) + "\""

		content = strings.Replace(content, oldStr, newStr, -1)
	}
	return content
}

func (m VirtualIDToIDMap) Merge(m2 VirtualIDToIDMap) {
	for k, v := range m2 {
		m[k] = v
	}
}

func GetCollections(ctx context.Context, p *project.Project, cIDs ...uint) ([]*Collection, error) {
	var collections []*Collection
	tx := model.DB(ctx)
	if len(cIDs) > 0 {
		tx = tx.Where("id in (?)", cIDs)
	}
	tx = tx.Where("project_id = ?", p.ID).Order("display_order asc")
	return collections, tx.Find(&collections).Error
}

func BatchDeleteCollections(ctx context.Context, DeletedBy uint, cIDs ...uint) error {
	return model.DB(ctx).Model(&Collection{}).Where("id IN (?)", cIDs).Updates(map[string]interface{}{
		"deleted_by": DeletedBy,
		"deleted_at": time.Now(),
	}).Error
}

func ExportCollections(ctx context.Context, pID string) spec.Collections {
	result := make(spec.Collections, 0)

	collections, err := GetCollections(ctx, &project.Project{ID: pID})
	if err != nil {
		slog.ErrorContext(ctx, "GetDefinitionResponses", "err", err)
		return result
	}

	return exportBuildcollectionsTree(ctx, collections, &Collection{ParentID: 0}, pID)
}

func exportBuildcollectionsTree(ctx context.Context, collections []*Collection, parentCollection *Collection, projectID string) []*spec.Collection {
	collectItems := []*spec.Collection{}

	for _, collection := range collections {
		if collection.ParentID == parentCollection.ID {
			collectItem := &spec.Collection{
				ID:       collection.ID,
				ParentID: collection.ParentID,
				Title:    collection.Title,
				Type:     spec.CollectionType(collection.Type),
			}

			// 将父级的分类名称也加入Tags中
			if parentCollection.ID > 0 {
				if !collectItem.HasTag(parentCollection.Title) {
					collectItem.Tags = append(collectItem.Tags, parentCollection.Title)
				}
			}

			if tags := TagExport(ctx, collection.ID); len(tags) > 0 {
				collectItem.Tags = append(collectItem.Tags, tags...)
			}

			if collection.Type != CategoryType {
				content := []*spec.NodeProxy{}
				if json.Unmarshal([]byte(collection.Content), &content) == nil {
					collectItem.Content = content
				}
			}

			collectItem.Items = exportBuildcollectionsTree(ctx, collections, collection, projectID)
			collectItems = append(collectItems, collectItem)
		}
	}

	return collectItems
}

func GetCollectionHistories(ctx context.Context, c *Collection, start, end time.Time) ([]*CollectionHistory, error) {
	var list []*CollectionHistory
	tx := model.DB(ctx)
	if !start.IsZero() && !end.IsZero() {
		tx = tx.Where("created_at BETWEEN ? AND ?", start, end)
	}
	err := tx.Where("collection_id = ?", c.ID).Order("created_at desc").Find(&list).Error
	return list, err
}

// GetDeletedCollections 获取删除了的 collection
func GetDeletedCollections(ctx context.Context, projectID string) ([]*Collection, error) {
	// 计算三十天前的时间
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var list []*Collection
	tx := model.DB(ctx).Unscoped()
	return list, tx.Where("project_id = ? AND deleted_at >= ? AND type <> ?", projectID, thirtyDaysAgo, CategoryType).Order("deleted_at desc").Find(&list).Error
}

// GetDeletedCollectionsByIDs 通过 ID 获取删除了的 collection
func GetDeletedCollectionsByIDs(ctx context.Context, projectID string, ids []uint) ([]*Collection, error) {
	var list []*Collection
	tx := model.DB(ctx).Unscoped()
	return list, tx.Where("id in ? AND project_id = ? AND  type <> ?", ids, projectID, CategoryType).Find(&list).Error
}

// RestoreCollections 恢复删除的 collection
func RestoreCollections(ctx context.Context, member *team.TeamMember, collections []*Collection) ([]uint, error) {
	var restoreIDs []uint
	if len(collections) == 0 {
		return restoreIDs, nil
	}

	var collectionParentIDs []uint
	for _, c := range collections {
		collectionParentIDs = append(collectionParentIDs, c.ParentID)
	}

	var existParentCollectionIDs []uint
	if err := model.DB(ctx).Model(&Collection{}).Where("id in ?", collectionParentIDs).Pluck("id", &existParentCollectionIDs).Error; err != nil {
		return restoreIDs, err
	}

	var parentIdMap = make(map[uint]bool)
	for _, id := range existParentCollectionIDs {
		parentIdMap[id] = true
	}

	for _, c := range collections {
		updateData := map[string]interface{}{
			"display_order": 0,
			"updated_by":    member.ID,
			"updated_at":    time.Now(),
			"deleted_by":    gorm.Expr("NULL"),
			"deleted_at":    gorm.Expr("NULL"),
		}
		updateFields := []string{"display_order", "updated_by", "updated_at", "deleted_by", "deleted_at"}
		if _, exist := parentIdMap[c.ParentID]; !exist {
			updateData["parent_id"] = 0
			updateFields = append(updateFields, "parent_id")
		}

		ret := model.DB(ctx).Unscoped().Model(&c).Select(updateFields).Updates(updateData)
		if ret.Error == nil {
			restoreIDs = append(restoreIDs, c.ID)
		}
	}

	return restoreIDs, nil
}

func GetByName(ctx context.Context, projectID string, name string) (*Tag, error) {
	t := &Tag{}
	err := model.DB(ctx).Where("project_id = ? and name = ?", projectID, name).Take(t).Error
	return t, err
}

func TagImport(ctx context.Context, projectID string, collectionID uint, tags []string) {
	if len(tags) > 0 {
		for _, tag := range tags {
			t, err := GetByName(ctx, projectID, tag)
			if err != nil {
				t.ProjectID = projectID
				t.Name = tag
				err = t.Create(ctx)
			}

			if err == nil {
				var ttc TagToCollection
				ttc.TagID = t.ID
				ttc.CollectionID = collectionID
				ttc.Create(ctx)
			}
		}
	}
}

func TagExport(ctx context.Context, collectionID uint) []string {
	var tagNames []string

	tagIDs := CollectionToTagID(ctx, collectionID)
	if len(tagIDs) > 0 {
		var tags []Tag
		if err := model.DB(ctx).Where("id IN ?", tagIDs).Find(&tags).Error; err == nil {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}
	}

	return tagNames
}

func CollectionToTagID(ctx context.Context, collectionID uint) []uint {
	var (
		tagIDs  []uint
		records []TagToCollection
	)
	if err := model.DB(ctx).Where("collection_id = ?", collectionID).Find(&records).Error; err != nil {
		for _, v := range records {
			tagIDs = append(tagIDs, v.TagID)
		}
	}

	return tagIDs
}

// TODO 移除表关系后将次方法移动到backend/service/collection_relations/collection_relations.go中。collection_relations中引用了model.collection导致现在移动不了
func GetCollectionContentSpec(ctx context.Context, content string) ([]*spec.NodeProxy, error) {
	specContent := []*spec.NodeProxy{}
	if err := json.Unmarshal([]byte(content), &specContent); err != nil {
		return nil, errors.New("GetCollectionContentSpec unmarshal error")
	}

	return specContent, nil
}

// TODO 移除表关系后将次方法移动到backend/service/collection_relations/collection_relations.go中。collection_relations中引用了model.collection导致现在移动不了
func GetCollectionURLNode(ctx context.Context, content string) (*spec.HTTPNode[spec.HTTPURLNode], error) {
	url := &spec.HTTPNode[spec.HTTPURLNode]{}
	if content == "" {
		return url, nil
	}

	var specContent []*spec.NodeProxy
	if err := json.Unmarshal([]byte(content), &specContent); err != nil {
		return url, errors.New("unmarshal error")
	}

	for _, i := range specContent {
		switch nx := i.Node.(type) {
		case *spec.HTTPNode[spec.HTTPURLNode]:
			url = nx
		}
	}
	if url == nil {
		return url, errors.New("parsing error")
	}

	return url, nil
}

// TODO 移除表关系后将次方法移动到backend/service/collection_relations/collection_relations.go中。collection_relations中引用了model.collection导致现在移动不了
func GetCollectionRequestNode(ctx context.Context, content string) (*spec.HTTPNode[spec.HTTPRequestNode], error) {
	request := &spec.HTTPNode[spec.HTTPRequestNode]{}
	if content == "" {
		return request, nil
	}

	var specContent []*spec.NodeProxy
	if err := json.Unmarshal([]byte(content), &specContent); err != nil {
		return request, errors.New("unmarshal error")
	}

	for _, i := range specContent {
		switch nx := i.Node.(type) {
		case *spec.HTTPNode[spec.HTTPRequestNode]:
			request = nx
		}
	}
	if request == nil {
		return request, errors.New("parsing error")
	}

	return request, nil
}

// TODO 移除表关系后将次方法移动到backend/service/collection_relations/collection_relations.go中。collection_relations中引用了model.collection导致现在移动不了
func GetCollectionResponseNode(ctx context.Context, content string) (*spec.HTTPNode[spec.HTTPResponsesNode], error) {
	response := &spec.HTTPNode[spec.HTTPResponsesNode]{}
	if content == "" {
		return response, nil
	}

	var specContent []*spec.NodeProxy
	if err := json.Unmarshal([]byte(content), &specContent); err != nil {
		return response, errors.New("unmarshal error")
	}

	for _, i := range specContent {
		switch nx := i.Node.(type) {
		case *spec.HTTPNode[spec.HTTPResponsesNode]:
			response = nx
		}
	}
	if response == nil {
		return response, errors.New("parsing error")
	}

	return response, nil
}

// GetTestCases 获取测试用例列表
func GetTestCases(ctx context.Context, projectID string, collectionID uint) ([]*TestCase, error) {
	var list []*TestCase
	return list, model.DB(ctx).Where("project_id = ? and collection_id = ?", projectID, collectionID).Find(&list).Error
}

func DelAllTestCases(ctx context.Context, projectID string) error {
	return model.DB(ctx).Where("project_id = ?", projectID).Delete(&TestCase{}).Error
}

func MemberInfo(ctx context.Context, memberID uint, unscoped bool) (*team.TeamMember, error) {
	var (
		tm *team.TeamMember
		tx *gorm.DB
	)

	if unscoped {
		tx = model.DB(ctx).Unscoped()
	} else {
		tx = model.DB(ctx)
	}

	if err := tx.First(&tm, memberID).Error; err != nil {
		return nil, err
	} else {
		return tm, nil
	}
}

func UserInfo(ctx context.Context, memberID uint, unscoped bool) (*user.User, error) {
	if tm, err := MemberInfo(ctx, memberID, unscoped); err != nil {
		return nil, err
	} else {
		return tm.UserInfo(ctx, unscoped)
	}
}
