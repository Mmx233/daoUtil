package daoUtil

import (
	"gorm.io/gorm"
)

type ServicePackage struct {
	Tx *gorm.DB
}

func (ServicePackage) Begin() ServicePackage {
	return ServicePackage{
		Tx: Begin(),
	}
}

func (ServicePackage) BeginWith(tx *gorm.DB) ServicePackage {
	if tx == nil {
		tx = c.DB
	}
	return ServicePackage{
		Tx: tx,
	}
}

func (a *ServicePackage) committed() bool {
	ed, ok := a.Tx.Get("committed")
	return ok && ed.(bool) == true
}

func (a *ServicePackage) mark() {
	a.Tx.Set("committed", true)
}

func (a *ServicePackage) RollBack() {
	if !a.committed() {
		a.mark()
		a.Tx.Rollback()
	}
}

func (a *ServicePackage) Commit() {
	if !a.committed() {
		a.mark()
		a.Tx.Commit()
	}
}
