package models

import "gorm.io/gorm"

// Todo represents a todo item in the system
// @Description Todo information
type Todo struct {
	gorm.Model
	Title       string `json:"title" example:"Learn Go" binding:"required"`
	Description string `json:"description" example:"Study Go programming language"`
	Completed   bool   `json:"completed" example:"false"`
	UserID      uint   `json:"user_id" example:"1"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
}
