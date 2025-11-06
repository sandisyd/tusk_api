package controller

import (
	"net/http"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (u *UserController) Login(c *gin.Context) {
	user := models.User{}

	errBind := c.ShouldBindJSON(&user)
	if errBind != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return;
	}

	password := user.Password

	errDB := u.DB.Where("email = ? ", user.Email).Take(&user).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return;
	}

	errHash := bcrypt.CompareHashAndPassword( []byte(user.Password) ,[]byte(password))
	if errHash != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or Password is wrong"})
		return;
	}

	c.JSON(http.StatusOK, user)
}
func (u *UserController) Register(c *gin.Context) {
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
func (u *UserController) Delete(c *gin.Context) {
	id := c.Param("id")
	

	errDB := u.DB.Delete(&models.User{},id).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error() })
		return;
	}

	c.JSON(http.StatusOK, gin.H{"message":"Success deleted user"})
}