package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
type Users struct {
	Id uint `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Role string `gorm:"type:varchar(10)" json:"role"`
	Name string `gorm:"type:varchar(255)" json:"name"`
	Email string `gorm:"type:varchar(50)" json:"email"`
	Password string `gorm:"type:varchar(255)" json:"password"`
	CreatedAt time.Time `json:"cretedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Task []Tasks `gorm:"constraint:OnDelete:CASCADE" json:"tasks, omitempty"`
}

func (u *Users) AfterDeleted(tx *gorm.DB)(err error)  {
	tx.Clauses(clause.Returning{}).Where("user_id = ?", u.Id).Delete(&Tasks{})
	return
}

