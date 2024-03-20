package user_relations

import (
	protouserresponse "github.com/apicat/apicat/backend/route/proto/user/response"

	"github.com/apicat/apicat/backend/model/user"

	"github.com/gin-gonic/gin"
)

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
