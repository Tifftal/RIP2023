package delivery

import (
	"MSRM/internal/app/ds"
	"MSRM/internal/app/repository"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка образцов с фильтрацией по имени и типу горной породы
// @Description Получение списка образцов с возможностью фильтрации по имени и/или типу горной породы
// @Accept json
// @Produce json
// @Tags Samples
// @Param name query string false "Фильтр по имени образца"
// @Param rockType query string false "Фильтр по типу горной породы"
// @Success 200 {object} []ds.Samples "Список образцов"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/sample [get]
func GetAllSamples(repository *repository.Repository, c *gin.Context, user_id int) {
	var sample []ds.Samples

	name := c.DefaultQuery("name", "")
	rockType := c.DefaultQuery("rockType", "")

	sample, draftMission_id, err := repository.GetAllSamples(name, rockType, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"samples": sample, "draftMission_id": draftMission_id})
}

// @Summary Получение образца по ID
// @Description Получение информации о конкретном образце по его идентификатору
// @Accept json
// @Produce json
// @Tags Samples
// @Param id path integer true "Идентификатор образца"
// @Success 200 {object} ds.Samples "Информация об образце"
// @Failure 400 {object} string "Неверный запрос или ошибка получения образца"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/sample/{id} [get]
func GetSampleByID(repository *repository.Repository, c *gin.Context) {
	var sample *ds.Samples
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

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

// @Summary Создание образца
// @Description Создание нового образца с указанными данными
// @Accept json
// @Security User_id
// @Produce json
// @Tags Samples
// @Param sampleData body ds.Samples true "Данные для создания образца"
// @Success 200 {object} string "Образец успешно создан"
// @Failure 400 {object} string "Неверный запрос или ошибка создания образца"
// @Router /api/sample/create [post]
// @Security JwtAuth
func CreateSample(repository *repository.Repository, c *gin.Context, user_id int) {
	var sample ds.Samples

	if err := c.BindJSON(&sample); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sample.Name == "" || sample.Type == "" || sample.Sol_Sealed == 0 || sample.Current_Location == "" || sample.Sample_status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка, указаны неправильные данные для образца"})
		return
	}

	createdSampleID, err := repository.CreateSample(&sample, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Образец успешно создан!", "id": createdSampleID})
}

// @Summary Удаление образца по ID
// @Description Удаление образца с указанным идентификатором пользователя и ID образца
// @Accept json
// @Produce json
// @Tags Samples
// @Param id path integer true "Идентификатор образца"
// @Success 200 {string} string "Образец успешно удален"
// @Failure 400 {object} string "Неверный запрос или ошибка удаления образца"
// @Router /api/sample/delete/{id} [delete]
// @Security JwtAuth
func DeleteSampleByID(repository *repository.Repository, c *gin.Context, user_id int) {
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

	err = repository.DeleteSampleByID(id, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Образец успешно удален!")
}

// @Summary Обновление образца
// @Description Обновление данных образца с указанным идентификатором пользователя и ID образца
// @Accept json
// @Produce json
// @Tags Samples
// @Param id path integer true "Идентификатор образца"
// @Param sampleData body classes.Sample_Update true "Данные для обновления образца"
// @Success 200 {string} string "Образец успешно обновлен"
// @Failure 400 {object} string "Неверный запрос или ошибка обновления образца"
// @Router /api/sample/update/{id} [put]
// @Security JwtAuth
func UpdateSample(repository *repository.Repository, c *gin.Context, user_id int) {
	Id_sample, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if Id_sample < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка id<0",
		})
		return
	}

	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Извлечение значений из JSON
	Name, nameOk := jsonData["Name"].(string)
	Type, typeOk := jsonData["Type"].(string)
	Rock_Type, rockTypeOk := jsonData["Rock_Type"].(string)
	Current_Location, locationOk := jsonData["Current_Location"].(string)
	Sample_status, statusOk := jsonData["Sample_status"].(string)

	// Получение объекта образца из БД
	candidate, err := repository.GetSampleByID(int(Id_sample))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Обновление значений, если они указаны
	if nameOk {
		candidate.Name = Name
	}
	if typeOk {
		candidate.Type = Type
	}
	if rockTypeOk {
		candidate.Rock_Type = Rock_Type
	}
	if locationOk {
		candidate.Current_Location = Current_Location
	}
	if statusOk {
		candidate.Sample_status = Sample_status
	}

	// Сохранение обновленного образца в БД
	err = repository.UpdateSample(candidate, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Образец успешно обновлен",
	})
}

// @Summary Добавление образца в миссию
// @Description Добавление образца в последнюю миссию с Mission_status "Draft"
// @Accept json
// @Produce json
// @Tags Samples
// @Param id path integer true "Идентификатор образца"
// @Success 200 {string} string "Образец успешно добавлен в миссию"
// @Failure 400 {object} string "Неверный запрос или ошибка добавления образца в миссию"
// @Failure 401 {object} string "Неавторизованный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/sample/to_mission/{id} [put]
// @Security JwtAuth
func AddSampleToMission(repository *repository.Repository, c *gin.Context, user_id int) {
	// Получаем Id_sample из параметра запроса
	sampleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id_sample"})
		return
	}
	// Вызываем функцию для добавления образца в последнюю миссию с Mission_status = "Draft"
	mission, samples, err := repository.AddSampleToLastDraftMission(sampleID, user_id)
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

// @Summary Добавление изображения к образцу
// @Description Загружает изображение и добавляет его к указанному образцу
// @Accept mpfd
// @Produce json
// @Tags Samples
// @Param id path integer true "Идентификатор образца"
// @Param image formData file true "Изображение для загрузки"
// @Success 200 {string} string "Изображение усспешно загружено"
// @Failure 400 {object} string "Неверный запрос или ошибка загрузки изображения"
// @Failure 401 {object} string "Неавторизованный запрос"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/sample/{id}/image [post]
// @Security JwtAuth
func AddImageToSample(repository *repository.Repository, c *gin.Context, user_id int) {
	sampleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопсутимый ИД багажа"})
		return
	}

	// Чтение изображения из запроса
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
		return
	}

	// Чтение содержимого изображения в байтах
	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть изображение"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать изображение в байтах"})
		return
	}
	// Получение Content-Type из заголовков запроса
	contentType := image.Header.Get("Content-Type")

	// Вызов функции репозитория для добавления изображения
	err = repository.AddSampleImage(sampleID, imageBytes, contentType, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение усспешно загружено"})
}
