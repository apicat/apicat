package collection

import (
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/model/collection"
	"apicat-cloud/backend/model/project"
	"apicat-cloud/backend/route/middleware/access"
	"apicat-cloud/backend/route/middleware/jwt"
	protobase "apicat-cloud/backend/route/proto/base"
	collection_proto "apicat-cloud/backend/route/proto/collection"
	"apicat-cloud/backend/route/proto/collection/base"
	"apicat-cloud/backend/route/proto/collection/request"
	"apicat-cloud/backend/route/proto/collection/response"
	"apicat-cloud/backend/service/ai"
	collectionrelations "apicat-cloud/backend/service/collection_relations"
	"log/slog"
	"net/http"

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
	sc, err := collectionrelations.CollectionDerefWithSpec(ctx, c)
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

	sc, err := collectionrelations.CollectionDerefWithSpec(ctx, c)
	if err != nil {
		slog.ErrorContext(ctx, "collectionDerefWithApiCatSpec", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}

	apiSummary, err := ai.APISummarize(ctx, sc)
	if err != nil {
		slog.ErrorContext(ctx, "ai.APISummarize", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}

	var result string
	if opt.Prompt != "" {
		result, err = ai.TestCaseDetailRegenerate(t, jwt.GetUser(ctx).Language, apiSummary, opt.Prompt)
		if err != nil {
			slog.ErrorContext(ctx, "ai.TestCaseDetailRegenerate", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
		}
	} else {
		result, err = ai.TestCaseDetailGenerate(jwt.GetUser(ctx).Language, apiSummary, t.Title)
		if err != nil {
			slog.ErrorContext(ctx, "ai.TestCaseDetailGenerate", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
		}
	}

	if result == "" {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("testCase.RegenerationFailed"))
	}

	if err := t.Update(ctx, t.Title, result); err != nil {
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
