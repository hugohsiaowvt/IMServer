package groupservice

import (

	"rz/mysql"
	"rz/restfulapi"
	"rz/util"

	"github.com/gin-gonic/gin"

)

func CreateGroup(c *gin.Context) {

	res := CreateGroupResponse{}
	creator := c.Param("creator")

	// 綁定輸入參數
	var input CreateGroupInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, err.Error())
		return
	}

	// 開啟Transactions
	tx := mysql.Instance().Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	// 將輸入參數轉成Group
	group := input.ToModel()

	// 產生group id
	groupId := util.MD5(GROUP_ID_PREFIX + util.GetRandomCode(10))
	group.GroupId = groupId

	members := removeDumpMembers(input.Members)

	// 新增群組
	if err := InsertGroup(tx, group); err != nil {
		tx.Rollback()
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct {}{}, err.Error())
		return
	}

	// 先新增群主
	groupMember := &GroupMember{}
	groupMember.GroupId = groupId
	groupMember.OpenId = creator
	groupMember.Role = ROLE_GROUP_MANAGER
	groupMember.Status = GROUP_MEMBER_STATUS_OK

	if err := InsertGroupMembers(tx, groupMember); err != nil {
		tx.Rollback()
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct {}{}, err.Error())
		return
	}

	// 新增群組人員
	groupMember.Role = ROLE_MEMBER
	for _, m := range members {
		groupMember.OpenId = m.OpenId
		if err := InsertGroupMembers(tx, groupMember); err != nil {
			tx.Rollback()
			restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct {}{}, err.Error())
			return
		}
	}

	// 提交
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_TRANSACTION_ERROR_CODE, struct {}{}, restfulapi.ERROR_MYSQL_TRANSACTION_ERROR_MSG)
		return
	}

	res.GroupId = groupId

	restfulapi.Response(c, restfulapi.SUCCESS_CODE, res, "")

}




// 移除重複的群組成員
func removeDumpMembers(members []GroupMemberInput) []GroupMemberInput {
	var membersMap = map[string]GroupMemberInput{}
	if members!=nil && len(members) > 0 {
		for _ ,member := range members {
			groupMemberInput := membersMap[member.OpenId]
			if groupMemberInput.OpenId == "" {
				membersMap[member.OpenId] = member
			}
		}
	}
	var newGroupMembers = make([]GroupMemberInput,0)
	for _ ,member := range membersMap {
		newGroupMembers = append(newGroupMembers, member)
	}
	return newGroupMembers
}
