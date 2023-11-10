package route

import (
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/route/api/ai"
	"github.com/apicat/apicat/backend/route/api/collection"
	"github.com/apicat/apicat/backend/route/api/config"
	"github.com/apicat/apicat/backend/route/api/definition"
	"github.com/apicat/apicat/backend/route/api/doc"
	"github.com/apicat/apicat/backend/route/api/global"
	"github.com/apicat/apicat/backend/route/api/iteration"
	"github.com/apicat/apicat/backend/route/api/mock"
	"github.com/apicat/apicat/backend/route/api/project"
	"github.com/apicat/apicat/backend/route/api/server"
	"github.com/apicat/apicat/backend/route/api/trash"
	"github.com/apicat/apicat/backend/route/api/user"
	"github.com/apicat/apicat/backend/route/middleware/check"
	"github.com/apicat/apicat/backend/route/middleware/db"
	"github.com/apicat/apicat/backend/route/middleware/jwt"
	"github.com/apicat/apicat/backend/route/middleware/log"
	"github.com/apicat/apicat/frontend"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(r *gin.Engine) {

	r.Use(
		log.RequestIDLog("/assets/", "/static/"),
		db.CheckDBConnStatus("/assets/", "/static/"),
		translator.UseValidatori18n(),
		gin.Recovery(),
	)

	assets, err := fs.Sub(frontend.FrontDist, "dist/assets")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/assets", http.FS(assets))

	static, err := fs.Sub(frontend.FrontDist, "dist/static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(static))

	r.GET("/", func(ctx *gin.Context) {
		ctx.FileFromFS("dist/", http.FS(frontend.FrontDist))
	})

	mocksrv := mock.NewMockServer()

	r.Any("/mock/:id/*path", mocksrv.Handler)

	configs := r.Group("/config")
	{
		configs.GET("/db", config.GetDBConfig)
	}

	apiRouter := r.Group("/api")
	{
		configs := apiRouter.Group("/config")
		{
			configs.PUT("/db", config.SetDBConfig)
		}

		// 未登录状态下可访问的API
		notLogin := apiRouter.Group("")
		{
			account := notLogin.Group("/account")
			{
				account.POST("/login/email", user.EmailLogin)
				account.POST("/register/email", user.EmailRegister)
			}

			projects := notLogin.Group("/projects")
			projects.Use(check.CheckProject())
			{
				projects.GET("/:project-id/data", project.ProjectDataGet)
				projects.GET("/:project-id/share/status", check.CheckMemberHalfLogin(), project.ProjectShareStatus)
				projects.POST("/:project-id/share/check", project.ProjectShareSecretkeyCheck)
			}

			collections := notLogin.Group("/projects/:project-id/collections")
			collections.Use(check.CheckProject())
			{
				collections.GET("/:collection-id/data", check.CheckCollection(), collection.CollectionDataGet)
				collections.POST("/:collection-id/share/check", check.CheckCollection(), doc.DocShareCheck)
			}

			collectionShare := notLogin.Group("/collections")
			{
				collectionShare.GET("/:public_collection_id/share/status", doc.DocShareStatus)
			}
		}

		// 半登录状态下可访问的API。半登录：登录或不登录时都可访问，但响应的参数不同
		halfLogin := apiRouter.Group("")
		halfLogin.Use(check.CheckProject(), check.CheckMemberHalfLogin(), check.CheckProjectMemberHalfLogin(), mocksrv.ClearCache())
		{
			projects := halfLogin.Group("/projects/:project-id")
			{
				projects.GET("", project.ProjectsGet)
			}

			servers := halfLogin.Group("/projects/:project-id/servers")
			{
				servers.GET("", server.UrlList)
			}

			definitionSchemas := halfLogin.Group("/projects/:project-id/definition/schemas")
			{
				definitionSchemas.GET("/:schemas-id", check.CheckDefinitionSchema(), definition.DefinitionSchemasGet)
				definitionSchemas.GET("", definition.DefinitionSchemasList)
			}

			collections := halfLogin.Group("/projects/:project-id/collections")
			{
				collections.GET("", collection.CollectionsList)
				collections.GET("/:collection-id", check.CheckCollection(), collection.CollectionsGet)
			}

			globalParameters := halfLogin.Group("/projects/:project-id/global/parameters")
			{
				globalParameters.GET("", global.GlobalParametersList)
			}

			definitionResponses := halfLogin.Group("/projects/:project-id/definition/responses")
			{
				definitionResponses.GET("", definition.DefinitionResponsesList)
				definitionResponses.GET("/:response-id", definition.DefinitionResponsesDetail)
			}
		}

		// 仅登录状态下可访问的API。仅登录：仅登录了apicat便可访问，一般为项目外的操作
		onlyLogin := apiRouter.Group("")
		onlyLogin.Use(check.CheckMember())
		{
			users := onlyLogin.Group("/user")
			{
				users.GET("/self", user.GetUserInfo)
				users.PUT("/self", user.SetUserInfo)
				users.PUT("/self/password", user.ChangePassword)
			}

			members := onlyLogin.Group("/members")
			{
				members.GET("", project.GetMembers)
				members.POST("/", project.AddMember)
				members.PUT("/:user-id", project.SetMember)
				members.DELETE("/:user-id", project.DeleteMember)
			}

			projects := onlyLogin.Group("/projects")
			{
				projects.GET("", project.ProjectsList)
				projects.POST("", project.ProjectsCreate)
			}

			iterations := onlyLogin.Group("/iterations")
			{
				iterations.GET("", iteration.IterationsList)
				iterations.GET("/:iteration-id", iteration.IterationsDetails)
				iterations.POST("", iteration.IterationsCreate)
				iterations.PUT("/:iteration-id", iteration.IterationsUpdate)
				iterations.DELETE("/:iteration-id", iteration.IterationsDelete)
			}

			projectGroup := onlyLogin.Group("/project_group")
			{
				projectGroup.GET("", project.ProjectGroupList)
				projectGroup.POST("", project.ProjectGroupCreate)
				projectGroup.PUT("/:group_id/rename", project.ProjectGroupRename)
				projectGroup.DELETE("/:group_id", project.ProjectGroupDelete)
				projectGroup.PUT("/order", project.ProjectGroupOrder)
			}
		}

		// 项目内部操作
		projects := apiRouter.Group("/projects/:project-id")
		projects.Use(jwt.JWTAuthMiddleware(), check.CheckProject(), check.CheckProjectMember(), mocksrv.ClearCache())
		{
			projects := projects.Group("")
			{
				projects.PUT("", project.ProjectsUpdate)
				projects.DELETE("", project.ProjectsDelete)
				projects.DELETE("/exit", project.ProjectExit)
				projects.PUT("/transfer", project.ProjectTransfer)
				projects.GET("/share", project.ProjectShareDetails)
				projects.PUT("/share/switch", project.ProjectSharingSwitch)
				projects.PUT("/share/reset", project.ProjectShareReset)
				projects.POST("/follow", project.ProjectFollow)
				projects.DELETE("/follow", project.ProjectUnFollow)
				projects.PUT("/change_group", project.ProjectChangeGroup)
			}

			definitionSchemas := projects.Group("/definition/schemas")
			{
				definitionSchemas.POST("", definition.DefinitionSchemasCreate)
				definitionSchemas.PUT("/:schemas-id", check.CheckDefinitionSchema(), definition.DefinitionSchemasUpdate)
				definitionSchemas.DELETE("/:schemas-id", check.CheckDefinitionSchema(), definition.DefinitionSchemasDelete)
				definitionSchemas.POST("/:schemas-id", check.CheckDefinitionSchema(), definition.DefinitionSchemasCopy)
				definitionSchemas.PUT("/movement", definition.DefinitionSchemasMove)
			}

			servers := projects.Group("/servers")
			{
				servers.PUT("", server.UrlSettings)
			}

			globalParameters := projects.Group("/global/parameters")
			{
				globalParameters.POST("", global.GlobalParametersCreate)
				globalParameters.PUT("/:parameter-id", global.GlobalParametersUpdate)
				globalParameters.DELETE("/:parameter-id", global.GlobalParametersDelete)
			}

			definitionResponses := projects.Group("/definition/responses")
			{
				definitionResponses.POST("", definition.DefinitionResponsesCreate)
				definitionResponses.PUT("/:response-id", definition.DefinitionResponsesUpdate)
				definitionResponses.DELETE("/:response-id", definition.DefinitionResponsesDelete)
			}

			collections := projects.Group("/collections")
			{
				collections.POST("", collection.CollectionsCreate)
				collections.PUT("/:collection-id", check.CheckCollection(), collection.CollectionsUpdate)
				collections.POST("/:collection-id", check.CheckCollection(), collection.CollectionsCopy)
				collections.PUT("/movement", collection.CollectionsMovement)
				collections.DELETE("/:collection-id", check.CheckCollection(), collection.CollectionsDelete)
				collections.GET("/:collection-id/share", check.CheckCollection(), doc.DocShareDetails)
				collections.PUT("/:collection-id/share/switch", check.CheckCollection(), doc.DocShareSwitch)
				collections.PUT("/:collection-id/share/reset", check.CheckCollection(), doc.DocShareReset)
			}

			trashs := projects.Group("/trashs")
			{
				trashs.GET("", trash.TrashsList)
				trashs.PUT("", trash.TrashsRecover)
			}

			ais := projects.Group("/ai")
			{
				ais.GET("/collections/name", ai.AICreateApiNames)
				ais.POST("/collections", ai.AICreateCollection)
				ais.POST("/schemas", ai.AICreateSchema)
			}

			projectMember := projects.Group("/members")
			{
				projectMember.GET("", project.ProjectMembersList)
				projectMember.POST("", project.ProjectMembersCreate)
				projectMember.PUT("/authority/:user-id", project.ProjectMembersAuthUpdate)
				projectMember.DELETE("/:user-id", project.ProjectMembersDelete)
				projectMember.GET("/without", project.ProjectMembersWithout)
			}

			collectionHistories := projects.Group("/collections/:collection-id/histories")
			collectionHistories.Use(check.CheckCollection())
			{
				collectionHistories.GET("", doc.CollectionHistoryList)
				collectionHistories.GET("/:history-id", doc.CollectionHistoryDetails)
				collectionHistories.GET("/diff", doc.CollectionHistoryDiff)
				collectionHistories.PUT("/:history-id/restore", doc.CollectionHistoryRestore)

			}

			definitionSchemaHistories := projects.Group("/definition/schemas/:schemas-id/histories")
			definitionSchemaHistories.Use(check.CheckDefinitionSchema())
			{
				definitionSchemaHistories.GET("", definition.DefinitionSchemaHistoryList)
				definitionSchemaHistories.GET("/:history-id", definition.DefinitionSchemaHistoryDetails)
				definitionSchemaHistories.GET("/diff", definition.DefinitionSchemaHistoryDiff)
				definitionSchemaHistories.PUT("/:history-id/restore", definition.DefinitionSchemaHistoryRestore)
			}
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.FileFromFS("dist/", http.FS(frontend.FrontDist))
	})
}
