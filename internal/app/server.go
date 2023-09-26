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

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"name": "Main website",
			"css":  "/styles/home.css",
		})
	})

	router.GET("/", func(context *gin.Context) {
		sample, err := a.repository.GetSampleByID(2)
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}

		context.JSON(200, sample)
	})

	router.GET("/services", func(c *gin.Context) {
		sample, err := a.repository.GetAllSamples()
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}

		// c.JSON(200, sample)
		for i := 0; i < len(sample); i++ {
			sample[i].Date_Sealed = sample[i].Date_Sealed[:10]
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
		sample.Date_Sealed = sample.Date_Sealed[:10]

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

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	}

	log.Println("Server down")
}
