package router

import (
	"io/fs"
	"net/http"

	"github.com/apicat/apicat/app/api"
	"github.com/apicat/apicat/app/middleware"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/frontend"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(r *gin.Engine) {

	r.Use(
		middleware.RequestIDLog("/assets/", "/static/"),
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

	mocksrv := api.NewMockServer()

	r.Any("/mock/:id/*path", mocksrv.Handler)

	apiRouter := r.Group("/api")
	apiRouter.Use(translator.UseValidatori18n())
	{
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
				project.POST("/:project-id/share/secretkey_check", api.ProjectShareSecretkeyCheck)

			}

			collection := notLogin.Group("/projects/:project-id/collections")
			collection.Use(middleware.CheckProject())
			{
				collection.GET("/:collection-id/data", api.CollectionDataGet)
				collection.POST("/:collection-id/share/secretkey_check", api.DocShareSecretkeyCheck)
			}

			share := notLogin.Group("/share")
			{
				share.GET("/collections/:public_collection_id/status", api.DocShareStatus)
			}
		}

		// 半登录状态下可访问的API。半登录：登录或不登录时都可访问，但响应的参数不同
		halfLogin := apiRouter.Group("")
		halfLogin.Use(middleware.CheckProject(), middleware.CheckMemberHalfLogin(), middleware.CheckProjectMemberHalfLogin(), mocksrv.ClearCache())
		{
			project := halfLogin.Group("/projects/:project-id")
			{
				project.GET("", api.ProjectsGet)
				project.GET("/status", api.ProjectStatus)
			}

			servers := halfLogin.Group("/projects/:project-id/servers")
			{
				servers.GET("", api.UrlList)
			}

			definitionSchemas := halfLogin.Group("/projects/:project-id/definition/schemas")
			{
				definitionSchemas.GET("", api.DefinitionSchemasList)
			}

			collections := halfLogin.Group("/projects/:project-id/collections")
			{
				collections.GET("", api.CollectionsList)
				collections.GET("/:collection-id", api.CollectionsGet)
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
				projects.PUT("/share", api.ProjectSharingSwitch)
				projects.PUT("/share/reset_share_secretkey", api.ProjectShareResetSecretKey)
			}

			definitionSchemas := project.Group("/definition/schemas")
			{
				definitionSchemas.POST("", api.DefinitionSchemasCreate)
				definitionSchemas.PUT("/:schemas-id", api.DefinitionSchemasUpdate)
				definitionSchemas.DELETE("/:schemas-id", api.DefinitionSchemasDelete)
				definitionSchemas.GET("/:schemas-id", api.DefinitionSchemasGet)
				definitionSchemas.POST("/:schemas-id", api.DefinitionSchemasCopy)
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
				collections.PUT("/:collection-id", api.CollectionsUpdate)
				collections.POST("/:collection-id", api.CollectionsCopy)
				collections.PUT("/movement", api.CollectionsMovement)
				collections.DELETE("/:collection-id", api.CollectionsDelete)
				collections.GET("/:collection-id/share", api.DocStatus)
				collections.PUT("/:collection-id/share", api.DocSharingSwitch)
				collections.PUT("/:collection-id/share/reset_share_secretkey", api.DocShareResetSecretKey)
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
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.FileFromFS("dist/", http.FS(frontend.FrontDist))
	})
}
