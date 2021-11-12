package daoUtil

import (
	"gorm.io/gorm"
)

type ServicePackage struct {
	Tx    *gorm.DB
	ended bool
}

type modelInterface interface {
	Lock(tx *gorm.DB) error
}

func (ServicePackage) Begin(model modelInterface) (*ServicePackage, error) {
	tx, e := Begin()
	if e != nil {
		return nil, e
	}
	if model != nil {
		e = model.Lock(tx)
		if e != nil {
			tx.Rollback()
		}
	}
	return &ServicePackage{
		Tx: tx,
	}, e
}

func (ServicePackage) BeginWith(tx *gorm.DB) *ServicePackage {
	if tx == nil {
		tx = c.DB
	}
	return &ServicePackage{
		Tx: tx,
	}
}

func (ServicePackage) LockAndBeginWith(tx *gorm.DB, model modelInterface) (*ServicePackage, error) {
	return &ServicePackage{
		Tx: tx,
	}, model.Lock(tx)
}

func (a *ServicePackage) end(e func() *gorm.DB) error {
	if a.ended {
		return nil
	}
	a.ended = true
	return e().Error
}

// RollBack 回滚，使用行锁时必须defer
func (a *ServicePackage) RollBack() error {
	return a.end(a.Tx.Rollback)
}

func (a *ServicePackage) Commit() error {
	return a.end(a.Tx.Commit)
}
