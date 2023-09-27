package app

import "C"
import (
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = a.repository.DeleteSampleByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sample, err := a.repository.GetAllSamples()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data := gin.H{
			"css":      "/styles/employee_mode.css",
			"Services": sample,
		}
		c.HTML(http.StatusOK, "employee_mode.tmpl", data)
	})

	router.POST("/return", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("q", ""))
		log.Print(c.DefaultQuery("q", ""))
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = a.repository.ReturnSampleByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sample, err := a.repository.GetAllSamples()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data := gin.H{
			"css":      "/styles/employee_mode.css",
			"Services": sample,
		}
		c.HTML(http.StatusOK, "employee_mode.tmpl", data)
	})

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	}

	log.Println("Server down")
}
