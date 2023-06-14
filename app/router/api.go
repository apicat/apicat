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
		account := apiRouter.Group("/account")
		{
			account.POST("/login/email", api.EmailLogin)
			account.POST("/register/email", api.EmailRegister)
		}

		projects := apiRouter.Group("/projects")
		{
			projects.GET("", middleware.JWTAuthMiddleware(), api.ProjectsList)
			projects.POST("", middleware.JWTAuthMiddleware(), api.ProjectsCreate)
			projects.GET("/:project-id", middleware.JWTAuthMiddleware(), api.ProjectsGet)
			projects.GET("/:project-id/data", api.ProjectDataGet)
		}

		members := apiRouter.Group("/members")
		members.Use(middleware.JWTAuthMiddleware())
		{
			members.GET("", api.GetMembers)
			members.POST("/", api.AddMember)
			members.PUT("/:user-id", api.SetMember)
			members.DELETE("/:user-id", api.DeleteMember)
		}

		user := apiRouter.Group("/user")
		user.Use(middleware.JWTAuthMiddleware())
		{
			user.GET("/self", api.GetUserInfo)
			user.PUT("/self", api.SetUserInfo)
			user.PUT("/self/password", api.ChangePassword)
		}

		project := apiRouter.Group("/projects/:project-id")
		project.Use(middleware.JWTAuthMiddleware(), middleware.CheckProject(), middleware.CheckProjectMember(), mocksrv.ClearCache())
		{
			projects := project.Group("")
			{
				projects.PUT("", api.ProjectsUpdate)
				projects.DELETE("", api.ProjectsDelete)
				projects.DELETE("/exit", api.ProjectExit)
				projects.PUT("/transfer", api.ProjectTransfer)
			}

			definitionSchemas := project.Group("/definition/schemas")
			{
				definitionSchemas.GET("", api.DefinitionSchemasList)
				definitionSchemas.POST("", api.DefinitionSchemasCreate)
				definitionSchemas.PUT("/:schemas-id", api.DefinitionSchemasUpdate)
				definitionSchemas.DELETE("/:schemas-id", api.DefinitionSchemasDelete)
				definitionSchemas.GET("/:schemas-id", api.DefinitionSchemasGet)
				definitionSchemas.POST("/:schemas-id", api.DefinitionSchemasCopy)
				definitionSchemas.PUT("/movement", api.DefinitionSchemasMove)
			}

			servers := project.Group("/servers")
			{
				servers.GET("", api.UrlList)
				servers.PUT("", api.UrlSettings)
			}

			globalParameters := project.Group("/global/parameters")
			{
				globalParameters.GET("", api.GlobalParametersList)
				globalParameters.POST("", api.GlobalParametersCreate)
				globalParameters.PUT("/:parameter-id", api.GlobalParametersUpdate)
				globalParameters.DELETE("/:parameter-id", api.GlobalParametersDelete)
			}

			definitionResponses := project.Group("/definition/responses")
			{
				definitionResponses.GET("", api.DefinitionResponsesList)
				definitionResponses.GET("/:response-id", api.DefinitionResponsesDetail)
				definitionResponses.POST("", api.DefinitionResponsesCreate)
				definitionResponses.PUT("/:response-id", api.DefinitionResponsesUpdate)
				definitionResponses.DELETE("/:response-id", api.DefinitionResponsesDelete)
			}

			collections := project.Group("/collections")
			{
				collections.GET("", api.CollectionsList)
				collections.GET("/:collection-id", api.CollectionsGet)
				collections.POST("", api.CollectionsCreate)
				collections.PUT("/:collection-id", api.CollectionsUpdate)
				collections.POST("/:collection-id", api.CollectionsCopy)
				collections.PUT("/movement", api.CollectionsMovement)
				collections.DELETE("/:collection-id", api.CollectionsDelete)
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
