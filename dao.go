package daoUtil

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Begin() (*gorm.DB, error) {
	tx := c.DB.Begin()
	return tx, tx.Error
}

func EnablePrepareStmt(tx *gorm.DB) *gorm.DB {
	return tx.Session(&gorm.Session{
		PrepareStmt: true,
	})
}

func DefaultInsert(a interface{}) error {
	return DefaultInsertTx(c.DB, a)
}

func DefaultInsertTx(tx *gorm.DB, a interface{}) error {
	return tx.Create(a).Error
}

func DefaultDelete(a interface{}) error {
	return DefaultDeleteTx(c.DB, a)
}

func DefaultDeleteTx(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Delete(a).Error
}

func DefaultFind(a interface{}) error {
	return DefaultFindTx(c.DB, a)
}

func DefaultFindTx(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Find(a).Error
}

func DefaultExist(a interface{}) bool {
	return DefaultExistTx(c.DB, a)
}

func DefaultExistTx(tx *gorm.DB, a interface{}) bool {
	var t bool
	tx.Model(a).Select("1").Where(a).Find(&t)
	return t
}

func DefaultGet(a interface{}) error {
	return DefaultGetTx(c.DB, a)
}

func DefaultGetTx(tx *gorm.DB, a interface{}) error {
	return tx.Find(a).Error
}

func DefaultGetWhitQuery(a interface{}, t interface{}) error {
	return DefaultGetWhitQueryTx(c.DB, a, t)
}

func DefaultGetWhitQueryTx(tx *gorm.DB, a interface{}, t interface{}) error {
	return tx.Where(a).Find(t).Error
}

func DefaultCounter(t interface{}) (int64, error) {
	return DefaultCounterTx(c.DB, t)
}

func DefaultCounterTx(tx *gorm.DB, t interface{}) (int64, error) {
	var n int64
	return n, tx.Model(t).Where(t).Count(&n).Error
}

func DefaultLock(tx *gorm.DB, t interface{}) (bool, error) {
	var r bool
	return r, tx.Select("1").Model(t).Where(t).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&r).Error
}

func MultiLock(tx *gorm.DB, t interface{}) (bool, error) {
	var r bool
	return r, tx.Model(t).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&r).Error
}
