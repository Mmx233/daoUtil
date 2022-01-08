package daoUtil

import (
	"gorm.io/gorm"
)

type ServicePackage struct {
	Tx    *gorm.DB
	ended bool
	es    []func(success bool)
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

func (ServicePackage) LockWith(tx *gorm.DB, model modelInterface) (*ServicePackage, error) {
	return &ServicePackage{
		Tx: tx,
	}, model.Lock(tx)
}

// Hook 添加在事务结束时执行的函数
func (a *ServicePackage) Hook(e func(success bool)) {
	a.es = append(a.es, e)
}

func (a *ServicePackage) end(success bool) error {
	if a.ended {
		return nil
	}
	a.ended = true
	var e error
	if success {
		e = a.Tx.Commit().Error
	} else {
		e = a.Tx.Rollback().Error
	}
	success = success && e == nil
	for _, e := range a.es {
		e(success)
	}
	return e
}

// RollBack 回滚，使用行锁时必须defer
func (a *ServicePackage) RollBack() error {
	return a.end(false)
}

func (a *ServicePackage) Commit() error {
	return a.end(true)
}
