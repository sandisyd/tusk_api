package controller

import (
	"net/http"
	"os"
	"strconv"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

// function login
func (u *TaskController) CreateTask(c *gin.Context) {
	task := models.Task{}

	errBind := c.ShouldBindJSON(&task)
	if errBind != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return;
	}

	errDB := u.DB.Create(&task).Error
	
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return;
	}

	c.JSON(http.StatusOK, task)
}

// function delete
func (u *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")
	task := models.Task{}

	// cek data di db

	if err := u.DB.First(&task,id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error() })
		return;
	}

	//  delete data
	errDB := u.DB.Delete(&models.Task{},id).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error() })
		return;
	}

	// cek attachment
	if task.Attachment != "" {
		os.Remove("attachment/" + task.Attachment)
	}

	c.JSON(http.StatusOK, gin.H{"message":"Success deleted task"})
}

// update attachment dan submited date
func (u *TaskController) PatchData(c *gin.Context) {
	id := c.Param("id")
	task := models.Task{}
	submitedDate := c.PostForm("submitedDate")
	file, errFile := c.FormFile("attachment")

	// cek error file
	if errFile != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errFile.Error() })
		return;
	}

	// cek data di db

	if err := u.DB.First(&task,id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data not found" })
		return;
	}

	// remove old attachment
	attachment := task.Attachment
	fileInfo, _ := os.Stat("attachment/" + attachment)
	if fileInfo != nil {
		// found
		os.Remove("attachment/" + attachment)
	}

	// create new file attachment
	attachment = file.Filename
	errSave := c.SaveUploadedFile(file, "attachment/" + attachment)

	if errSave != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errSave.Error() })
		return;
	}
	//  delete data
	errDB := u.DB.Where("id=?", id).Updates(models.Task{
		Status: "Review",
		SubmitDate: submitedDate,
		Attachment: attachment,
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, gin.H{"message":"Success Update data"})
}

// rejected task
func (u *TaskController) RejectedTask(c *gin.Context) {
	id := c.Param("id")
	task := models.Task{}
	reason := c.PostForm("reason")
	rejectedDate := c.PostForm("rejectedDate")


	// cek data di db

	if err := u.DB.First(&task,id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data not found" })
		return;
	}
	//  delete data
	errDB := u.DB.Where("id=?", id).Updates(models.Task{
		Status: "Rejected",
		Reason: reason,
		RejectedDate: rejectedDate,
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, gin.H{"message":"You're task was been rejected"})
}

// Fix to Queue
func (u *TaskController) FixedTask(c *gin.Context) {
	id := c.Param("id")
	revision, errConv := strconv.Atoi(c.PostForm("revision"))

	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errConv.Error() })
		return;
	}
	// cek data di db

	if err := u.DB.First(&models.Task{},id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data not found" })
		return;
	}
	//  delete data
	errDB := u.DB.Where("id=?", id).Updates(models.Task{
		Status: "Queue",
		Revision: int8(revision),
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, gin.H{"message":"Change Status to Queue"})
}

// approved task
func (u *TaskController) ApprovedTask(c *gin.Context) {
	id := c.Param("id")
	approvedDate := c.PostForm("approvedDate")
	// cek data di db

	if err := u.DB.First(&models.Task{},id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data not found" })
		return;
	}
	//  delete data
	errDB := u.DB.Where("id=?", id).Updates(models.Task{
		Status: "Approved",
		ApprovedDate: approvedDate,
	}).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, gin.H{"message":"Change Status to Approved"})
}

// function find task by id
func (u *TaskController) TaskbyId(c *gin.Context) {
	task := models.Task{}
	id := c.Param("id")
	
	if errData := u.DB.First(&models.Task{}, id).Error; errData != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found" })
		return;
	}
	errDB := u.DB.Preload("User").Find(&task, id).Error
	if errDB != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, task)
}

// task review
func (u *TaskController) TaskReview(c *gin.Context) {
	tasks := []models.Task{}
	errDB := u.DB.Preload("User").Where("status=?", "Review").Order("submit_date ASC").Limit(2).Find(&tasks).Error
	if errDB != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, tasks)
}

// task progres
func (u *TaskController) TaskProgress(c *gin.Context) {
	tasks := []models.Task{}
	userId := c.Param("userId")
	errDB := u.DB.Where("(status!=? AND user_id=?) OR (revision!=? AND user_id=?)", "Queue", userId, 0, userId).Order("updated_at DESC").Limit(5).Find(&tasks).Error
	if errDB != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, tasks)
}
// task statistik
func (u *TaskController) TaskStatistik(c *gin.Context) {
	userId := c.Param("userId")

	stat_progress := []map[string]interface{}{}
	errDB := u.DB.Model(models.Task{}).Select("status, count(status) as total").Where("user_id=?", userId).Group("status").Find(&stat_progress).Error
	if errDB != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, stat_progress)
}

// task by user and status
func (u *TaskController) TaskStatus(c *gin.Context) {
	tasks := []models.Task{}
	userId := c.Param("userId")
	status := c.Param("status")
	errDB := u.DB.Where("(status=? AND user_id=?)", status, userId).Find(&tasks).Error
	if errDB != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, tasks)
}



