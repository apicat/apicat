package route

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/module/storage"
	"github.com/apicat/apicat/v2/backend/route/middleware/dump"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"
	"github.com/apicat/apicat/v2/backend/route/middleware/log"
	"github.com/apicat/apicat/v2/frontend"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

func Init() error {
	conf := config.Get()
	if !conf.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()
	e.ContextWithFallback = true
	e.Use(
		gin.Recovery(),
		log.AccessLog("/assets"),
	)

	// index 使用模版 主要是为了自定义标题
	// 为了解决分享的时候app抓取网页title不正确的问题
	t, _ := template.ParseFS(frontend.Dist, "dist/templates/*.tmpl", "dist/index.html")
	e.SetHTMLTemplate(t)

	// 前端单页路由以及静态资源
	frontendHandle(e)

	if s := config.Get().Storage; s.Driver == storage.LOCAL {
		e.Static("/uploads/", s.LocalDisk.Path)
	}

	// api接口
	g := e.Group("/api",
		ginrpc.ReponseRender(dump.Response),
		ginrpc.RequestBeforeHook(dump.Request),
		// 需要登录 如果无需验证请添加对应的路由前缀
		jwt.JwtUser(jwt.NotAbortPathList{
			// 账号相关API
			{Method: []string{"all"}, Path: "/api/account/*"},
			// 修改邮箱
			{Method: []string{http.MethodPut}, Path: "/api/user/email/:code"},
			// 检查团队邀请token
			{Method: []string{http.MethodGet}, Path: "/api/teams/invitation-tokens/:invitationToken"},
			// 项目详情
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID"},
			// 项目server列表
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/servers"},
			// 公共参数列表
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/global/parameters"},
			// 集合列表、详情
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/collections"},
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/collections/:collectionID"},
			// 定义模型列表、详情
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/definition/schemas"},
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/definition/schemas/:schemaID"},
			// 定义响应列表、详情
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/definition/responses"},
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/definition/responses/:responseID"},
			// 项目分享状态
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/share/status"},
			// 检查项目分享密钥
			{Method: []string{http.MethodPost}, Path: "/api/projects/:projectID/share/check"},
			// 集合分享状态
			{Method: []string{http.MethodGet}, Path: "/api/collections/:collectionPublicID/share/status"},
			// 检查集合分享密钥
			{Method: []string{http.MethodPost}, Path: "/api/projects/:projectID/collections/:collectionID/share/check"},
			// 导出项目
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/export/:code"},
			// 导出集合
			{Method: []string{http.MethodGet}, Path: "/api/projects/:projectID/collections/:collectionID/export/:code"},
			// mock
			{Method: []string{"all"}, Path: "/api/mock/:projectID/*"},
			// Get GitHub client id
			{Method: []string{http.MethodGet}, Path: "/api/sysconfigs/github"},
		}),
	)

	registerUser(g)
	registerAccount(g)
	registerTeam(g)
	registerTeamMember(g)
	registerProjectGroup(g)
	registerProject(g)
	registerProjectShare(g)
	registerProjectGlobalParameter(g)
	registerProjectServer(g)
	registerProjectMember(g)
	registerProjectDefinitionSchema(g)
	registerProjectDefinitionSchemaHistory(g)
	registerProjectDefinitionResponse(g)
	registerCollection(g)
	registerCollectionMock(g)
	registerCollectionShare(g)
	registerCollectionHistory(g)
	registerTestCase(g)
	registerIteration(g)
	registerOauthSysconfig(g)
	registerServiceSysconfig(g)
	registerStorageSysconfig(g)
	registerEmailSysconfig(g)
	registerModelSysconfig(g)
	registerJsonSchema(g)

	slog.Info("init router", "bind", conf.App.AppServerBind)
	return e.Run(conf.App.AppServerBind)
}
