package access

import (
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/route/middleware/jwt"
	"net/http"

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
