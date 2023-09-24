package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Sample struct {
	Id              int
	Name            string
	Type            string
	SolSealed       string
	DateSealed      string
	RockType        string
	SampleHeight    string
	CurrentLocation string
	Url             string
	UrlVideo        string
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	samples := []Sample{
		{
			Id:              1,
			Name:            "Roubion",
			Type:            "Atmospheric",
			DateSealed:      "Aug. 6, 2021",
			SolSealed:       "164",
			RockType:        "n/a",
			SampleHeight:    "n/a",
			CurrentLocation: "Sample Depot",
			Url:             "../../imgSample/no1.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27520/",
		},
		{
			Id:              2,
			Name:            "Montdenier",
			Type:            "Rock Core",
			DateSealed:      "Sept. 6, 2021",
			SolSealed:       "194",
			RockType:        "Igneous",
			SampleHeight:    "5.98 cm/2.35 in",
			CurrentLocation: "Sample Depot",
			Url:             "../../imgSample/no2.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27530/",
		},
		{
			Id:              3,
			Name:            "Montagnac",
			Type:            "Rock Core",
			DateSealed:      "Sept. 8, 2021",
			SolSealed:       "196",
			RockType:        "Igneous",
			SampleHeight:    "6.14 cm/2.42 in",
			CurrentLocation: "Perseverance Rover",
			Url:             "../../imgSample/no3.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27530/",
		},
		{
			Id:              4,
			Name:            "Salette",
			Type:            "Rock Core",
			DateSealed:      "Nov. 15, 2021",
			SolSealed:       "262",
			RockType:        "Igneous",
			SampleHeight:    "6.28 cm/2.47 in",
			CurrentLocation: "Perseverance Rover",
			Url:             "../../imgSample/no4.jpg",
			UrlVideo:        "https://mars.nasa.gov/embed/27542/",
		},
		{
			Id:              5,
			Name:            "Coulettes",
			Type:            "Rock Core",
			DateSealed:      "Nov. 24, 2021",
			SolSealed:       "271",
			RockType:        "Igneous",
			SampleHeight:    "3.3 cm/1.30 in",
			CurrentLocation: "Sample Depot",
			Url:             "../../imgSample/no5.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27542/",
		},
		{
			Id:              6,
			Name:            "Robine",
			Type:            "Rock Core",
			DateSealed:      "Dec. 22, 2021",
			SolSealed:       "298",
			RockType:        "Igneous",
			SampleHeight:    "6.08 cm/2.39 in",
			CurrentLocation: "Perseverance Rover",
			Url:             "../../imgSample/no6.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27556/",
		},
		{
			Id:              7,
			Name:            "Malay",
			Type:            "Rock Core",
			DateSealed:      "Jan. 31, 2022",
			SolSealed:       "337",
			RockType:        "Igneous",
			SampleHeight:    "3.07 cm/1.21 in",
			CurrentLocation: "Sample Depot",
			Url:             "../../imgSample/no7.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27556/",
		},
	}
	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"name": "Main website",
			"css":  "/styles/home.css",
		})
	})

	r.Static("/styles", "./resources/styles")
	r.Static("/imgSample", "./resources/imgSample")

	r.GET("/services", func(c *gin.Context) {
		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"css":      "/styles/services.css",
			"Services": samples,
		})
	})

	// r.GET("/services/:id", func(c *gin.Context) {
	// 	id, err := strconv.Atoi(c.Param("id"))
	// 	if err != nil || id < 1 || id > len(samples) {
	// 		log.Println("error")
	// 		c.HTML(http.StatusOK, "info.tmpl", gin.H{
	// 			"css":    "/styles/info.css",
	// 			"Sample": nil,
	// 		})
	// 		return
	// 	}

	// 	sample := samples[id-1]
	// 	var next Sample
	// 	if id < len(samples) {
	// 		next = samples[id]
	// 	}

	// 	c.HTML(http.StatusOK, "info.tmpl", gin.H{
	// 		"css":    "/styles/info.css",
	// 		"Sample": sample,
	// 		"Next":   next,
	// 	})
	// })
	r.GET("/services/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			return
		}

		if id <= 0 || id > len(samples) {
			log.Println("error")
			c.HTML(http.StatusOK, "info.tmpl", gin.H{
				"css":    "/styles/info.css",
				"Sample": nil,
				"Prev":   nil, // Добавили новую переменную "Prev"
			})
			return
		}

		sample := samples[id-1]

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
			"Prev":   prev, // Добавили предыдущий объект "Prev"
		})
	})

	r.GET("/search", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("q", "")
		var foundSample []Sample
		for _, sample := range samples {
			if strings.HasPrefix(strings.ToLower(sample.Name), strings.ToLower(searchQuery)) {
				foundSample = append(foundSample, sample)
			}
		}
		data := gin.H{
			"css":    "/styles/test.css",
			"Sample": foundSample,
		}
		c.HTML(http.StatusOK, "test.tmpl", data)
	})

	r.Run()

	log.Println("Server down")
}
