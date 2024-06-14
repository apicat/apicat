package collection

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	collection_proto "github.com/apicat/apicat/v2/backend/route/proto/collection"
	"github.com/apicat/apicat/v2/backend/route/proto/collection/base"
	"github.com/apicat/apicat/v2/backend/route/proto/collection/request"
	"github.com/apicat/apicat/v2/backend/route/proto/collection/response"
	"github.com/apicat/apicat/v2/backend/service/ai"
	"github.com/apicat/apicat/v2/backend/service/relations"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type testCaseApiImpl struct{}

func NewTestCaseApi() collection_proto.TestCaseApi {
	return &testCaseApiImpl{}
}

func (ts *testCaseApiImpl) Generate(ctx *gin.Context, opt *request.GenerateTestCaseOption) (*ginrpc.Empty, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfPM.ProjectID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.GenerationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("collection.DoesNotExist"))
	}
	sc, err := relations.CollectionDerefWithSpec(ctx, c)
	if err != nil {
		slog.ErrorContext(ctx, "collectionDerefWithApiCatSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.GenerationFailed"))
	}

	apiSummary, err := ai.APISummarize(ctx, sc)
	if err != nil {
		slog.ErrorContext(ctx, "ai.APISummarize", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.GenerationFailed"))
	}

	TestCaseRecordList := make([]string, 0)
	if opt.Regenerate {
		collection.DelAllTestCases(ctx, selfPM.ProjectID)
	} else {
		if testCases, err := collection.GetTestCases(ctx, selfPM.ProjectID, opt.CollectionID); err == nil {
			for _, testCase := range testCases {
				TestCaseRecordList = append(TestCaseRecordList, testCase.Title)
			}
		}
	}

	testCaseList, err := ai.TestCaseListGenerate(ctx, apiSummary, TestCaseRecordList, opt.Prompt)
	if err != nil {
		slog.ErrorContext(ctx, "ai.TestCaseListGenerate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.GenerationFailed"))
	}

	go ai.TestCaseGenerate(jwt.GetUser(ctx).Language, selfPM.ProjectID, apiSummary, c.ID, testCaseList)

	return nil, nil
}

func (ts *testCaseApiImpl) List(ctx *gin.Context, opt *base.ProjectCollectionIDOption) (*response.TestCaseList, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	testCaseList := make([]*response.TestCase, 0)
	if testCases, err := collection.GetTestCases(ctx, selfPM.ProjectID, opt.CollectionID); err == nil {
		for _, testCase := range testCases {
			testCaseList = append(testCaseList, &response.TestCase{
				Title: testCase.Title,
				IdCreateTimeInfo: protobase.IdCreateTimeInfo{
					ID:        testCase.ID,
					CreatedAt: testCase.CreatedAt.Unix(),
				},
			})
		}
	}
	return &response.TestCaseList{
		Generating: ai.TestCaseGeneratingStatus(selfPM.ProjectID, opt.CollectionID),
		Records:    testCaseList,
	}, nil
}

func (ts *testCaseApiImpl) Get(ctx *gin.Context, opt *request.GetTestCaseOption) (*response.TestCaseDetail, error) {
	t := &collection.TestCase{
		ID: opt.TestCaseID,
	}
	exist, err := t.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "t.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("testCase.DoesNotExist"))
	}
	return &response.TestCaseDetail{
		ID:      t.ID,
		Title:   t.Title,
		Content: t.Content,
	}, nil
}

func (ts *testCaseApiImpl) Regenerate(ctx *gin.Context, opt *request.RegenerateTestCaseOption) (*response.TestCaseDetail, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfPM.ProjectID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	t := &collection.TestCase{
		ID: opt.TestCaseID,
	}
	exist, err = t.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "t.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("testCase.DoesNotExist"))
	}

	sc, err := relations.CollectionDerefWithSpec(ctx, c)
	if err != nil {
		slog.ErrorContext(ctx, "collectionDerefWithApiCatSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}

	apiSummary, err := ai.APISummarize(ctx, sc)
	if err != nil {
		slog.ErrorContext(ctx, "ai.APISummarize", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}

	result, err := ai.TestCaseDetailRegenerate(t, jwt.GetUser(ctx).Language, apiSummary, opt.Prompt)
	if err != nil {
		slog.ErrorContext(ctx, "ai.TestCaseDetailRegenerate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}

	content := fmt.Sprintf(
		"`%s`<br>\n>%s\n### %s\n%s\n### %s\n%s\n### %s\n%s",
		result.Type,
		result.Description,
		i18n.NewTran("testCase.Steps").Translate(ctx),
		result.Steps,
		i18n.NewTran("testCase.Input").Translate(ctx),
		result.Input,
		i18n.NewTran("testCase.Output").Translate(ctx),
		result.Output,
	)

	if err := t.Update(ctx, result.Purpose, content); err != nil {
		slog.ErrorContext(ctx, "t.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}

	return &response.TestCaseDetail{
		ID:      t.ID,
		Title:   t.Title,
		Content: t.Content,
	}, nil
}

func (ts *testCaseApiImpl) Delete(ctx *gin.Context, opt *request.DeleteTestCaseOption) (*ginrpc.Empty, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	t := &collection.TestCase{
		ID: opt.TestCaseID,
	}
	exist, err := t.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "t.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.FailedToGet"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("testCase.DoesNotExist"))
	}

	if err := t.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "t.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.FailedToDelete"))
	}
	return &ginrpc.Empty{}, nil
}
