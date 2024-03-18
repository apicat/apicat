package dump

import (
	"log/slog"
	"net/http"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

func Request(ctx *gin.Context, in any, lastBindErr error) error {
	if ginrpc.IsEmpty(in) {
		return nil
	}
	slog.DebugContext(ctx, "dump", "in", in)
	return lastBindErr
}

func Response(ctx *gin.Context, respdata any, err *ginrpc.Error) {
	statusCode := http.StatusOK
	if err != nil {
		resbody := make(map[string]any)
		if !ginrpc.IsEmpty(err.Attrs) {
			for k := range err.Attrs {
				resbody[k] = err.Attrs[k]
			}
		}
		resbody["message"] = transValidErr(ctx, err.Err).Error()
		if err.Code > 0 {
			statusCode = err.Code
		}
		ctx.JSON(statusCode, resbody)
		return
	}

	slog.DebugContext(ctx, "dump", "out", respdata)
	switch ctx.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		statusCode = http.StatusCreated
	case http.MethodDelete:
		statusCode = http.StatusNoContent
	}

	if ginrpc.IsEmpty(respdata) {
		ctx.Status(statusCode)
		return
	}
	ctx.JSON(statusCode, respdata)
}
