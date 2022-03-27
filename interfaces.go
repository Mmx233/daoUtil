package daoUtil

import "gorm.io/gorm"

type Service interface {
	fill(tx *gorm.DB)
}
