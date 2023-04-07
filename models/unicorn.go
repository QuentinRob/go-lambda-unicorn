package models

import (
	"gorm.io/gorm"
)

type Unicorn struct {
	gorm.Model
	Name  string `json:"name"`
	Color string `json:"color"`
}
