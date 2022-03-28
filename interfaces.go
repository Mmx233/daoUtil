package daoUtil

import "gorm.io/gorm"

type Model interface {
	Lock(tx *gorm.DB) (bool, error)
}

type Service interface {
	fill(tx *gorm.DB)
}
