package relations

import (
	"github.com/apicat/apicat/v2/backend/model/user"
	protouserresponse "github.com/apicat/apicat/v2/backend/route/proto/user/response"

	"github.com/gin-gonic/gin"
)

// ConvertModelUser 将model.user转换为proto.user
func ConvertModelUser(ctx *gin.Context, src *user.User) protouserresponse.User {
	var dest protouserresponse.User
	dest.ID = src.ID
	dest.Name = src.Name
	dest.Email = src.Email
	dest.Avatar = src.Avatar
	dest.Role = src.Role
	dest.Language = src.Language

	oauths, err := src.Oauths(ctx, "github")
	if err == nil && len(oauths) > 0 {
		dest.Github = true
	}
	return dest
}
