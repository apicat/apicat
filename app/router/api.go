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

	apiRouter := r.Group("/api").Use(translator.UseValidatori18n())
	{
		projects := apiRouter.(*gin.RouterGroup).Group("/projects", middleware.JWTAuthMiddleware())
		{
			projects.GET("", api.ProjectsList)
			projects.GET("/:id", api.ProjectsGet)
			projects.GET("/:id/data", api.ProjectDataGet)
			projects.POST("/", api.ProjectsCreate)
			projects.PUT("/:id", api.ProjectsUpdate)
			projects.DELETE("/:id", api.ProjectsDelete)
		}

		account := apiRouter.(*gin.RouterGroup).Group("/account")
		{
			account.POST("/login/email", api.EmailLogin)
			account.POST("/register/email", api.EmailRegister)
		}

		user := apiRouter.(*gin.RouterGroup).Group("/user", middleware.JWTAuthMiddleware())
		{
			user.GET("/", api.GetUserInfo)
			user.PUT("/", api.SetUserInfo)
			user.PUT("/password", api.ChangePassword)
		}

		project := apiRouter.(*gin.RouterGroup).Group("/projects/:id", middleware.JWTAuthMiddleware()).Use(middleware.CheckProject())
		{
			definitions := project.(*gin.RouterGroup).Group("/definitions")
			{
				definitions.GET("", api.DefinitionSchemasList)
				definitions.POST("", api.DefinitionSchemasCreate)
				definitions.PUT("/:definition-id", api.DefinitionSchemasUpdate)
				definitions.DELETE("/:definition-id", api.DefinitionSchemasDelete)
				definitions.GET("/:definition-id", api.DefinitionSchemasGet)
				definitions.POST("/:definition-id", api.DefinitionSchemasCopy)
				definitions.PUT("/movement", api.DefinitionSchemasMove)
			}

			servers := project.(*gin.RouterGroup).Group("/servers")
			{
				servers.GET("", api.UrlList)
				servers.PUT("", api.UrlSettings)
			}

			globalParameters := project.(*gin.RouterGroup).Group("/global/parameters")
			{
				globalParameters.GET("", api.GlobalParametersList)
				globalParameters.POST("", api.GlobalParametersCreate)
				globalParameters.PUT("/:parameter-id", api.GlobalParametersUpdate)
				globalParameters.DELETE("/:parameter-id", api.GlobalParametersDelete)
			}

			definitionsResponses := project.(*gin.RouterGroup).Group("/definitions/responses")
			{
				definitionsResponses.GET("", api.CommonResponsesList)
				definitionsResponses.GET("/:response-id", api.CommonResponsesDetail)
				definitionsResponses.POST("", api.CommonResponsesCreate)
				definitionsResponses.PUT("/:response-id", api.CommonResponsesUpdate)
				definitionsResponses.DELETE("/:response-id", api.CommonResponsesDelete)
			}

			collections := project.(*gin.RouterGroup).Group("/collections")
			{
				collections.GET("", api.CollectionsList)
				collections.GET("/:collection-id", api.CollectionsGet)
				collections.POST("", api.CollectionsCreate)
				collections.PUT("/:collection-id", api.CollectionsUpdate)
				collections.POST("/:collection-id", api.CollectionsCopy)
				collections.PUT("/movement", api.CollectionsMovement)
				collections.DELETE("/:collection-id", api.CollectionsDelete)
			}

			trashs := project.(*gin.RouterGroup).Group("/trashs")
			{
				trashs.GET("", api.TrashsList)
				trashs.PUT("", api.TrashsRecover)
			}

			ai := project.(*gin.RouterGroup).Group("/ai")
			{
				ai.GET("/collections/name", api.AICreateApiNames)
				ai.POST("/collections", api.AICreateCollection)
				ai.POST("/schemas", api.AICreateSchema)
			}
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.FileFromFS("dist/", http.FS(frontend.FrontDist))
	})
}
