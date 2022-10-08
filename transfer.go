package daoUtil

import "gorm.io/gorm"

type ServiceOpt func(tx *gorm.DB) *gorm.DB
