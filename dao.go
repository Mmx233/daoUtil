package daoUtil

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DaoUtil struct {
	DB *gorm.DB
}

func (s *DaoUtil) Begin() (*gorm.DB, error) {
	tx := s.DB.Begin()
	return tx, tx.Error
}

func (s *DaoUtil) BeginService(p Service) (Service, error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return p, tx.Error
	}
	p.fill(tx)
	return p, nil
}

func (s *DaoUtil) NewServicePackage(tx *gorm.DB) *ServicePackage {
	if tx == nil {
		tx = s.DB
	}
	return &ServicePackage{Tx: tx}
}

func (s *DaoUtil) EnablePrepareStmt(tx *gorm.DB) *gorm.DB {
	return tx.Session(&gorm.Session{
		PrepareStmt: true,
	})
}

func (s *DaoUtil) DefaultInsert(a interface{}) error {
	return s.DefaultInsertTx(s.DB, a)
}

func (s *DaoUtil) DefaultInsertTx(tx *gorm.DB, a interface{}) error {
	return tx.Create(a).Error
}

func (s *DaoUtil) DefaultDelete(a interface{}) error {
	return s.DefaultDeleteTx(s.DB, a)
}

func (s *DaoUtil) DefaultDeleteTx(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Delete(a).Error
}

func (s *DaoUtil) DefaultFind(a interface{}) error {
	return s.DefaultFindTx(s.DB, a)
}

func (s *DaoUtil) DefaultFindTx(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Find(a).Error
}

func (s *DaoUtil) DefaultExist(a interface{}) bool {
	return s.DefaultExistTx(s.DB, a)
}

func (s *DaoUtil) DefaultExistTx(tx *gorm.DB, a interface{}) bool {
	var t bool
	tx.Model(a).Select("1").Where(a).Find(&t)
	return t
}

func (s *DaoUtil) DefaultGet(a interface{}) error {
	return s.DefaultGetTx(s.DB, a)
}

func (s *DaoUtil) DefaultGetTx(tx *gorm.DB, a interface{}) error {
	return tx.Find(a).Error
}

func (s *DaoUtil) DefaultGetWhitQuery(a interface{}, t interface{}) error {
	return s.DefaultGetWhitQueryTx(s.DB, a, t)
}

func (s *DaoUtil) DefaultGetWhitQueryTx(tx *gorm.DB, a interface{}, t interface{}) error {
	return tx.Where(a).Find(t).Error
}

func (s *DaoUtil) DefaultCounter(t interface{}) (int64, error) {
	return s.DefaultCounterTx(s.DB, t)
}

func (s *DaoUtil) DefaultCounterTx(tx *gorm.DB, t interface{}) (int64, error) {
	var n int64
	return n, tx.Model(t).Where(t).Count(&n).Error
}

func (s *DaoUtil) DefaultLock(tx *gorm.DB, t interface{}) (bool, error) {
	var r bool
	return r, tx.Select("1").Model(t).Where(t).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&r).Error
}

func (s *DaoUtil) MultiLock(tx *gorm.DB, t interface{}) (bool, error) {
	var r bool
	return r, tx.Model(t).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&r).Error
}
