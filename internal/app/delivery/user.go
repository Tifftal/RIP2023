package delivery

import (
	"MSRM/internal/app/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteUserByID(repository *repository.Repository, c *gin.Context) {
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

	err = repository.DeleteUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Deleted!")
}

func EditUser(repository *repository.Repository, c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Id_user, idOk := jsonData["Id_user"].(float64)
	Name, nameOk := jsonData["Name"].(string)
	Status, statusOk := jsonData["User_status"].(string)

	fmt.Println(Id_user, Name, Status)
	if !idOk || Id_user <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing Id_user"})
		return
	}

	candidate, err := repository.GetUserByID(int(Id_user))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if nameOk {
		candidate.Name = Name
	}
	if statusOk {
		candidate.User_status = Status
	}
	err = repository.EditUser(candidate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User edited successfully",
	})
}

func GetUserByRole(repository *repository.Repository, c *gin.Context) {
	// Parse status from the request parameters
	role := c.Param("role")

	// Call the repository function to get missions by status
	users, err := repository.GetUserByRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve missions"})
		return
	}

	// Return the missions in the response
	c.JSON(http.StatusOK, users)
}
