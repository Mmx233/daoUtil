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
)
