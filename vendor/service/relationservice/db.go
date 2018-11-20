package relationservice

import (
	"rz/mysql"
)

func QueryFriendRecord(relationship *RelationShips, count *int, userOne, userTwo string) error {
	return mysql.Instance().Table(TABLE_RELATIONSHIPS).
		Where("user_one_id = ? AND user_two_id = ?", userOne, userTwo).Find(relationship).Count(count).Error
}

func QueryFriendRecordByRequestToken(relationship *RelationShips, count *int, token string) error {
	return mysql.Instance().Table(TABLE_RELATIONSHIPS).
		Where("request_token = ?", token).Find(relationship).Count(count).Error
}

func QueryFriendSyncList(relationship *[]RelationShips, count *int, userId, time string) error {
	return mysql.Instance().Debug().Table(TABLE_RELATIONSHIPS).
		Where("(user_one_id = ? OR user_two_id = ?) AND time > ?", userId, userId, time).
		Find(relationship).Count(count).Error
}

func InsertRelation(relationship *RelationShips) error {
	return mysql.Instance().Table(TABLE_RELATIONSHIPS).Create(relationship).Error
}

func UpdateRelation(r *RelationShips) error {
	return mysql.Instance().Table(TABLE_RELATIONSHIPS).
		Where("user_one_id=? AND user_two_id=?", r.UserOneId, r.UserTwoId).
		Updates(r).Error
	//Save(r).Error
}

