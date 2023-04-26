package daoUtil

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	EnablePrepareStmt = ServiceOpt(func(tx *gorm.DB) *gorm.DB {
		return tx.Session(&gorm.Session{
			PrepareStmt: true,
		})
	})
	LockForUpdate = ServiceOpt(func(tx *gorm.DB) *gorm.DB {
		return tx.Clauses(clause.Locking{Strength: "UPDATE"})
	})
	LockForShare = ServiceOpt(func(tx *gorm.DB) *gorm.DB {
		return tx.Clauses(clause.Locking{Strength: "SHARE", Table: clause.Table{Name: clause.CurrentTable}})
	})
	SelectAssociations = ServiceOpt(func(tx *gorm.DB) *gorm.DB {
		return tx.Select(clause.Associations)
	})
	UnScoped = ServiceOpt(func(tx *gorm.DB) *gorm.DB {
		return tx.Unscoped()
	})
)

func TxOpts(DB *gorm.DB, opts ...ServiceOpt) *gorm.DB {
	for _, f := range opts {
		DB = f(DB)
	}
	return DB
}
