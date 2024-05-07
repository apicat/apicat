package collection

import (
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/project"
	collectionrequest "github.com/apicat/apicat/v2/backend/route/proto/collection/request"
	"github.com/apicat/apicat/v2/backend/service/relations"

	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec"

	"github.com/gin-gonic/gin"
)

func Mock(ctx *gin.Context) {
	// 解析和校验 URI 中的参数
	opt := &collectionrequest.GetMockOption{}
	if err := ctx.ShouldBindUri(&opt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	p := &project.Project{ID: opt.ProjectID}
	exist, err := p.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "p.Get", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": i18n.NewErr("mock.FailedToMock").Error()})
		return
	}
	if !exist {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": i18n.NewErr("project.DoesNotExist").Error()})
		return
	}

	methodDict := map[string]string{
		"GET":     "get",
		"POST":    "post",
		"PUT":     "put",
		"PATCH":   "patch",
		"DELETE":  "delete",
		"OPTIONS": "options",
	}
	method, ok := methodDict[ctx.Request.Method]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": i18n.NewErr("mock.FailedToMock").Error()})
		return
	}

	c := &collection.Collection{ProjectID: opt.ProjectID, Path: fmt.Sprintf("/%s", strings.TrimPrefix(opt.Path, "/")), Method: method}
	exist, err = c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": i18n.NewErr("mock.FailedToMock").Error()})
		return
	}
	if !exist {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": i18n.NewErr("collection.DoesNotExist").Error()})
		return
	}

	collectionSpec, err := relations.CollectionDerefWithSpec(ctx, c)
	if err != nil {
		slog.ErrorContext(ctx, "collectionDerefWithSpec", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": i18n.NewErr("mock.FailedToMock").Error()})
		return
	}

	var resp *spec.HTTPNode[spec.HTTPResponsesNode]
	for _, i := range collectionSpec.Content {
		switch nx := i.Node.(type) {
		case *spec.HTTPNode[spec.HTTPResponsesNode]:
			resp = nx
		}
	}
	if resp == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": i18n.NewErr("mock.FailedToMock").Error()})
		return
	}

	ctx.JSON(http.StatusOK, &resp.Attrs.List)
}
