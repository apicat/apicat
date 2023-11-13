package enum

const (
	// 登录令牌无效或不正确
	InvalidOrIncorrectLoginToken = 101
	// 访问令牌（分享密钥）无效或不正确
	InvalidOrIncorrectAccessToken = 201

	// 成员权限不足
	MemberInsufficientPermissionsCode = 101
	// 项目成员权限不足
	ProjectMemberInsufficientPermissionsCode = 201
	// 目标项目成员权限不足
	TargetProjectMemberInsufficientPermissionsCode = 202
	// 跳转403页面
	Redirect403Page = 301

	// 跳转404页面
	Redirect404Page = 101
	// 显示404错误信息
	Display404ErrorMessage = 201
)
