package daoUtil

import "gorm.io/gorm"

type Config struct {
	DB *gorm.DB
}

var c *Config

func Init(config *Config) {
	c = config
}
