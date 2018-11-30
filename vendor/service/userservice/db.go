package userservice

import (
	"rz/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestDB(user *User, input *User, count *int) error {
	return mysql.Instance().Table(TABLE_USERS).
		Where(input).Find(user).Count(count).Error
}

func QueryUser(user, input *User, count *int) error {
	return mysql.Instance().Table(TABLE_USERS).
		Where(input).Find(user).Count(count).Error
}

func QueryUserByOpenId(user *User, openId string) error {
	return mysql.Instance().Table(TABLE_USERS).
		Where("open_id=?", openId).Find(user).Error
}

func QueryUserPublicInfoByOpenIds(user *[]User, count *int, openIds []string) error {
	return mysql.Instance().Table(TABLE_USERS).
		Select("open_id, user_name, zone, mobile, sex, email, status").
		Where("open_id IN (?)", openIds).
		Find(user).Count(count).Error
}

func InsertUsers(tx *gorm.DB, input *User) error {
	return tx.Table(TABLE_USERS).Create(input).Error
}

func UpdatePassword(openId, password string) error {
	return mysql.Instance().Table(TABLE_USERS).
		Where("open_id=?", openId).
		UpdateColumn("password", password).Error
}