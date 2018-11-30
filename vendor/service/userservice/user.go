package userservice

import (

	"fmt"
	"errors"

	"rz/util"
	"rz/mysql"
	"rz/redis"
	"rz/restfulapi"
	"service/smsservice"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"

)

// PUBLIC METHOD

func Test(c *gin.Context) {
	//user, query := &User{}, &User{}
	//var count int
	//query.OpenId = "875c5080f1f7e5eb7cbd4c893290994f"
	//query.Password = "ss"
	v, err := redis.Instance().LRange("mylist", 0, -1)
	fmt.Println(v, err)

	restfulapi.Response(c, 0, v, "")

}

func CreateUser(c *gin.Context) {

	res := LoginResponse{}
	var loginInput LoginInput

	vcode := c.Param("code")

	// 綁定輸入參數
	var input RegisterInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, restfulapi.ERROR_MISSING_PARAMETER_MSG)
		return
	}

	// 輸入參數確認
	if err := input.CheckInput(); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, err.Error())
		return
	}

	// 驗證碼確認
	if code, err := smsservice.Check(input.Zone + input.Mobile, vcode); err != nil {
		restfulapi.Response(c, code, struct{}{}, err.Error())
		return
	}

	// 將輸入參數轉成User
	user := input.ToModel()

	// 查詢用戶存不存在
	var count int = -1
	if err := QueryUser(&User{}, user, &count); err != nil {
		if count != 0 {
			restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
			return
		}
	} else if count == 1 {
		restfulapi.Response(c, restfulapi.ERROR_PHONE_ISEXIST_CODE, struct{}{}, restfulapi.ERROR_PHONE_ISEXIST_MSG)
		return
	}

	// 新增一個Open id
	user.OpenId = util.MD5(OPEN_ID_PREFIX + input.Zone + input.Mobile + OPEN_ID_KEY)
	user.Status = USER_STATUS_OK

	// 開啟Transactions
	tx := mysql.Instance().Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	// 新增用戶
	if err := InsertUsers(tx, user); err != nil {
		tx.Rollback()
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct {}{}, err.Error())
		return
	}

	// 提交
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_TRANSACTION_ERROR_CODE, struct {}{}, restfulapi.ERROR_MYSQL_TRANSACTION_ERROR_MSG)
		return
	}

	// 登入
	loginInput.Zone = input.Zone
	loginInput.Mobile = input.Mobile
	u := User{}
	if token, im_token, code, err := login(&u, loginInput); err != nil {
		restfulapi.Response(c, code, struct{}{}, err.Error())
		return
	} else {
		res.Token = token
		res.IMToken = im_token
		res.User = u
	}

	restfulapi.Response(c, restfulapi.SUCCESS_CODE, res, "")
}

func Login(c *gin.Context) {

	res := LoginResponse{}
	user := User{}

	// 綁定輸入參數
	var input LoginInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, restfulapi.ERROR_MISSING_PARAMETER_MSG)
		return
	}

	// 輸入參數確認
	if err := input.CheckInput(); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, err.Error())
		return
	}

	// 登入
	if token, im_token, code, err := login(&user, input); err != nil {
		restfulapi.Response(c, code, struct{}{}, err.Error())
		return
	} else {
		res.Token = token
		res.IMToken = im_token
		res.User = user
	}

	restfulapi.Response(c, restfulapi.SUCCESS_CODE, res, "")
}

func SetPassword(c *gin.Context) {

	user := &User{}

	// 綁定輸入參數
	var input UpdatePasswordInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, restfulapi.ERROR_MISSING_PARAMETER_MSG)
		return
	}

	// 輸入參數確認
	if err := input.CheckInput(); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, err.Error())
		return
	}

	// 查詢用戶
	if err := QueryUserByOpenId(user, input.OpenId); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
		return
	}

	// 檢查密碼
	if user.Password != "" {
		if user.Password != input.OldPassword {
			restfulapi.Response(c, restfulapi.ERROR_WRONG_PASSWORD_CODE, struct{}{}, restfulapi.ERROR_WRONG_PASSWORD_MSG)
			return
		}
	}

	// 更新密碼
	if err := UpdatePassword(input.OpenId, input.NewPassword); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MYSQL_ERROR_CODE, struct{}{}, err.Error())
		return
	}

	restfulapi.Response(c, restfulapi.SUCCESS_CODE, input, "")
}

// PRIVATE METHOD

// 確認用戶存不存在
func QueryUserExistByMobile(zone, mobile string) (User, bool) {
	user := User{}
	var count int
	query := &User{
		Zone: zone,
		Mobile: mobile,
	}
	QueryUser(&user, query, &count)
	if count != 0 {
		return user, true
	}
	return user, false
}

// 確認用戶存不存在
func QueryUserExistByOpenId(openId string) (User, bool) {
	user := User{}
	var count int
	query := &User{
		OpenId: openId,
	}
	QueryUser(&user, query, &count)
	if count != 0 {
		return user, true
	}
	return user, false
}

// 登入
func login(user *User, input LoginInput) (string, string, int32, error) {

	// 查詢用戶
	var count int
	query := &User{}
	query.Zone = input.Zone
	query.Mobile = input.Mobile
	query.Password = input.Password
	if err := QueryUser(user, query, &count); err != nil {
		if (count != 0) {
			return "", "", restfulapi.ERROR_MYSQL_ERROR_CODE, err
		} else {
			return "", "", restfulapi.ERROR_WRONG_INPUT_CODE, errors.New(restfulapi.ERROR_WRONG_INPUT_MSG)
		}
	}

	if user.Status != USER_STATUS_OK {
		return "", "", restfulapi.ERROR_USER_LOCK_CODE, errors.New(restfulapi.ERROR_USER_LOCK_MSG)
	}

	// 製作登入token
	key := LOGIN_TOKEN_REDIS_PREFIX + user.OpenId
	token := uuid.New()
	if _, err := redis.Instance().Set(key, token); err != nil {
		return "", "", restfulapi.ERROR_REDIS_ERROR_CODE, err
	}

	// 取得im token
	userId := user.OpenId
	if imToken, code, tokenError := getIMToken(userId); tokenError != nil {
		return "", "", code, tokenError
	} else {
		return token, imToken, restfulapi.SUCCESS_CODE, nil
	}

}

// 取得IM Token
func getIMToken(openId string) (string, int32, error) {
	token := uuid.New()
	key := IM_TOKEN_REDIS_PREFIX + openId
	if _, err := redis.Instance().Set(key, token); err != nil {
		return "", restfulapi.ERROR_REDIS_ERROR_CODE, err
	}
	return token, restfulapi.SUCCESS_CODE, nil
}

// 確認IM Token
func CheckIMToken(openId, token string) bool {
	key := IM_TOKEN_REDIS_PREFIX + openId
	if v, err := redis.Instance().Get(key); err != nil {
		return false
	} else {
		if v == token {
			return true
		} else {
			return false
		}
	}
}