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

	name := c.DefaultQuery("name", "")

	sample, err := repository.GetAllSamples(name)
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
	// user_id, err := strconv.Atoi(c.Param("user_id"))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }
	user_id := 2
	var sample *ds.Samples

	if err := c.BindJSON(&sample); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(sample)

	if sample.Name == "" || sample.Type == "" || sample.Sol_Sealed == 0 || sample.Current_Location == "" || sample.Sample_status == "" {
		c.JSON(http.StatusBadRequest, "Ошибка, указаны неправильные данные для образца")
	}

	if err := repository.CreateSample(sample, user_id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Образец успешно создан!")
}

func DeleteSampleByID(repository *repository.Repository, c *gin.Context) {
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
	fmt.Println(id)

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

func UpdateSample(repository *repository.Repository, c *gin.Context) {
	// user_id, err := strconv.Atoi(c.Param("user_id"))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }
	user_id := 2
	Id_sample, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(Id_sample)

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

	Name, nameOk := jsonData["Name"].(string)
	Type, typeOk := jsonData["Type"].(string)

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
	err = repository.UpdateSample(candidate, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Образец успешно обновлен",
	})
}

func AddSampleToMission(repository *repository.Repository, c *gin.Context) {
	// user_id, err := strconv.Atoi(c.Param("user_id"))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }
	user_id := 7
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
