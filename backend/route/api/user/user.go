package user

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/user"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/oauth2"
	"github.com/apicat/apicat/v2/backend/module/oauth2/github"
	"github.com/apicat/apicat/v2/backend/module/storage"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protouser "github.com/apicat/apicat/v2/backend/route/proto/user"
	protouserbase "github.com/apicat/apicat/v2/backend/route/proto/user/base"
	protouserrequest "github.com/apicat/apicat/v2/backend/route/proto/user/request"
	protouserresponse "github.com/apicat/apicat/v2/backend/route/proto/user/response"
	"github.com/apicat/apicat/v2/backend/service/relations"
	imgutil "github.com/apicat/apicat/v2/backend/utils/image"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type userApiImpl struct{}

func NewUserApi() protouser.UserApi {
	return &userApiImpl{}
}

func (*userApiImpl) GetList(ctx *gin.Context, opt *protouserrequest.UserListOption) (*protouserresponse.UserList, error) {
	if opt.Page <= 0 {
		opt.Page = 1
	}
	if opt.PageSize <= 0 {
		opt.PageSize = 15
	}

	items, err := user.GetUsers(ctx, opt.Page, opt.PageSize, opt.Keywords)
	if err != nil {
		slog.ErrorContext(ctx, "user.GetUsers", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.FailedToGetList"))
	}

	count, err := user.GetUserCount(ctx, opt.Keywords)
	if err != nil {
		slog.ErrorContext(ctx, "user.GetUserCount", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.FailedToGetList"))
	}

	var list = &protouserresponse.UserList{
		PaginationInfo: protobase.PaginationInfo{
			Count:       int(count),
			TotalPage:   int(math.Ceil(float64(count) / float64(opt.PageSize))),
			CurrentPage: opt.Page,
		},
		Items: make([]protouserresponse.User, len(items)),
	}
	for k, v := range items {
		list.Items[k] = relations.ConvertModelUser(ctx, v)
	}
	return list, nil
}

func (ua *userApiImpl) ChangePasswordByAdmin(ctx *gin.Context, opt *protouserrequest.ChangePasswordOption) (*ginrpc.Empty, error) {
	u := user.User{ID: opt.UserID}
	exist, err := u.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "u.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordUpdateFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.DoesNotExist"))
	}

	if u.IsSysAdmin(ctx) {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.PasswordUpdateFailed"))
	}

	u.Password = opt.Password
	if err := u.UpdatePassword(ctx); err != nil {
		slog.ErrorContext(ctx, "u.UpdatePassword", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordUpdateFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (ua *userApiImpl) DelUser(ctx *gin.Context, opt *protouserrequest.UserIDOption) (*ginrpc.Empty, error) {
	u := user.User{ID: opt.UserID}
	exist, err := u.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "u.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.FailedToDelete"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.DoesNotExist"))
	}

	if u.IsSysAdmin(ctx) {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.FailedToDelete"))
	}

	if err := u.Delete(ctx); err != nil {
		slog.ErrorContext(ctx, "u.Delete", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.FailedToDelete"))
	}

	return &ginrpc.Empty{}, nil
}

// GetSelf 当前登录的用户信息
func (*userApiImpl) GetSelf(ctx *gin.Context, _ *ginrpc.Empty) (*protouserresponse.User, error) {
	u := jwt.GetUser(ctx)
	usr := relations.ConvertModelUser(ctx, u)
	return &usr, nil
}

// ChangePassword 修改密码
func (ua *userApiImpl) ChangePassword(ctx *gin.Context, opt *protouserrequest.ChangePwdOption) (*ginrpc.Empty, error) {
	u := jwt.GetUser(ctx)

	ucache, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordUpdateFailed"))
	}

	// 按照用户id设置最大重试次数
	changePasswordTimesKey := fmt.Sprintf("changePassword-%d", u.ID)
	ts, ok, _ := ucache.Get(changePasswordTimesKey)
	var number int
	if ok {
		var err error
		number, err = strconv.Atoi(ts)
		if err != nil {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordUpdateFailed"))
		}
		if number > 10 {
			return nil, ginrpc.NewError(http.StatusTooManyRequests, i18n.NewErr("common.TooManyOperations"))
		}
	}

	_ = ucache.Set(changePasswordTimesKey, strconv.Itoa(number+1), time.Hour)

	if !u.CheckPassword(opt.Password) {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.OriginalPasswordWrong"))
	}
	u.Password = opt.NewPassword
	if err := u.UpdatePassword(ctx); err != nil {
		slog.ErrorContext(ctx, "u.UpdatePassword", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordUpdateFailed"))
	}

	_ = ucache.Del(changePasswordTimesKey)

	return &ginrpc.Empty{}, nil
}

// SetSelf 设置当前用户自身的信息
func (*userApiImpl) SetSelf(ctx *gin.Context, opt *protouserrequest.SetUserSelfOption) (*ginrpc.Empty, error) {
	u := jwt.GetUser(ctx)

	u.Name = opt.Name

	if _, exist := user.SupportedLanguages[opt.Language]; exist {
		u.Language = opt.Language
	}

	if err := u.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "u.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ModificationFailed"))
	}

	return &ginrpc.Empty{}, nil
}

