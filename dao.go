package daoUtil

import "gorm.io/gorm"

var db *gorm.DB

func Init(DB *gorm.DB) {
	db = DB
}

func Begin() *gorm.DB {
	return db.Begin()
}

func DefaultInsert(a interface{}) error {
	return DefaultInsertTx(db, a)
}

func DefaultInsertTx(tx *gorm.DB, a interface{}) error {
	return tx.Create(a).Error
}

func DefaultDelete(a interface{}) error {
	return DefaultDeleteTx(db, a)
}

func DefaultDeleteTx(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Delete(a).Error
}

func DefaultFind(a interface{}) error {
	return DefaultFindTx(db, a)
}

func DefaultFindTx(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Find(a).Error
}

func DefaultExist(a interface{}) bool {
	return DefaultExistTx(db, a)
}

func DefaultExistTx(tx *gorm.DB, a interface{}) bool {
	var t struct {
		ID uint
	}
	tx.Model(a).Where(a).Find(&t)
	return t.ID != 0
}

func DefaultGet(a interface{}) error {
	return DefaultGetTx(db, a)
}

func DefaultGetTx(tx *gorm.DB, a interface{}) error {
	return tx.Find(a).Error
}

func DefaultGetWhitQuery(a interface{}, t interface{}) error {
	return DefaultGetWhitQueryTx(db, a, t)
}

func DefaultGetWhitQueryTx(tx *gorm.DB, a interface{}, t interface{}) error {
	return tx.Where(a).Find(t).Error
}

func DefaultCounter(t interface{}) (int64, error) {
	return DefaultCounterTx(db, t)
}

func DefaultCounterTx(tx *gorm.DB, t interface{}) (int64, error) {
	var n int64
	return n, tx.Where(t).Count(&n).Error
}

type ServicePackage struct {
	Tx        *gorm.DB
	committed bool
}

func (ServicePackage) Begin() ServicePackage {
	return ServicePackage{
		Tx: Begin(),
	}
}

func (a *ServicePackage) RollBack() {
	if !a.committed {
		a.committed = true
		a.Tx.Rollback()
	}
}

func (a *ServicePackage) Commit() {
	if !a.committed {
		a.committed = true
		a.Tx.Commit()
	}
}
