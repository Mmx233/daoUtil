package daoUtil

import (
	"gorm.io/gorm"
	"sync"
)

var lock sync.Mutex

type Config struct {
	DB *gorm.DB
	//强制非nil事务串行化
	Serializable bool
}

var c *Config

func Init(config *Config) {
	c = config
}

func Begin() *gorm.DB {
	return c.DB.Begin()
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

type ServicePackage struct {
	Tx      *gorm.DB
	locking bool
}

func (ServicePackage) Begin() ServicePackage {
	if c.Serializable {
		lock.Lock()
	}
	return ServicePackage{
		Tx:      Begin(),
		locking: c.Serializable,
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
		if a.locking {
			lock.Unlock()
		}
	}
}

func (a *ServicePackage) Commit() {
	if !a.committed() {
		a.mark()
		a.Tx.Commit()
		if a.locking {
			lock.Unlock()
		}
	}
}
