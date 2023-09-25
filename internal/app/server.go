package app

import "C"
import (
	"log"
	"net/http"

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
		for i := 0; i < len(sample); i++ {
			date := sample[i].Date_Sealed
			date.Format("January 02, 2006")
			sample[i].Date_Sealed = date
		}
		// c.JSON(200, sample)
		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"css":      "/styles/services.css",
			"Services": sample,
		})
	})

	// router.GET("/services/:id", func(c *gin.Context) {

	// 	sample, err := a.repository.GetSampleByID(id)
	// 	if err != nil {
	// 		log.Println("Error with running\nServer down")
	// 		return
	// 	}
	// 	c.JSON(200, sample)
	// 	// c.HTML(http.StatusOK, "services.tmpl", gin.H{
	// 	// 	"css":      "/styles/services.css",
	// 	// 	"Services": sample,
	// 	// })
	// })

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