func (ua *userApiImpl) ChangeEmail(ctx *gin.Context, opt *protouserbase.EmailOption) (*ginrpc.Empty, error) {
	u := jwt.GetUser(ctx)
	if u.Email == opt.Email {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.EmailNotChanged"))
	}

	ucache, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.EmailUpdateFailed"))
	}

	// 按照用户id设置最大重试次数
	changeEmailTimesKey := fmt.Sprintf("changeEmail-%d", u.ID)
	ts, ok, _ := ucache.Get(changeEmailTimesKey)
	var number int
	if ok {
		var err error
		number, err = strconv.Atoi(ts)
		if err != nil {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.EmailUpdateFailed"))
		}
		if number > 10 {
			return nil, ginrpc.NewError(http.StatusTooManyRequests, i18n.NewErr("common.TooManyOperations"))
		}
	}

	_ = ucache.Set(changeEmailTimesKey, strconv.Itoa(number+1), time.Hour)

	usr := &user.User{
		Email: opt.Email,
	}
	exist, err := usr.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "usr.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.EmailUpdateFailed"))
	}
	if exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.EmailHasBeenUsed"))
	}

	u.Email = opt.Email
	if err := u.UpdateEmail(ctx); err != nil {
		slog.ErrorContext(ctx, "u.UpdateEmail", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.EmailUpdateFailed"))
	}
	return &ginrpc.Empty{}, nil
}

// UploadAvatar 上传头像
func (*userApiImpl) UploadAvatar(ctx *gin.Context, opt *protouserrequest.UploadAvatarOption) (*protouserbase.AvatarOption, error) {
	u := jwt.GetUser(ctx)

	if opt.Avatar.Size > 1024*1024*2 {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.ImageTooLarge"))
	}

	img, fileExt, err := imgutil.FileHeaderToImage(opt.Avatar)
	if err != nil {
		slog.ErrorContext(ctx, "imageOpt.FileHeaderToImage", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.ImageUploadFailed"))
	}

	fileName := fmt.Sprintf("%s/%x%s", "avatars", md5.Sum([]byte(fmt.Sprintf("%d_%d", u.ID, time.Now().Unix()))), fileExt)

	// 裁剪图片
	croppedFileBytes, err := imgutil.Cropping(img, opt.CroppedX, opt.CroppedY, opt.CroppedWidth, opt.CroppedHeight)
	if err != nil {
		slog.ErrorContext(ctx, "imageOpt.Cropping", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("common.ImageUploadFailed"))
	}

	helper := storage.NewStorage(config.Get().Storage.ToCfg())
	contentType := http.DetectContentType(croppedFileBytes)
	path, err := helper.PutObject(fileName, croppedFileBytes, contentType)
	if err != nil {
		slog.ErrorContext(ctx, "helper.PutObject", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ImageUploadFailed"))
	}

	u.Avatar = path
	if err := u.Update(ctx); err != nil {
		slog.ErrorContext(ctx, "u.Update", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.ImageUploadFailed"))
	}

	return &protouserbase.AvatarOption{
		Avatar: path,
	}, nil
}

func (ui *userApiImpl) OauthConnect(ctx *gin.Context, opt *protouserrequest.OauthOption) (*ginrpc.Empty, error) {
	var (
		oauthUser *oauth2.AuthUser
		err       error
	)

	oauthMap := config.Get().Oauth2
	cfg, ok := oauthMap[opt.Type]
	if !ok {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("user.NotSupportOauth", opt.Type))
	} else {
		// Now there is only github
		oauthObj := oauth2.NewObject(cfg, &github.Github{})
		oauthUser, err = oauthObj.GetUserByState(ctx, opt.Code)
		if err != nil {
			slog.ErrorContext(ctx, "oauthObj.GetUserByState", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.OauthConnectFailed", opt.Type))
		}
	}

	usr, err := user.GetUserByOauth(ctx, oauthUser.ID, opt.Type)
	if err != nil {
		slog.ErrorContext(ctx, "user.GetUserByOauth", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.OauthConnectFailed", opt.Type))
	}
	if usr != nil {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.OauthConnectRepeat", opt.Type))
	}

	selfUser := jwt.GetUser(ctx)
	if err := selfUser.BindOrRecoverOauth(ctx, opt.Type, oauthUser.ID); err != nil {
		slog.ErrorContext(ctx, "selfUser.BindOrRecoverOauth", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.OauthConnectFailed", opt.Type))
	}

	return &ginrpc.Empty{}, nil
}

func (ui *userApiImpl) OauthDisconnect(ctx *gin.Context, opt *protouserbase.OauthTypeOption) (*ginrpc.Empty, error) {
	oauthMap := config.Get().Oauth2
	_, ok := oauthMap[opt.Type]
	if !ok {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("user.OauthDisconnectFailed", opt.Type))
	}

	u := jwt.GetUser(ctx)
	if err := u.UnBindOauth(ctx, opt.Type); err != nil {
		slog.ErrorContext(ctx, "u.UnBindOauth", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.OauthDisconnectFailed", opt.Type))
	}
	return &ginrpc.Empty{}, nil
}
