package delivery

import (
	"MSRM/internal/app/ds"
	"MSRM/internal/app/repository"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка миссий
// @Description Возвращает список всех миссий или миссий в указанном диапазоне дат
// @Accept json
// @Produce json
// @Tags Missions
// @Param start_date query string false "Дата начала в формате YYYY-MM-DD (необязательно)"
// @Param end_date query string false "Дата окончания в формате YYYY-MM-DD (необязательно)"
// @Success 200 {object} string "Список миссий"
// @Failure 400 {object} string "Неверный запрос или ошибка в формате дат"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/mission/ [get]
// @Security JwtAuth
func GetAllMissions(repository *repository.Repository, c *gin.Context, user_id int) {
	startDateString := c.Query("start_date")
	endDateString := c.Query("end_date")
	missionStatus := c.Query("mission_status")

	var startDate, endDate time.Time
	var err error

	if startDateString != "" {
		if startDate, err = time.Parse("2006-01-02", startDateString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала"})
			return
		}
	}

	if endDateString != "" {
		if endDate, err = time.Parse("2006-01-02", endDateString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты конца"})
			return
		}
	}

	// Если start_date и end_date не указаны, вызывайте функцию для получения всех миссий
	if startDateString == "" && endDateString == "" {
		var missions []ds.MissionWithUser
		var err error

		// Если статус не указан, выводим все миссии без фильтрации по статусу
		if missionStatus == "" {
			missions, err = repository.GetAllMissions(user_id, "")
		} else {
			fmt.Println("HERE")
			missions, err = repository.GetAllMissions(user_id, missionStatus)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := gin.H{"missions": missions}
		c.JSON(http.StatusOK, response)
		return
	}

	// Иначе вызывайте функцию для получения миссий с фильтрацией по дате и статусу
	missions, err := repository.GetAllMissionsByDateRange(startDate, endDate, user_id, missionStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{"missions": missions}
	c.JSON(http.StatusOK, response)
}

// @Summary Получение деталей миссии по идентификатору
// @Description Получение информации о миссии и связанных образцах
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор миссии"
// @Security JwtAuth
// @Success 200 {object} string "Детали миссии и образцы"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/mission/{id} [get]
func GetMissionDetailByID(repository *repository.Repository, c *gin.Context, user_id int) {
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

	mission, samples, err := repository.GetMissioninDetailByID(id, user_id)
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

// @Summary Удаление миссии по идентификатору
// @Description Удаление миссии по указанному идентификатору
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор миссии"
// @Security JwtAuth
// @Success 200 {string} string "Миссия успешно удалена"
// @Failure 400 {object} string "Ошибка запроса"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/mission/delete/{id} [delete]
func DeleteMissionByID(repository *repository.Repository, c *gin.Context, user_id int) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка id<0",
		})
		return
	}

	err = repository.DeleteMissionByID(id, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Миссия успешно удалена!")
}

// @Summary Обновление информации о миссии
// @Description Обновление данных миссии по её идентификатору
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор миссии"
// @Param missionData body classes.UpdateMission true "Данные для редактирования миссии"
// @Security JwtAuth
// @Success 200 {object} string "Миссия успешно обновлена"
// @Failure 400 {object} string "Ошибка запроса"
// @Router /api/mission/update/{id} [put]
func UpdateMission(repository *repository.Repository, c *gin.Context, user_id int) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidate, err := repository.GetMissionByID(int(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if Name, nameOk := jsonData["Name"].(string); nameOk {
		candidate.Name = Name
	}

	if formationDateStr, formationOk := jsonData["Formation_date"].(string); formationOk {
		parsedTime, parseErr := time.Parse("2006-01-02", formationDateStr)
		if parseErr != nil {
			// Обработка ошибки парсинга времени
			fmt.Println("Ошибка парсинга времени:", parseErr)
		} else {
			// Присваиваем его полю в вашей структуре
			candidate.Formation_date = &parsedTime
		}
	}

	if completionDateStr, completionOk := jsonData["Completion_date"].(string); completionOk {
		parsedTime, parseErr := time.Parse("2006-01-02", completionDateStr)
		if parseErr != nil {
			// Обработка ошибки парсинга времени
			fmt.Println("Ошибка парсинга Completion_date:", parseErr)
		} else {
			// Присваиваем его полю в вашей структуре
			candidate.Completion_date = &parsedTime
		}
	}

	err = repository.UpdateMission(candidate, id, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Миссия успешно изменена",
	})
}

// @Summary Обновление статуса миссии пользователем
// @Description Обновляет статус миссии по её идентификатору
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор миссии"
// @Param body body classes.UpdateMissionStatus true "Данные для обновления статуса миссии"
// @Security JwtAuth
// @Success 200 {object} string "Статус миссии успешно изменен"
// @Failure 400 {object} string "Ошибка запроса"
// @Router /api/mission/status_by_user/{id} [put]
func UpdateMissionStatusByUser(repository *repository.Repository, c *gin.Context, user_id int) {
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

	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStatus, statusOk := jsonData["Mission_status"].(string)

	// Проверяем валидность нового статуса
	if !statusOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный статус миссии"})
		return
	}

	// Вызываем метод для обновления статуса миссии
	err = repository.UpdateMissionStatusByUser(id, newStatus, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус миссии успешно изменен"})
}

// @Summary Обновление статуса миссии модератором
// @Description Обновляет статус миссии по её идентификатору
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор миссии"
// @Param body body classes.UpdateMissionStatus true "Данные для обновления статуса миссии"
// @Security JwtAuth
// @Success 200 {object} string "Статус миссии успешно изменен"
// @Failure 400 {object} string "Ошибка запроса"
// @Router /api/mission/status_by_moderator/{id} [put]
func UpdateMissionStatusByModerator(repository *repository.Repository, c *gin.Context, user_id int) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка (id < 0)",
		})
		return
	}

	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStatus, statusOk := jsonData["Mission_status"].(string)

	// Проверяем валидность нового статуса
	if !statusOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный статус миссии"})
		return
	}

	// Вызываем метод для обновления статуса миссии
	err = repository.UpdateMissionStatusByModerator(id, newStatus, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус миссии успешно изменен"})
}

// @Summary Удаление образца из последней черновой миссии
// @Description Удаление образца из последней черновой миссии пользователя
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор образца"
// @Security JwtAuth
// @Success 200 {object} map[string]interface{} "Информация о миссии и образцах"
// @Failure 400 {object} string "Ошибка запроса"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/mission/delete_from_last/{id} [delete]
func RemoveSampleFromLastDraftMission(repository *repository.Repository, c *gin.Context, user_id int) {
	sampleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id_sample"})
		return
	}
	mission, samples, err := repository.RemoveSampleFromLastDraftMission(sampleID, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Возвращаем JSON-ответ с деталями миссии и образцами
	c.JSON(http.StatusOK, gin.H{
		"mission": mission,
		"samples": samples,
	})
}
