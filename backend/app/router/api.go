package router

import (
	"github.com/apicat/apicat/backend/app/api"
	"github.com/apicat/apicat/backend/app/middleware"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/frontend"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(r *gin.Engine) {

	r.Use(
		middleware.RequestIDLog("/assets/", "/static/"),
		translator.UseValidatori18n(),
		middleware.CheckDBConnStatus("/assets/", "/static/"),
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

	t, _ := template.ParseFS(frontend.FrontDist, "dist/templates/*.tmpl")
	r.SetHTMLTemplate(t)

	r.GET("/", func(ctx *gin.Context) {
		ctx.FileFromFS("dist/", http.FS(frontend.FrontDist))
	})

	mocksrv := api.NewMockServer()

	r.Any("/mock/:id/*path", mocksrv.Handler)

	config := r.Group("/config")
	{
		config.GET("/db", api.GetDBConfig)
	}

	apiRouter := r.Group("/api")
	{
		config := apiRouter.Group("/config")
		{
			config.PUT("/db", api.SetDBConfig)
		}

		// 未登录状态下可访问的API
		notLogin := apiRouter.Group("")
		{
			account := notLogin.Group("/account")
			{
				account.POST("/login/email", api.EmailLogin)
				account.POST("/register/email", api.EmailRegister)
			}

			project := notLogin.Group("/projects")
			project.Use(middleware.CheckProject())
			{
				project.GET("/:project-id/data", api.ProjectDataGet)
				project.GET("/:project-id/share/status", middleware.CheckMemberHalfLogin(), api.ProjectShareStatus)
				project.POST("/:project-id/share/check", api.ProjectShareSecretkeyCheck)
			}

			collection := notLogin.Group("/projects/:project-id/collections")
			collection.Use(middleware.CheckProject())
			{
				collection.GET("/:collection-id/data", middleware.CheckCollection(), api.CollectionDataGet)
				collection.POST("/:collection-id/share/check", middleware.CheckCollection(), api.DocShareCheck)
			}

			collection_share := notLogin.Group("/collections")
			{
				collection_share.GET("/:public_collection_id/share/status", api.DocShareStatus)
			}
		}

		// 半登录状态下可访问的API。半登录：登录或不登录时都可访问，但响应的参数不同
		halfLogin := apiRouter.Group("")
		halfLogin.Use(middleware.CheckProject(), middleware.CheckMemberHalfLogin(), middleware.CheckProjectMemberHalfLogin(), mocksrv.ClearCache())
		{
			project := halfLogin.Group("/projects/:project-id")
			{
				project.GET("", api.ProjectsGet)
			}

			servers := halfLogin.Group("/projects/:project-id/servers")
			{
				servers.GET("", api.UrlList)
			}

			definitionSchemas := halfLogin.Group("/projects/:project-id/definition/schemas")
			{
				definitionSchemas.GET("/:schemas-id", middleware.CheckDefinitionSchema(), api.DefinitionSchemasGet)
				definitionSchemas.GET("", api.DefinitionSchemasList)
			}

			collections := halfLogin.Group("/projects/:project-id/collections")
			{
				collections.GET("", api.CollectionsList)
				collections.GET("/:collection-id", middleware.CheckCollection(), api.CollectionsGet)
			}

			globalParameters := halfLogin.Group("/projects/:project-id/global/parameters")
			{
				globalParameters.GET("", api.GlobalParametersList)
			}

			definitionResponses := halfLogin.Group("/projects/:project-id/definition/responses")
			{
				definitionResponses.GET("", api.DefinitionResponsesList)
				definitionResponses.GET("/:response-id", api.DefinitionResponsesDetail)
			}
		}

		// 仅登录状态下可访问的API。仅登录：仅登录了apicat便可访问，一般为项目外的操作
		onlyLogin := apiRouter.Group("")
		onlyLogin.Use(middleware.CheckMember())
		{
			user := onlyLogin.Group("/user")
			{
				user.GET("/self", api.GetUserInfo)
				user.PUT("/self", api.SetUserInfo)
				user.PUT("/self/password", api.ChangePassword)
			}

			members := onlyLogin.Group("/members")
			{
				members.GET("", api.GetMembers)
				members.POST("/", api.AddMember)
				members.PUT("/:user-id", api.SetMember)
				members.DELETE("/:user-id", api.DeleteMember)
			}

			project := onlyLogin.Group("/projects")
			{
				project.GET("", api.ProjectsList)
				project.POST("", api.ProjectsCreate)
			}

			iteration := onlyLogin.Group("/iterations")
			{
				iteration.GET("", api.IterationsList)
				iteration.GET("/:iteration-id", api.IterationsDetails)
				iteration.POST("", api.IterationsCreate)
				iteration.PUT("/:iteration-id", api.IterationsUpdate)
				iteration.DELETE("/:iteration-id", api.IterationsDelete)
			}

			projectGroup := onlyLogin.Group("/project_group")
			{
				projectGroup.GET("", api.ProjectGroupList)
				projectGroup.POST("", api.ProjectGroupCreate)
				projectGroup.PUT("/:group_id/rename", api.ProjectGroupRename)
				projectGroup.DELETE("/:group_id", api.ProjectGroupDelete)
				projectGroup.PUT("/order", api.ProjectGroupOrder)
			}
		}

		// 项目内部操作
		project := apiRouter.Group("/projects/:project-id")
		project.Use(middleware.JWTAuthMiddleware(), middleware.CheckProject(), middleware.CheckProjectMember(), mocksrv.ClearCache())
		{
			projects := project.Group("")
			{
				projects.PUT("", api.ProjectsUpdate)
				projects.DELETE("", api.ProjectsDelete)
				projects.DELETE("/exit", api.ProjectExit)
				projects.PUT("/transfer", api.ProjectTransfer)
				projects.GET("/share", api.ProjectShareDetails)
				projects.PUT("/share/switch", api.ProjectSharingSwitch)
				projects.PUT("/share/reset", api.ProjectShareReset)
				projects.POST("/follow", api.ProjectFollow)
				projects.DELETE("/follow", api.ProjectUnFollow)
				projects.PUT("/change_group", api.ProjectChangeGroup)
			}

			definitionSchemas := project.Group("/definition/schemas")
			{
				definitionSchemas.POST("", api.DefinitionSchemasCreate)
				definitionSchemas.PUT("/:schemas-id", middleware.CheckDefinitionSchema(), api.DefinitionSchemasUpdate)
				definitionSchemas.DELETE("/:schemas-id", middleware.CheckDefinitionSchema(), api.DefinitionSchemasDelete)
				definitionSchemas.POST("/:schemas-id", middleware.CheckDefinitionSchema(), api.DefinitionSchemasCopy)
				definitionSchemas.PUT("/movement", api.DefinitionSchemasMove)
			}

			servers := project.Group("/servers")
			{
				servers.PUT("", api.UrlSettings)
			}

			globalParameters := project.Group("/global/parameters")
			{
				globalParameters.POST("", api.GlobalParametersCreate)
				globalParameters.PUT("/:parameter-id", api.GlobalParametersUpdate)
				globalParameters.DELETE("/:parameter-id", api.GlobalParametersDelete)
			}

			definitionResponses := project.Group("/definition/responses")
			{
				definitionResponses.POST("", api.DefinitionResponsesCreate)
				definitionResponses.PUT("/:response-id", api.DefinitionResponsesUpdate)
				definitionResponses.DELETE("/:response-id", api.DefinitionResponsesDelete)
			}

			collections := project.Group("/collections")
			{
				collections.POST("", api.CollectionsCreate)
				collections.PUT("/:collection-id", middleware.CheckCollection(), api.CollectionsUpdate)
				collections.POST("/:collection-id", middleware.CheckCollection(), api.CollectionsCopy)
				collections.PUT("/movement", api.CollectionsMovement)
				collections.DELETE("/:collection-id", middleware.CheckCollection(), api.CollectionsDelete)
				collections.GET("/:collection-id/share", middleware.CheckCollection(), api.DocShareDetails)
				collections.PUT("/:collection-id/share/switch", middleware.CheckCollection(), api.DocShareSwitch)
				collections.PUT("/:collection-id/share/reset", middleware.CheckCollection(), api.DocShareReset)
			}

			trashs := project.Group("/trashs")
			{
				trashs.GET("", api.TrashsList)
				trashs.PUT("", api.TrashsRecover)
			}

			ai := project.Group("/ai")
			{
				ai.GET("/collections/name", api.AICreateApiNames)
				ai.POST("/collections", api.AICreateCollection)
				ai.POST("/schemas", api.AICreateSchema)
			}

			projectMember := project.Group("/members")
			{
				projectMember.GET("", api.ProjectMembersList)
				projectMember.POST("", api.ProjectMembersCreate)
				projectMember.PUT("/authority/:user-id", api.ProjectMembersAuthUpdate)
				projectMember.DELETE("/:user-id", api.ProjectMembersDelete)
				projectMember.GET("/without", api.ProjectMembersWithout)
			}

			collectionHistories := project.Group("/collections/:collection-id/histories")
			collectionHistories.Use(middleware.CheckCollection())
			{
				collectionHistories.GET("", api.CollectionHistoryList)
				collectionHistories.GET("/:history-id", api.CollectionHistoryDetails)
				collectionHistories.GET("/diff", api.CollectionHistoryDiff)
				collectionHistories.PUT("/:history-id/restore", api.CollectionHistoryRestore)

			}

			definitionSchemaHistories := project.Group("/definition/schemas/:schemas-id/histories")
			definitionSchemaHistories.Use(middleware.CheckDefinitionSchema())
			{
				definitionSchemaHistories.GET("", api.DefinitionSchemaHistoryList)
				definitionSchemaHistories.GET("/:history-id", api.DefinitionSchemaHistoryDetails)
				definitionSchemaHistories.GET("/diff", api.DefinitionSchemaHistoryDiff)
				definitionSchemaHistories.PUT("/:history-id/restore", api.DefinitionSchemaHistoryRestore)
			}
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.FileFromFS("dist/", http.FS(frontend.FrontDist))
	})
}
