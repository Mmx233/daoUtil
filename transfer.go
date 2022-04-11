package daoUtil

import "gorm.io/gorm"

const packageKey = "DaoServicePackage"

type Context struct {
	Ended bool
	ES    []func() error
	EF    []func()
}

type Model interface {
	Lock(tx *gorm.DB) (bool, error)
}

type Service interface {
	db() *gorm.DB
	fill(tx *gorm.DB)
	LockOrRoll(m Model) (bool, error)
	Hook(e func(success bool))
	RollBack() error
	Commit() error
}

type ServiceOpt func(tx *gorm.DB) *gorm.DB
