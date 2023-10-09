package app

import (
	"MSRM/internal/app/delivery"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/imgSample", "./resources/imgSample")
	router.Static("/js", "./resources/js")

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
			//http://localhost:8080/api/user/delete_user/1										DONE
			user.DELETE("/delete_user/:id", func(c *gin.Context) {
				delivery.DeleteUserByID(a.repository, c)
			})

			//http://localhost:8080/api/user/edit_info											DONE
			/*
				{
					"Id_user": 2,
					"Name": "TESTEST",
					"User_status": "Active"
				}
			*/
			user.PUT("edit_info", func(c *gin.Context) {
				delivery.EditUser(a.repository, c)
			})
		}

		mission := api.Group("/mission")
		{
			//http://localhost:8080/api/mission/update_mission									DONE
			/*
				{
					"Id_mission": 2,
					"Name": "TESTEST",
					"Mission_status": "testest"
				}
			*/
			mission.PUT("/update_mission", func(c *gin.Context) {
				delivery.UpdateMission(a.repository, c)
			})

			//http://localhost:8080/api/mission/delete_mission/1								DONE
			mission.DELETE("/delete_mission/:id", func(c *gin.Context) {
				delivery.DeleteMissionByID(a.repository, c)
			})

			//http://localhost:8080/api/mission/get_all_missions								DONE
			mission.GET("/get_all_missions", func(c *gin.Context) {
				delivery.GetAllMissiions(a.repository, c)
			})

			//http://localhost:8080/api/mission/get_mission/2 									DONE
			mission.GET("/get_mission/:id", func(c *gin.Context) {
				delivery.GetMissionByID(a.repository, c)
			})
		}

		sample := api.Group("/sample")
		{
			//http://localhost:8080/api/sample/create_sample   									DONE
			/*
				{
					"Name": "TEST1",
					"Type": "test",
					"Date_Sealed": "2021-08-06T00:00:00Z",
					"Sol_Sealed": 100,
					"Rock_Type": "n/a",
					"Height": "n/a",
					"Current_Location": "Sample Depot",
					"Image": "../../imgSample/no1.png",
					"Video": "https://mars.nasa.gov/embed/27520/",
					"Sample_status": "Active"
				}
			*/
			sample.POST("create_sample", func(c *gin.Context) {
				delivery.CreateSample(a.repository, c)
			})

			//http://localhost:8080/api/sample/delete_sample/8									DONE
			sample.DELETE("/delete_sample/:id", func(c *gin.Context) {
				delivery.DeleteSampleByID(a.repository, c)
			})

			//http://localhost:8080/api/sample/get_all_samples                      			DONE
			sample.GET("/get_all_samples", func(c *gin.Context) {
				delivery.GetAllSamples(a.repository, c)
			})

			//http://localhost:8080/api/sample/get_sample/5										DONE
			sample.GET("/get_sample/:id", func(c *gin.Context) {
				delivery.GetSampleByID(a.repository, c)
			})

			//http://localhost:8080/api/sample/update_sample									DONE
			/*
				{
					"Id_sample": 9,
					"Name": "TESTEST",
					"Type": "testest"
				}
			*/
			sample.PUT("/update_sample", func(c *gin.Context) {
				delivery.UpdateSample(a.repository, c)
			})
		}
	}

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"name": "Main website",
			"css":  "/styles/home.css",
		})
	})

	router.GET("/services", func(c *gin.Context) {
		sample, err := a.repository.GetAllSamples()
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}

		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"css":      "/styles/services.css",
			"Services": sample,
		})
	})

	router.GET("/services/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			return
		}
		samples, err := a.repository.GetAllSamples()
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}

		sample, err := a.repository.GetSampleByID(id)
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}

		if id <= 0 || id > len(samples) {
			log.Println("error")
			c.HTML(http.StatusOK, "info.tmpl", gin.H{
				"css":    "/styles/info.css",
				"Sample": nil,
				"Prev":   nil,
			})
			return
		}
		// sample.Date_Sealed = sample.Date_Sealed[:10]

		nextID := id + 1
		if nextID > len(samples) {
			nextID = 1
		}
		next := samples[nextID-1]

		prevID := id - 1
		if prevID < 1 {
			prevID = len(samples)
		}
		prev := samples[prevID-1]

		c.HTML(http.StatusOK, "info.tmpl", gin.H{
			"css":    "/styles/info.css",
			"Sample": sample,
			"Next":   next,
			"Prev":   prev,
		})

	})

	router.GET("/search", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("q", "")

		samples, err := a.repository.GetSampleByName(searchQuery)
		if err != nil {
			log.Println("Error with running\nServer down")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if searchQuery == "" {
			samples, err = a.repository.GetAllSamples()
			if err != nil {
				log.Println("Error with running\nServer down")
				return
			}
		}
		// for i := 0; i < len(samples); i++ {
		// 	samples[i].Date_Sealed = samples[i].Date_Sealed[:10]
		// }

		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"css":      "/styles/services.css",
			"Services": samples,
			"Search":   searchQuery,
		})

	})

	router.GET("/employee_mode", func(c *gin.Context) {
		sample, err := a.repository.GetAllSamples()
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}
		// c.JSON(200, sample)
		// for i := 0; i < len(sample); i++ {
		// 	sample[i].Date_Sealed = sample[i].Date_Sealed[:10]
		// }

		c.HTML(http.StatusOK, "employee_mode.tmpl", gin.H{
			"css":      "/styles/employee_mode.css",
			"Services": sample,
		})
	})

	router.GET("/search_empl", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("q", "")

		samples, err := a.repository.GetSampleByName(searchQuery)
		if err != nil {
			log.Println("Error with running\nServer down")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if searchQuery == "" {
			samples, err = a.repository.GetAllSamples()
			if err != nil {
				log.Println("Error with running\nServer down")
				return
			}
		}

		// for i := 0; i < len(samples); i++ {
		// 	samples[i].Date_Sealed = samples[i].Date_Sealed[:10]
		// }

		c.HTML(http.StatusOK, "employee_mode.tmpl", gin.H{
			"css":      "/styles/employee_mode.css",
			"Services": samples,
			"Search":   searchQuery,
		})

	})

	router.POST("/delete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("q", ""))
		log.Print(c.DefaultQuery("q", ""))
		if err != nil {
			log.Print(err)
		}
		err = a.repository.DeleteSampleByID(id)
		if err != nil {
			log.Print(err)
		}
		sample, err := a.repository.GetAllSamples()
		if err != nil {
			log.Print(err)
		}
		data := gin.H{
			"css":      "/styles/employee_mode.css",
			"Services": sample,
		}
		c.HTML(http.StatusOK, "employee_mode.tmpl", data)
	})

	router.GET("/test", func(c *gin.Context) {
		mission, err := a.repository.GetAllMissiions()
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}
		c.JSON(http.StatusOK, mission)
	})

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	}

	log.Println("Server down")
}
