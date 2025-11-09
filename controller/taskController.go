package controller

import (
	"net/http"
	"os"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

// function create users
func (u *TaskController) Register(c *gin.Context) {
	user := models.User{}

	errBind := c.ShouldBindJSON(&user)
	if errBind != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return;
	}

	password := user.Password

	emailExist := u.DB.Where("email = ? ", user.Email).First(&user).RowsAffected != 0
	if emailExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exist"})
		return;
	}

	hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	user.Password = string(hashedPassword)
	user.Role = "Employee"

	errDB := u.DB.Create(&user).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, user)
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

// function list users employee
func (u *TaskController) GetListUsersEmployee(c *gin.Context) {
	users := []models.User{}
	

	errDB := u.DB.Where("role = ?", "Employee").Find(&users).Error
	if errDB != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, users)
}