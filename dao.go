package daoUtil

import (
	"context"
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
	tx := s.DB.Begin().WithContext(context.WithValue(context.Background(), packageKey, &Context{}))
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

func (s *DaoUtil) DefaultInsert(tx *gorm.DB, a interface{}) error {
	return tx.Create(a).Error
}

func (s *DaoUtil) DefaultDelete(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Delete(a).Error
}

func (s *DaoUtil) DefaultFind(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Find(a).Error
}

func (s *DaoUtil) DefaultExist(tx *gorm.DB, a interface{}) bool {
	var t bool
	tx.Model(a).Select("1").Where(a).Find(&t)
	return t
}

func (s *DaoUtil) DefaultGet(tx *gorm.DB, a interface{}) error {
	return tx.Find(a).Error
}

func (s *DaoUtil) DefaultGetWhitQuery(tx *gorm.DB, a interface{}, t interface{}) error {
	return tx.Where(a).Find(t).Error
}

func (s *DaoUtil) DefaultCounter(tx *gorm.DB, t interface{}) (int64, error) {
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
