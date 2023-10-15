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
	mission, err := repository.GetAllMissions()
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

func GetMissionDetailByID(repository *repository.Repository, c *gin.Context) {
	var mission *ds.Missions
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID (id < 0)",
		})
		return
	}

	mission, samples, err := repository.GetMissioninDetailByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return JSON response with mission details and associated samples
	c.JSON(http.StatusOK, gin.H{
		"mission": mission,
		"samples": samples,
	})
}

func GetMissionByUserID(repository *repository.Repository, c *gin.Context) {
	// Parse user ID from the request parameters
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the repository function to get missions for the user
	missions, err := repository.GetMissionByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve missions"})
		return
	}

	// Return the missions in the response
	c.JSON(http.StatusOK, missions)
}

func GetMissionByModeratorID(repository *repository.Repository, c *gin.Context) {
	// Parse user ID from the request parameters
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the repository function to get missions for the user
	missions, err := repository.GetMissionByModeratorID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve missions"})
		return
	}

	// Return the missions in the response
	c.JSON(http.StatusOK, missions)
}

func GetMissionByStatus(repository *repository.Repository, c *gin.Context) {
	// Parse status from the request parameters
	status := c.Param("status")

	// Call the repository function to get missions by status
	missions, err := repository.GetMissionByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve missions"})
		return
	}

	// Return the missions in the response
	c.JSON(http.StatusOK, missions)
}
