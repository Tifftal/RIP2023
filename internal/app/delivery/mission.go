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

func GetAllMissiions(repository *repository.Repository, c *gin.Context) {
	startDateString := c.Query("start_date")
	endDateString := c.Query("end_date")

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
		mission, err := repository.GetAllMissions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, mission)
		return
	}

	// Иначе вызывайте функцию для получения миссий с фильтрацией по дате
	mission, err := repository.GetAllMissionsByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, mission)
}

func DeleteMissionByID(repository *repository.Repository, c *gin.Context) {
	user_id := 2
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(id)

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

func UpdateMission(repository *repository.Repository, c *gin.Context) {
	// user_id, err := strconv.Atoi(c.Param("user_id"))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }
	user_id := 2

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

	if moderator, moderatorOk := jsonData["Moderator_id"].(float64); moderatorOk {
		candidate.Moderator_id = int(moderator)
	}

	if formationDateStr, formationOk := jsonData["Formation_date"].(string); formationOk {
		parsedTime, parseErr := time.Parse("2006-01-02", formationDateStr)
		if parseErr != nil {
			// Обработка ошибки парсинга времени
			fmt.Println("Ошибка парсинга времени:", parseErr)
		} else {
			// Присваиваем его полю в вашей структуре
			candidate.Formation_date = parsedTime
		}
	}

	if creationDateStr, creationOk := jsonData["Creation_date"].(string); creationOk {
		parsedTime, parseErr := time.Parse("2006-01-02", creationDateStr)
		if parseErr != nil {
			// Обработка ошибки парсинга времени
			fmt.Println("Ошибка парсинга Creation_date:", parseErr)
		} else {
			// Присваиваем его полю в вашей структуре
			candidate.Creation_date = parsedTime
		}
	}

	if completionDateStr, completionOk := jsonData["Completion_date"].(string); completionOk {
		parsedTime, parseErr := time.Parse("2006-01-02", completionDateStr)
		if parseErr != nil {
			// Обработка ошибки парсинга времени
			fmt.Println("Ошибка парсинга Completion_date:", parseErr)
		} else {
			// Присваиваем его полю в вашей структуре
			candidate.Completion_date = parsedTime
		}
	}

	fmt.Println(candidate)

	err = repository.UpdateMission(candidate, id, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Миссия успешно изменена",
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

// UpdateMissionStatus обновляет статус миссии.
func UpdateMissionStatusByUser(repository *repository.Repository, c *gin.Context) {
	user_id := 1
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
func UpdateMissionStatusByModerator(repository *repository.Repository, c *gin.Context) {
	user_id := 2
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

func RemoveSampleFromMission(repository *repository.Repository, c *gin.Context) {
	user_id := 2
	// Получаем Id_mission и Id_sample из параметров запроса
	missionID, err := strconv.Atoi(c.Param("mission_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id_mission"})
		return
	}

	sampleID, err := strconv.Atoi(c.Param("sample_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id_sample"})
		return
	}

	mission, samples, err := repository.RemoveSampleFromMission(uint(missionID), uint(sampleID), user_id)
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

func RemoveSampleFromLastDraftMission(repository *repository.Repository, c *gin.Context) {
	user_id := 1
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
