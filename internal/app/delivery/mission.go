package delivery

import (
	"MSRM/internal/app/ds"
	"MSRM/internal/app/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllMissiions(repository *repository.Repository, c *gin.Context) {
	var mission []ds.Missions
	mission, err := repository.GetAllMissiions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mission)
}

func GetMissionByID(repository *repository.Repository, c *gin.Context) {
	var mission *ds.Missions
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
	mission, err = repository.GetMissionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mission)
}

func DeleteMissionByID(repository *repository.Repository, c *gin.Context) {
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

	err = repository.DeleteMissionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Deleted!")
}

func UpdateMission(repository *repository.Repository, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Id_mission, idOk := jsonData["Id_mission"].(float64)
	Name, nameOk := jsonData["Name"].(string)
	Status, statusOk := jsonData["Mission_status"].(string)

	fmt.Println(Id_mission, Name, Status)
	if !idOk || Id_mission <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing Id_mission"})
		return
	}

	candidate, err := repository.GetMissionByID(int(Id_mission))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if nameOk {
		candidate.Name = Name
	}
	if statusOk {
		candidate.Mission_status = Status
	}
	err = repository.UpdateMission(candidate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mission updated successfully",
	})
}
