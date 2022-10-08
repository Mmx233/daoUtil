package daoUtil

import (
	"gorm.io/gorm"
)

func Insert(tx *gorm.DB, a interface{}) error {
	return tx.Create(a).Error
}

func Delete(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Delete(a).Error
}

func Find(tx *gorm.DB, a interface{}) error {
	return tx.Where(a).Find(a).Error
}

func Exist(tx *gorm.DB, a interface{}) (bool, error) {
	var t bool
	return t, tx.Model(a).Select("1").Where(a).Find(&t).Error
}

func Get(tx *gorm.DB, a interface{}) error {
	return tx.Find(a).Error
}

func GetWhitQuery(tx *gorm.DB, a interface{}, t interface{}) error {
	return tx.Where(a).Find(t).Error
}

func Counter(tx *gorm.DB, t interface{}) (int64, error) {
	var n int64
	return n, tx.Model(t).Where(t).Count(&n).Error
}
