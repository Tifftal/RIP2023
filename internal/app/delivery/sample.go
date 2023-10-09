package delivery

import (
	"MSRM/internal/app/ds"
	"MSRM/internal/app/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllSamples(repository *repository.Repository, c *gin.Context) {
	var sample []ds.Samples
	sample, err := repository.GetAllSamples()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, sample)
}

func GetSampleByID(repository *repository.Repository, c *gin.Context) {
	var sample *ds.Samples
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(id)

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Oshibochka id<0",
		})
		return
	}
	sample, err = repository.GetSampleByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, sample)
}

func CreateSample(repository *repository.Repository, c *gin.Context) {
	var sample *ds.Samples
	// Достаем данные из JSON'а из запроса
	if err := c.BindJSON(&sample); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(sample)

	if sample.Name == "" || sample.Type == "" || sample.Sol_Sealed == 0 || sample.Current_Location == "" || sample.Sample_status == "" {
		c.JSON(http.StatusBadRequest, "Oshibochka")
	}

	if err := repository.CreateSample(sample); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Created",
	})
}

func DeleteSampleByID(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(id)

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Oshibochka id<0",
		})
		return
	}

	err = repository.DeleteSampleByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Deleted!")
}

func UpdateSample(repository *repository.Repository, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Id_sample, idOk := jsonData["Id_sample"].(float64)
	Name, nameOk := jsonData["Name"].(string)
	Type, typeOk := jsonData["Type"].(string)

	fmt.Println(Id_sample, Name, Type)
	if !idOk || Id_sample <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing Id_sample"})
		return
	}

	candidate, err := repository.GetSampleByID(int(Id_sample))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if nameOk {
		candidate.Name = Name
	}
	if typeOk {
		candidate.Type = Type
	}
	err = repository.UpdateSample(candidate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sample updated successfully",
	})
}
