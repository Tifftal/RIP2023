package app

import (
	"MSRM/internal/app/delivery"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func (a *Application) StartServer() {
	router := gin.Default()

	user := router.Group("/user")
	{
		user.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "ping",
			})
		})

		user.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "ping",
			})
		})
	}

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.DELETE("/delete_user/:id", func(c *gin.Context) {
				delivery.DeleteUserByID(a.repository, c)
			})

			user.PUT("edit_info", func(c *gin.Context) {
				delivery.EditUser(a.repository, c)
			})

			user.GET("get_user_by_role/:role", func(c *gin.Context) {
				delivery.GetUserByRole(a.repository, c)
			})
		}

		mission := api.Group("/mission")
		{
			//DONE
			mission.GET("/get_all_missions", func(c *gin.Context) {
				delivery.GetAllMissiions(a.repository, c)
			})

			//DONE
			mission.GET("/get_mission/:id", func(c *gin.Context) {
				delivery.GetMissionDetailByID(a.repository, c)
			})

			//DONE
			mission.PUT("/update_mission/:id", func(c *gin.Context) {
				delivery.UpdateMission(a.repository, c)
			})

			//DONE
			mission.DELETE("/delete_mission/:id", func(c *gin.Context) {
				delivery.DeleteMissionByID(a.repository, c)
			})

			//DONE
			mission.DELETE("/delete_sample_from_last_mission/:id", func(c *gin.Context) {
				delivery.RemoveSampleFromLastDraftMission(a.repository, c)
			})

			//DONE
			mission.DELETE("/delete_sample_from_mission/:mission_id/:sample_id", func(c *gin.Context) {
				delivery.RemoveSampleFromMission(a.repository, c)
			})

			//DONE
			mission.PUT("/update_mission_status_by_user/:id", func(c *gin.Context) {
				delivery.UpdateMissionStatusByUser(a.repository, c)
			})

			//DONE
			mission.PUT("/update_mission_status_by_moderator/:id", func(c *gin.Context) {
				delivery.UpdateMissionStatusByModerator(a.repository, c)
			})
		}

		sample := api.Group("/sample")
		{
			//DONE
			sample.GET("/get_sample/:id", func(c *gin.Context) {
				delivery.GetSampleByID(a.repository, c)
			})

			//DONE
			sample.POST("create_sample", func(c *gin.Context) {
				delivery.CreateSample(a.repository, c)
			})

			//DONE
			sample.DELETE("/delete_sample/:id", func(c *gin.Context) {
				delivery.DeleteSampleByID(a.repository, c)
			})

			//DONE
			sample.GET("/get_all_samples", func(c *gin.Context) {
				delivery.GetAllSamples(a.repository, c)
			})

			//DONE
			sample.PUT("/update_sample/:id", func(c *gin.Context) {
				delivery.UpdateSample(a.repository, c)
			})

			//DONE
			sample.PUT("/add_sample_to_mission/:id", func(c *gin.Context) {
				delivery.AddSampleToMission(a.repository, c)
			})

			// sample.GET("get_all_samples_order_type", func(c *gin.Context) {
			// 	delivery.GetAllSamplesOrderByType(a.repository, c)
			// })

			// sample.GET("get_all_samples_order_date", func(c *gin.Context) {
			// 	delivery.GetAllSamplesOrderByDate(a.repository, c)
			// })

			// sample.GET("get_all_samples_active", func(c *gin.Context) {
			// 	delivery.GetAllSamplesStatusActive(a.repository, c)
			// })

			// sample.GET("get_all_samples_deleted", func(c *gin.Context) {
			// 	delivery.GetAllSamplesStatusDaleted(a.repository, c)
			// })
		}
	}

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	}

	log.Println("Server down")
}
