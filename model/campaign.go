package model

import "gorm.io/gorm"

type Campaign struct {
	gorm.Model

	Content string
}
