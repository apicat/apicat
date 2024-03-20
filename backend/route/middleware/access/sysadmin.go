package access

import (
	"net/http"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func SysAdmin() func(*gin.Context) {
	return func(ctx *gin.Context) {
		u := jwt.GetUser(ctx)
		if u == nil || !u.IsSysAdmin(ctx) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": i18n.NewTran("common.PermissionDenied").Translate(ctx)})
			return
		}
	}
}
