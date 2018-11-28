package relationservice

import (

	"strings"
	"time"

	"rz/util"
	"rz/restfulapi"
	"service/userservice"

	"github.com/gin-gonic/gin"

)

func SyncFriendList(c *gin.Context) {

	var relationship = []RelationShips{}
	var count int = -1
	var friendIds []string
	var friendList []userservice.User

	userId := c.Param("userId")
	time := c.Param("time")

	if err := QueryFriendSyncList(&relationship, &count, userId, time); err != nil {
		if (count != 0) {
			restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
			return
		}
	}

	for _, elem := range relationship {
		if (elem.UserOneId == userId) {
			friendIds = append(friendIds, elem.UserTwoId)
		} else {
			friendIds = append(friendIds, elem.UserOneId)
		}
	}

	if err := userservice.QueryUserPublicInfoByOpenIds(&friendList, &count, friendIds); err != nil {
		if (count != 0) {
			restfulapi.Response(c, restfulapi.ERROR_WRONG_OPEN_ID_CODE, struct{}{}, restfulapi.ERROR_WRONG_OPEN_ID_MSG)
			return
		}
	}

	restfulapi.Response(c, restfulapi.SUCCESS_CODE, relationship, "")

}

func InviteFriend(c *gin.Context) {
	var relationship = &RelationShips{}
	var count int = -1

	fromId := c.Param("fromId")

	// 綁定輸入參數
	var input InviteFriendInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, restfulapi.ERROR_MISSING_PARAMETER_MSG)
		return
	}

	// 判斷ID大小
	smallId := fromId
	bigId := input.ToId
	if (strings.Compare(fromId, input.ToId) == 1) {
		smallId	= input.ToId
		bigId = fromId
	}

	//查詢用戶資料
	var users []userservice.User
	var openIds = []string {
		smallId,
		bigId,
	}
	if err := userservice.QueryUserPublicInfoByOpenIds(&users, &count, openIds); err != nil {
		if (count != 2) {
			restfulapi.Response(c, restfulapi.ERROR_WRONG_OPEN_ID_CODE, struct{}{}, restfulapi.ERROR_WRONG_OPEN_ID_MSG)
			return
		}
	}

	// 查詢好友邀請清單
	if err := QueryFriendRecord(relationship, &count, smallId, bigId); err != nil {
		if (count != 0) {
			restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
			return
		}
	}

	// 已有發送邀請
	if count == 1 {
		// 之前被忽略過
		if relationship.Status == None {
			relationship.ActionUserId = fromId
			relationship.Status = Pending
			if err := UpdateRelation(relationship); err != nil {
				restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
				return
			}
			restfulapi.Response(c, restfulapi.SUCCESS_CODE, relationship, "")
		} else if fromId == relationship.ActionUserId {
			restfulapi.Response(c, restfulapi.ERROR_SENDED_FRIEND_REQUEST_CODE, struct{}{}, restfulapi.ERROR_SENDED_FRIEND_REQUEST_MSG)
		} else if input.ToId == relationship.ActionUserId {
			restfulapi.Response(c, restfulapi.ERROR_UNHANDLED_FRIEND_REQUEST_CODE, struct{}{}, restfulapi.ERROR_UNHANDLED_FRIEND_REQUEST_MSG)
		} else {
			restfulapi.Response(c, restfulapi.ERROR_UNKNOW_CODE, struct{}{}, restfulapi.ERROR_UNKNOW_MSG)
		}
		return
	}

	// 添加好友邀請
	relationship.UserOneId = smallId
	relationship.UserTwoId = bigId
	relationship.Time = time.Now().UnixNano()
	relationship.Status = Pending
	relationship.ActionUserId = fromId
	relationship.RequestToken = util.GenerateKey(REQUEST_TOKEN_PREFIX)
	if err := InsertRelation(relationship); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
		return
	}
	restfulapi.Response(c, restfulapi.SUCCESS_CODE, relationship, "")
}

func AcceptFriend(c *gin.Context) {

	var relationship = &RelationShips{}
	var count int = -1

	fromId := c.Param("fromId")

	// 綁定輸入參數
	var input QueryFriendRecordInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, restfulapi.ERROR_MISSING_PARAMETER_MSG)
		return
	}

	// 查詢好友邀請清單
	if err := QueryFriendRecordByRequestToken(relationship, &count, input.RequestToken); err != nil {
		if (count != 0) {
			restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
			return
		}
	}

	// 沒有邀請資訊
	if count == 0 {
		restfulapi.Response(c, restfulapi.ERROR_UNKNOW_CODE, struct{}{}, restfulapi.ERROR_UNKNOW_MSG)
		return
	}

	// 成為好友
	switch relationship.Status {
	case None:
		restfulapi.Response(c, restfulapi.ERROR_NONE_FRIEND_REQUEST_CODE, struct{}{}, restfulapi.ERROR_NONE_FRIEND_REQUEST_MSG)
		break
	case Pending:

		// 發送與接受邀請是同個人
		if (relationship.ActionUserId == fromId) {
			restfulapi.Response(c, restfulapi.ERROR_SENDED_FRIEND_REQUEST_CODE, struct{}{}, restfulapi.ERROR_SENDED_FRIEND_REQUEST_MSG)
			return
		}

		relationship.Time = time.Now().UnixNano()
		relationship.Status = Acception
		relationship.ActionUserId = fromId
		if err := UpdateRelation(relationship); err != nil {
			restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
			return
		}
		restfulapi.Response(c, restfulapi.SUCCESS_CODE, relationship, "")
		break
	case Acception:
		restfulapi.Response(c, restfulapi.ERROR_FRIEND_ALREADY_CODE, struct{}{}, restfulapi.ERROR_FRIEND_ALREADY_MSG)
		break
	default:
		restfulapi.Response(c, restfulapi.ERROR_UNKNOW_CODE, struct{}{}, restfulapi.ERROR_UNKNOW_MSG)
		break
	}

}

func IgnoreFriend(c *gin.Context) {

	var relationship = &RelationShips{}
	var count int = -1

	fromId := c.Param("fromId")

	// 綁定輸入參數
	var input QueryFriendRecordInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, restfulapi.ERROR_MISSING_PARAMETER_MSG)
		return
	}

	// 查詢好友邀請清單
	if err := QueryFriendRecordByRequestToken(relationship, &count, input.RequestToken); err != nil {
		if (count != 0) {
			restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
			return
		}
	}

	// 沒有邀請資訊
	if count == 0 {
		restfulapi.Response(c, restfulapi.ERROR_UNKNOW_CODE, struct{}{}, restfulapi.ERROR_UNKNOW_MSG)
		return
	} else {
		// 發送與接受邀請是同個人
		if fromId == relationship.ActionUserId {
			restfulapi.Response(c, restfulapi.ERROR_SENDED_FRIEND_REQUEST_CODE, struct{}{}, restfulapi.ERROR_SENDED_FRIEND_REQUEST_MSG)
			return
		}
	}

	// 修改status到none
	relationship.Time = time.Now().UnixNano()
	relationship.Status = None
	relationship.ActionUserId = fromId
	if err := UpdateRelation(relationship); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
		return
	}
	restfulapi.Response(c, restfulapi.SUCCESS_CODE, relationship, "")
}