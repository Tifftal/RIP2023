package app

import (
	"MSRM/docs"
	"MSRM/internal/app/delivery"
	"MSRM/internal/app/pkg"
	"log"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// @SecurityDefinitions.apikey JwtAuth
// @in header
// @name Authorization
func (a *Application) StartServer() {
	router := gin.Default()

	docs.SwaggerInfo.Title = "Mars Sample Return Mission"
	docs.SwaggerInfo.Description = "API endpoint'ы для сервиса доставки пробирок с образцами марсианского грунта"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // List of allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Enable credentials (e.g., cookies)
	}))

	user := router.Group("/user")
	{
		user.POST("/register", func(c *gin.Context) {
			delivery.Register(a.repository, c)
		})

		user.POST("/login", func(c *gin.Context) {
			delivery.Login(a.repository, c)
		})

		user.POST("/logout", a.RoleMiddleware(pkg.Moderator, pkg.User), func(c *gin.Context) {
			delivery.Logout(a.repository, c)
		})
	}

	api := router.Group("/api")
	{

		// user := api.Group("/user")
		// {
		// 	user.DELETE("/delete_user/:id", func(c *gin.Context) {
		// 		delivery.DeleteUserByID(a.repository, c)
		// 	})

		// 	user.PUT("edit_info", func(c *gin.Context) {
		// 		delivery.EditUser(a.repository, c)
		// 	})

		// 	user.GET("get_user_by_role/:role", func(c *gin.Context) {
		// 		delivery.GetUserByRole(a.repository, c)
		// 	})
		// }

		mission := api.Group("/mission")
		{
			mission.GET("/", a.RoleMiddleware(pkg.Moderator, pkg.User), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.GetAllMissiions(a.repository, c, user_id)
			})

			mission.GET("/:id", a.RoleMiddleware(pkg.Moderator, pkg.User), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.GetMissionDetailByID(a.repository, c, user_id)
			})

			mission.PUT("/update/:id", a.RoleMiddleware(pkg.Moderator), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.UpdateMission(a.repository, c, user_id)
			})

			// mission.DELETE("/delete/:id", a.RoleMiddleware(pkg.User), func(c *gin.Context) {
			// 	user_id := c.MustGet("User_id").(int)
			// 	delivery.DeleteMissionByID(a.repository, c, user_id)
			// })

			mission.DELETE("/delete_from_last/:id", a.RoleMiddleware(pkg.User), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.RemoveSampleFromLastDraftMission(a.repository, c, user_id)
			})

			mission.PUT("/status_by_user/:id", a.RoleMiddleware(pkg.User), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.UpdateMissionStatusByUser(a.repository, c, user_id)
			})

			mission.PUT("/status_by_moderator/:id", a.RoleMiddleware(pkg.Moderator), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.UpdateMissionStatusByModerator(a.repository, c, user_id)
			})
		}

		sample := api.Group("/sample")
		{
			sample.GET("/:id", func(c *gin.Context) {
				delivery.GetSampleByID(a.repository, c)
			})

			sample.POST("/create", a.RoleMiddleware(pkg.Moderator), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.CreateSample(a.repository, c, user_id)
			})

			sample.DELETE("/delete/:id", a.RoleMiddleware(pkg.Moderator), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.DeleteSampleByID(a.repository, c, user_id)
			})

			sample.GET("/", a.Guest(pkg.User), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.GetAllSamples(a.repository, c, user_id)
			})

			sample.PUT("/update/:id", a.RoleMiddleware(pkg.Moderator), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.UpdateSample(a.repository, c, user_id)
			})

			sample.PUT("/to_mission/:id", a.RoleMiddleware(pkg.User), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.AddSampleToMission(a.repository, c, user_id)
			})

			sample.POST("/:id/image", a.RoleMiddleware(pkg.Moderator), func(c *gin.Context) {
				user_id := c.MustGet("User_id").(int)
				delivery.AddImageToSample(a.repository, c, user_id)
			})
		}
	}

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	}

	log.Println("Server down")
}
