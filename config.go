package daoUtil

import "gorm.io/gorm"

type Config struct {
	DB *gorm.DB
}

func New(config *Config) *DaoUtil {
	return &DaoUtil{
		DB: config.DB,
	}
}
