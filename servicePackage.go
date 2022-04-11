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

// HookSuccess 在事务提交前执行，error不为nil时，中止提交并抛出错误
func (a *ServicePackage) HookSuccess(e func() error) {
	context := a.context()
	context.ES = append(context.ES, e)
}

// HookFail 在事务回滚后执行
func (a *ServicePackage) HookFail(e func()) {
	context := a.context()
	context.EF = append(context.EF, e)
}

func (a *ServicePackage) end(success bool) error {
	context := a.context()
	if context.Ended {
		return nil
	}
	context.Ended = true
	var e error
	if success {
		for _, es := range context.ES {
			if e = es(); e != nil {
				return e
			}
		}
		e = a.Tx.Commit().Error
	} else {
		e = a.Tx.Rollback().Error
		for _, ef := range context.EF {
			ef()
		}
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

// Transaction 对子事务的补充，仅在非事务中使用
func (a *ServicePackage) Transaction(e func(tx *gorm.DB) error) error {
	tx := a.Tx.Begin()
	defer tx.Rollback()
	err := e(tx)
	if err == nil {
		return tx.Commit().Error
	}
	return err
}
