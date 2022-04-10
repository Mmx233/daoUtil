package daoUtil

import (
	"context"
	"gorm.io/gorm"
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
	p.fill(tx)
	return p, tx.Error
}

func (s *DaoUtil) NewZeroService(p Service) Service {
	p.fill(s.DB)
	return p
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

func (s *DaoUtil) DefaultExist(tx *gorm.DB, a interface{}) (bool, error) {
	var t bool
	return t, tx.Model(a).Select("1").Where(a).Find(&t).Error
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
	return r, LockForUpdate(tx).Select("1").Model(t).Where(t).Find(&r).Error
}
