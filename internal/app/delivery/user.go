package delivery

import (
	"MSRM/internal/app/ds"
	"MSRM/internal/app/pkg"
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

// @Summary Регистрация нового пользователя
// @Description Регистрация нового пользователя в системе
// @Accept json
// @Produce json
// @Tags Users
// @Param userJSON body ds.Users true "Данные нового пользователя"
// @Success 200 {string} string "Пользователь успешно зарегистрирован"
// @Failure 400 {object} string "Неверный запрос или ошибка регистрации пользователя"
// @Router /user/register [post]
func Register(repository *repository.Repository, c *gin.Context) {
	var userJSON ds.Users

	if err := c.ShouldBindJSON(&userJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userJSON.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пустой пароль"})
		return
	}

	if userJSON.Password != userJSON.RepeatPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароли не совпадают"})
		return
	}

	candidate, err := repository.GetUserByEmail(userJSON.Email_address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if candidate.Email_address == userJSON.Email_address {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким адресом уже существует"})
		return
	}

	if userJSON.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя пользователя пустое"})
		return
	}

	if err := repository.CreateUser(&userJSON); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Пользователь успешно зарегестрирован!")

}

// @Summary Вход пользователя
// @Description Аутентификация пользователя и генерация JWT-токена
// @Accept json
// @Produce json
// @Tags Users
// @Param userJSON body classes.Login true "Данные пользователя для входа"
// @Success 200 {string} string "JWT-токен успешно сгенерирован"
// @Failure 400 {object} string "Неверный запрос или ошибка аутентификации"
// @Router /user/login [post]
func Login(repository *repository.Repository, c *gin.Context) {
	var userJSON ds.Users

	if err := c.ShouldBindJSON(&userJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidate, err := repository.GetUserByEmail(userJSON.Email_address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if candidate.Password != userJSON.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный пароль"})
		return
	}

	token, err := pkg.GenerateToken(uint(candidate.Id_user), candidate.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = repository.SaveJWTToken(uint(candidate.Id_user), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, token)
}
