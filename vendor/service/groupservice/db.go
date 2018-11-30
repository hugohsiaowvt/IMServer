package groupservice

import (

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

func InsertGroup(tx *gorm.DB, input *Group) error {
	return tx.Table(TABLE_GROUPS).Create(input).Error
}

func InsertGroupMembers(tx *gorm.DB, input *GroupMember) error {
	return tx.Table(TABLE_GROUP_MEMBERS).Create(input).Error
}