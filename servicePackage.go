package daoUtil

import (
	"gorm.io/gorm"
)

type ServicePackage struct {
	Tx *gorm.DB
}

func (a *ServicePackage) NewFromTx(tx *gorm.DB) {
	a.fill(tx)
}

func (a *ServicePackage) NewFromServ(p Service) {
	a.fill(p.db())
}

func (a *ServicePackage) LockOrRoll(m Model) (bool, error) {
	ok, e := m.Lock(a.Tx)
	if e != nil {
		_ = a.RollBack()
	}
	return ok, e
}

func (a *ServicePackage) WithOpts(opts ...ServiceOpt) *gorm.DB {
	tx := a.Tx
	for _, opt := range opts {
		tx = opt(a.Tx)
	}
	return tx
}

func (a *ServicePackage) db() *gorm.DB {
	return a.Tx
}

func (a *ServicePackage) fill(tx *gorm.DB) {
	a.Tx = tx
}

func (a *ServicePackage) context() *Context {
	return a.Tx.Statement.Context.Value(packageKey).(*Context)
}

// Hook 添加在事务结束时执行的函数
func (a *ServicePackage) Hook(e func(success bool)) {
	context := a.context()
	context.ES = append(context.ES, e)
}

func (a *ServicePackage) end(success bool) error {
	context := a.context()
	if context.Ended {
		return nil
	}
	context.Ended = true
	var e error
	if success {
		e = a.Tx.Commit().Error
	} else {
		e = a.Tx.Rollback().Error
	}
	success = success && e == nil
	for _, e := range context.ES {
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
