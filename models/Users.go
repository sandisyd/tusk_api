package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
type User struct {
	Id int `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Role string `gorm:"type:varchar(10)" json:"role"`
	Name string `gorm:"type:varchar(255)" json:"name"`
	Email string `gorm:"type:varchar(50)" json:"email"`
	Password string `gorm:"type:varchar(255)" json:"password"`
	CreatedAt time.Time `json:"cretedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Task []Task `gorm:"constraint:OnDelete:CASCADE" json:"tasks, omitempty"`
}

func (u *User) AfterDeleted(tx *gorm.DB)(err error)  {
	tx.Clauses(clause.Returning{}).Where("user_id = ?", u.Id).Delete(&Task{})
	return
}

