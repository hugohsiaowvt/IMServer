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
			restfulapi.Error(c, 1, err.Error())
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
			restfulapi.Error(c, 1, "open_id錯誤")
			return
		}
	}

	restfulapi.Success(c, relationship)

}

func InviteFriend(c *gin.Context) {
	var relationship = &RelationShips{}
	var count int = -1

	fromId := c.Param("fromId")

	// 綁定輸入參數
	var input InviteFriendInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Error(c, 1, "缺少參數")
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
			restfulapi.Error(c, 1, "open_id錯誤")
			return
		}
	}

	// 查詢好友邀請清單
	if err := QueryFriendRecord(relationship, &count, smallId, bigId); err != nil {
		if (count != 0) {
			restfulapi.Error(c, 1, err.Error())
			return
		}
	}

	// 添加好友邀請
	relationship.UserOneId = smallId
	relationship.UserTwoId = bigId
	relationship.Time = time.Now().UnixNano()
	relationship.Status = Pending
	relationship.ActionUserId = fromId
	relationship.RequestToken = util.GenerateKey("request_token:")
	if err := InsertRelation(relationship); err != nil {
		restfulapi.Error(c, 1, err.Error())
		return
	}
	restfulapi.Success(c, relationship)
}

func AcceptFriend(c *gin.Context) {

	var relationship = &RelationShips{}
	var count int = -1

	fromId := c.Param("fromId")

	// 綁定輸入參數
	var input AcceptFriendInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Error(c, 1, "缺少參數")
		return
	}

	// 查詢好友邀請清單
	if err := QueryFriendRecordByRequestToken(relationship, &count, input.RequestToken); err != nil {
		restfulapi.Error(c, 1, err.Error())
		return
	}

	// 成為好友
	switch relationship.Status {
	case Pending:
		if (relationship.ActionUserId == fromId) {
			restfulapi.Error(c, 1, "已發送邀請")
			return
		}
		relationship.Time = time.Now().UnixNano()
		relationship.Status = Acception
		relationship.ActionUserId = fromId
		if err := UpdateRelation(relationship); err != nil {
			restfulapi.Error(c, 1, err.Error())
			return
		}
		restfulapi.Success(c, relationship)
		break
	case Acception:
		restfulapi.Error(c, 1, "已是好友")
		break
	default:
		restfulapi.Error(c, 1, "未知錯誤")
		break
	}

}