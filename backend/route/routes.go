package route

import (
	"github.com/apicat/apicat/v2/backend/route/api/collection"
	"github.com/apicat/apicat/v2/backend/route/api/iteration"
	"github.com/apicat/apicat/v2/backend/route/api/project"
	"github.com/apicat/apicat/v2/backend/route/api/sysconfig"
	"github.com/apicat/apicat/v2/backend/route/api/team"
	"github.com/apicat/apicat/v2/backend/route/api/user"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

func registerUser(g *gin.RouterGroup) {
	srv := user.NewUserApi()
	g.GET("/users", access.SysAdmin(), ginrpc.Handle(srv.GetList))
	g.PATCH("/users/:userID", access.SysAdmin(), ginrpc.Handle(srv.ChangePasswordByAdmin))
	g.DELETE("/users/:userID", access.SysAdmin(), ginrpc.Handle(srv.DelUser))

	r := g.Group("/user")
	r.GET("", ginrpc.Handle(srv.GetSelf))
	r.PUT("", ginrpc.Handle(srv.SetSelf))
	r.PUT("/password", ginrpc.Handle(srv.ChangePassword))
	r.POST("/email", ginrpc.Handle(srv.SendChangeEmail))
	r.PUT("/email/:code", ginrpc.Handle(srv.ChangeEmailFire))
	r.POST("/avatar", ginrpc.Handle(srv.UploadAvatar))
	r.POST("/oauth/:type/connect", ginrpc.Handle(srv.OauthConnect))
	r.DELETE("/oauth/:type/disconnect", ginrpc.Handle(srv.OauthDisconnect))
}

func registerAccount(g *gin.RouterGroup) {
	srv := user.NewAccountApi()
	r := g.Group("/account")
	r.POST("/login", ginrpc.Handle(srv.Login))
	r.POST("/register", ginrpc.Handle(srv.Register))
	r.POST("/oauth/:type/login", ginrpc.Handle(srv.LoginWithOauthCode))
	r.PUT("/email-verification/:code", ginrpc.Handle(srv.RegisterFire))
	r.POST("/retrieve-password", ginrpc.Handle(srv.SendResetPasswordMail))
	r.GET("/reset-password/check/:code", ginrpc.Handle(srv.ResetPasswordCheck))
	r.PUT("/reset-password/:code", ginrpc.Handle(srv.ResetPassword))
}

func registerTeam(g *gin.RouterGroup) {
	srv := team.NewTeamApi()
	r := g.Group("/teams")
	r.POST("", ginrpc.Handle(srv.Create))
	r.GET("", ginrpc.Handle(srv.TeamList))
	r.GET("/current", ginrpc.Handle(srv.Current))
	r.POST("/join", ginrpc.Handle(srv.Join))
	r.GET("/invitation-tokens/:invitationToken", ginrpc.Handle(srv.CheckInvitationToken))

	t := g.Group("/teams/:teamID", access.BelongToTeam())
	t.GET("", ginrpc.Handle(srv.Get))
	t.PUT("/switch", ginrpc.Handle(srv.Switch))
	t.GET("/invitation-tokens", ginrpc.Handle(srv.GetInvitationToken))
	t.PUT("/invitation-tokens", ginrpc.Handle(srv.ResetInvitationToken))
	t.PUT("/setting", ginrpc.Handle(srv.Setting))
	t.PUT("/transfer", ginrpc.Handle(srv.Transfer))
	t.DELETE("", ginrpc.Handle(srv.Delete))
}

func registerTeamMember(g *gin.RouterGroup) {
	srv := team.NewTeamMemberApi()
	r := g.Group("/teams/:teamID/members", access.BelongToTeam())
	r.GET("", ginrpc.Handle(srv.MemberList))
	r.PUT("/:memberID", ginrpc.Handle(srv.UpdateMember))
	r.DELETE("/:memberID", ginrpc.Handle(srv.DeleteMember))
	r.DELETE("", ginrpc.Handle(srv.Quit))
}

func registerProjectGroup(g *gin.RouterGroup) {
	srv := project.NewProjectGroupApi()
	r := g.Group("/teams/:teamID/project-groups", access.BelongToTeam())
	r.POST("", ginrpc.Handle(srv.Create))
	r.GET("", ginrpc.Handle(srv.List))
	r.PUT("/sort", ginrpc.Handle(srv.Sort))

	gp := g.Group("/project-groups/:groupID", access.BelongToTeam())
	gp.DELETE("", ginrpc.Handle(srv.Delete))
	gp.PUT("", ginrpc.Handle(srv.Rename))
}

func registerProject(g *gin.RouterGroup) {
	srv := project.NewProjectApi()

	r := g.Group("/teams/:teamID/projects", access.BelongToTeam())
	r.POST("", ginrpc.Handle(srv.Create))
	r.GET("", ginrpc.Handle(srv.List))

	noAuth := g.Group("/projects/:projectID")
	noAuth.GET("", access.AllowGuestByShareCode(), ginrpc.Handle(srv.Get))
	// 导出项目内容，需返回不同的 Content-Type，单独处理
	noAuth.GET("/export/:code", project.Export)

	p := g.Group("/projects/:projectID", access.BelongToTeam(), access.BelongToProject())
	p.PUT("", ginrpc.Handle(srv.Setting))
	p.DELETE("", ginrpc.Handle(srv.Delete))
	p.PUT("/group", ginrpc.Handle(srv.ChangeGroup))
	p.POST("/follow", ginrpc.Handle(srv.Follow))
	p.DELETE("/follow", ginrpc.Handle(srv.UnFollow))
	p.PUT("/transfer", ginrpc.Handle(srv.Transfer))
	p.DELETE("/exit", ginrpc.Handle(srv.Exit))
	p.GET("/export", ginrpc.Handle(srv.GetExportPath))
}

func registerProjectShare(g *gin.RouterGroup) {
	srv := project.NewProjectShareApi()

	noAuth := g.Group("/projects/:projectID/share")
	noAuth.GET("/status", access.AllowGuestByShareStatus(), ginrpc.Handle(srv.Status))
	noAuth.POST("/check", ginrpc.Handle(srv.Check))

	r := g.Group("/projects/:projectID/share", access.BelongToTeam(), access.BelongToProject())
	r.GET("", ginrpc.Handle(srv.Detail))
	r.PUT("", ginrpc.Handle(srv.Switch))
	r.PUT("/reset", ginrpc.Handle(srv.Reset))
}

func registerProjectGlobalParameter(g *gin.RouterGroup) {
	srv := project.NewGlobalParameterAPI()

	noAuth := g.Group("/projects/:projectID/global/parameters", access.AllowGuestByShareCode())
	noAuth.GET("", ginrpc.Handle(srv.List))

	r := g.Group("/projects/:projectID/global/parameters", access.BelongToTeam(), access.BelongToProject())
	r.POST("", ginrpc.Handle(srv.Create))
	r.PUT("/:parameterID", ginrpc.Handle(srv.Update))
	r.DELETE("/:parameterID", ginrpc.Handle(srv.Delete))
	r.PUT("/sort", ginrpc.Handle(srv.Sort))
}

func registerProjectServer(g *gin.RouterGroup) {
	srv := project.NewProjectServerApi()

	noAuth := g.Group("/projects/:projectID/servers", access.AllowGuestByShareCode())
	noAuth.GET("", ginrpc.Handle(srv.List))

	r := g.Group("/projects/:projectID/servers", access.BelongToTeam(), access.BelongToProject())
	r.POST("", ginrpc.Handle(srv.Create))
	r.PUT("/:serverID", ginrpc.Handle(srv.Update))
	r.DELETE("/:serverID", ginrpc.Handle(srv.Delete))
	r.PUT("/sort", ginrpc.Handle(srv.Sort))
}

func registerProjectMember(g *gin.RouterGroup) {
	srv := project.NewProjectMemberApi()
	r := g.Group("/projects/:projectID/members", access.BelongToTeam(), access.BelongToProject())
	r.POST("", ginrpc.Handle(srv.Create))
	r.GET("", ginrpc.Handle(srv.Members))
	r.GET("/notin", ginrpc.Handle(srv.NotInProjectMembers))
	r.PUT("/:memberID", ginrpc.Handle(srv.Update))
	r.DELETE("/:memberID", ginrpc.Handle(srv.Delete))
}

func registerProjectDefinitionSchema(g *gin.RouterGroup) {
	srv := project.NewDefinitionSchemaApi()

	noAuth := g.Group("/projects/:projectID/definition/schemas", access.AllowGuestByShareCode())
	noAuth.GET("", ginrpc.Handle(srv.List))
	noAuth.GET("/:schemaID", ginrpc.Handle(srv.Get))

	g.POST("/projects/:projectID/definition/ai/schemas", access.BelongToTeam(), access.BelongToProject(), ginrpc.Handle(srv.AIGenerate))
	r := g.Group("/projects/:projectID/definition/schemas", access.BelongToTeam(), access.BelongToProject())
	r.POST("", ginrpc.Handle(srv.Create))
	r.PUT("/:schemaID", ginrpc.Handle(srv.Update))
	r.DELETE("/:schemaID", ginrpc.Handle(srv.Delete))
	r.PUT("/move", ginrpc.Handle(srv.Move))
	r.POST("/:schemaID/copy", ginrpc.Handle(srv.Copy))
}

func registerProjectDefinitionSchemaHistory(g *gin.RouterGroup) {
	srv := project.NewDefinitionSchemaHistoryApi()

	r := g.Group("/projects/:projectID/definition/schemas/:schemaID/histories", access.BelongToTeam(), access.BelongToProject())
	r.GET("", ginrpc.Handle(srv.List))
	r.GET("/:historyID", ginrpc.Handle(srv.Get))
	r.PUT("/:historyID/restore", ginrpc.Handle(srv.Restore))
	r.GET("/diff", ginrpc.Handle(srv.Diff))
}

func registerProjectDefinitionResponse(g *gin.RouterGroup) {
	srv := project.NewDefinitionResponseApi()

	noAuth := g.Group("/projects/:projectID/definition/responses", access.AllowGuestByShareCode())
	noAuth.GET("", ginrpc.Handle(srv.List))
	noAuth.GET("/:responseID", ginrpc.Handle(srv.Get))

	r := g.Group("/projects/:projectID/definition/responses", access.BelongToTeam(), access.BelongToProject())
	r.POST("", ginrpc.Handle(srv.Create))
	r.PUT("/:responseID", ginrpc.Handle(srv.Update))
	r.DELETE("/:responseID", ginrpc.Handle(srv.Delete))
	r.PUT("/move", ginrpc.Handle(srv.Move))
	r.POST("/:responseID/copy", ginrpc.Handle(srv.Copy))
}

func registerCollection(g *gin.RouterGroup) {
	srv := collection.NewCollectionApi()
	noAuth := g.Group("/projects/:projectID/collections")
	noAuth.GET("", access.AllowGuestByShareCode(), ginrpc.Handle(srv.List))
	noAuth.GET("/:collectionID", access.AllowGuestByShareCode(), ginrpc.Handle(srv.Get))
	// 导出集合内容，需返回不同的 Content-Type，单独处理
	noAuth.GET("/:collectionID/export/:code", collection.Export)

	g.POST("/projects/:projectID/ai/collections", access.BelongToTeam(), access.BelongToProject(), ginrpc.Handle(srv.AIGenerate))
	r := g.Group("/projects/:projectID/collections", access.BelongToTeam(), access.BelongToProject())
	r.POST("", ginrpc.Handle(srv.Create))
	r.PUT("/:collectionID", ginrpc.Handle(srv.Update))
	r.DELETE("/:collectionID", ginrpc.Handle(srv.Delete))
	r.PUT("/move", ginrpc.Handle(srv.Move))
	r.POST("/:collectionID/copy", ginrpc.Handle(srv.Copy))
	r.GET("/trashes", ginrpc.Handle(srv.Trashes))
	r.PUT("/restore", ginrpc.Handle(srv.Restore))
	r.GET("/:collectionID/export", ginrpc.Handle(srv.GetExportPath))
}

func registerCollectionMock(g *gin.RouterGroup) {
	g.Any("/mock/:projectID/*path", collection.Mock)
}

func registerCollectionShare(g *gin.RouterGroup) {
	srv := collection.NewCollectionShareApi()

	g.GET("/collections/:collectionPublicID/share/status", ginrpc.Handle(srv.Status))
	g.POST("/projects/:projectID/collections/:collectionID/share/check", ginrpc.Handle(srv.Check))

	r := g.Group("/projects/:projectID/collections/:collectionID/share", access.BelongToTeam(), access.BelongToProject())
	r.GET("", ginrpc.Handle(srv.Detail))
	r.PUT("", ginrpc.Handle(srv.Switch))
	r.PUT("/reset", ginrpc.Handle(srv.Reset))
}

func registerCollectionHistory(g *gin.RouterGroup) {
	srv := collection.NewCollectionHistoryApi()

	r := g.Group("/projects/:projectID/collections/:collectionID/histories", access.BelongToTeam(), access.BelongToProject())
	r.GET("", ginrpc.Handle(srv.List))
	r.GET("/:historyID", ginrpc.Handle(srv.Get))
	r.PUT("/:historyID/restore", ginrpc.Handle(srv.Restore))
	r.GET("/diff", ginrpc.Handle(srv.Diff))
}

func registerTestCase(g *gin.RouterGroup) {
	srv := collection.NewTestCaseApi()

	r := g.Group("/projects/:projectID/collections/:collectionID/testcases", access.BelongToTeam(), access.BelongToProject())
	r.POST("", ginrpc.Handle(srv.Generate))
	r.GET("", ginrpc.Handle(srv.List))
	r.GET("/:testCaseID", ginrpc.Handle(srv.Get))
	r.PUT("/:testCaseID", ginrpc.Handle(srv.Regenerate))
	r.DELETE("/:testCaseID", ginrpc.Handle(srv.Delete))
}

func registerIteration(g *gin.RouterGroup) {
	srv := iteration.NewIterationApi()

	r := g.Group("/teams/:teamID/iterations", access.BelongToTeam())
	r.POST("", ginrpc.Handle(srv.Create))
	r.GET("", ginrpc.Handle(srv.List))

	i := g.Group("/iterations/:iterationID", access.BelongToTeam(), access.BelongToProject())
	i.GET("", ginrpc.Handle(srv.Get))
	i.PUT("", ginrpc.Handle(srv.Update))
	i.DELETE("", ginrpc.Handle(srv.Delete))
}

func registerOauthSysconfig(g *gin.RouterGroup) {
	srv := sysconfig.NewOauthApi()
	g.GET("/sysconfigs/github", ginrpc.Handle(srv.GetGithubClientID))
	g.GET("/sysconfigs/oauth", access.SysAdmin(), ginrpc.Handle(srv.Get))
	g.PUT("/sysconfigs/oauth", access.SysAdmin(), ginrpc.Handle(srv.Update))
}

func registerServiceSysconfig(g *gin.RouterGroup) {
	srv := sysconfig.NewServiceApi()
	g.GET("/sysconfigs/service", access.SysAdmin(), ginrpc.Handle(srv.Get))
	g.PUT("/sysconfigs/service", access.SysAdmin(), ginrpc.Handle(srv.Update))
}

func registerStorageSysconfig(g *gin.RouterGroup) {
	srv := sysconfig.NewStorageApi()
	g.GET("/sysconfigs/storages", access.SysAdmin(), ginrpc.Handle(srv.Get))
	g.PUT("/sysconfigs/storages/disk", access.SysAdmin(), ginrpc.Handle(srv.UpdateDisk))
	g.PUT("/sysconfigs/storages/cloudflare", access.SysAdmin(), ginrpc.Handle(srv.UpdateCloudflare))
	g.PUT("/sysconfigs/storages/qiniu", access.SysAdmin(), ginrpc.Handle(srv.UpdateQiniu))
}

func registerEmailSysconfig(g *gin.RouterGroup) {
	srv := sysconfig.NewEmailApi()
	g.GET("/sysconfigs/emails", access.SysAdmin(), ginrpc.Handle(srv.Get))
	g.PUT("/sysconfigs/emails/smtp", access.SysAdmin(), ginrpc.Handle(srv.UpdateSMTP))
	g.PUT("/sysconfigs/emails/sendcloud", access.SysAdmin(), ginrpc.Handle(srv.UpdateSendCloud))
}

func registerModelSysconfig(g *gin.RouterGroup) {
	srv := sysconfig.NewModelApi()
	g.GET("/sysconfigs/models", access.SysAdmin(), ginrpc.Handle(srv.Get))
	g.PUT("/sysconfigs/models/openai", access.SysAdmin(), ginrpc.Handle(srv.UpdateOpenAI))
	g.PUT("/sysconfigs/models/azure-openai", access.SysAdmin(), ginrpc.Handle(srv.UpdateAzureOpenAI))
}
