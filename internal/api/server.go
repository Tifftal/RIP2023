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
			Type:            "Атмосферный",
			DateSealed:      "6 августа 2021",
			SolSealed:       "164",
			RockType:        "n/a",
			SampleHeight:    "n/a",
			CurrentLocation: "Хранилище образцов",
			Url:             "../../imgSample/no1.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27520/",
		},
		{
			Id:              2,
			Name:            "Montdenier",
			Type:            "Горный образец",
			DateSealed:      "6 сентября 2021",
			SolSealed:       "194",
			RockType:        "Магматический",
			SampleHeight:    "5.98 см/2.35 дюйма",
			CurrentLocation: "Хранилище образцов",
			Url:             "../../imgSample/no2.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27530/",
		},
		{
			Id:              3,
			Name:            "Montagnac",
			Type:            "Горный образец",
			DateSealed:      "8 сентября 2021",
			SolSealed:       "196",
			RockType:        "Магматический",
			SampleHeight:    "6.14 см/2.42 дюйма",
			CurrentLocation: "Ровер Perseverance",
			Url:             "../../imgSample/no3.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27530/",
		},
		{
			Id:              4,
			Name:            "Salette",
			Type:            "Горный образец",
			DateSealed:      "15 ноября 2021",
			SolSealed:       "262",
			RockType:        "Магматический",
			SampleHeight:    "6.28 см/2.47 дюйма",
			CurrentLocation: "Ровер Perseverance",
			Url:             "../../imgSample/no4.jpg",
			UrlVideo:        "https://mars.nasa.gov/embed/27542/",
		},
		{
			Id:              5,
			Name:            "Coulettes",
			Type:            "Горный образец",
			DateSealed:      "24 ноября 2021",
			SolSealed:       "271",
			RockType:        "Магматический",
			SampleHeight:    "3.3 см/1.30 дюйма",
			CurrentLocation: "Хранилище образцов",
			Url:             "../../imgSample/no5.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27542/",
		},
		{
			Id:              6,
			Name:            "Robine",
			Type:            "Горный образец",
			DateSealed:      "22 декабря 2021",
			SolSealed:       "298",
			RockType:        "Магматический",
			SampleHeight:    "6.08 см/2.39 дюйма",
			CurrentLocation: "Ровер Perseverance",
			Url:             "../../imgSample/no6.png",
			UrlVideo:        "https://mars.nasa.gov/embed/27556/",
		},
		{
			Id:              7,
			Name:            "Malay",
			Type:            "Горный образец",
			DateSealed:      "31 января 2022",
			SolSealed:       "337",
			RockType:        "Магматический",
			SampleHeight:    "3.07 см/1.21 дюйма",
			CurrentLocation: "Хранилище образцов",
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
		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"css":      "/styles/services.css",
			"Services": foundSample,
			"Search":   searchQuery,
		})

	})

	r.Run()

	log.Println("Server down")
}
